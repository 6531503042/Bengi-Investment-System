import React, { useState } from 'react'
import { StyleSheet, Image, TouchableOpacity } from 'react-native'
import { Text, XStack, YStack, View } from 'tamagui'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

interface HoldingItemProps {
    symbol: string
    name: string
    logoUrl?: string
    quantity: number
    avgCost: number
    currentPrice: number
    allocation: number // percentage of portfolio
    onPress?: () => void
}

// Generate consistent color from symbol
const getSymbolColor = (symbol: string): string => {
    const colors = ['#4CAF50', '#2196F3', '#9C27B0', '#FF9800', '#E91E63', '#00BCD4', '#FF5722', '#3F51B5']
    let hash = 0
    for (let i = 0; i < symbol.length; i++) {
        hash = symbol.charCodeAt(i) + ((hash << 5) - hash)
    }
    return colors[Math.abs(hash) % colors.length]
}

export const HoldingItem: React.FC<HoldingItemProps> = ({
    symbol,
    name,
    logoUrl,
    quantity,
    avgCost,
    currentPrice,
    allocation,
    onPress,
}) => {
    const [imageError, setImageError] = useState(false)

    const totalValue = quantity * currentPrice
    const totalCost = quantity * avgCost
    const pnlAmount = totalValue - totalCost
    const pnlPercent = totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0
    const isProfit = pnlPercent >= 0
    const symbolColor = getSymbolColor(symbol)

    const formatCurrency = (value: number) => {
        if (value >= 1000000) {
            return `$${(value / 1000000).toFixed(2)}M`
        }
        if (value >= 1000) {
            return `$${(value / 1000).toFixed(2)}K`
        }
        return `$${value.toFixed(2)}`
    }

    const formatUSD = (value: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 2,
        }).format(value)
    }

    const showPlaceholder = !logoUrl || imageError

    return (
        <TouchableOpacity
            style={styles.container}
            activeOpacity={0.7}
            onPress={onPress}
        >
            {/* Left: Logo + Info */}
            <XStack alignItems="center" flex={1} gap="$3">
                {/* Logo */}
                <View style={styles.logoContainer}>
                    {!showPlaceholder ? (
                        <Image
                            source={{ uri: logoUrl }}
                            style={styles.logo}
                            onError={() => setImageError(true)}
                        />
                    ) : (
                        <View style={[styles.logoPlaceholder, { backgroundColor: symbolColor + '25' }]}>
                            <Text color={symbolColor} fontSize={18} fontWeight="bold">
                                {symbol.charAt(0)}
                            </Text>
                        </View>
                    )}
                </View>

                {/* Symbol + Name */}
                <YStack flex={1}>
                    <XStack alignItems="center" gap="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                            {symbol}
                        </Text>
                        <View style={[styles.allocationBadge, { backgroundColor: symbolColor + '20' }]}>
                            <Ionicons name="pie-chart" size={10} color={symbolColor} />
                            <Text color={symbolColor} fontSize={10} fontWeight="600" marginLeft={2}>
                                {allocation.toFixed(1)}%
                            </Text>
                        </View>
                    </XStack>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12} numberOfLines={1}>
                        {name}
                    </Text>
                </YStack>
            </XStack>

            {/* Right: Value + P&L */}
            <YStack alignItems="flex-end">
                <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                    {formatCurrency(totalValue)}
                </Text>
                <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                    ≈ {formatUSD(totalValue)}
                </Text>
                <XStack alignItems="center" gap="$1" marginTop="$1">
                    <Ionicons
                        name={isProfit ? "trending-up" : "trending-down"}
                        size={12}
                        color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                    />
                    <Text
                        color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                        fontSize={13}
                        fontWeight="600"
                    >
                        {isProfit ? '+' : ''}{pnlPercent.toFixed(2)}%
                    </Text>
                </XStack>
                <Text
                    color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                    fontSize={11}
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
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        marginHorizontal: 16,
        marginBottom: 8,
        borderRadius: 12,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    logoContainer: {
        width: 44,
        height: 44,
        borderRadius: 22,
        overflow: 'hidden',
    },
    logo: {
        width: 44,
        height: 44,
        borderRadius: 22,
    },
    logoPlaceholder: {
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
        paddingVertical: 2,
        borderRadius: 4,
    },
})
