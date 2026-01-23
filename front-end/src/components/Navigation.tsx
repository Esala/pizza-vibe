'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function Navigation() {
  const pathname = usePathname();

  const navItems = [
    { href: '/', label: 'Order' },
    { href: '/kitchen', label: 'Kitchen' },
    { href: '/delivery', label: 'Delivery' },
    { href: '/management', label: 'Management' },
  ];

  return (
    <nav className="nav">
      <div className="nav-brand">Pizza Vibe</div>
      <ul className="nav-links">
        {navItems.map((item) => (
          <li key={item.href}>
            <Link
              href={item.href}
              className={pathname === item.href ? 'nav-link active' : 'nav-link'}
            >
              {item.label}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
