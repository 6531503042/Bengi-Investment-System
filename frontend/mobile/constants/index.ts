export const API_CONFIG = {
    baseUrl: process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/api/v1',
    wsUrl: process.env.EXPO_PUBLIC_WS_URL ?? 'ws://localhost:8080/ws/prices',
    timeout: 10_000,
} as const

export const THEME = {
    colors: {
        primary: '#22c55e',
        secondary: '#3b82f6',
        success: '#22c55e',
        danger: '#ef4444',
        warning: '#f59e0b',
        background: {
            light: '#ffffff',
            dark: '#0f0f0f',
        },
        text: {
            light: '#000000',
            dark: '#ffffff',
        },
    },
    spacing: {
        xs: 4,
        sm: 8,
        md: 16,
        lg: 24,
        xl: 32,
    },
} as const

export const DEFAULT_WATCHLIST = ['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'TSLA'] as const

export const ORDER_TYPES = ['MARKET', 'LIMIT', 'STOP'] as const
export const ORDER_SIDES = ['BUY', 'SELL'] as const
