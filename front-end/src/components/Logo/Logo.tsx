import Image from 'next/image';
import styles from './Logo.module.css';

export default function Logo({ className }: { className?: string }) {
  return (
    <div className={`${styles.logo}${className ? ` ${className}` : ''}`}>
      <Image
        src="/images/logo-icon.svg"
        alt=""
        width={62}
        height={62}
        className={styles.icon}
      />
      <Image
        src="/images/logo-text.svg"
        alt="PizzaVibe"
        width={289}
        height={53}
        className={styles.text}
      />
    </div>
  );
}
