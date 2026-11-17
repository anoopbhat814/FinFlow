import { create } from 'zustand';

interface SessionState {
  token: string | null;
  userId: string | null;
  setSession: (token: string, userId: string) => void;
  clear: () => void;
}

export const useSession = create<SessionState>((set) => ({
  token: null,
  userId: null,
  setSession: (token, userId) => set({ token, userId }),
  clear: () => set({ token: null, userId: null }),
}));
