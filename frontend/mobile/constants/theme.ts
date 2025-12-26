// Dime-style theme colors and constants
export const dimeTheme = {
    // Core colors
    colors: {
        // Background
        background: '#0D0D0D',
        backgroundSecondary: '#1A1A1A',
        surface: '#1F1F1F',

        // Primary (Dime green)
        primary: '#00E676',
        primaryDark: '#00C853',
        primaryLight: '#69F0AE',

        // Status colors
        profit: '#00E676',
        loss: '#FF3B30',
        warning: '#FF9500',

        // Text
        textPrimary: '#FFFFFF',
        textSecondary: '#8E8E93',
        textTertiary: '#636366',

        // Border
        border: '#2C2C2E',
        borderLight: '#3A3A3C',

        // Accent
        blue: '#007AFF',
        purple: '#AF52DE',
        teal: '#5AC8FA',
    },

    // Spacing
    spacing: {
        xs: 4,
        sm: 8,
        md: 16,
        lg: 24,
        xl: 32,
    },

    // Border radius
    radius: {
        sm: 8,
        md: 12,
        lg: 16,
        xl: 20,
        full: 9999,
    },

    // Typography
    font: {
        light: 300,
        regular: 400,
        medium: 500,
        semibold: 600,
        bold: 700,
    },

    // Chart colors
    chart: {
        candleUp: '#00E676',
        candleDown: '#FF3B30',
        ma5: '#FFA726',
        ma10: '#42A5F5',
        ma20: '#AB47BC',
        volume: '#8E8E93',
        grid: '#2C2C2E',
    },
} as const

export type DimeTheme = typeof dimeTheme

// Export Tamagui-compatible tokens
export const tamaguiDimeConfig = {
    colors: {
        ...dimeTheme.colors,
    },
}
