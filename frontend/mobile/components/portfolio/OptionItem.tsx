import React, { useState } from 'react'
import { StyleSheet, Image, TouchableOpacity } from 'react-native'
import { Text, XStack, YStack, View } from 'tamagui'
import { Ionicons } from '@expo/vector-icons'
import { dimeTheme } from '@/constants/theme'

interface OptionItemProps {
    symbol: string
    name: string
    logoUrl?: string
    type: 'Call' | 'Put'
    strike: number
    expiry: string // e.g., "Jan 9, 2025"
    contracts: number
    premium: number // Price paid per share
    currentPrice: number // Current option price
    delta?: number
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
    onPress,
}) => {
    const [imageError, setImageError] = useState(false)

    const sharesPerContract = 100
    const totalCost = contracts * sharesPerContract * premium
    const totalValue = contracts * sharesPerContract * currentPrice
    const pnlAmount = totalValue - totalCost
    const pnlPercent = totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0
    const isProfit = pnlPercent >= 0

    const typeColor = type === 'Call' ? dimeTheme.colors.profit : dimeTheme.colors.loss
    const symbolColor = getSymbolColor(symbol)
    const showPlaceholder = !logoUrl || imageError

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
            {/* Expiry Badge */}
            <View style={styles.expiryBadge}>
                <Ionicons name="time-outline" size={10} color={dimeTheme.colors.textSecondary} />
                <Text color={dimeTheme.colors.textSecondary} fontSize={10} marginLeft={2}>
                    Expires: {expiry}
                </Text>
            </View>

            <XStack alignItems="center" flex={1} gap="$3" marginTop="$2">
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

                {/* Symbol + Type */}
                <YStack flex={1}>
                    <XStack alignItems="center" gap="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                            {symbol}
                        </Text>
                        <View style={[styles.typeBadge, { backgroundColor: typeColor + '20' }]}>
                            <Text color={typeColor} fontSize={11} fontWeight="bold">
                                {type}
                            </Text>
                        </View>
                    </XStack>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12}>
                        Strike: ${strike.toFixed(2)}
                    </Text>
                    <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                        {contracts} contract{contracts > 1 ? 's' : ''} = {contracts * 100} shares
                    </Text>
                </YStack>

                {/* Right: Value + P&L */}
                <YStack alignItems="flex-end">
                    <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize={15}>
                        {formatUSD(totalValue)}
                    </Text>
                    <XStack alignItems="center" gap="$1" marginTop="$1">
                        <Ionicons
                            name={isProfit ? "trending-up" : "trending-down"}
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
                    >
                        ({isProfit ? '+' : ''}{formatUSD(pnlAmount)})
                    </Text>
                </YStack>
            </XStack>

            {/* Bottom Stats */}
            <XStack marginTop="$3" paddingTop="$2" borderTopWidth={1} borderColor={dimeTheme.colors.border} justifyContent="space-between">
                <YStack>
                    <Text color={dimeTheme.colors.textTertiary} fontSize={10}>Premium Paid</Text>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12}>${premium.toFixed(2)}</Text>
                </YStack>
                <YStack alignItems="center">
                    <Text color={dimeTheme.colors.textTertiary} fontSize={10}>Current Price</Text>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12}>${currentPrice.toFixed(2)}</Text>
                </YStack>
                <YStack alignItems="center">
                    <Text color={dimeTheme.colors.textTertiary} fontSize={10}>Total Cost</Text>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12}>{formatUSD(totalCost)}</Text>
                </YStack>
                <YStack alignItems="flex-end">
                    <Text color={dimeTheme.colors.textTertiary} fontSize={10}>Delta</Text>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12}>{delta.toFixed(2)}</Text>
                </YStack>
            </XStack>
        </TouchableOpacity>
    )
}

const styles = StyleSheet.create({
    container: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        marginHorizontal: 16,
        marginBottom: 8,
        borderRadius: 12,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    expiryBadge: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.background,
        paddingHorizontal: 8,
        paddingVertical: 4,
        borderRadius: 4,
        alignSelf: 'flex-start',
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
    typeBadge: {
        paddingHorizontal: 8,
        paddingVertical: 3,
        borderRadius: 4,
    },
})
