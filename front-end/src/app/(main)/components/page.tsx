'use client';

import { useState } from 'react';
import Logo from '@/components/Logo';
import Header from '@/components/Header';
import Footer from '@/components/Footer';
import Button from '@/components/Button';
import Tabs from '@/components/Tabs';
import Icon from '@/components/Icon';
import QuantitySelector from '@/components/QuantitySelector';
import CartItem from '@/components/CartItem';
import EmptyBlock from '@/components/EmptyBlock';
import PizzaItem from '@/components/PizzaItem';
import Chat, { ChatMessage } from '@/components/Chat';
import Message from '@/components/Chat/Message';
import MessageTurn from '@/components/Chat/MessageTurn';
import ChatInput from '@/components/Chat/ChatInput';

export default function ComponentsShowcase() {
  const [qty1, setQty1] = useState(2);
  const [qty2, setQty2] = useState(2);
  const [chatInput, setChatInput] = useState('');
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([
    { id: '1', content: 'Welcome to Pizza Vibe! What kind of pizza are you in the mood for today?', type: 'bot' },
    { id: '2', content: 'Following message from bot', type: 'bot' },
    { id: '3', content: 'I want a Margherita!', type: 'user' },
    { id: '4', content: 'Great choice! Adding a Margherita to your cart.', type: 'bot' },
  ]);

  const handleChatSubmit = () => {
    if (!chatInput.trim()) return;
    setChatMessages(prev => [...prev, { id: String(Date.now()), content: chatInput, type: 'user' }]);
    setChatInput('');
  };

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

      {/* Logo (Isotype) */}
      <section>
        <h2>Logo (Isotype)</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
          <Logo type="isotype" />
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

      {/* Icon */}
      <section>
        <h2>Icon</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', gap: '20px', alignItems: 'center' }}>
          <Icon name="minus" />
          <Icon name="add" />
          <Icon name="delete" />
          <Icon name="send" />
          <Icon name="check" />
          <Icon name="order" />
          <Icon name="kitchen" />
          <Icon name="delivery" />
        </div>
      </section>

      {/* Quantity Selector */}
      <section>
        <h2>Quantity Selector</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', gap: '20px', alignItems: 'center' }}>
          <QuantitySelector
            value={qty1}
            min={0}
            onIncrement={() => setQty1(qty1 + 1)}
            onDecrement={() => setQty1(qty1 - 1)}
          />
          <QuantitySelector
            value={qty2}
            min={1}
            deleteAtMin
            onIncrement={() => setQty2(qty2 + 1)}
            onDecrement={() => setQty2(qty2 - 1)}
            onDelete={() => setQty2(0)}
          />
        </div>
      </section>

      {/* Cart Item */}
      <section>
        <h2>Cart Item</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
          <CartItem name="Margherita" unitPrice={10} quantity={2} />
        </div>
      </section>

      {/* Empty Block */}
      <section>
        <h2>Empty Block</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
          <EmptyBlock />
        </div>
      </section>

      {/* Pizza Item */}
      <section>
        <h2>Pizza Item</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', gap: '20px' }}>
          <PizzaItem
            name="Margherita"
            price={10}
            description="San Marzano tomatoes, mozzarella cheese, fresh basil, salt, and extra-virgin olive oil"
            image="/images/pizza-margherita.svg"
          />
          <PizzaItem
            name="Pepperoni"
            price={15}
            description="Mozzarella cheese, pepperoni slices, olive oil, salt, and pepper"
            image="/images/pizza-pepperoni.svg"
          />
          <PizzaItem
            name="Hawaiian"
            price={15}
            description="Tomato sauce, mozzarella cheese, cooked ham, pineapple"
            image="/images/pizza-hawaiian.svg"
          />
          <PizzaItem
            name="Vegan"
            price={12}
            description="Vegan cheese, tomato sauce, mushrooms, onions, green peppers, and black olives"
            image="/images/pizza-vegan.svg"
          />
        </div>
      </section>

      {/* Tabs */}
      <section>
        <h2>Tabs</h2>
        <div style={{ marginTop: '16px', border: '1px solid #ccc', borderRadius: '8px', overflow: 'hidden', padding: '20px' }}>
          <Tabs
            tabs={[
              { label: 'Tab Item', value: 'tab1' },
              { label: 'Tab Item', value: 'tab2' },
              { label: 'Tab Item', value: 'tab3' },
            ]}
            defaultValue="tab1"
          />
        </div>
      </section>

      {/* Message */}
      <section>
        <h2>Message</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', flexDirection: 'column', gap: '12px', background: '#08b869' }}>
          <Message message="Welcome to Pizza Vibe! What kind of pizza are you in the mood for today?" type="bot" />
          <Message message="I want a Margherita!" type="user" />
        </div>
      </section>

      {/* MessageTurn */}
      <section>
        <h2>MessageTurn</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', flexDirection: 'column', gap: '20px', background: '#08b869' }}>
          <MessageTurn
            messages={['Welcome to Pizza Vibe!', 'What kind of pizza are you in the mood for today?']}
            type="bot"
          />
          <MessageTurn
            messages={['I want a Margherita!', 'Make it a large please.']}
            type="user"
          />
        </div>
      </section>

      {/* ChatInput */}
      <section>
        <h2>ChatInput</h2>
        <div style={{ marginTop: '16px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'flex', flexDirection: 'column', gap: '16px', background: '#08b869' }}>
          <ChatInput value={chatInput} onChange={setChatInput} onSubmit={handleChatSubmit} />
          <ChatInput value="" onChange={() => {}} onSubmit={() => {}} disabled />
        </div>
      </section>

      {/* Chat */}
      <section>
        <h2>Chat</h2>
        <div style={{ marginTop: '16px', border: '1px solid #ccc', borderRadius: '8px', overflow: 'hidden', height: '600px' }}>
          <Chat
            messages={chatMessages}
            inputValue={chatInput}
            onInputChange={setChatInput}
            onSubmit={handleChatSubmit}
          />
        </div>
      </section>
    </div>
  );
}
