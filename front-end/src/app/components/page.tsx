'use client';

import Logo from '@/components/Logo';
import Header from '@/components/Header';
import Footer from '@/components/Footer';
import Button from '@/components/Button';

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

      {/* Logo (Small) */}
      <section>
        <h2>Logo (Small)</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
          <Logo size="small" />
        </div>
      </section>

      {/* Header */}
      <section>
        <h2>Header</h2>
        <div style={{ marginTop: '16px', border: '1px solid #ccc', borderRadius: '8px', overflow: 'hidden' }}>
          <Header />
        </div>
      </section>

      {/* Footer */}
      <section>
        <h2>Footer</h2>
        <div style={{ marginTop: '16px', border: '1px solid #ccc', borderRadius: '8px', overflow: 'hidden' }}>
          <Footer />
        </div>
      </section>

      {/* Button */}
      <section>
        <h2>Button</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', gap: '20px', alignItems: 'center' }}>
          <Button>Button</Button>
          <Button disabled>Button</Button>
        </div>
      </section>
    </div>
  );
}
