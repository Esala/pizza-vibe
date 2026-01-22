'use client';

import { useState, FormEvent } from 'react';

export default function Home() {
  const [pizzaType, setPizzaType] = useState('Margherita');
  const [quantity, setQuantity] = useState(1);
  const [message, setMessage] = useState<string | null>(null);
  const [isError, setIsError] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setMessage(null);
    setIsError(false);

    try {
      const response = await fetch('/api/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          items: [{ pizza_type: pizzaType, quantity: quantity }],
        }),
      });

      if (response.ok) {
        setMessage('Order placed successfully!');
        setIsError(false);
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
        <button type="submit">Place Order</button>
      </form>
      {message && (
        <p role="status" style={{ color: isError ? 'red' : 'green' }}>
          {message}
        </p>
      )}
    </main>
  );
}
