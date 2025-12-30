import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, TextInput, Alert, TouchableOpacity, Switch, ActivityIndicator } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams, useRouter } from 'expo-router'
import { useMarketStore } from '@/stores/market'
import { useDemoStore } from '@/stores/demo'
import { usePortfolioStore } from '@/stores/portfolio'
import { orderService } from '@/services/api'
import { dimeTheme } from '@/constants/theme'
import { Ionicons } from '@expo/vector-icons'
import type { CreateOrderInput } from '@/types/trade'

type OrderSide = 'buy' | 'sell'
type OrderMode = 'spot' | 'leverage' | 'options'
type OrderType = 'market' | 'limit' | 'scheduled'

const DEMO_SYMBOLS = ['AAPL', 'GOOGL', 'TSLA', 'BTC', 'ETH', 'MSFT', 'AMZN', 'NVDA', 'PLTR', 'META']

const MOCK_PRICES: Record<string, { price: number; change: number }> = {
    AAPL: { price: 254.49, change: 2.35 },
    GOOGL: { price: 193.42, change: -0.85 },
    TSLA: { price: 421.06, change: 4.12 },
    BTC: { price: 94250.00, change: 1.25 },
    ETH: { price: 3380.50, change: 2.80 },
    MSFT: { price: 425.50, change: 1.15 },
    AMZN: { price: 227.05, change: -0.45 },
    NVDA: { price: 134.82, change: 3.75 },
    PLTR: { price: 75.14, change: 5.20 },
    META: { price: 585.50, change: 1.85 },
}

const LEVERAGE_OPTIONS = [1, 2, 5, 10, 20, 50, 100]

export default function TradeScreen() {
    const router = useRouter()
    const params = useLocalSearchParams<{ symbol?: string; side?: string }>()
    const { quotes } = useMarketStore()
    const { account: demoAccount, fetchDemo } = useDemoStore()
    const { activePortfolio } = usePortfolioStore()
    const [isSubmitting, setIsSubmitting] = useState(false)

    const [symbol, setSymbol] = useState(params.symbol ?? 'AAPL')
    const [side, setSide] = useState<OrderSide>((params.side === 'short' ? 'sell' : 'buy') as OrderSide)
    const [orderMode, setOrderMode] = useState<OrderMode>('spot')
    const [orderType, setOrderType] = useState<OrderType>('market')
    const [amount, setAmount] = useState('')
    const [shares, setShares] = useState('')
    const [leverage, setLeverage] = useState(10)
    const [limitPrice, setLimitPrice] = useState('')

    // SL/TP
    const [enableSLTP, setEnableSLTP] = useState(false)
    const [stopLoss, setStopLoss] = useState('')
    const [takeProfit, setTakeProfit] = useState('')

    // Scheduled order
    const [scheduledDate, setScheduledDate] = useState(new Date())
    const [showDatePicker, setShowDatePicker] = useState(false)

    useEffect(() => {
        fetchDemo()
    }, [])

    const mockQuote = MOCK_PRICES[symbol] ?? { price: 100, change: 0 }
    const quote = quotes[symbol] ?? mockQuote
    const currentPrice = quote?.price ?? mockQuote.price
    const changePercent = quote?.changePercent ?? mockQuote.change
    const isPositive = changePercent >= 0

    const balance = demoAccount?.balance ?? 50000
    const inputAmount = Number(amount) || 0
    const inputShares = Number(shares) || 0

    // Calculate based on mode
    const getPositionDetails = () => {
        if (orderMode === 'spot') {
            const qty = inputShares || (inputAmount / currentPrice)
            return { shares: qty, value: qty * currentPrice, margin: qty * currentPrice }
        } else if (orderMode === 'leverage') {
            const margin = inputAmount
            return { shares: (margin * leverage) / currentPrice, value: margin * leverage, margin }
        }
        return { shares: 0, value: 0, margin: 0 }
    }

    const position = getPositionDetails()

    const handleTrade = async () => {
        // Validation
        if (position.margin <= 0) {
            Alert.alert('Invalid Amount', 'Please enter a valid trade amount')
            return
        }
        if (position.margin > balance) {
            Alert.alert('Insufficient Balance', `You only have $${balance.toFixed(2)} available`)
            return
        }
        if (!demoAccount?.accountId) {
            Alert.alert('No Account', 'Please set up an account first')
            return
        }
        if (!activePortfolio?.id) {
            Alert.alert('No Portfolio', 'Please create a portfolio first')
            return
        }

        const orderDetails = orderMode === 'spot'
            ? `Shares: ${position.shares.toFixed(4)}\nTotal: $${position.value.toFixed(2)}`
            : `Margin: $${position.margin.toFixed(2)}\nLeverage: ${leverage}x\nPosition: $${position.value.toFixed(2)}`

        const slTpDetails = enableSLTP
            ? `\n\nStop Loss: ${stopLoss ? `$${stopLoss}` : 'Not set'}\nTake Profit: ${takeProfit ? `$${takeProfit}` : 'Not set'}`
            : ''

        Alert.alert(
            `Confirm ${side.toUpperCase()} Order`,
            `${orderMode.toUpperCase()} ${orderType.toUpperCase()}\n\n` +
            `Symbol: ${symbol}\n` +
            `Price: $${currentPrice.toFixed(2)}\n` +
            orderDetails + slTpDetails,
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Place Order',
                    onPress: async () => {
                        setIsSubmitting(true)
                        try {
                            const orderInput: CreateOrderInput = {
                                accountId: demoAccount.accountId,
                                portfolioId: activePortfolio.id,
                                symbol: symbol,
                                side: side.toUpperCase() as 'BUY' | 'SELL',
                                type: orderType === 'limit' ? 'LIMIT' : 'MARKET',
                                quantity: position.shares,
                                price: orderType === 'limit' ? Number(limitPrice) || currentPrice : undefined,
                            }

                            const order = await orderService.create(orderInput)

                            // Refresh demo account balance
                            await fetchDemo()

                            Alert.alert(
                                'Order Placed! 🎉',
                                `Your ${side.toUpperCase()} order for ${position.shares.toFixed(4)} ${symbol} has been placed!\n\nOrder ID: ${order.id.slice(0, 8)}...`,
                                [{
                                    text: 'OK', onPress: () => {
                                        // Clear form
                                        setAmount('')
                                        setShares('')
                                    }
                                }]
                            )
                        } catch (error: any) {
                            const message = error?.response?.data?.message || error?.message || 'Failed to place order'
                            Alert.alert('Order Failed', message)
                        } finally {
                            setIsSubmitting(false)
                        }
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
                    {/* Header */}
                    <XStack justifyContent="space-between" alignItems="center" padding="$4" paddingBottom="$2">
                        <YStack>
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                                Trade
                            </Text>
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                Balance: ${balance.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                            </Text>
                        </YStack>
                        <View style={styles.demoBadge}>
                            <Text color={dimeTheme.colors.primary} fontSize="$1" fontWeight="bold">
                                DEMO
                            </Text>
                        </View>
                    </XStack>

                    {/* Symbol Selector */}
                    <YStack paddingHorizontal="$4" marginBottom="$3">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Select Symbol
                        </Text>
                        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                            <XStack gap="$2">
                                {DEMO_SYMBOLS.map((s) => (
                                    <TouchableOpacity
                                        key={s}
                                        style={[styles.symbolChip, symbol === s && styles.symbolChipActive]}
                                        onPress={() => setSymbol(s)}
                                    >
                                        <Text color={symbol === s ? '#fff' : dimeTheme.colors.textSecondary} fontWeight="600">
                                            {s}
                                        </Text>
                                    </TouchableOpacity>
                                ))}
                            </XStack>
                        </ScrollView>
                    </YStack>

                    {/* Current Price Card */}
                    <View style={styles.priceCard}>
                        <XStack justifyContent="space-between" alignItems="center">
                            <YStack>
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">Current Price</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontSize="$7" fontWeight="bold">
                                    ${currentPrice.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                                </Text>
                            </YStack>
                            <View style={[styles.changeBadge, { backgroundColor: isPositive ? 'rgba(0,200,83,0.15)' : 'rgba(255,82,82,0.15)' }]}>
                                <Ionicons name={isPositive ? "arrow-up" : "arrow-down"} size={14} color={isPositive ? dimeTheme.colors.profit : dimeTheme.colors.loss} />
                                <Text color={isPositive ? dimeTheme.colors.profit : dimeTheme.colors.loss} fontWeight="bold">
                                    {isPositive ? '+' : ''}{changePercent.toFixed(2)}%
                                </Text>
                            </View>
                        </XStack>
                    </View>

                    {/* Buy/Sell Toggle */}
                    <XStack paddingHorizontal="$4" marginBottom="$3" gap="$2">
                        <TouchableOpacity
                            style={[styles.sideButton, side === 'buy' && styles.buyButtonActive]}
                            onPress={() => setSide('buy')}
                        >
                            <Ionicons name="arrow-up-circle" size={20} color={side === 'buy' ? '#fff' : dimeTheme.colors.profit} />
                            <Text color={side === 'buy' ? '#fff' : dimeTheme.colors.profit} fontWeight="bold" marginLeft={8}>
                                BUY
                            </Text>
                        </TouchableOpacity>
                        <TouchableOpacity
                            style={[styles.sideButton, side === 'sell' && styles.sellButtonActive]}
                            onPress={() => setSide('sell')}
                        >
                            <Ionicons name="arrow-down-circle" size={20} color={side === 'sell' ? '#fff' : dimeTheme.colors.loss} />
                            <Text color={side === 'sell' ? '#fff' : dimeTheme.colors.loss} fontWeight="bold" marginLeft={8}>
                                SELL
                            </Text>
                        </TouchableOpacity>
                    </XStack>

                    {/* Order Mode: Spot / Leverage / Options */}
                    <YStack paddingHorizontal="$4" marginBottom="$3">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Order Mode
                        </Text>
                        <XStack gap="$2">
                            {(['spot', 'leverage', 'options'] as OrderMode[]).map((mode) => (
                                <TouchableOpacity
                                    key={mode}
                                    style={[styles.modeButton, orderMode === mode && styles.modeButtonActive]}
                                    onPress={() => setOrderMode(mode)}
                                >
                                    <Ionicons
                                        name={mode === 'spot' ? 'cash' : mode === 'leverage' ? 'trending-up' : 'options'}
                                        size={16}
                                        color={orderMode === mode ? '#fff' : dimeTheme.colors.textSecondary}
                                    />
                                    <Text
                                        color={orderMode === mode ? '#fff' : dimeTheme.colors.textSecondary}
                                        fontSize="$2"
                                        fontWeight="600"
                                        textTransform="capitalize"
                                        marginLeft={4}
                                    >
                                        {mode}
                                    </Text>
                                </TouchableOpacity>
                            ))}
                        </XStack>
                    </YStack>

                    {/* Order Type: Market / Limit / Scheduled */}
                    <YStack paddingHorizontal="$4" marginBottom="$3">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            Order Type
                        </Text>
                        <XStack gap="$2">
                            {(['market', 'limit', 'scheduled'] as OrderType[]).map((type) => (
                                <TouchableOpacity
                                    key={type}
                                    style={[styles.typeChip, orderType === type && styles.typeChipActive]}
                                    onPress={() => setOrderType(type)}
                                >
                                    <Text
                                        color={orderType === type ? '#fff' : dimeTheme.colors.textSecondary}
                                        fontSize="$2"
                                        fontWeight="600"
                                        textTransform="capitalize"
                                    >
                                        {type}
                                    </Text>
                                </TouchableOpacity>
                            ))}
                        </XStack>
                    </YStack>

                    {/* Leverage Slider (only for leverage mode) */}
                    {orderMode === 'leverage' && (
                        <YStack paddingHorizontal="$4" marginBottom="$3">
                            <XStack justifyContent="space-between" alignItems="center" marginBottom="$2">
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">Leverage</Text>
                                <Text color={dimeTheme.colors.primary} fontWeight="bold">{leverage}x</Text>
                            </XStack>
                            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                                <XStack gap="$2">
                                    {LEVERAGE_OPTIONS.map((lev) => (
                                        <TouchableOpacity
                                            key={lev}
                                            style={[styles.leverageChip, leverage === lev && styles.leverageChipActive]}
                                            onPress={() => setLeverage(lev)}
                                        >
                                            <Text color={leverage === lev ? '#fff' : dimeTheme.colors.textSecondary} fontWeight="600">
                                                {lev}x
                                            </Text>
                                        </TouchableOpacity>
                                    ))}
                                </XStack>
                            </ScrollView>
                        </YStack>
                    )}

                    {/* Amount Input */}
                    <YStack paddingHorizontal="$4" marginBottom="$3">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                            {orderMode === 'spot' ? 'Amount (USD) or Shares' : 'Margin (USD)'}
                        </Text>
                        <XStack gap="$2">
                            <View style={[styles.inputContainer, { flex: 1 }]}>
                                <Text color={dimeTheme.colors.textTertiary}>$</Text>
                                <TextInput
                                    style={styles.input}
                                    value={amount}
                                    onChangeText={setAmount}
                                    placeholder="0.00"
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                    keyboardType="decimal-pad"
                                />
                            </View>
                            {orderMode === 'spot' && (
                                <View style={[styles.inputContainer, { flex: 1 }]}>
                                    <Ionicons name="layers" size={16} color={dimeTheme.colors.textTertiary} />
                                    <TextInput
                                        style={styles.input}
                                        value={shares}
                                        onChangeText={setShares}
                                        placeholder="Shares"
                                        placeholderTextColor={dimeTheme.colors.textTertiary}
                                        keyboardType="decimal-pad"
                                    />
                                </View>
                            )}
                        </XStack>
                    </YStack>

                    {/* Limit Price (for limit orders) */}
                    {orderType === 'limit' && (
                        <YStack paddingHorizontal="$4" marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                                Limit Price
                            </Text>
                            <View style={styles.inputContainer}>
                                <Text color={dimeTheme.colors.textTertiary}>$</Text>
                                <TextInput
                                    style={styles.input}
                                    value={limitPrice}
                                    onChangeText={setLimitPrice}
                                    placeholder={currentPrice.toFixed(2)}
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                    keyboardType="decimal-pad"
                                />
                            </View>
                        </YStack>
                    )}

                    {/* Scheduled Date (for scheduled orders) */}
                    {orderType === 'scheduled' && (
                        <YStack paddingHorizontal="$4" marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">
                                Schedule Time
                            </Text>
                            <TouchableOpacity style={styles.dateButton} onPress={() => setShowDatePicker(true)}>
                                <Ionicons name="calendar" size={20} color={dimeTheme.colors.primary} />
                                <Text color={dimeTheme.colors.textPrimary} marginLeft={8}>
                                    {scheduledDate.toLocaleString()}
                                </Text>
                            </TouchableOpacity>
                        </YStack>
                    )}

                    {/* SL/TP Toggle */}
                    <XStack paddingHorizontal="$4" marginBottom="$2" justifyContent="space-between" alignItems="center">
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                            Stop Loss / Take Profit
                        </Text>
                        <Switch
                            value={enableSLTP}
                            onValueChange={setEnableSLTP}
                            trackColor={{ false: dimeTheme.colors.surface, true: dimeTheme.colors.primary }}
                            thumbColor="#fff"
                        />
                    </XStack>

                    {enableSLTP && (
                        <XStack paddingHorizontal="$4" marginBottom="$3" gap="$2">
                            <View style={[styles.inputContainer, styles.slInput, { flex: 1 }]}>
                                <Text color={dimeTheme.colors.loss} fontSize="$1" fontWeight="600">SL</Text>
                                <TextInput
                                    style={[styles.input, { textAlign: 'right' }]}
                                    value={stopLoss}
                                    onChangeText={setStopLoss}
                                    placeholder="Stop Loss"
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                    keyboardType="decimal-pad"
                                />
                            </View>
                            <View style={[styles.inputContainer, styles.tpInput, { flex: 1 }]}>
                                <Text color={dimeTheme.colors.profit} fontSize="$1" fontWeight="600">TP</Text>
                                <TextInput
                                    style={[styles.input, { textAlign: 'right' }]}
                                    value={takeProfit}
                                    onChangeText={setTakeProfit}
                                    placeholder="Take Profit"
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                    keyboardType="decimal-pad"
                                />
                            </View>
                        </XStack>
                    )}

                    {/* Order Summary */}
                    <View style={styles.summaryCard}>
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2" marginBottom="$2">Order Summary</Text>
                        <XStack justifyContent="space-between" marginBottom="$1">
                            <Text color={dimeTheme.colors.textTertiary}>Shares</Text>
                            <Text color={dimeTheme.colors.textPrimary} fontWeight="600">{position.shares.toFixed(4)}</Text>
                        </XStack>
                        <XStack justifyContent="space-between" marginBottom="$1">
                            <Text color={dimeTheme.colors.textTertiary}>Position Value</Text>
                            <Text color={dimeTheme.colors.textPrimary} fontWeight="600">${position.value.toFixed(2)}</Text>
                        </XStack>
                        {orderMode === 'leverage' && (
                            <XStack justifyContent="space-between" marginBottom="$1">
                                <Text color={dimeTheme.colors.textTertiary}>Required Margin</Text>
                                <Text color={dimeTheme.colors.primary} fontWeight="600">${position.margin.toFixed(2)}</Text>
                            </XStack>
                        )}
                    </View>

                    {/* Place Order Button */}
                    <View style={{ paddingHorizontal: 16, paddingBottom: 32 }}>
                        <TouchableOpacity
                            style={[
                                styles.orderButton,
                                { backgroundColor: side === 'buy' ? dimeTheme.colors.profit : dimeTheme.colors.loss }
                            ]}
                            onPress={handleTrade}
                        >
                            <Text color="#fff" fontSize="$5" fontWeight="bold">
                                {side === 'buy' ? '🚀 Place Buy Order' : '📉 Place Sell Order'}
                            </Text>
                        </TouchableOpacity>
                    </View>

                    <View style={{ height: 100 }} />
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
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
    },
    symbolChip: {
        paddingHorizontal: 16,
        paddingVertical: 10,
        borderRadius: 20,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    symbolChipActive: {
        backgroundColor: dimeTheme.colors.primary,
        borderColor: dimeTheme.colors.primary,
    },
    priceCard: {
        marginHorizontal: 16,
        marginBottom: 16,
        padding: 16,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 16,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    changeBadge: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: 12,
        paddingVertical: 8,
        borderRadius: 12,
        gap: 4,
    },
    sideButton: {
        flex: 1,
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        paddingVertical: 14,
        borderRadius: 12,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    buyButtonActive: {
        backgroundColor: dimeTheme.colors.profit,
        borderColor: dimeTheme.colors.profit,
    },
    sellButtonActive: {
        backgroundColor: dimeTheme.colors.loss,
        borderColor: dimeTheme.colors.loss,
    },
    modeButton: {
        flex: 1,
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        paddingVertical: 12,
        borderRadius: 10,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    modeButtonActive: {
        backgroundColor: dimeTheme.colors.primary,
        borderColor: dimeTheme.colors.primary,
    },
    typeChip: {
        paddingHorizontal: 16,
        paddingVertical: 10,
        borderRadius: 20,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    typeChipActive: {
        backgroundColor: dimeTheme.colors.primary,
        borderColor: dimeTheme.colors.primary,
    },
    leverageChip: {
        paddingHorizontal: 16,
        paddingVertical: 10,
        borderRadius: 8,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    leverageChipActive: {
        backgroundColor: dimeTheme.colors.primary,
        borderColor: dimeTheme.colors.primary,
    },
    inputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 12,
        paddingHorizontal: 16,
        paddingVertical: 14,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        gap: 8,
    },
    input: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 16,
        fontWeight: '600',
    },
    slInput: {
        borderColor: 'rgba(255, 82, 82, 0.3)',
    },
    tpInput: {
        borderColor: 'rgba(0, 200, 83, 0.3)',
    },
    dateButton: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 12,
        paddingHorizontal: 16,
        paddingVertical: 14,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    summaryCard: {
        marginHorizontal: 16,
        marginBottom: 16,
        padding: 16,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 16,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    orderButton: {
        alignItems: 'center',
        justifyContent: 'center',
        paddingVertical: 18,
        borderRadius: 16,
    },
})
