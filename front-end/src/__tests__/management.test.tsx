import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import ManagementPage from '@/app/management/page';

// Mock fetch
global.fetch = jest.fn();

// Mock data matching store/models.go types
const mockOrders = [
  {
    orderId: '123e4567-e89b-12d3-a456-426614174000',
    orderItems: [{ pizzaType: 'Margherita', quantity: 2 }],
    orderData: 'Test order 1',
    orderStatus: 'pending',
  },
  {
    orderId: '223e4567-e89b-12d3-a456-426614174001',
    orderItems: [{ pizzaType: 'Pepperoni', quantity: 1 }],
    orderData: 'Test order 2',
    orderStatus: 'COOKED',
  },
];

const mockEvents = [
  { orderId: '123e4567-e89b-12d3-a456-426614174000', status: 'cooking', source: 'kitchen' },
  { orderId: '123e4567-e89b-12d3-a456-426614174000', status: 'in oven', source: 'kitchen' },
  { orderId: '123e4567-e89b-12d3-a456-426614174000', status: 'DONE', source: 'kitchen' },
];

describe('Management Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the management page title', () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);
    expect(screen.getByRole('heading', { name: /management/i })).toBeInTheDocument();
  });

  it('displays status for all three services', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      expect(screen.getByText(/store service/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/kitchen service/i)).toBeInTheDocument();
    expect(screen.getByText(/delivery service/i)).toBeInTheDocument();
  });

  it('shows healthy status when services are up', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      const healthyStatuses = screen.getAllByText(/healthy/i);
      expect(healthyStatuses.length).toBeGreaterThan(0);
    });
  });

  it('shows unhealthy status when services are down', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: false,
          status: 500,
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      const unhealthyStatuses = screen.getAllByText(/unhealthy/i);
      expect(unhealthyStatuses.length).toBeGreaterThan(0);
    });
  });

  it('displays loading state while checking services', () => {
    (global.fetch as jest.Mock).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    render(<ManagementPage />);

    expect(screen.getByText(/checking services/i)).toBeInTheDocument();
  });

  it('handles fetch errors gracefully', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Network error'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      const unhealthyStatuses = screen.getAllByText(/unhealthy/i);
      expect(unhealthyStatuses.length).toBeGreaterThan(0);
    });
  });

  it('displays orders list section', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockOrders,
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /orders/i })).toBeInTheDocument();
    });
  });

  it('displays all orders with their status', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockOrders,
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      expect(screen.getByText(/pending/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/COOKED/i)).toBeInTheDocument();
    expect(screen.getByText(/Margherita/i)).toBeInTheDocument();
    expect(screen.getByText(/Pepperoni/i)).toBeInTheDocument();
  });

  it('shows events when an order is selected', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockOrders,
        });
      }
      if (url.includes('/api/events')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockEvents,
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    // Wait for orders to load
    await waitFor(() => {
      expect(screen.getByText(/pending/i)).toBeInTheDocument();
    });

    // Click on the first order
    const orderRows = screen.getAllByTestId('order-row');
    fireEvent.click(orderRows[0]);

    // Wait for events to load
    await waitFor(() => {
      expect(screen.getByText(/cooking/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/in oven/i)).toBeInTheDocument();
    expect(screen.getByText(/DONE/i)).toBeInTheDocument();
  });

  it('displays empty message when no orders exist', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => [],
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      expect(screen.getByText(/no orders/i)).toBeInTheDocument();
    });
  });

  it('displays events section heading when order is selected', async () => {
    (global.fetch as jest.Mock).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve({
          ok: true,
          json: async () => ({ status: 'healthy' }),
        });
      }
      if (url.includes('/api/orders')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockOrders,
        });
      }
      if (url.includes('/api/events')) {
        return Promise.resolve({
          ok: true,
          json: async () => mockEvents,
        });
      }
      return Promise.reject(new Error('Unknown endpoint'));
    });

    render(<ManagementPage />);

    await waitFor(() => {
      expect(screen.getByText(/pending/i)).toBeInTheDocument();
    });

    const orderRows = screen.getAllByTestId('order-row');
    fireEvent.click(orderRows[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /events/i })).toBeInTheDocument();
    });
  });
});
