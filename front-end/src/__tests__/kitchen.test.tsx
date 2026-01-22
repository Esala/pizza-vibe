import { render, screen, waitFor } from '@testing-library/react';
import KitchenPage from '@/app/kitchen/page';

// Mock fetch
global.fetch = jest.fn();

describe('Kitchen Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the kitchen page title', () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orders: [] }),
    });

    render(<KitchenPage />);
    expect(screen.getByRole('heading', { name: /kitchen/i })).toBeInTheDocument();
  });

  it('displays a list of orders from the kitchen service', async () => {
    const mockOrders = [
      {
        orderId: '123e4567-e89b-12d3-a456-426614174000',
        orderItems: [
          { pizzaType: 'Margherita', quantity: 2 },
          { pizzaType: 'Pepperoni', quantity: 1 },
        ],
        status: 'cooking',
      },
      {
        orderId: '123e4567-e89b-12d3-a456-426614174001',
        orderItems: [
          { pizzaType: 'Hawaiian', quantity: 1 },
        ],
        status: 'ready',
      },
    ];

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orders: mockOrders }),
    });

    render(<KitchenPage />);

    await waitFor(() => {
      expect(screen.getByText('123e4567-e89b-12d3-a456-426614174000')).toBeInTheDocument();
    });

    expect(screen.getByText('123e4567-e89b-12d3-a456-426614174001')).toBeInTheDocument();
    expect(screen.getByText(/Margherita/)).toBeInTheDocument();
    expect(screen.getByText(/Pepperoni/)).toBeInTheDocument();
    expect(screen.getByText(/Hawaiian/)).toBeInTheDocument();
  });

  it('displays order status for each order', async () => {
    const mockOrders = [
      {
        orderId: '123e4567-e89b-12d3-a456-426614174000',
        orderItems: [{ pizzaType: 'Margherita', quantity: 1 }],
        status: 'cooking',
      },
    ];

    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orders: mockOrders }),
    });

    render(<KitchenPage />);

    await waitFor(() => {
      expect(screen.getByText(/cooking/i)).toBeInTheDocument();
    });
  });

  it('displays loading state while fetching orders', () => {
    (global.fetch as jest.Mock).mockImplementationOnce(
      () => new Promise(() => {}) // Never resolves
    );

    render(<KitchenPage />);

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('displays error message when fetch fails', async () => {
    (global.fetch as jest.Mock).mockRejectedValueOnce(new Error('Failed to fetch'));

    render(<KitchenPage />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });

  it('displays empty state when no orders are available', async () => {
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orders: [] }),
    });

    render(<KitchenPage />);

    await waitFor(() => {
      expect(screen.getByText(/no orders/i)).toBeInTheDocument();
    });
  });
});
