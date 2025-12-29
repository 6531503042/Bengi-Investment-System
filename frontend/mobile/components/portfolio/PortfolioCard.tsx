import React from 'react'
import { StyleSheet, View, TouchableOpacity, Dimensions } from 'react-native'
import { Text, XStack, YStack } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { BlurView } from 'expo-blur'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

const { width: screenWidth } = Dimensions.get('window')

interface PortfolioCardProps {
    totalValue: number
    cashBalance: number
    investedValue: number
    dailyChange: number
    dailyChangePercent: number
    totalPnl: number
    totalPnlPercent: number
    currency?: string
    onDetailsPress?: () => void
}

export const PortfolioCard: React.FC<PortfolioCardProps> = ({
    totalValue,
    cashBalance,
    investedValue,
    dailyChange,
    dailyChangePercent,
    totalPnl,
    totalPnlPercent,
    currency = 'USD',
    onDetailsPress
}) => {
    const isProfit = totalPnl >= 0
    const isDailyProfit = dailyChange >= 0

    // Premium gradient colors based on overall P&L
    const gradientColors: readonly [string, string, string] = isProfit
        ? ['#0D4A2B', '#0A3821', '#061F14'] // Rich green gradient
        : ['#4A1A2B', '#3D1020', '#1F0810'] // Deep red gradient

    const formatCurrency = (value: number, showSign = false) => {
        const formatted = new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency,
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(Math.abs(value))

        if (showSign && value !== 0) {
            return value > 0 ? `+${formatted}` : `-${formatted}`
        }
        return formatted
    }

    const formatPercent = (value: number, showSign = true) => {
        const sign = showSign && value > 0 ? '+' : ''
        return `${sign}${value.toFixed(2)}%`
    }

    // Calculate allocation percentages
    const cashPercent = totalValue > 0 ? (cashBalance / totalValue) * 100 : 0
    const investedPercent = totalValue > 0 ? (investedValue / totalValue) * 100 : 0

    return (
        <View style={styles.container}>
            <LinearGradient
                colors={gradientColors}
                start={{ x: 0, y: 0 }}
                end={{ x: 1, y: 1 }}
                style={styles.gradient}
            >
                {/* Glassmorphism overlay */}
                <View style={styles.glassOverlay}>
                    {/* Header Row */}
                    <XStack justifyContent="space-between" alignItems="center" marginBottom="$2">
                        <XStack alignItems="center" gap="$2">
                            <View style={styles.portfolioIcon}>
                                <Ionicons name="briefcase" size={16} color={dimeTheme.colors.primary} />
                            </View>
                            <Text color={dimeTheme.colors.textSecondary} fontSize={13} fontWeight="500">
                                Portfolio Value
                            </Text>
                        </XStack>

                        <TouchableOpacity onPress={onDetailsPress} style={styles.detailsButton}>
                            <Text color={dimeTheme.colors.primary} fontSize={12} fontWeight="600">
                                Details
                            </Text>
                            <Ionicons name="chevron-forward" size={14} color={dimeTheme.colors.primary} />
                        </TouchableOpacity>
                    </XStack>

                    {/* Main Value */}
                    <YStack marginBottom="$4">
                        <Text
                            color={dimeTheme.colors.textPrimary}
                            fontSize={36}
                            fontWeight="bold"
                            letterSpacing={-1}
                        >
                            {formatCurrency(totalValue)}
                        </Text>

                        {/* Daily Change Pill */}
                        <XStack alignItems="center" gap="$2" marginTop="$1">
                            <View style={[
                                styles.changePill,
                                { backgroundColor: isDailyProfit ? 'rgba(0, 200, 83, 0.15)' : 'rgba(255, 82, 82, 0.15)' }
                            ]}>
                                <Ionicons
                                    name={isDailyProfit ? "trending-up" : "trending-down"}
                                    size={12}
                                    color={isDailyProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                />
                                <Text
                                    color={isDailyProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                    fontSize={12}
                                    fontWeight="600"
                                >
                                    {formatCurrency(dailyChange, true)} ({formatPercent(dailyChangePercent)})
                                </Text>
                            </View>
                            <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                                today
                            </Text>
                        </XStack>
                    </YStack>

                    {/* Allocation Bar */}
                    <View style={styles.allocationContainer}>
                        <View style={styles.allocationBar}>
                            <View style={[styles.investedBar, { width: `${investedPercent}%` }]} />
                            <View style={[styles.cashBar, { width: `${cashPercent}%` }]} />
                        </View>
                    </View>

                    {/* Stats Row */}
                    <XStack justifyContent="space-between" marginTop="$3">
                        {/* Invested */}
                        <YStack flex={1}>
                            <XStack alignItems="center" gap="$1" marginBottom="$1">
                                <View style={styles.legendDot} />
                                <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                                    Invested
                                </Text>
                            </XStack>
                            <Text color={dimeTheme.colors.textPrimary} fontSize={15} fontWeight="600">
                                {formatCurrency(investedValue)}
                            </Text>
                        </YStack>

                        {/* Cash */}
                        <YStack flex={1} alignItems="center">
                            <XStack alignItems="center" gap="$1" marginBottom="$1">
                                <View style={[styles.legendDot, { backgroundColor: dimeTheme.colors.textSecondary }]} />
                                <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                                    Cash
                                </Text>
                            </XStack>
                            <Text color={dimeTheme.colors.textPrimary} fontSize={15} fontWeight="600">
                                {formatCurrency(cashBalance)}
                            </Text>
                        </YStack>

                        {/* Total P&L */}
                        <YStack flex={1} alignItems="flex-end">
                            <Text color={dimeTheme.colors.textTertiary} fontSize={11} marginBottom="$1">
                                Total P&L
                            </Text>
                            <XStack alignItems="center" gap="$1">
                                <Text
                                    color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                    fontSize={15}
                                    fontWeight="600"
                                >
                                    {formatCurrency(totalPnl, true)}
                                </Text>
                            </XStack>
                            <Text
                                color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                fontSize={11}
                            >
                                {formatPercent(totalPnlPercent)}
                            </Text>
                        </YStack>
                    </XStack>
                </View>
            </LinearGradient>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        marginHorizontal: 16,
        borderRadius: 20,
        overflow: 'hidden',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 8 },
        shadowOpacity: 0.3,
        shadowRadius: 16,
        elevation: 10,
    },
    gradient: {
        borderRadius: 20,
    },
    glassOverlay: {
        padding: 20,
        backgroundColor: 'rgba(255, 255, 255, 0.03)',
        borderWidth: 1,
        borderColor: 'rgba(255, 255, 255, 0.08)',
        borderRadius: 20,
    },
    portfolioIcon: {
        width: 28,
        height: 28,
        borderRadius: 8,
        backgroundColor: 'rgba(0, 200, 83, 0.15)',
        alignItems: 'center',
        justifyContent: 'center',
    },
    detailsButton: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 4,
        paddingVertical: 6,
        paddingHorizontal: 10,
        backgroundColor: 'rgba(0, 200, 83, 0.1)',
        borderRadius: 8,
    },
    changePill: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 4,
        paddingVertical: 4,
        paddingHorizontal: 8,
        borderRadius: 12,
    },
    allocationContainer: {
        marginTop: 4,
    },
    allocationBar: {
        height: 6,
        backgroundColor: 'rgba(255, 255, 255, 0.1)',
        borderRadius: 3,
        flexDirection: 'row',
        overflow: 'hidden',
    },
    investedBar: {
        height: '100%',
        backgroundColor: dimeTheme.colors.primary,
        borderRadius: 3,
    },
    cashBar: {
        height: '100%',
        backgroundColor: dimeTheme.colors.textSecondary,
        opacity: 0.5,
    },
    legendDot: {
        width: 8,
        height: 8,
        borderRadius: 4,
        backgroundColor: dimeTheme.colors.primary,
    },
})
