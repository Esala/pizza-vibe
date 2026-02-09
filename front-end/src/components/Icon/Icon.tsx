import styles from './Icon.module.css';

const icons = {
  minus: {
    viewBox: '0 0 24 4',
    width: 24,
    height: 4,
    path: 'M24 0V4H0V0H24Z',
  },
  add: {
    viewBox: '0 0 24 24',
    width: 24,
    height: 24,
    path: 'M14 10H24V14H14V24H10V14H0V10H10V0H14V10Z',
  },
  delete: {
    viewBox: '0 0 20 20',
    width: 20,
    height: 20,
    path: 'M16 4H20V8H17.5L16 20H4L2.5 8H0V4H4V0H16V4Z',
  },
} as const;

export type IconName = keyof typeof icons;

interface IconProps {
  name: IconName;
  className?: string;
}

export default function Icon({ name, className }: IconProps) {
  const icon = icons[name];

  return (
    <span
      className={`${styles.icon}${className ? ` ${className}` : ''}`}
      role="img"
      aria-label={name}
    >
      <svg
        viewBox={icon.viewBox}
        width={icon.width}
        height={icon.height}
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path d={icon.path} fill="currentColor" />
      </svg>
    </span>
  );
}
