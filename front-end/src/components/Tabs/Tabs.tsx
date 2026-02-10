'use client';

import { useState } from 'react';
import styles from './Tabs.module.css';
import TabItem from './TabItem';

interface Tab {
  label: string;
  value: string;
}

interface TabsProps {
  tabs: Tab[];
  defaultValue?: string;
  onTabChange?: (value: string) => void;
  className?: string;
}

export default function Tabs({ tabs, defaultValue, onTabChange, className }: TabsProps) {
  const [activeTab, setActiveTab] = useState(defaultValue ?? tabs[0]?.value);

  const handleTabClick = (value: string) => {
    setActiveTab(value);
    onTabChange?.(value);
  };

  return (
    <div
      className={`${styles.tabs}${className ? ` ${className}` : ''}`}
      role="tablist"
    >
      <div className={styles.buttonGroup}>
        <div className={styles.buttonWrapper}>
          {tabs.map((tab) => (
            <TabItem
              key={tab.value}
              active={activeTab === tab.value}
              onClick={() => handleTabClick(tab.value)}
            >
              {tab.label}
            </TabItem>
          ))}
        </div>
      </div>
    </div>
  );
}
