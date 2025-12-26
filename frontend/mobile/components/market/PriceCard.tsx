import type { FC } from 'react'
import { XStack, YStack, Text, Card } from 'tamagui'
import { useMarketStore } from '@/stores/market'

interface PriceCardProps {
    symbol: string
    onPress?: () => void
}

export const PriceCard: FC<PriceCardProps> = ({ symbol, onPress }) => {
    const quote = useMarketStore((s) => s.quotes[symbol])
    const isPositive = (quote?.change ?? 0) >= 0
    const changeColor = isPositive ? '$green10' : '$red10'

    return (
        <Card
            elevate
            bordered
            animation="bouncy"
            pressStyle={{ scale: 0.98 }}
            onPress={onPress}
            marginBottom="$2"
            padding="$3"
        >
            <XStack justifyContent="space-between" alignItems="center">
                <YStack>
                    <Text fontWeight="bold" fontSize="$5">{symbol}</Text>
                    <Text color="$gray10" fontSize="$2">{quote ? 'Live' : 'Loading...'}</Text>
                </YStack>

                <YStack alignItems="flex-end">
                    <Text fontWeight="bold" fontSize="$5">
                        ${quote?.price?.toFixed(2) ?? '---'}
                    </Text>
                    <XStack gap="$2">
                        <Text color={changeColor} fontSize="$3">
                            {isPositive ? '+' : ''}{quote?.change?.toFixed(2) ?? '0.00'}
                        </Text>
                        <Text color={changeColor} fontSize="$3">
                            ({isPositive ? '+' : ''}{quote?.changePercent?.toFixed(2) ?? '0.00'}%)
                        </Text>
                    </XStack>
                </YStack>
            </XStack>
        </Card>
    )
}
