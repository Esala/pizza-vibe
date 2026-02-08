'use client';

import Logo from '@/components/Logo';
import Header from '@/components/Header';

export default function ComponentsShowcase() {
  return (
    <div style={{ padding: '40px', display: 'flex', flexDirection: 'column', gap: '60px' }}>
      <h1>Components Showcase</h1>

      {/* Logo */}
      <section>
        <h2>Logo</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
          <Logo />
        </div>
      </section>

      {/* Header */}
      <section>
        <h2>Header</h2>
        <div style={{ marginTop: '16px', border: '1px solid #ccc', borderRadius: '8px', overflow: 'hidden' }}>
          <Header />
        </div>
      </section>
    </div>
  );
}
