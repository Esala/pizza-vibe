'use client';

import { useState, FormEvent } from 'react';
import { useOrder } from '@/context/OrderContext';

interface OrderItem {
  pizzaType: string;
  quantity: number;
}

export default function Home() {
  const [pizzaType, setPizzaType] = useState('Margherita');
  const [quantity, setQuantity] = useState(1);
  const [cart, setCart] = useState<OrderItem[]>([]);
  const [message, setMessage] = useState<string | null>(null);
  const [isError, setIsError] = useState(false);

  const { orderId, setOrderId, events, setEvents, wsConnected, connectWebSocket } = useOrder();

  const handleAddToCart = () => {
    setCart((prevCart) => {
      const existingIndex = prevCart.findIndex(
        (item) => item.pizzaType === pizzaType
      );
      if (existingIndex >= 0) {
        const updated = [...prevCart];
        updated[existingIndex] = {
          ...updated[existingIndex],
          quantity: updated[existingIndex].quantity + quantity,
        };
        return updated;
      }
      return [...prevCart, { pizzaType, quantity }];
    });
  };

  const handleRemoveFromCart = (pizzaTypeToRemove: string) => {
    setCart((prevCart) =>
      prevCart.filter((item) => item.pizzaType !== pizzaTypeToRemove)
    );
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (cart.length === 0) return;

    setMessage(null);
    setIsError(false);
    setOrderId(null);
    setEvents([]);

    try {
      // Connect WebSocket before placing the order so no events are missed
      await connectWebSocket();

      const response = await fetch('/api/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          orderItems: cart,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setOrderId(data.orderId);
        setMessage('Order placed successfully!');
        setIsError(false);
        setCart([]);
      } else {
        setMessage('Failed to place order');
        setIsError(true);
      }
    } catch {
      setMessage('Failed to place order');
      setIsError(true);
    }
  };

  const kitchenEvents = events.filter((e) => e.source === 'kitchen');
  const deliveryEvents = events.filter((e) => e.source === 'delivery');
  const isCooked = kitchenEvents.some((e) => e.status === 'COOKED');
  const isDelivered = deliveryEvents.some((e) => e.status === 'DELIVERED');

  return (
    <main>
      <h1>Pizza Vibe</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="pizzaType">Pizza Type</label>
          <select
            id="pizzaType"
            value={pizzaType}
            onChange={(e) => setPizzaType(e.target.value)}
          >
            <option value="Margherita">Margherita</option>
            <option value="Pepperoni">Pepperoni</option>
            <option value="Hawaiian">Hawaiian</option>
            <option value="Veggie">Veggie</option>
          </select>
        </div>
        <div>
          <label htmlFor="quantity">Quantity</label>
          <input
            id="quantity"
            type="number"
            min="1"
            value={quantity}
            onChange={(e) => setQuantity(parseInt(e.target.value, 10) || 1)}
          />
        </div>
        <button type="button" onClick={handleAddToCart}>Add to Cart</button>
        {cart.length > 0 && (
          <table data-testid="cart">
            <thead>
              <tr>
                <th>Pizza Type</th>
                <th>Quantity</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {cart.map((item) => (
                <tr key={item.pizzaType}>
                  <td>{item.pizzaType}</td>
                  <td>{item.quantity}</td>
                  <td>
                    <button
                      type="button"
                      onClick={() => handleRemoveFromCart(item.pizzaType)}
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <button type="submit" disabled={cart.length === 0}>Place Order</button>
      </form>
      {message && (
        <p role="status" style={{ color: isError ? 'red' : 'green' }}>
          {message}
        </p>
      )}
      {orderId && (
        <p data-testid="order-id">Order ID: {orderId}</p>
      )}
      {orderId && (
        <p data-testid="ws-status">
          WebSocket: {wsConnected ? 'Connected' : 'Disconnected'}
        </p>
      )}
      {events.length > 0 && (() => {
        const latestOvenProgress = kitchenEvents
          .filter((e) => e.status === 'oven_progress')
          .slice(-1)[0];
        const ovenPercent = latestOvenProgress?.message
          ? parseInt(latestOvenProgress.message.match(/(\d+)% complete/)?.[1] || '0', 10)
          : null;
        const progressPercent = isCooked ? 100 : (ovenPercent !== null ? Math.min(99, ovenPercent) : Math.min(99, kitchenEvents.length * 20));
        return (
          <p data-testid="cooking-progress">
            Cooking progress: {progressPercent}%
          </p>
        );
      })()}
      {isCooked && (
        <p data-testid="delivery-progress">
          Delivery: {isDelivered ? 'Delivered' : `In progress (${deliveryEvents.filter((e) => e.status !== 'DELIVERED').length} updates)`}
        </p>
      )}
      {events.length > 0 && (
        <table data-testid="events-table">
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Status</th>
              <th>Source</th>
              <th>Timestamp</th>
              <th>Details</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event, index) => (
              <tr key={index}>
                <td>{event.orderId}</td>
                <td>{event.status}</td>
                <td>{event.source}</td>
                <td>{event.timestamp}</td>
                <td>
                  {event.message && <span>{event.message}</span>}
                  {event.toolName && <span> [{event.toolName}]</span>}
                  {event.toolInput && <span> {event.toolInput}</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </main>
  );
}
