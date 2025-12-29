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
    allocation: number
    onPress?: () => void
}

// Generate consistent color from symbol
const getSymbolColor = (symbol: string): string => {
    const colors = ['#E91E63', '#9C27B0', '#673AB7', '#3F51B5', '#2196F3', '#00BCD4', '#009688', '#4CAF50', '#FF9800', '#FF5722']
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

    const showPlaceholder = !logoUrl || imageError

    return (
        <TouchableOpacity
            style={styles.container}
            activeOpacity={0.7}
            onPress={onPress}
        >
            {/* Logo - Round with colored background */}
            <View style={styles.logoWrapper}>
                {!showPlaceholder ? (
                    <Image
                        source={{ uri: logoUrl }}
                        style={styles.logo}
                        onError={() => setImageError(true)}
                    />
                ) : (
                    <View style={[styles.logoPlaceholder, { backgroundColor: symbolColor }]}>
                        <Text color="#fff" fontSize={16} fontWeight="bold">
                            {symbol.charAt(0)}
                        </Text>
                    </View>
                )}
            </View>

            {/* Symbol + Name + Allocation */}
            <YStack flex={1} marginLeft={12}>
                <XStack alignItems="center" gap={8}>
                    <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                        {symbol}
                    </Text>
                    {/* Allocation Badge */}
                    <View style={[styles.allocationBadge, { backgroundColor: symbolColor + '25' }]}>
                        <Ionicons name="pie-chart" size={10} color={symbolColor} />
                        <Text color={symbolColor} fontSize={10} fontWeight="600" marginLeft={3}>
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
        backgroundColor: dimeTheme.colors.surface,
        paddingVertical: 14,
        paddingHorizontal: 16,
        marginHorizontal: 16,
        marginBottom: 8,
        borderRadius: 14,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    logoWrapper: {
        width: 48,
        height: 48,
        borderRadius: 24,
        overflow: 'hidden',
        backgroundColor: '#fff',
        alignItems: 'center',
        justifyContent: 'center',
        borderWidth: 2,
        borderColor: 'rgba(255,255,255,0.1)',
    },
    logo: {
        width: 32,
        height: 32,
        resizeMode: 'contain',
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
        paddingVertical: 3,
        borderRadius: 10,
    },
})
