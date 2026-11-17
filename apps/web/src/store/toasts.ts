import { create } from 'zustand';

export interface ToastItem {
  id: number;
  message: string;
  kind: 'success' | 'error';
}

interface ToastState {
  items: ToastItem[];
  push: (message: string, kind?: ToastItem['kind']) => void;
  dismiss: (id: number) => void;
}

let nextId = 1;

export const useToasts = create<ToastState>((set) => ({
  items: [],
  push: (message, kind = 'success') =>
    set((s) => ({ items: [...s.items, { id: nextId++, message, kind }] })),
  dismiss: (id) => set((s) => ({ items: s.items.filter((t) => t.id !== id) })),
}));
