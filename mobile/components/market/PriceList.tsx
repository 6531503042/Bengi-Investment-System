import type { FC } from 'react'
import { FlatList } from 'react-native'
import { YStack, Text, Spinner } from 'tamagui'
import { useMarketStore } from '@/stores/market'
import { PriceCard } from './PriceCard'

interface PriceListProps {
    onSymbolPress?: (symbol: string) => void
}

export const PriceList: FC<PriceListProps> = ({ onSymbolPress }) => {
    const { watchedSymbols, isLoading } = useMarketStore()

    if (isLoading) {
        return (
            <YStack flex={1} justifyContent="center" alignItems="center">
                <Spinner size="large" color="$green10" />
            </YStack>
        )
    }

    if (watchedSymbols.length === 0) {
        return (
            <YStack flex={1} justifyContent="center" alignItems="center" padding="$4">
                <Text color="$gray10" textAlign="center">
                    No symbols in your watchlist.{'\n'}Add some to get started.
                </Text>
            </YStack>
        )
    }

    return (
        <FlatList
            data={watchedSymbols}
            keyExtractor={(item) => item}
            renderItem={({ item }) => (
                <PriceCard symbol={item} onPress={() => onSymbolPress?.(item)} />
            )}
            contentContainerStyle={{ padding: 16 }}
            showsVerticalScrollIndicator={false}
        />
    )
}
