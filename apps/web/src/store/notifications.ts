import { create } from 'zustand';
import { eventUserFacing, type DomainEvent } from '@finflow/types';

interface InboxEntry {
  id: string;
  text: string;
  receivedAt: number;
  read: boolean;
}

interface NotificationState {
  entries: InboxEntry[];
  receive: (event: DomainEvent) => void;
  markAllRead: () => void;
  unreadCount: () => number;
}

export const useNotifications = create<NotificationState>((set, get) => ({
  entries: [],
  receive: (event) =>
    set((s) => ({
      entries: [
        { id: crypto.randomUUID(), text: eventUserFacing(event), receivedAt: Date.now(), read: false },
        ...s.entries,
      ].slice(0, 50),
    })),
  markAllRead: () => set((s) => ({ entries: s.entries.map((e) => ({ ...e, read: true })) })),
  unreadCount: () => get().entries.filter((e) => !e.read).length,
}));
