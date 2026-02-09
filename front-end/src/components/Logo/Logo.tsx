import Image from 'next/image';
import styles from './Logo.module.css';

interface LogoProps {
  className?: string;
  size?: 'default' | 'small';
}

export default function Logo({ className, size = 'default' }: LogoProps) {
  const sizeClass = size === 'small' ? styles.small : '';

  return (
    <div className={`${styles.logo} ${sizeClass}${className ? ` ${className}` : ''}`}>
      <Image
        src="/images/logo-icon.svg"
        alt=""
        width={size === 'small' ? 31 : 62}
        height={size === 'small' ? 31 : 62}
        className={styles.icon}
      />
      <Image
        src="/images/logo-text.svg"
        alt="PizzaVibe"
        width={size === 'small' ? 144 : 289}
        height={size === 'small' ? 26 : 53}
        className={styles.text}
      />
    </div>
  );
}
