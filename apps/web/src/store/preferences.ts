import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Currency } from '@finflow/types';

interface PreferencesState {
  locale: string;
  displayCurrency: Currency;
  compactNumbers: boolean;
  setLocale: (locale: string) => void;
  setDisplayCurrency: (currency: Currency) => void;
  toggleCompact: () => void;
}

export const usePreferences = create<PreferencesState>()(
  persist(
    (set) => ({
      locale: 'en-US',
      displayCurrency: 'USD',
      compactNumbers: false,
      setLocale: (locale) => set({ locale }),
      setDisplayCurrency: (displayCurrency) => set({ displayCurrency }),
      toggleCompact: () => set((s) => ({ compactNumbers: !s.compactNumbers })),
    }),
    { name: 'finflow.preferences' },
  ),
);
