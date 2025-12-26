import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, TextInput } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams } from 'expo-router'
import { useMarketStore } from '@/stores/market'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { Ionicons } from '@expo/vector-icons'

type OrderSide = 'buy' | 'sell'
type OrderType = 'market' | 'limit'

export default function TradeScreen() {
    const params = useLocalSearchParams<{ symbol?: string; side?: string }>()
    const { quotes, watchedSymbols } = useMarketStore()

    const [symbol, setSymbol] = useState(params.symbol ?? 'AAPL')
    const [side, setSide] = useState<OrderSide>((params.side as OrderSide) ?? 'buy')
    const [orderType, setOrderType] = useState<OrderType>('market')
    const [quantity, setQuantity] = useState('')
    const [limitPrice, setLimitPrice] = useState('')

    const quote = quotes[symbol]
    const currentPrice = quote?.price ?? 0
    const estimatedTotal = Number(quantity) * (orderType === 'limit' ? Number(limitPrice) : currentPrice)

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                <ScrollView showsVerticalScrollIndicator={false}>
                    {/* Header */}
                    <YStack padding="$4" paddingBottom="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                            Trade
                        </Text>
                    </YStack>

                    {/* Symbol Selector */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Symbol
                        </Text>
                        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                            <XStack gap="$2">
                                {watchedSymbols.slice(0, 5).map(s => (
                                    <Button
                                        key={s}
                                        size="$3"
                                        backgroundColor={symbol === s ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                        onPress={() => setSymbol(s)}
                                    >
                                        <Text
                                            color={symbol === s ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                            fontWeight="600"
                                        >
                                            {s}
                                        </Text>
                                    </Button>
                                ))}
                            </XStack>
                        </ScrollView>
                    </YStack>

                    {/* Price Display */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.priceCard}>
                            <XStack justifyContent="space-between" alignItems="center">
                                <YStack>
                                    <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                        {symbol}
                                    </Text>
                                    <Text color={dimeTheme.colors.textPrimary} fontSize="$7" fontWeight="bold">
                                        ${currentPrice.toFixed(2)}
                                    </Text>
                                </YStack>
                                <PriceChip value={quote?.changePercent ?? 0} size="lg" />
                            </XStack>
                        </View>
                    </YStack>

                    {/* Buy/Sell Toggle */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <XStack gap="$3">
                            <Button
                                flex={1}
                                size="$5"
                                backgroundColor={side === 'buy' ? dimeTheme.colors.profit : dimeTheme.colors.surface}
                                borderWidth={side === 'buy' ? 0 : 1}
                                borderColor={dimeTheme.colors.border}
                                onPress={() => setSide('buy')}
                            >
                                <Text
                                    color={side === 'buy' ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                    fontWeight="bold"
                                    fontSize="$4"
                                >
                                    Buy
                                </Text>
                            </Button>
                            <Button
                                flex={1}
                                size="$5"
                                backgroundColor={side === 'sell' ? dimeTheme.colors.loss : dimeTheme.colors.surface}
                                borderWidth={side === 'sell' ? 0 : 1}
                                borderColor={dimeTheme.colors.border}
                                onPress={() => setSide('sell')}
                            >
                                <Text
                                    color={side === 'sell' ? dimeTheme.colors.textPrimary : dimeTheme.colors.textPrimary}
                                    fontWeight="bold"
                                    fontSize="$4"
                                >
                                    Sell
                                </Text>
                            </Button>
                        </XStack>
                    </YStack>

                    {/* Order Type */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Order Type
                        </Text>
                        <XStack gap="$3">
                            <Button
                                flex={1}
                                size="$4"
                                backgroundColor={orderType === 'market' ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                onPress={() => setOrderType('market')}
                            >
                                <Text
                                    color={orderType === 'market' ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                    fontWeight="600"
                                >
                                    Market
                                </Text>
                            </Button>
                            <Button
                                flex={1}
                                size="$4"
                                backgroundColor={orderType === 'limit' ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                onPress={() => setOrderType('limit')}
                            >
                                <Text
                                    color={orderType === 'limit' ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                    fontWeight="600"
                                >
                                    Limit
                                </Text>
                            </Button>
                        </XStack>
                    </YStack>

                    {/* Quantity Input */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Quantity
                        </Text>
                        <XStack alignItems="center" gap="$3">
                            <Button
                                size="$4"
                                circular
                                backgroundColor={dimeTheme.colors.surface}
                                onPress={() => setQuantity(String(Math.max(0, Number(quantity) - 1)))}
                            >
                                <Ionicons name="remove" size={20} color={dimeTheme.colors.textPrimary} />
                            </Button>
                            <View style={styles.quantityInput}>
                                <TextInput
                                    style={styles.input}
                                    value={quantity}
                                    onChangeText={setQuantity}
                                    keyboardType="numeric"
                                    placeholder="0"
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                />
                            </View>
                            <Button
                                size="$4"
                                circular
                                backgroundColor={dimeTheme.colors.surface}
                                onPress={() => setQuantity(String(Number(quantity) + 1))}
                            >
                                <Ionicons name="add" size={20} color={dimeTheme.colors.textPrimary} />
                            </Button>
                        </XStack>
                    </YStack>

                    {/* Limit Price (if applicable) */}
                    {orderType === 'limit' && (
                        <YStack paddingHorizontal="$4" marginBottom="$4">
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                                Limit Price
                            </Text>
                            <View style={styles.priceInput}>
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$4">$</Text>
                                <TextInput
                                    style={styles.input}
                                    value={limitPrice}
                                    onChangeText={setLimitPrice}
                                    keyboardType="decimal-pad"
                                    placeholder={currentPrice.toFixed(2)}
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                />
                            </View>
                        </YStack>
                    )}

                    {/* Order Summary */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.summaryCard}>
                            <XStack justifyContent="space-between" marginBottom="$2">
                                <Text color={dimeTheme.colors.textSecondary}>Estimated Total</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize="$5">
                                    ${isNaN(estimatedTotal) ? '0.00' : estimatedTotal.toFixed(2)}
                                </Text>
                            </XStack>
                            <XStack justifyContent="space-between">
                                <Text color={dimeTheme.colors.textSecondary}>Commission</Text>
                                <Text color={dimeTheme.colors.profit} fontWeight="600">
                                    $0.00
                                </Text>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Submit Button */}
                    <YStack paddingHorizontal="$4" paddingBottom="$4">
                        <Button
                            size="$5"
                            backgroundColor={side === 'buy' ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            disabled={!quantity || Number(quantity) <= 0}
                            opacity={!quantity || Number(quantity) <= 0 ? 0.5 : 1}
                        >
                            <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                {side === 'buy' ? 'Buy' : 'Sell'} {quantity || 0} {symbol}
                            </Text>
                        </Button>
                    </YStack>
                </ScrollView>
            </SafeAreaView>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: dimeTheme.colors.background,
    },
    safeArea: {
        flex: 1,
    },
    priceCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 20,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    quantityInput: {
        flex: 1,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: dimeTheme.radius.lg,
        padding: 16,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    priceInput: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: dimeTheme.radius.lg,
        padding: 16,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        gap: 8,
    },
    input: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 18,
        textAlign: 'center',
    },
    summaryCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
