import styles from './TabItem.module.css';

interface TabItemProps {
  children: React.ReactNode;
  active?: boolean;
  onClick?: () => void;
  className?: string;
}

export default function TabItem({ children, active = false, onClick, className }: TabItemProps) {
  return (
    <button
      className={`${styles.tabItem}${active ? ` ${styles.active}` : ''}${className ? ` ${className}` : ''}`}
      onClick={onClick}
      role="tab"
      aria-selected={active}
    >
      {children}
    </button>
  );
}
