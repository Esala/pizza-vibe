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

// Helper: add an item to the cart
async function addItemToCart(
  user: ReturnType<typeof userEvent.setup>,
  pizzaType: string,
  quantity: number
) {
  const pizzaSelect = screen.getByLabelText(/pizza type/i);
  const quantityInput = screen.getByLabelText(/quantity/i);
  const addToCartButton = screen.getByRole('button', { name: /add to cart/i });

  await user.selectOptions(pizzaSelect, pizzaType);
  await user.tripleClick(quantityInput);
  await user.keyboard(String(quantity));
  await user.click(addToCartButton);
}

describe('Home Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the page title', () => {
    render(<Home />);
    expect(screen.getByRole('heading', { name: /pizza vibe/i })).toBeInTheDocument();
  });

  it('displays an order form with pizza type selection', () => {
    render(<Home />);
    expect(screen.getByLabelText(/pizza type/i)).toBeInTheDocument();
  });

  it('displays an order form with quantity input', () => {
    render(<Home />);
    expect(screen.getByLabelText(/quantity/i)).toBeInTheDocument();
  });

  it('displays a submit button to place the order', () => {
    render(<Home />);
    expect(screen.getByRole('button', { name: /place order/i })).toBeInTheDocument();
  });

  it('submits the order to the store service when form is submitted', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 2);

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

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 2);

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/order placed successfully/i)).toBeInTheDocument();
    });

    // Verify order ID is displayed
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

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 1);

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

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 1);

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    // Wait for order to be placed (WebSocket connected first, then order placed)
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

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 1);

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

  it('displays error message when order fails', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      status: 500,
    });

    render(<Home />);

    // Add item to cart first
    await addItemToCart(user, 'Margherita', 1);

    const submitButton = screen.getByRole('button', { name: /place order/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/failed to place order/i)).toBeInTheDocument();
    });
  });
});
