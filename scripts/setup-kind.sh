#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="pizza-vibe"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Pizza Vibe - KIND Setup ==="
echo "Project root: $PROJECT_ROOT"
echo ""

# -------------------------------------------------------
# Pre-flight: Require ANTHROPIC_API_KEY
# -------------------------------------------------------
if [ -z "${ANTHROPIC_API_KEY:-}" ]; then
  echo "ERROR: ANTHROPIC_API_KEY environment variable is not set."
  echo "The agent services require an Anthropic API key to function."
  echo ""
  echo "Set it before running this script:"
  echo "  export ANTHROPIC_API_KEY=<YOUR_KEY>"
  exit 1
fi

# -------------------------------------------------------
# 1. Create KIND cluster
# -------------------------------------------------------
echo "--- Step 1: Creating KIND cluster '$CLUSTER_NAME' ---"
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "Cluster '$CLUSTER_NAME' already exists, skipping creation."
else
  kind create cluster --name "$CLUSTER_NAME"
fi
kubectl cluster-info --context "kind-${CLUSTER_NAME}"
echo ""

# -------------------------------------------------------
# 2. Install Dapr
# -------------------------------------------------------
echo "--- Step 2: Installing Dapr ---"
helm repo add dapr https://dapr.github.io/helm-charts/ 2>/dev/null || true
helm repo update
if helm status dapr -n dapr-system &>/dev/null; then
  echo "Dapr is already installed, skipping."
else
  helm install dapr dapr/dapr --namespace dapr-system --create-namespace --wait
fi
echo "Dapr pods:"
kubectl get pods -n dapr-system
echo ""

# -------------------------------------------------------
# 3. Install PostgreSQL
# -------------------------------------------------------
echo "--- Step 3: Installing PostgreSQL ---"
helm repo add bitnami https://charts.bitnami.com/bitnami 2>/dev/null || true
helm repo update
if helm status postgresql &>/dev/null; then
  echo "PostgreSQL is already installed, skipping."
else
  helm install postgresql bitnami/postgresql \
    --set auth.postgresPassword=postgres \
    --set auth.database=dapr_store \
    --wait
fi
echo "PostgreSQL pods:"
kubectl get pods -l app.kubernetes.io/name=postgresql
echo ""

# -------------------------------------------------------
# 4. Create secrets
# -------------------------------------------------------
echo "--- Step 4: Creating secrets ---"
kubectl create secret generic anthropic-secret \
  --from-literal=api-key="$ANTHROPIC_API_KEY" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "Secret 'anthropic-secret' created."
echo ""

# -------------------------------------------------------
# 5. Install Jaeger
# -------------------------------------------------------
echo "--- Step 5: Installing Jaeger ---"
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts 2>/dev/null || true
helm repo update
if helm status jaeger &>/dev/null; then
  echo "Jaeger is already installed, skipping."
else
  helm install jaeger jaegertracing/jaeger --version 3.4.1 -f "$PROJECT_ROOT/jaeger/values.yaml" --wait
fi
echo "Jaeger pods:"
kubectl get pods -l app.kubernetes.io/name=jaeger
echo ""

# -------------------------------------------------------
# 6. Create OpenTelemetry namespace and Dash0 secrets
# -------------------------------------------------------
echo "--- Step 6: Creating OpenTelemetry namespace and Dash0 secrets ---"
kubectl create namespace opentelemetry --dry-run=client -o yaml | kubectl apply -f -

if [ -n "${DASH0_AUTH_TOKEN:-}" ]; then
  DASH0_ENDPOINT_OTLP_GRPC_HOSTNAME="${DASH0_ENDPOINT_OTLP_GRPC_HOSTNAME:-ingress.eu-west-1.aws.dash0.com}"
  DASH0_ENDPOINT_OTLP_GRPC_PORT="${DASH0_ENDPOINT_OTLP_GRPC_PORT:-4317}"

  kubectl create secret generic dash0-secrets \
    --from-literal=dash0-authorization-token="$DASH0_AUTH_TOKEN" \
    --from-literal=dash0-grpc-hostname="$DASH0_ENDPOINT_OTLP_GRPC_HOSTNAME" \
    --from-literal=dash0-grpc-port="$DASH0_ENDPOINT_OTLP_GRPC_PORT" \
    --namespace=opentelemetry \
    --dry-run=client -o yaml | kubectl apply -f -
  echo "Secret 'dash0-secrets' created in opentelemetry namespace."
else
  echo "DASH0_AUTH_TOKEN not set, skipping Dash0 secrets. Only Jaeger will receive telemetry."
fi
echo ""

# -------------------------------------------------------
# 7. Install OpenTelemetry Collector
# -------------------------------------------------------
echo "--- Step 7: Installing OpenTelemetry Collector ---"
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts 2>/dev/null || true
helm repo update
if helm status otel-collector -n opentelemetry &>/dev/null; then
  echo "OpenTelemetry Collector is already installed, skipping."
else
  helm install otel-collector open-telemetry/opentelemetry-collector \
    --namespace opentelemetry \
    -f "$PROJECT_ROOT/collector/config.yaml" \
    --wait
fi
echo "OpenTelemetry Collector pods:"
kubectl get pods -n opentelemetry -l app.kubernetes.io/name=opentelemetry-collector
echo ""

# -------------------------------------------------------
# 8. Install cert-manager and OpenTelemetry Operator
# -------------------------------------------------------
echo "--- Step 8: Installing cert-manager ---"
helm repo add jetstack https://charts.jetstack.io --force-update
helm repo update
if helm status cert-manager -n cert-manager &>/dev/null; then
  echo "cert-manager is already installed, skipping."
else
  helm upgrade --install cert-manager jetstack/cert-manager \
    --namespace cert-manager --create-namespace \
    --set crds.enabled=true \
    --wait
fi
echo "cert-manager pods:"
kubectl get pods -n cert-manager
echo ""

echo "--- Step 8b: Installing OpenTelemetry Operator ---"
if helm status opentelemetry-operator -n opentelemetry &>/dev/null; then
  echo "OpenTelemetry Operator is already installed, skipping."
else
  helm upgrade --install opentelemetry-operator open-telemetry/opentelemetry-operator \
    --namespace opentelemetry \
    --set manager.extraArgs='{--enable-go-instrumentation}' \
    --wait
fi
echo "OpenTelemetry Operator pods:"
kubectl get pods -n opentelemetry -l app.kubernetes.io/name=opentelemetry-operator
echo ""

# -------------------------------------------------------
# 9. Apply OpenTelemetry Instrumentation resource
# -------------------------------------------------------
echo "--- Step 9: Applying OpenTelemetry Instrumentation ---"
kubectl apply -f "$PROJECT_ROOT/instrumentation/instrumentation.yaml"
echo "Instrumentation resource applied."
echo ""

# -------------------------------------------------------
# 10-14. Build, load images, and deploy
# -------------------------------------------------------
export CLUSTER_NAME
"$SCRIPT_DIR/rebuild-and-deploy.sh"

echo "Access the application with:"
echo "  kubectl port-forward svc/store 8080:8080"
echo "Then open http://localhost:8080"
echo ""
echo "To access Jaeger UI:"
echo "  kubectl port-forward svc/jaeger-query 16686"
echo "Then open http://localhost:16686"
echo ""
echo "To access PostgreSQL from outside the cluster:"
echo "  kubectl port-forward svc/postgresql 5432:5432"
echo "  psql postgres://postgres:postgres@localhost:5432/dapr_store"
