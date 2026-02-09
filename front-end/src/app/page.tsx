'use client';

import { useState, useRef, useCallback, useEffect, FormEvent } from 'react';

interface OrderItem {
  pizzaType: string;
  quantity: number;
}

interface WebSocketEvent {
  orderId: string;
  status: string;
  source: string;
  timestamp: string;
}

function generateClientId(): string {
  return 'client-' + Math.random().toString(36).substring(2, 15);
}

export default function Home() {
  const [pizzaType, setPizzaType] = useState('Margherita');
  const [quantity, setQuantity] = useState(1);
  const [cart, setCart] = useState<OrderItem[]>([]);
  const [message, setMessage] = useState<string | null>(null);
  const [isError, setIsError] = useState(false);
  const [orderId, setOrderId] = useState<string | null>(null);
  const [wsConnected, setWsConnected] = useState(false);
  const [events, setEvents] = useState<WebSocketEvent[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const clientIdRef = useRef<string>(generateClientId());

  const connectWebSocket = useCallback((): Promise<void> => {
    return new Promise((resolve, reject) => {
      const storeWsUrl = process.env.NEXT_PUBLIC_STORE_WS_URL || 'ws://localhost:8080';
      const wsUrl = `${storeWsUrl}/ws?clientId=${clientIdRef.current}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        setWsConnected(true);
        resolve();
      };

      ws.onmessage = (event: MessageEvent) => {
        const data: WebSocketEvent = JSON.parse(event.data);
        setEvents((prev) => [...prev, data]);
      };

      ws.onclose = () => {
        setWsConnected(false);
      };

      ws.onerror = () => {
        setWsConnected(false);
        reject(new Error('WebSocket connection failed'));
      };

      wsRef.current = ws;
    });
  }, []);

  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

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
      {events.length > 0 && (
        <table data-testid="events-table">
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Status</th>
              <th>Source</th>
              <th>Timestamp</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event, index) => (
              <tr key={index}>
                <td>{event.orderId}</td>
                <td>{event.status}</td>
                <td>{event.source}</td>
                <td>{event.timestamp}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </main>
  );
}
