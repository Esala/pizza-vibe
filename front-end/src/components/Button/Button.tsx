import styles from './Button.module.css';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  className?: string;
}

export default function Button({ children, className, ...props }: ButtonProps) {
  return (
    <button
      className={`${styles.button}${className ? ` ${className}` : ''}`}
      {...props}
    >
      {children}
    </button>
  );
}
