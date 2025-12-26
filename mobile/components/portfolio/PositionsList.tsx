import type { FC } from 'react'
import { FlatList } from 'react-native'
import { YStack, XStack, Text, Card, Spinner } from 'tamagui'
import { usePortfolioStore } from '@/stores/portfolio'
import type { Position } from '@/types/portfolio'

interface PositionCardProps {
    position: Position
}

const PositionCard: FC<PositionCardProps> = ({ position }) => {
    const isPositive = position.unrealizedPL >= 0
    const changeColor = isPositive ? '$green10' : '$red10'

    return (
        <Card elevate bordered marginBottom="$2" padding="$3">
            <XStack justifyContent="space-between" alignItems="center">
                <YStack>
                    <Text fontWeight="bold" fontSize="$4">{position.symbol}</Text>
                    <Text color="$gray10" fontSize="$2">
                        {position.quantity} shares @ ${position.averageCost.toFixed(2)}
                    </Text>
                </YStack>

                <YStack alignItems="flex-end">
                    <Text fontWeight="bold" fontSize="$4">
                        ${position.marketValue.toFixed(2)}
                    </Text>
                    <XStack gap="$1">
                        <Text color={changeColor} fontSize="$3">
                            {isPositive ? '+' : ''}${position.unrealizedPL.toFixed(2)}
                        </Text>
                        <Text color={changeColor} fontSize="$3">
                            ({isPositive ? '+' : ''}{position.unrealizedPLPercent.toFixed(2)}%)
                        </Text>
                    </XStack>
                </YStack>
            </XStack>
        </Card>
    )
}

export const PositionsList: FC = () => {
    const { portfolios, selectedPortfolioId, isLoading } = usePortfolioStore()
    const portfolio = portfolios.find((p) => p.id === selectedPortfolioId)

    if (isLoading) {
        return (
            <YStack flex={1} justifyContent="center" alignItems="center">
                <Spinner size="large" color="$green10" />
            </YStack>
        )
    }

    if (!portfolio || portfolio.positions.length === 0) {
        return (
            <YStack flex={1} justifyContent="center" alignItems="center" padding="$4">
                <Text color="$gray10" textAlign="center">
                    No positions yet.{'\n'}Start trading to build your portfolio!
                </Text>
            </YStack>
        )
    }

    return (
        <FlatList<Position>
            data={portfolio.positions}
            keyExtractor={(item) => item.id}
            renderItem={({ item }) => <PositionCard position={item} />}
            contentContainerStyle={{ padding: 16 }}
            ListHeaderComponent={
                <YStack marginBottom="$4">
                    <Text fontSize="$6" fontWeight="bold">
                        ${portfolio.totalValue?.toFixed(2) ?? '0.00'}
                    </Text>
                    <Text color={portfolio.totalPL >= 0 ? '$green10' : '$red10'}>
                        {portfolio.totalPL >= 0 ? '+' : ''}${portfolio.totalPL?.toFixed(2) ?? '0.00'} ({portfolio.totalPLPercent?.toFixed(2) ?? '0.00'}%)
                    </Text>
                </YStack>
            }
        />
    )
}
