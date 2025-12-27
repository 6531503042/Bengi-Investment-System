import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, TextInput, Alert } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams } from 'expo-router'
import { useMarketStore } from '@/stores/market'
import { useDemoStore } from '@/stores/demo'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { Ionicons } from '@expo/vector-icons'

type OrderSide = 'long' | 'short'
type OrderType = 'market' | 'limit'

// Available symbols for demo trading
const DEMO_SYMBOLS = ['AAPL', 'GOOGL', 'TSLA', 'BTC', 'ETH', 'MSFT', 'AMZN', 'NVDA']

// Mock prices for demo (until we connect to real price feeds)
const MOCK_PRICES: Record<string, { price: number; change: number }> = {
    AAPL: { price: 178.50, change: 2.35 },
    GOOGL: { price: 142.30, change: -0.85 },
    TSLA: { price: 248.75, change: 4.12 },
    BTC: { price: 43250.00, change: 1.25 },
    ETH: { price: 2280.50, change: 2.80 },
    MSFT: { price: 375.20, change: 1.15 },
    AMZN: { price: 155.80, change: -0.45 },
    NVDA: { price: 495.50, change: 3.75 },
}

export default function TradeScreen() {
    const params = useLocalSearchParams<{ symbol?: string; side?: string }>()
    const { quotes } = useMarketStore()
    const { account: demoAccount, fetchDemo } = useDemoStore()

    const [symbol, setSymbol] = useState(params.symbol ?? 'AAPL')
    const [side, setSide] = useState<OrderSide>((params.side as OrderSide) ?? 'long')
    const [orderType, setOrderType] = useState<OrderType>('market')
    const [amount, setAmount] = useState('')
    const [leverage, setLeverage] = useState(10)
    const [limitPrice, setLimitPrice] = useState('')

    useEffect(() => {
        fetchDemo()
    }, [])

    const mockQuote = MOCK_PRICES[symbol] ?? { price: 100, change: 0 }
    const quote = quotes[symbol] ?? mockQuote
    const currentPrice = quote?.price ?? mockQuote.price
    const changePercent = quote?.changePercent ?? mockQuote.change

    const balance = demoAccount?.balance ?? 0
    const margin = Number(amount) || 0
    const positionSize = margin * leverage

    const handleTrade = () => {
        if (margin <= 0) {
            Alert.alert('Invalid Amount', 'Please enter a valid trade amount')
            return
        }
        if (margin > balance) {
            Alert.alert('Insufficient Balance', `You only have $${balance.toFixed(2)} available`)
            return
        }

        Alert.alert(
            'Confirm Trade',
            `Open ${side.toUpperCase()} position:\n\n` +
            `Symbol: ${symbol}\n` +
            `Margin: $${margin.toFixed(2)}\n` +
            `Leverage: ${leverage}x\n` +
            `Position Size: $${positionSize.toFixed(2)}\n` +
            `Entry Price: $${currentPrice.toFixed(2)}`,
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Confirm',
                    onPress: () => {
                        // TODO: Connect to backend leverage trading API
                        Alert.alert('Coming Soon', 'Leverage trading will be available in the next update!')
                    }
                },
            ]
        )
    }

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                <ScrollView showsVerticalScrollIndicator={false}>
                    {/* Header with Balance */}
                    <XStack justifyContent="space-between" alignItems="center" padding="$4" paddingBottom="$2">
                        <YStack>
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                                Trade
                            </Text>
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                Demo Balance: ${balance.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                            </Text>
                        </YStack>
                        <View style={styles.demoBadge}>
                            <Text color={dimeTheme.colors.primary} fontSize="$1" fontWeight="bold">
                                DEMO
                            </Text>
                        </View>
                    </XStack>

                    {/* Symbol Selector */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Symbol
                        </Text>
                        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                            <XStack gap="$2">
                                {DEMO_SYMBOLS.map(s => (
                                    <Button
                                        key={s}
                                        size="$3"
                                        backgroundColor={symbol === s ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                        borderWidth={symbol === s ? 0 : 1}
                                        borderColor={dimeTheme.colors.border}
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
                                    <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                                        ${currentPrice.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                                    </Text>
                                </YStack>
                                <PriceChip value={changePercent} size="lg" />
                            </XStack>
                        </View>
                    </YStack>

                    {/* Long/Short Toggle */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <XStack gap="$3">
                            <Button
                                flex={1}
                                size="$5"
                                backgroundColor={side === 'long' ? dimeTheme.colors.profit : dimeTheme.colors.surface}
                                borderWidth={side === 'long' ? 0 : 1}
                                borderColor={dimeTheme.colors.border}
                                onPress={() => setSide('long')}
                            >
                                <XStack alignItems="center" gap="$2">
                                    <Ionicons
                                        name="trending-up"
                                        size={20}
                                        color={side === 'long' ? dimeTheme.colors.background : dimeTheme.colors.profit}
                                    />
                                    <Text
                                        color={side === 'long' ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                        fontWeight="bold"
                                        fontSize="$4"
                                    >
                                        Long
                                    </Text>
                                </XStack>
                            </Button>
                            <Button
                                flex={1}
                                size="$5"
                                backgroundColor={side === 'short' ? dimeTheme.colors.loss : dimeTheme.colors.surface}
                                borderWidth={side === 'short' ? 0 : 1}
                                borderColor={dimeTheme.colors.border}
                                onPress={() => setSide('short')}
                            >
                                <XStack alignItems="center" gap="$2">
                                    <Ionicons
                                        name="trending-down"
                                        size={20}
                                        color={side === 'short' ? dimeTheme.colors.background : dimeTheme.colors.loss}
                                    />
                                    <Text
                                        color={side === 'short' ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                        fontWeight="bold"
                                        fontSize="$4"
                                    >
                                        Short
                                    </Text>
                                </XStack>
                            </Button>
                        </XStack>
                    </YStack>

                    {/* Leverage Selector */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <XStack justifyContent="space-between" alignItems="center" marginBottom="$2">
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                Leverage
                            </Text>
                            <View style={styles.leverageBadge}>
                                <Text color={dimeTheme.colors.primary} fontWeight="bold">
                                    {leverage}x
                                </Text>
                            </View>
                        </XStack>
                        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                            <XStack gap="$2">
                                {[1, 2, 5, 10, 25, 50, 100].map(lev => (
                                    <Button
                                        key={lev}
                                        size="$3"
                                        backgroundColor={leverage === lev ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                        borderWidth={leverage === lev ? 0 : 1}
                                        borderColor={dimeTheme.colors.border}
                                        onPress={() => setLeverage(lev)}
                                    >
                                        <Text
                                            color={leverage === lev ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                            fontWeight="600"
                                        >
                                            {lev}x
                                        </Text>
                                    </Button>
                                ))}
                            </XStack>
                        </ScrollView>
                    </YStack>

                    {/* Margin Amount Input */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Margin Amount (USD)
                        </Text>
                        <View style={styles.amountInput}>
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$5">$</Text>
                            <TextInput
                                style={styles.input}
                                value={amount}
                                onChangeText={setAmount}
                                keyboardType="decimal-pad"
                                placeholder="0.00"
                                placeholderTextColor={dimeTheme.colors.textTertiary}
                            />
                        </View>
                        <XStack gap="$2" marginTop="$2">
                            {[10, 25, 50, 100].map(pct => (
                                <Button
                                    key={pct}
                                    flex={1}
                                    size="$2"
                                    backgroundColor={dimeTheme.colors.surface}
                                    onPress={() => setAmount(((balance * pct) / 100).toFixed(2))}
                                >
                                    <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                        {pct}%
                                    </Text>
                                </Button>
                            ))}
                        </XStack>
                    </YStack>

                    {/* Order Summary */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.summaryCard}>
                            <XStack justifyContent="space-between" marginBottom="$3">
                                <Text color={dimeTheme.colors.textSecondary}>Margin</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    ${margin.toFixed(2)}
                                </Text>
                            </XStack>
                            <XStack justifyContent="space-between" marginBottom="$3">
                                <Text color={dimeTheme.colors.textSecondary}>Position Size</Text>
                                <Text color={dimeTheme.colors.primary} fontWeight="bold" fontSize="$5">
                                    ${positionSize.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                                </Text>
                            </XStack>
                            <XStack justifyContent="space-between">
                                <Text color={dimeTheme.colors.textSecondary}>Entry Price</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    ${currentPrice.toFixed(2)}
                                </Text>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Submit Button */}
                    <YStack paddingHorizontal="$4" paddingBottom="$4">
                        <Button
                            size="$5"
                            backgroundColor={side === 'long' ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            disabled={margin <= 0}
                            opacity={margin <= 0 ? 0.5 : 1}
                            onPress={handleTrade}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons
                                    name={side === 'long' ? 'trending-up' : 'trending-down'}
                                    size={20}
                                    color={dimeTheme.colors.background}
                                />
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                    Open {side.toUpperCase()} • {leverage}x
                                </Text>
                            </XStack>
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
    demoBadge: {
        backgroundColor: 'rgba(0, 230, 118, 0.15)',
        paddingHorizontal: 8,
        paddingVertical: 4,
        borderRadius: 6,
        borderWidth: 1,
        borderColor: dimeTheme.colors.primary,
    },
    priceCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 20,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    leverageBadge: {
        backgroundColor: 'rgba(0, 230, 118, 0.15)',
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 8,
    },
    amountInput: {
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
        fontSize: 24,
        fontWeight: 'bold',
    },
    summaryCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
