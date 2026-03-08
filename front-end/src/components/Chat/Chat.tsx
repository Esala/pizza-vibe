'use client';

import { useRef, useEffect } from 'react';
import Logo from '@/components/Logo';
import ChatInput from './ChatInput';
import MessageTurn from './MessageTurn';
import styles from './Chat.module.css';

export interface ChatMessage {
  id: string;
  content: string;
  type: 'bot' | 'user';
}

interface ChatTurn {
  type: 'bot' | 'user';
  messages: string[];
}

interface ChatProps {
  messages?: ChatMessage[];
  inputValue?: string;
  onInputChange?: (value: string) => void;
  onSubmit?: () => void;
  inputDisabled?: boolean;
}

function groupIntoTurns(messages: ChatMessage[]): ChatTurn[] {
  const turns: ChatTurn[] = [];
  for (const msg of messages) {
    const last = turns[turns.length - 1];
    if (last && last.type === msg.type) {
      last.messages.push(msg.content);
    } else {
      turns.push({ type: msg.type, messages: [msg.content] });
    }
  }
  return turns;
}

export default function Chat({
  messages = [],
  inputValue = '',
  onInputChange,
  onSubmit,
  inputDisabled = false,
}: ChatProps) {
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (contentRef.current) {
      contentRef.current.scrollTop = contentRef.current.scrollHeight;
    }
  }, [messages]);

  const turns = groupIntoTurns(messages);

  return (
    <div className={styles.chat}>
      <div className={styles.header}>
        <Logo type="isotype" />
        <div className={styles.headerInfo}>
          <p className={styles.title}>Pizza Vibe Assistant</p>
          <p className={styles.subtitle}>Lets get you a pizza</p>
        </div>
      </div>
      <div className={styles.content} ref={contentRef}>
        <div className={styles.spacer} />
        {turns.map((turn, index) => (
          <MessageTurn key={index} messages={turn.messages} type={turn.type} />
        ))}
      </div>
      <div className={styles.inputArea}>
        <ChatInput
          value={inputValue}
          onChange={onInputChange ?? (() => {})}
          onSubmit={onSubmit ?? (() => {})}
          disabled={inputDisabled}
        />
      </div>
    </div>
  );
}
