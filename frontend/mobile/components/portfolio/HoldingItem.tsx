import React from 'react'
import { StyleSheet, TouchableOpacity } from 'react-native'
import { Text, XStack, YStack, View } from 'tamagui'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

interface HoldingItemProps {
    symbol: string
    name: string
    logoUrl?: string // ignored - we use letter icons
    quantity: number
    avgCost: number
    currentPrice: number
    allocation: number
    onPress?: () => void
}

// Curated brand colors for popular stocks (Dime-style)
const BRAND_COLORS: Record<string, { bg: string; text: string }> = {
    // Tech
    AAPL: { bg: '#000000', text: '#fff' },
    MSFT: { bg: '#00A4EF', text: '#fff' },
    GOOGL: { bg: '#4285F4', text: '#fff' },
    GOOG: { bg: '#4285F4', text: '#fff' },
    AMZN: { bg: '#FF9900', text: '#000' },
    META: { bg: '#0668E1', text: '#fff' },
    NVDA: { bg: '#76B900', text: '#fff' },
    TSLA: { bg: '#E82127', text: '#fff' },
    AMD: { bg: '#ED1C24', text: '#fff' },
    INTC: { bg: '#0071C5', text: '#fff' },
    NFLX: { bg: '#E50914', text: '#fff' },
    PLTR: { bg: '#000000', text: '#fff' },

    // Finance
    JPM: { bg: '#0A6EBD', text: '#fff' },
    V: { bg: '#1A1F71', text: '#fff' },
    MA: { bg: '#EB001B', text: '#fff' },
    COIN: { bg: '#0052FF', text: '#fff' },

    // Crypto
    BTC: { bg: '#F7931A', text: '#fff' },
    ETH: { bg: '#627EEA', text: '#fff' },
    SOL: { bg: '#9945FF', text: '#fff' },

    // Default
    DEFAULT: { bg: '#374151', text: '#fff' },
}

const getSymbolStyle = (symbol: string) => {
    return BRAND_COLORS[symbol.toUpperCase()] || BRAND_COLORS.DEFAULT
}

export const HoldingItem: React.FC<HoldingItemProps> = ({
    symbol,
    name,
    quantity,
    avgCost,
    currentPrice,
    allocation,
    onPress,
}) => {
    const totalValue = quantity * currentPrice
    const totalCost = quantity * avgCost
    const pnlAmount = totalValue - totalCost
    const pnlPercent = totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0
    const isProfit = pnlPercent >= 0
    const brandStyle = getSymbolStyle(symbol)

    const formatValue = (value: number) => {
        if (value >= 1000000) return `${(value / 1000000).toFixed(2)}M`
        if (value >= 1000) return `$${(value / 1000).toFixed(2)}K`
        return `$${value.toFixed(2)}`
    }

    const formatUSD = (value: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
        }).format(value)
    }

    return (
        <TouchableOpacity
            style={styles.container}
            activeOpacity={0.7}
            onPress={onPress}
        >
            {/* Logo - Dime style: colored circle with letter */}
            <View style={[styles.logoContainer, { backgroundColor: brandStyle.bg }]}>
                <Text color={brandStyle.text} fontSize={18} fontWeight="bold">
                    {symbol.charAt(0)}
                </Text>
            </View>

            {/* Symbol + Name + Allocation */}
            <YStack flex={1} marginLeft={12}>
                <XStack alignItems="center" gap={8}>
                    <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                        {symbol}
                    </Text>
                    {/* Allocation Badge */}
                    <View style={styles.allocationBadge}>
                        <Ionicons name="pie-chart" size={10} color={dimeTheme.colors.primary} />
                        <Text color={dimeTheme.colors.primary} fontSize={10} fontWeight="600" marginLeft={3}>
                            {allocation.toFixed(1)}%
                        </Text>
                    </View>
                </XStack>
                <Text
                    color={dimeTheme.colors.textSecondary}
                    fontSize={12}
                    numberOfLines={1}
                    marginTop={2}
                >
                    {name}
                </Text>
            </YStack>

            {/* Value Column */}
            <YStack alignItems="flex-end" minWidth={90}>
                <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                    {formatValue(totalValue)}
                </Text>
                <Text color={dimeTheme.colors.textTertiary} fontSize={11} marginTop={2}>
                    ≈ {formatUSD(totalValue)}
                </Text>
            </YStack>

            {/* P&L Column */}
            <YStack alignItems="flex-end" minWidth={80} marginLeft={8}>
                <XStack alignItems="center" gap={3}>
                    <Ionicons
                        name={isProfit ? "arrow-up" : "arrow-down"}
                        size={12}
                        color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                    />
                    <Text
                        color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                        fontSize={14}
                        fontWeight="bold"
                    >
                        {isProfit ? '+' : ''}{pnlPercent.toFixed(2)}%
                    </Text>
                </XStack>
                <Text
                    color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                    fontSize={11}
                    marginTop={2}
                >
                    ({isProfit ? '+' : ''}{formatUSD(pnlAmount)})
                </Text>
            </YStack>
        </TouchableOpacity>
    )
}

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingVertical: 14,
        paddingHorizontal: 16,
        backgroundColor: dimeTheme.colors.surface,
        marginHorizontal: 16,
        marginBottom: 8,
        borderRadius: 14,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    logoContainer: {
        width: 44,
        height: 44,
        borderRadius: 22,
        alignItems: 'center',
        justifyContent: 'center',
    },
    allocationBadge: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: 6,
        paddingVertical: 3,
        borderRadius: 10,
        backgroundColor: 'rgba(0, 230, 118, 0.15)',
    },
})
