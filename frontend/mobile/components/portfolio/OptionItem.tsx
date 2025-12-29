import React from 'react'
import { StyleSheet, TouchableOpacity } from 'react-native'
import { Text, XStack, YStack, View } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

interface OptionItemProps {
    symbol: string
    name: string
    logoUrl?: string
    type: 'Call' | 'Put'
    strike: number
    expiry: string
    contracts: number
    premium: number
    currentPrice: number
    delta?: number
    theta?: number
    iv?: number // Implied volatility
    onPress?: () => void
}

// Curated brand colors for popular stocks (Dime-style)
const BRAND_COLORS: Record<string, { bg: string; text: string }> = {
    AAPL: { bg: '#000000', text: '#fff' },
    MSFT: { bg: '#00A4EF', text: '#fff' },
    GOOGL: { bg: '#4285F4', text: '#fff' },
    AMZN: { bg: '#FF9900', text: '#000' },
    META: { bg: '#0668E1', text: '#fff' },
    NVDA: { bg: '#76B900', text: '#fff' },
    TSLA: { bg: '#E82127', text: '#fff' },
    AMD: { bg: '#ED1C24', text: '#fff' },
    NFLX: { bg: '#E50914', text: '#fff' },
    PLTR: { bg: '#000000', text: '#fff' },
    DEFAULT: { bg: '#374151', text: '#fff' },
}

const getSymbolStyle = (symbol: string) => {
    return BRAND_COLORS[symbol.toUpperCase()] || BRAND_COLORS.DEFAULT
}

export const OptionItem: React.FC<OptionItemProps> = ({
    symbol,
    name,
    logoUrl,
    type,
    strike,
    expiry,
    contracts,
    premium,
    currentPrice,
    delta = 0,
    theta = 0,
    iv = 0,
    onPress,
}) => {
    const sharesPerContract = 100
    const totalCost = contracts * sharesPerContract * premium
    const totalValue = contracts * sharesPerContract * currentPrice
    const pnlAmount = totalValue - totalCost
    const pnlPercent = totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0
    const isProfit = pnlPercent >= 0

    const isCall = type === 'Call'
    const typeColor = isCall ? '#00C853' : '#FF5252'
    const brandStyle = getSymbolStyle(symbol)

    // Gradient based on option type
    const gradientColors: readonly [string, string] = isCall
        ? ['rgba(0, 200, 83, 0.08)', 'rgba(0, 200, 83, 0.02)']
        : ['rgba(255, 82, 82, 0.08)', 'rgba(255, 82, 82, 0.02)']

    const formatUSD = (value: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
        }).format(value)
    }

    // Calculate days until expiry
    const getDaysUntilExpiry = () => {
        const expiryDate = new Date(expiry)
        const today = new Date()
        const diffTime = expiryDate.getTime() - today.getTime()
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
        return diffDays
    }

    const daysLeft = getDaysUntilExpiry()
    const isExpiringSoon = daysLeft <= 7

    return (
        <TouchableOpacity activeOpacity={0.8} onPress={onPress}>
            <LinearGradient
                colors={gradientColors}
                start={{ x: 0, y: 0 }}
                end={{ x: 1, y: 1 }}
                style={styles.container}
            >
                {/* Top: Type Badge + Expiry */}
                <XStack justifyContent="space-between" alignItems="center" marginBottom="$3">
                    {/* Option Type Badge */}
                    <View style={[styles.typeBadge, { backgroundColor: typeColor }]}>
                        <Ionicons
                            name={isCall ? "arrow-up" : "arrow-down"}
                            size={12}
                            color="#fff"
                        />
                        <Text color="#fff" fontSize={12} fontWeight="bold" marginLeft={4}>
                            {type.toUpperCase()}
                        </Text>
                    </View>

                    {/* Expiry Badge */}
                    <View style={[
                        styles.expiryBadge,
                        isExpiringSoon && styles.expiryBadgeUrgent
                    ]}>
                        <Ionicons
                            name="time-outline"
                            size={12}
                            color={isExpiringSoon ? '#FF9800' : dimeTheme.colors.textSecondary}
                        />
                        <Text
                            color={isExpiringSoon ? '#FF9800' : dimeTheme.colors.textSecondary}
                            fontSize={11}
                            fontWeight={isExpiringSoon ? "600" : "400"}
                            marginLeft={4}
                        >
                            {daysLeft > 0 ? `${daysLeft}d left` : 'Expired'}
                        </Text>
                    </View>
                </XStack>

                {/* Main Row */}
                <XStack alignItems="center" gap="$3">
                    {/* Logo with Option Indicator - Dime style */}
                    <View style={styles.logoWrapper}>
                        <View style={[styles.logoContainer, { backgroundColor: brandStyle.bg }]}>
                            <Text color={brandStyle.text} fontSize={18} fontWeight="bold">
                                {symbol.charAt(0)}
                            </Text>
                        </View>
                        {/* Option type indicator ring */}
                        <View style={[styles.optionRing, { borderColor: typeColor }]} />
                    </View>

                    {/* Symbol + Strike Info */}
                    <YStack flex={1}>
                        <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={16}>
                            {symbol}
                        </Text>
                        <XStack alignItems="center" gap="$2" marginTop="$1">
                            <View style={styles.strikeBox}>
                                <Text color={dimeTheme.colors.textSecondary} fontSize={11}>
                                    Strike
                                </Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600" fontSize={14}>
                                    ${strike.toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.contractsBox}>
                                <Text color={dimeTheme.colors.textSecondary} fontSize={11}>
                                    Contracts
                                </Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600" fontSize={14}>
                                    {contracts}
                                </Text>
                            </View>
                        </XStack>
                    </YStack>

                    {/* Right: Value + P&L */}
                    <YStack alignItems="flex-end">
                        <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={17}>
                            {formatUSD(totalValue)}
                        </Text>
                        <XStack
                            alignItems="center"
                            gap="$1"
                            paddingVertical="$1"
                            paddingHorizontal="$2"
                            backgroundColor={isProfit ? 'rgba(0, 200, 83, 0.15)' : 'rgba(255, 82, 82, 0.15)'}
                            borderRadius={6}
                            marginTop="$1"
                        >
                            <Ionicons
                                name={isProfit ? "caret-up" : "caret-down"}
                                size={14}
                                color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            />
                            <Text
                                color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                fontSize={13}
                                fontWeight="bold"
                            >
                                {isProfit ? '+' : ''}{pnlPercent.toFixed(1)}%
                            </Text>
                        </XStack>
                        <Text
                            color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            fontSize={11}
                            marginTop="$1"
                        >
                            {isProfit ? '+' : ''}{formatUSD(pnlAmount)}
                        </Text>
                    </YStack>
                </XStack>

                {/* Bottom Stats - Greeks & Price */}
                <View style={styles.statsContainer}>
                    <View style={styles.statItem}>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={10}>PREMIUM</Text>
                        <Text color={dimeTheme.colors.textPrimary} fontSize={13} fontWeight="600">
                            ${premium.toFixed(2)}
                        </Text>
                    </View>
                    <View style={styles.statDivider} />
                    <View style={styles.statItem}>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={10}>CURRENT</Text>
                        <Text color={dimeTheme.colors.textPrimary} fontSize={13} fontWeight="600">
                            ${currentPrice.toFixed(2)}
                        </Text>
                    </View>
                    <View style={styles.statDivider} />
                    <View style={styles.statItem}>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={10}>DELTA</Text>
                        <Text
                            color={delta >= 0 ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            fontSize={13}
                            fontWeight="600"
                        >
                            {delta >= 0 ? '+' : ''}{delta.toFixed(2)}
                        </Text>
                    </View>
                    <View style={styles.statDivider} />
                    <View style={styles.statItem}>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={10}>IV</Text>
                        <Text color={dimeTheme.colors.textSecondary} fontSize={13} fontWeight="600">
                            {(iv * 100).toFixed(0)}%
                        </Text>
                    </View>
                </View>
            </LinearGradient>
        </TouchableOpacity>
    )
}

const styles = StyleSheet.create({
    container: {
        marginHorizontal: 16,
        marginBottom: 12,
        borderRadius: 16,
        padding: 16,
        borderWidth: 1,
        borderColor: 'rgba(255, 255, 255, 0.08)',
    },
    typeBadge: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: 10,
        paddingVertical: 5,
        borderRadius: 20,
    },
    expiryBadge: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: 'rgba(255, 255, 255, 0.05)',
        paddingHorizontal: 10,
        paddingVertical: 5,
        borderRadius: 8,
    },
    expiryBadgeUrgent: {
        backgroundColor: 'rgba(255, 152, 0, 0.15)',
    },
    logoWrapper: {
        position: 'relative',
    },
    logoContainer: {
        width: 48,
        height: 48,
        borderRadius: 24,
        alignItems: 'center',
        justifyContent: 'center',
    },
    logo: {
        width: 48,
        height: 48,
        borderRadius: 24,
    },
    logoPlaceholder: {
        width: 48,
        height: 48,
        borderRadius: 24,
        alignItems: 'center',
        justifyContent: 'center',
    },
    optionRing: {
        position: 'absolute',
        top: -2,
        left: -2,
        right: -2,
        bottom: -2,
        borderRadius: 26,
        borderWidth: 2,
    },
    strikeBox: {
        backgroundColor: 'rgba(255, 255, 255, 0.05)',
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 6,
    },
    contractsBox: {
        backgroundColor: 'rgba(255, 255, 255, 0.05)',
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 6,
    },
    statsContainer: {
        flexDirection: 'row',
        marginTop: 16,
        paddingTop: 12,
        borderTopWidth: 1,
        borderTopColor: 'rgba(255, 255, 255, 0.08)',
        justifyContent: 'space-between',
    },
    statItem: {
        flex: 1,
        alignItems: 'center',
    },
    statDivider: {
        width: 1,
        backgroundColor: 'rgba(255, 255, 255, 0.08)',
    },
})
