'use client';

interface OrderItem {
  pizzaType: string;
  quantity: number;
}

interface DeliveryOrder {
  orderId: string;
  orderItems: OrderItem[];
  status: string;
  address: string;
}

export default function DeliveryPage() {
  // Mock data since delivery service is pending
  const mockOrders: DeliveryOrder[] = [
    {
      orderId: 'Order #1',
      orderItems: [
        { pizzaType: 'Margherita', quantity: 2 },
        { pizzaType: 'Pepperoni', quantity: 1 },
      ],
      status: 'In Transit',
      address: '123 Main St',
    },
    {
      orderId: 'Order #2',
      orderItems: [
        { pizzaType: 'Hawaiian', quantity: 1 },
      ],
      status: 'Delivered',
      address: '456 Oak Ave',
    },
  ];

  return (
    <main>
      <h1>Delivery</h1>
      <p>Note: Delivery service is pending - showing mock data</p>
      <div>
        {mockOrders.map((order) => (
          <div key={order.orderId}>
            <h3>{order.orderId}</h3>
            <p>Status: {order.status}</p>
            <p>Address: {order.address}</p>
            <ul>
              {order.orderItems.map((item, index) => (
                <li key={index}>
                  {item.quantity}x {item.pizzaType}
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </main>
  );
}
