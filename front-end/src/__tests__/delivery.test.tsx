import { render, screen } from '@testing-library/react';
import DeliveryPage from '@/app/delivery/page';

describe('Delivery Page', () => {
  it('renders the delivery page title', () => {
    render(<DeliveryPage />);
    expect(screen.getByRole('heading', { name: /delivery/i })).toBeInTheDocument();
  });

  it('displays a list of delivery orders', () => {
    render(<DeliveryPage />);

    // Check for mock order IDs
    expect(screen.getByText(/order #1/i)).toBeInTheDocument();
    expect(screen.getByText(/order #2/i)).toBeInTheDocument();
  });

  it('displays delivery status for each order', () => {
    render(<DeliveryPage />);

    expect(screen.getByText(/in transit/i)).toBeInTheDocument();
    expect(screen.getByText(/delivered/i)).toBeInTheDocument();
  });

  it('displays order items for each delivery', () => {
    render(<DeliveryPage />);

    expect(screen.getByText(/Margherita/i)).toBeInTheDocument();
    expect(screen.getByText(/Pepperoni/i)).toBeInTheDocument();
  });

  it('displays delivery address information', () => {
    render(<DeliveryPage />);

    expect(screen.getByText(/123 Main St/i)).toBeInTheDocument();
    expect(screen.getByText(/456 Oak Ave/i)).toBeInTheDocument();
  });

  it('shows a note that delivery service is pending', () => {
    render(<DeliveryPage />);

    expect(screen.getByText(/service is pending/i)).toBeInTheDocument();
  });
});
