/** Typed access to public environment variables. */
export const env = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080',
  appName: process.env.NEXT_PUBLIC_APP_NAME ?? 'FinFlow',
  analyticsEnabled: process.env.NEXT_PUBLIC_ANALYTICS === 'on',
} as const;
