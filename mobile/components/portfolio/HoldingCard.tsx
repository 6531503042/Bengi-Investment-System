import { type FC } from 'react'
import { TouchableOpacity } from 'react-native'
import { XStack, YStack, Text, View } from 'tamagui'
import { PriceChip } from '@/components/common/PriceChip'
import { dimeTheme } from '@/constants/theme'

interface HoldingCardProps {
    symbol: string
    name: string
    quantity: number
    avgCost: number
    currentPrice: number
    onPress?: () => void
}

export const HoldingCard: FC<HoldingCardProps> = ({
    symbol,
    name,
    quantity,
    avgCost,
    currentPrice,
    onPress,
}) => {
    const totalValue = quantity * currentPrice
    const totalCost = quantity * avgCost
    const profitLoss = totalValue - totalCost
    const profitLossPercent = totalCost > 0 ? ((profitLoss / totalCost) * 100) : 0
    const isProfit = profitLoss >= 0

    return (
        <TouchableOpacity onPress={onPress} activeOpacity={0.7}>
            <YStack
                backgroundColor={dimeTheme.colors.surface}
                padding="$4"
                borderRadius="$4"
                marginBottom="$3"
                borderWidth={1}
                borderColor={dimeTheme.colors.border}
            >
                {/* Top Row: Symbol & Value */}
                <XStack justifyContent="space-between" alignItems="center" marginBottom="$3">
                    <XStack alignItems="center" gap="$3">
                        {/* Symbol Icon */}
                        <View
                            width={40}
                            height={40}
                            borderRadius={20}
                            backgroundColor={dimeTheme.colors.backgroundSecondary}
                            alignItems="center"
                            justifyContent="center"
                        >
                            <Text
                                color={dimeTheme.colors.primary}
                                fontWeight="bold"
                                fontSize="$4"
                            >
                                {symbol.charAt(0)}
                            </Text>
                        </View>

                        <YStack>
                            <Text
                                color={dimeTheme.colors.textPrimary}
                                fontWeight="bold"
                                fontSize="$4"
                            >
                                {symbol}
                            </Text>
                            <Text
                                color={dimeTheme.colors.textSecondary}
                                fontSize="$2"
                            >
                                {quantity} shares
                            </Text>
                        </YStack>
                    </XStack>

                    <YStack alignItems="flex-end">
                        <Text
                            color={dimeTheme.colors.textPrimary}
                            fontWeight="bold"
                            fontSize="$5"
                        >
                            ${totalValue.toFixed(2)}
                        </Text>
                        <PriceChip value={profitLossPercent} size="sm" />
                    </YStack>
                </XStack>

                {/* Bottom Row: Details */}
                <XStack justifyContent="space-between" paddingTop="$2" borderTopWidth={1} borderTopColor={dimeTheme.colors.border}>
                    <YStack>
                        <Text color={dimeTheme.colors.textTertiary} fontSize="$1">
                            Avg Cost
                        </Text>
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                            ${avgCost.toFixed(2)}
                        </Text>
                    </YStack>

                    <YStack alignItems="center">
                        <Text color={dimeTheme.colors.textTertiary} fontSize="$1">
                            Current
                        </Text>
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                            ${currentPrice.toFixed(2)}
                        </Text>
                    </YStack>

                    <YStack alignItems="flex-end">
                        <Text color={dimeTheme.colors.textTertiary} fontSize="$1">
                            P/L
                        </Text>
                        <Text
                            color={isProfit ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            fontSize="$3"
                            fontWeight="600"
                        >
                            {isProfit ? '+' : ''}${profitLoss.toFixed(2)}
                        </Text>
                    </YStack>
                </XStack>
            </YStack>
        </TouchableOpacity>
    )
}
