'use client';

import { useState, useEffect } from 'react';

interface ServiceStatus {
  name: string;
  status: 'healthy' | 'unhealthy';
  url: string;
}

export default function ManagementPage() {
  const [services, setServices] = useState<ServiceStatus[]>([
    { name: 'Store Service', status: 'unhealthy', url: '/api/health/store' },
    { name: 'Kitchen Service', status: 'unhealthy', url: '/api/health/kitchen' },
    { name: 'Delivery Service', status: 'unhealthy', url: '/api/health/delivery' },
  ]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkServices = async () => {
      const updatedServices = await Promise.all(
        services.map(async (service) => {
          try {
            const response = await fetch(service.url);
            return {
              ...service,
              status: response.ok ? ('healthy' as const) : ('unhealthy' as const),
            };
          } catch {
            return {
              ...service,
              status: 'unhealthy' as const,
            };
          }
        })
      );

      setServices(updatedServices);
      setLoading(false);
    };

    checkServices();
  }, []);

  if (loading) {
    return (
      <main>
        <h1>Management</h1>
        <p>Checking services...</p>
      </main>
    );
  }

  return (
    <main>
      <h1>Management</h1>
      <h2>Service Status</h2>
      <div>
        {services.map((service) => (
          <div key={service.name}>
            <h3>{service.name}</h3>
            <p>Status: {service.status}</p>
          </div>
        ))}
      </div>
    </main>
  );
}
