import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Home from '@/app/page';

// Mock fetch
global.fetch = jest.fn();

// Helper: create a mock WebSocket that auto-fires onopen after construction
function createMockWebSocket() {
  const mockWs = {
    close: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    readyState: 1,
    onopen: null as ((ev: Event) => void) | null,
    onmessage: null as ((ev: MessageEvent) => void) | null,
    onclose: null as ((ev: CloseEvent) => void) | null,
    onerror: null as ((ev: Event) => void) | null,
  };
  const MockWebSocket = jest.fn(() => {
    // Simulate async connection: fire onopen on next microtask
    Promise.resolve().then(() => {
      if (mockWs.onopen) {
        mockWs.onopen(new Event('open'));
      }
    });
    return mockWs;
  });
  (global as unknown as Record<string, unknown>).WebSocket = MockWebSocket;
  return { mockWs, MockWebSocket };
}

// Helper: add a pizza to the cart by clicking on it
async function addPizzaToCart(
  user: ReturnType<typeof userEvent.setup>,
  pizzaName: string
) {
  const matches = screen.getAllByText(pizzaName);
  await user.click(matches[0]);
}

describe('Home Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders all pizza items', () => {
    render(<Home />);
    expect(screen.getByText('Margherita')).toBeInTheDocument();
    expect(screen.getByText('Pepperoni')).toBeInTheDocument();
    expect(screen.getByText('Hawaiian')).toBeInTheDocument();
    expect(screen.getByText('Vegan')).toBeInTheDocument();
  });

  it('displays a Place Order button', () => {
    render(<Home />);
    expect(screen.getByRole('button', { name: /place order/i })).toBeInTheDocument();
  });

  it('displays the Place Order button as disabled when cart is empty', () => {
    render(<Home />);
    const placeOrderButton = screen.getByRole('button', { name: /place order/i });
    expect(placeOrderButton).toBeDisabled();
  });

  it('submits the order to the store service when Place Order is clicked', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    // Add Margherita twice (quantity 2)
    await addPizzaToCart(user, 'Margherita');
    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith('/api/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          orderItems: [{ pizzaType: 'Margherita', quantity: 2 }],
        }),
      });
    });
  });

  it('displays success message and order ID after successful order', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    // After success, page auto-switches to "Your Orders" tab
    await waitFor(() => {
      expect(screen.getByText(/order placed successfully/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/test-order-id/i)).toBeInTheDocument();
  });

  it('connects to WebSocket with unique client ID before placing order', async () => {
    const user = userEvent.setup();
    const { MockWebSocket } = createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'ws-test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(MockWebSocket).toHaveBeenCalled();
    });

    // Verify WebSocket URL connects directly to the store service with clientId
    const wsUrl = MockWebSocket.mock.calls[0][0] as string;
    expect(wsUrl).toContain('/ws?clientId=');
    expect(wsUrl).toMatch(/^wss?:\/\//);

    // Verify fetch was called after WebSocket connected
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled();
    });
  });

  it('displays WebSocket connection indicator', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'indicator-test-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    // Wait for order to be placed — auto-switches to "Your Orders" tab
    await waitFor(() => {
      expect(screen.getByText(/indicator-test-id/i)).toBeInTheDocument();
    });

    expect(screen.getByTestId('ws-status')).toBeInTheDocument();
    expect(screen.getByText(/connected/i)).toBeInTheDocument();
  });

  it('displays incoming WebSocket events in a table', async () => {
    const user = userEvent.setup();
    const { mockWs } = createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'table-test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/table-test-order-id/i)).toBeInTheDocument();
    });

    // Simulate receiving WebSocket events
    await act(async () => {
      if (mockWs.onmessage) {
        mockWs.onmessage(new MessageEvent('message', {
          data: JSON.stringify({
            orderId: 'table-test-order-id',
            status: 'cooking',
            source: 'kitchen',
            timestamp: '2026-01-26T10:00:00Z',
          }),
        }));
      }
    });

    // Verify events table is displayed
    expect(screen.getByTestId('events-table')).toBeInTheDocument();

    // Verify table headers
    expect(screen.getByText('Order ID')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Source')).toBeInTheDocument();
    expect(screen.getByText('Timestamp')).toBeInTheDocument();

    // Verify event data in table
    expect(screen.getByText('cooking')).toBeInTheDocument();
    expect(screen.getByText('kitchen')).toBeInTheDocument();
    expect(screen.getByText('2026-01-26T10:00:00Z')).toBeInTheDocument();

    // Simulate a second event
    await act(async () => {
      if (mockWs.onmessage) {
        mockWs.onmessage(new MessageEvent('message', {
          data: JSON.stringify({
            orderId: 'table-test-order-id',
            status: 'COOKED',
            source: 'kitchen',
            timestamp: '2026-01-26T10:05:00Z',
          }),
        }));
      }
    });

    // Verify both events are in the table
    expect(screen.getByText('cooking')).toBeInTheDocument();
    expect(screen.getByText('COOKED')).toBeInTheDocument();

    // Verify table rows (header + 2 data rows)
    const rows = screen.getByTestId('events-table').querySelectorAll('tr');
    expect(rows.length).toBe(3); // 1 header + 2 data rows
  });

  it('disables New Order tab after placing order and re-enables after delivery', async () => {
    const user = userEvent.setup();
    const { mockWs } = createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'disable-tab-test-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    // Wait for order to be placed — auto-switches to "Your Orders" tab
    await waitFor(() => {
      expect(screen.getByTestId('order-id')).toBeInTheDocument();
    });

    // "New Order" tab should be disabled
    const newOrderTab = screen.getByRole('tab', { name: /new order/i });
    expect(newOrderTab).toBeDisabled();

    // Simulate DELIVERED event
    await act(async () => {
      if (mockWs.onmessage) {
        mockWs.onmessage(new MessageEvent('message', {
          data: JSON.stringify({
            orderId: 'disable-tab-test-id',
            status: 'DELIVERED',
            source: 'delivery',
            timestamp: '2026-01-26T10:10:00Z',
          }),
        }));
      }
    });

    // "New Order" tab should be re-enabled
    const newOrderTabAfter = screen.getByRole('tab', { name: /new order/i });
    expect(newOrderTabAfter).not.toBeDisabled();
  });

  it('displays error message when order fails', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      status: 500,
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/failed to place order/i)).toBeInTheDocument();
    });
  });
});
