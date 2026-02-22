#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="pizza-vibe"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Pizza Vibe - KIND Setup ==="
echo "Project root: $PROJECT_ROOT"
echo ""

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
# 3. Build agent services (Maven)
# -------------------------------------------------------
echo "--- Step 3: Building agent services with Maven ---"
AGENTS=(pizza-mcp cooking-agent delivery-agent store-mgmt-agent)
for agent in "${AGENTS[@]}"; do
  echo "Building agents/$agent ..."
  (cd "$PROJECT_ROOT/agents/$agent" && ./mvnw clean package -DskipTests)
done
echo ""

# -------------------------------------------------------
# 4. Build Docker images
# -------------------------------------------------------
echo "--- Step 4: Building Docker images ---"
cd "$PROJECT_ROOT"

# Go services
docker build -t pizza-vibe-store:latest -f store/Dockerfile .
docker build -t pizza-vibe-front-end:latest -f front-end/Dockerfile ./front-end
docker build -t pizza-vibe-inventory:latest -f inventory/Dockerfile .
docker build -t pizza-vibe-oven:latest -f oven/Dockerfile .
docker build -t pizza-vibe-kitchen:latest -f kitchen/Dockerfile .
docker build -t pizza-vibe-delivery:latest -f delivery/Dockerfile .
docker build -t pizza-vibe-bikes:latest -f bikes/Dockerfile .
docker build -t pizza-vibe-drinks-stock:latest -f drinks-stock/Dockerfile .

# Java/Quarkus agent services
docker build -t pizza-vibe-pizza-mcp:latest -f agents/pizza-mcp/src/main/docker/Dockerfile.jvm ./agents/pizza-mcp
docker build -t pizza-vibe-cooking-agent:latest -f agents/cooking-agent/src/main/docker/Dockerfile.jvm ./agents/cooking-agent
docker build -t pizza-vibe-delivery-agent:latest -f agents/delivery-agent/src/main/docker/Dockerfile.jvm ./agents/delivery-agent
docker build -t pizza-vibe-store-mgmt-agent:latest -f agents/store-mgmt-agent/src/main/docker/Dockerfile.jvm ./agents/store-mgmt-agent
echo ""

# -------------------------------------------------------
# 5. Load images into KIND
# -------------------------------------------------------
echo "--- Step 5: Loading images into KIND cluster ---"
IMAGES=(
  pizza-vibe-store
  pizza-vibe-front-end
  pizza-vibe-inventory
  pizza-vibe-oven
  pizza-vibe-kitchen
  pizza-vibe-delivery
  pizza-vibe-bikes
  pizza-vibe-drinks-stock
  pizza-vibe-pizza-mcp
  pizza-vibe-cooking-agent
  pizza-vibe-delivery-agent
  pizza-vibe-store-mgmt-agent
)
for image in "${IMAGES[@]}"; do
  echo "Loading $image:latest ..."
  kind load docker-image "$image:latest" --name "$CLUSTER_NAME"
done
echo ""

# -------------------------------------------------------
# 6. Create secrets
# -------------------------------------------------------
echo "--- Step 6: Creating secrets ---"
if [ -z "${OPENAI_API_KEY:-}" ]; then
  echo "WARNING: OPENAI_API_KEY environment variable is not set."
  echo "Set it and re-run, or create the secret manually:"
  echo "  kubectl create secret generic openai-secret --from-literal=api-key=<YOUR_KEY>"
else
  kubectl create secret generic openai-secret \
    --from-literal=api-key="$OPENAI_API_KEY" \
    --dry-run=client -o yaml | kubectl apply -f -
  echo "Secret 'openai-secret' created."
fi
echo ""

# -------------------------------------------------------
# 7. Deploy the application
# -------------------------------------------------------
echo "--- Step 7: Deploying application ---"
kubectl apply -f "$PROJECT_ROOT/k8s/"
echo ""

echo "--- Waiting for pods to be ready ---"
kubectl wait --for=condition=Ready pods --all --timeout=120s || true
kubectl get pods
echo ""

echo "=== Setup complete ==="
echo "Access the application with:"
echo "  kubectl port-forward svc/front-end 3000:3000"
echo "Then open http://localhost:3000"
