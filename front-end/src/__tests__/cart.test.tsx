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

// Helper: add a pizza to the cart by clicking on the PizzaItem
async function addPizzaToCart(
  user: ReturnType<typeof userEvent.setup>,
  pizzaName: string
) {
  const matches = screen.getAllByText(pizzaName);
  await user.click(matches[0]);
}

describe('Cart Functionality', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('shows empty cart message when no items are added', () => {
    render(<Home />);
    expect(screen.getByText(/your cart is empty/i)).toBeInTheDocument();
  });

  it('adds an item to the cart when a pizza is clicked', async () => {
    const user = userEvent.setup();
    render(<Home />);

    await addPizzaToCart(user, 'Pepperoni');

    // Cart should show the pizza count summary
    expect(screen.getByText(/1 pizzas in the cart/i)).toBeInTheDocument();
  });

  it('adds multiple different pizza types to the cart', async () => {
    const user = userEvent.setup();
    render(<Home />);

    await addPizzaToCart(user, 'Margherita');
    await addPizzaToCart(user, 'Hawaiian');

    // Should show 2 pizzas in cart
    expect(screen.getByText(/2 pizzas in the cart/i)).toBeInTheDocument();
  });

  it('increments quantity when same pizza is clicked again', async () => {
    const user = userEvent.setup();
    render(<Home />);

    // Click Margherita 3 times
    await addPizzaToCart(user, 'Margherita');
    await addPizzaToCart(user, 'Margherita');
    await addPizzaToCart(user, 'Margherita');

    // Should show 3 pizzas in cart
    expect(screen.getByText(/3 pizzas in the cart/i)).toBeInTheDocument();
    // Price should be 3 × $10 = $30
    expect(screen.getByText('$30')).toBeInTheDocument();
  });

  it('removes an item from the cart via delete button', async () => {
    const user = userEvent.setup();
    render(<Home />);

    // Add Margherita (qty 1 → delete button visible)
    await addPizzaToCart(user, 'Margherita');
    // Add Pepperoni
    await addPizzaToCart(user, 'Pepperoni');

    expect(screen.getByText(/2 pizzas in the cart/i)).toBeInTheDocument();

    // At qty 1, the minus button shows as "Delete item"
    const deleteButtons = screen.getAllByRole('button', { name: /delete item/i });
    // Click first delete button (Margherita)
    await user.click(deleteButtons[0]);

    // Should show 1 pizza in cart (Pepperoni remains)
    expect(screen.getByText(/1 pizzas in the cart/i)).toBeInTheDocument();
  });

  it('places order with all cart items', async () => {
    const user = userEvent.setup();
    createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'cart-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    // Add Margherita x2
    await addPizzaToCart(user, 'Margherita');
    await addPizzaToCart(user, 'Margherita');

    // Add Pepperoni x1
    await addPizzaToCart(user, 'Pepperoni');

    // Place order
    const placeOrderButton = screen.getByRole('button', { name: /place order/i });
    await user.click(placeOrderButton);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith('/api/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          orderItems: [
            { pizzaType: 'Margherita', quantity: 2 },
            { pizzaType: 'Pepperoni', quantity: 1 },
          ],
        }),
      });
    });
  });

  it('disables Place Order when cart is empty', () => {
    render(<Home />);
    const placeOrderButton = screen.getByRole('button', { name: /place order/i });
    expect(placeOrderButton).toBeDisabled();
  });

  it('clears the cart after placing order', async () => {
    const user = userEvent.setup();
    const { mockWs } = createMockWebSocket();

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'clear-cart-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    await addPizzaToCart(user, 'Margherita');

    // Place order
    const placeOrderButton = screen.getByRole('button', { name: /place order/i });
    await user.click(placeOrderButton);

    // Wait for order to be placed — auto-switches to "Your Orders" tab
    await waitFor(() => {
      expect(screen.getByTestId('order-id')).toBeInTheDocument();
    });

    // "New Order" tab is disabled until DELIVERED — simulate delivery
    await act(async () => {
      if (mockWs.onmessage) {
        mockWs.onmessage(new MessageEvent('message', {
          data: JSON.stringify({
            orderId: 'clear-cart-id',
            status: 'DELIVERED',
            source: 'delivery',
            timestamp: '2026-01-26T10:10:00Z',
          }),
        }));
      }
    });

    // Switch back to "New Order" tab to verify cart is cleared
    const newOrderTab = screen.getByRole('tab', { name: /new order/i });
    await user.click(newOrderTab);

    // Cart should be empty
    expect(screen.getByText(/your cart is empty/i)).toBeInTheDocument();
  });
});
