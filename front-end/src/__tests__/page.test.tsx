import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Home from '@/app/page';

// Mock fetch
global.fetch = jest.fn();

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
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    const pizzaSelect = screen.getByLabelText(/pizza type/i);
    const quantityInput = screen.getByLabelText(/quantity/i);
    const submitButton = screen.getByRole('button', { name: /place order/i });

    await user.selectOptions(pizzaSelect, 'Margherita');
    await user.tripleClick(quantityInput);
    await user.keyboard('2');
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

  it('displays success message after successful order', async () => {
    const user = userEvent.setup();
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ orderId: 'test-order-id', orderStatus: 'pending' }),
    });

    render(<Home />);

    const pizzaSelect = screen.getByLabelText(/pizza type/i);
    const quantityInput = screen.getByLabelText(/quantity/i);
    const submitButton = screen.getByRole('button', { name: /place order/i });

    await user.selectOptions(pizzaSelect, 'Margherita');
    await user.tripleClick(quantityInput);
    await user.keyboard('2');
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/order placed successfully/i)).toBeInTheDocument();
    });
  });

  it('displays error message when order fails', async () => {
    const user = userEvent.setup();
    (global.fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      status: 500,
    });

    render(<Home />);

    const pizzaSelect = screen.getByLabelText(/pizza type/i);
    const quantityInput = screen.getByLabelText(/quantity/i);
    const submitButton = screen.getByRole('button', { name: /place order/i });

    await user.selectOptions(pizzaSelect, 'Margherita');
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/failed to place order/i)).toBeInTheDocument();
    });
  });
});
