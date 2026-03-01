import styles from './Message.module.css';

interface MessageProps {
  message: string;
  type?: 'bot' | 'user';
}

export default function Message({ message, type = 'bot' }: MessageProps) {
  return (
    <div className={`${styles.message} ${type === 'user' ? styles.user : styles.bot}`}>
      <p className={styles.text}>{message}</p>
    </div>
  );
}
