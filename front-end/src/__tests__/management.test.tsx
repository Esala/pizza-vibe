import { render, screen, waitFor } from '@testing-library/react';
import ManagementPage from '@/app/management/page';

// Mock fetch
global.fetch = jest.fn();

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
    (global.fetch as jest.Mock).mockRejectedValue(new Error('Network error'));

    render(<ManagementPage />);

    await waitFor(() => {
      const unhealthyStatuses = screen.getAllByText(/unhealthy/i);
      expect(unhealthyStatuses.length).toBeGreaterThan(0);
    });
  });
});
