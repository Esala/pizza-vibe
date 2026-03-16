'use client';

import { useState, useEffect, useCallback } from 'react';
import styles from './page.module.css';
import DashboardPanels from '@/components/DashboardPanels';
import DashboardBlock from '@/components/DashboardBlock';
import InventoryItem from '@/components/InventoryItem';
import OvenItem from '@/components/OvenItem';
import BikeItem from '@/components/BikeItem';
import AgentBlock from '@/components/AgentBlock';
import Button from '@/components/Button';
import Message from '@/components/Chat/Message';
import MessageTurn from '@/components/Chat/MessageTurn';

interface InventoryData {
  item: string;
  quantity: number;
}

interface DrinkData {
  item: string;
  quantity: number;
}

interface Oven {
  id: string;
  status: string;
  user?: string;
  progress: number;
  updatedAt: string;
}

interface Bike {
  id: string;
  status: string;
  user?: string;
  orderId?: string;
  updatedAt: string;
}

interface AgentEvent {
  agentId: string;
  kind: 'request' | 'response' | 'error';
  text: string;
  timestamp: string;
}

interface AgentConfig {
  agentIds: string[];
  label: string;
  emoji: string;
}

const INVENTORY_EMOJI: Record<string, string> = {
  PizzaDough: '🫓',
  Mozzarella: '🧀',
  Sauce: '🍅',
  Pepperoni: '🥓',
  Pineapple: '🍍',
};

const DRINK_EMOJI: Record<string, string> = {
  Beer: '🍺',
  Coke: '🥤',
  Wine: '🍷',
  Water: '💧',
  Juice: '🧃',
  Coffee: '☕',
  Lemonade: '🍋',
};

const AGENTS: AgentConfig[] = [
  { agentIds: ['store-mgmt-agent', 'chat-agent', 'pizza-order-workflow'], label: 'Store Manager', emoji: '👩‍💼' },
  { agentIds: ['drinks-agent'], label: 'Drinks Agent', emoji: '🤵‍♂️' },
  { agentIds: ['cooking-agent'], label: 'Cooking Agent', emoji: '👨‍🍳' },
  { agentIds: ['delivery-agent'], label: 'Delivery Agent', emoji: '🚴' },
];

function getInventoryEmoji(item: string): string {
  return INVENTORY_EMOJI[item] || '📦';
}

function getDrinkEmoji(item: string): string {
  return DRINK_EMOJI[item] || '🥤';
}

type StatusType = 'active' | 'inactive' | 'failed';

function getServiceStatus(ok: boolean): StatusType {
  return ok ? 'active' : 'failed';
}

/** Group events into consecutive turns by kind (request = user, response/error = bot) */
function groupEventsIntoTurns(events: AgentEvent[]): { type: 'bot' | 'user'; messages: string[] }[] {
  const turns: { type: 'bot' | 'user'; messages: string[] }[] = [];
  for (const event of events) {
    const type = event.kind === 'request' ? 'user' : 'bot';
    const lastTurn = turns[turns.length - 1];
    if (lastTurn && lastTurn.type === type) {
      lastTurn.messages.push(event.text);
    } else {
      turns.push({ type, messages: [event.text] });
    }
  }
  return turns;
}

const STORE_SERVICE_URL = process.env.NEXT_PUBLIC_STORE_SERVICE_URL || '';

export default function AgentsDashPage() {
  const [inventory, setInventory] = useState<InventoryData[]>([]);
  const [drinks, setDrinks] = useState<DrinkData[]>([]);
  const [ovens, setOvens] = useState<Oven[]>([]);
  const [bikes, setBikes] = useState<Bike[]>([]);
  const [inventoryOk, setInventoryOk] = useState(true);
  const [drinksOk, setDrinksOk] = useState(true);
  const [ovensOk, setOvensOk] = useState(true);
  const [bikesOk, setBikesOk] = useState(true);
  const [eventsByAgent, setEventsByAgent] = useState<Record<string, AgentEvent[]>>({});

  const fetchAll = useCallback(async () => {
    const base = STORE_SERVICE_URL;

    try {
      const res = await fetch(`${base}/api/inventory`);
      if (res.ok) { setInventory(await res.json()); setInventoryOk(true); }
      else setInventoryOk(false);
    } catch { setInventoryOk(false); }

    try {
      const res = await fetch(`${base}/api/drinks-stock`);
      if (res.ok) { setDrinks(await res.json()); setDrinksOk(true); }
      else setDrinksOk(false);
    } catch { setDrinksOk(false); }

    try {
      const res = await fetch(`${base}/api/oven`);
      if (res.ok) { setOvens(await res.json()); setOvensOk(true); }
      else setOvensOk(false);
    } catch { setOvensOk(false); }

    try {
      const res = await fetch(`${base}/api/bikes`);
      if (res.ok) { setBikes(await res.json()); setBikesOk(true); }
      else setBikesOk(false);
    } catch { setBikesOk(false); }

    try {
      const res = await fetch(`${base}/api/agents-events`);
      if (!res.ok) return;
      const allEvents: AgentEvent[] = await res.json();
      const grouped: Record<string, AgentEvent[]> = {};
      for (const agent of AGENTS) {
        grouped[agent.label] = [];
      }
      for (const event of allEvents) {
        const agent = AGENTS.find((a) => a.agentIds.includes(event.agentId));
        if (agent) {
          grouped[agent.label].push(event);
        }
      }
      setEventsByAgent(grouped);
    } catch {
      // silently ignore
    }
  }, []);

  useEffect(() => {
    fetchAll();
    const interval = setInterval(fetchAll, 1000);
    return () => clearInterval(interval);
  }, [fetchAll]);

  const handleClearAll = async () => {
    const base = STORE_SERVICE_URL;
    try {
      await fetch(`${base}/api/agents-events`, { method: 'DELETE' });
      setEventsByAgent({});
    } catch {
      // silently ignore
    }
  };

  return (
    <main className={styles.page}>
      {/* Dashboard Panels */}
      <DashboardPanels>
        <DashboardBlock icon="drinks" title="Drinks Stock" status={getServiceStatus(drinksOk)}>
          {drinks.map((item) => (
            <InventoryItem key={item.item} emoji={getDrinkEmoji(item.item)} quantity={item.quantity} />
          ))}
        </DashboardBlock>
        <DashboardBlock icon="inventory" title="Inventory" status={getServiceStatus(inventoryOk)}>
          {inventory.map((item) => (
            <InventoryItem key={item.item} emoji={getInventoryEmoji(item.item)} quantity={item.quantity} />
          ))}
        </DashboardBlock>
        <DashboardBlock icon="kitchen" title="Ovens" status={getServiceStatus(ovensOk)}>
          {ovens.map((oven, i) => (
            <OvenItem
              key={oven.id}
              number={i + 1}
              status={oven.status === 'AVAILABLE' ? 'idle' : 'cooking'}
            />
          ))}
        </DashboardBlock>
        <DashboardBlock icon="bikes" title="Bikes" status={getServiceStatus(bikesOk)}>
          {bikes.map((bike, i) => (
            <BikeItem
              key={bike.id}
              number={i + 1}
              status={bike.status === 'AVAILABLE' ? 'idle' : 'delivering'}
            />
          ))}
        </DashboardBlock>
      </DashboardPanels>

      {/* Agents Section */}
      <div className={styles.agentsSection}>
        <div className={styles.agentsTitleRow}>
          <h2 className={styles.agentsTitle}>Agents</h2>
          <Button color="danger" onClick={handleClearAll}>Clear all</Button>
        </div>
        <div className={styles.agentPanels}>
          {AGENTS.map((agent) => {
            const events = eventsByAgent[agent.label] || [];
            const hasActivity = events.some((e) => e.kind === 'response');
            const turns = groupEventsIntoTurns(events);

            return (
              <AgentBlock
                key={agent.label}
                emoji={agent.emoji}
                title={agent.label}
                status={hasActivity ? 'talking' : 'default'}
              >
                {events.length === 0 ? (
                  <Message message="No events yet..." type="bot" size="small" />
                ) : (
                  turns.map((turn, idx) => (
                    <MessageTurn
                      key={idx}
                      messages={turn.messages}
                      type={turn.type}
                      size="small"
                    />
                  ))
                )}
              </AgentBlock>
            );
          })}
        </div>
      </div>
    </main>
  );
}
