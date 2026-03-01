import Message from '../Message';
import styles from './MessageTurn.module.css';

interface MessageTurnProps {
  messages: string[];
  type?: 'bot' | 'user';
}

export default function MessageTurn({ messages, type = 'bot' }: MessageTurnProps) {
  return (
    <div className={`${styles.turn} ${type === 'user' ? styles.user : styles.bot}`}>
      {messages.map((message, index) => (
        <Message key={index} message={message} type={type} />
      ))}
    </div>
  );
}
