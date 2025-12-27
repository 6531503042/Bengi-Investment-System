import React from 'react'
import { StyleSheet, View } from 'react-native'
import { Text, XStack, YStack } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

interface PortfolioCardProps {
    totalValue: number
    initialValue: number
    dailyChangePercent: number
    totalPnlPercent: number
    totalPnlAmount: number
    currency?: string
}

export const PortfolioCard: React.FC<PortfolioCardProps> = ({
    totalValue,
    initialValue,
    dailyChangePercent,
    totalPnlPercent,
    totalPnlAmount,
    currency = 'USD'
}) => {
    const isProfit = totalPnlPercent >= 0
    const isDailyProfit = dailyChangePercent >= 0
    const progress = Math.min(100, Math.max(0, (totalValue / initialValue) * 100))

    // Gradient colors based on P&L
    const gradientColors = isProfit
        ? ['#1a472a', '#0d3320', '#0a1f14'] as const
        : ['#4a1a1a', '#331010', '#1a0808'] as const

    const formatCurrency = (value: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency,
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(value)
    }

    const formatPercent = (value: number) => {
        const sign = value >= 0 ? '↗' : '↘'
        return `${sign} ${Math.abs(value).toFixed(2)}%`
    }

    return (
        <View style={styles.container}>
            <LinearGradient
                colors={gradientColors}
                start={{ x: 0, y: 0 }}
                end={{ x: 1, y: 1 }}
                style={styles.gradient}
            >
                {/* Total Value */}
                <YStack marginBottom="$3">
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12} marginBottom="$1">
                        Total Portfolio Value
                    </Text>
                    <XStack alignItems="baseline" gap="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontSize={32} fontWeight="bold">
                            {formatCurrency(totalValue)}
                        </Text>
                    </XStack>
                    <Text color={dimeTheme.colors.textTertiary} fontSize={12} marginTop="$1">
                        / {formatCurrency(initialValue)}
                    </Text>
                </YStack>

                {/* Progress Bar */}
                <View style={styles.progressContainer}>
                    <View style={[styles.progressBar, { width: `${progress}%` }]} />
                </View>
                <Text color={dimeTheme.colors.textSecondary} fontSize={11} marginTop="$1">
                    {progress.toFixed(0)}%
                </Text>

                {/* P&L Info */}
                <XStack marginTop="$3" gap="$4">
                    {/* Daily Change */}
                    <YStack>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                            Daily Change
                        </Text>
                        <Text
                            color={isDailyProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            fontSize={14}
                            fontWeight="600"
                        >
                            {formatPercent(dailyChangePercent)}
                        </Text>
                    </YStack>

                    {/* Total P&L */}
                    <YStack>
                        <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                            Total P&L
                        </Text>
                        <XStack alignItems="center" gap="$1">
                            <Text
                                color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                fontSize={14}
                                fontWeight="600"
                            >
                                {formatPercent(totalPnlPercent)}
                            </Text>
                            <Text
                                color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                fontSize={12}
                            >
                                ({isProfit ? '+' : ''}{formatCurrency(totalPnlAmount)})
                            </Text>
                        </XStack>
                    </YStack>
                </XStack>
            </LinearGradient>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        marginHorizontal: 16,
        borderRadius: 16,
        overflow: 'hidden',
    },
    gradient: {
        padding: 20,
        borderRadius: 16,
    },
    progressContainer: {
        height: 6,
        backgroundColor: 'rgba(255,255,255,0.1)',
        borderRadius: 3,
        overflow: 'hidden',
    },
    progressBar: {
        height: '100%',
        backgroundColor: dimeTheme.colors.primary,
        borderRadius: 3,
    },
})
