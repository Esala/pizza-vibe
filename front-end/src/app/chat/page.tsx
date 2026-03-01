'use client';

import { useState } from 'react';
import Chat, { ChatMessage } from '@/components/Chat';
import styles from './page.module.css';

export default function ChatPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([
    { id: '1', content: 'Welcome to Pizza Vibe! What kind of pizza are you in the mood for today?', type: 'bot' },
  ]);
  const [inputValue, setInputValue] = useState('');

  const handleSubmit = () => {
    if (!inputValue.trim()) return;
    setMessages(prev => [...prev, { id: String(Date.now()), content: inputValue, type: 'user' }]);
    setInputValue('');
  };

  return (
    <main className={styles.page}>
      <div className={styles.chatWrapper}>
        <Chat
          messages={messages}
          inputValue={inputValue}
          onInputChange={setInputValue}
          onSubmit={handleSubmit}
        />
      </div>
      <div className={styles.placeholder} />
    </main>
  );
}
