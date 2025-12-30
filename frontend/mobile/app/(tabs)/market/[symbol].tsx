import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, Dimensions, TextInput, Alert, Modal, TouchableOpacity, ActivityIndicator } from 'react-native'
import { YStack, XStack, Text, View, Button, Spinner } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams, useRouter } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { useMarketStore } from '@/stores/market'
import { useDemoStore } from '@/stores/demo'
import { usePortfolioStore } from '@/stores/portfolio'
import { instrumentService, orderService } from '@/services/api'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { TradingViewChart } from '@/components/chart/TradingViewChart'
import type { CandleData, Instrument, Quote } from '@/types/market'
import type { CreateOrderInput, OrderSide } from '@/types/trade'

const { width: screenWidth } = Dimensions.get('window')

// Time period options with their resolution and days
// Yahoo Finance free: Daily data only (no intraday for free tier)
const TIME_PERIODS = [
    { label: '1W', resolution: 'D', days: 7 },       // Daily for 1 week
    { label: '1M', resolution: 'D', days: 30 },
    { label: '3M', resolution: 'D', days: 90 },
    { label: '1Y', resolution: 'D', days: 365 },
    { label: '5Y', resolution: 'W', days: 365 * 5 },  // Weekly for 5 years
    { label: 'ALL', resolution: 'M', days: 365 * 20 }, // Monthly for max
] as const

type TimePeriod = typeof TIME_PERIODS[number]['label']

export default function SymbolDetailScreen() {
    const { symbol: rawSymbol } = useLocalSearchParams<{ symbol: string }>()
    const router = useRouter()
    const { instruments } = useMarketStore()
    const { account: demoAccount, fetchDemo } = useDemoStore()
    const { activePortfolio, fetchPortfolios } = usePortfolioStore()

    // Decode URL-encoded symbol
    const symbol = rawSymbol ? decodeURIComponent(rawSymbol) : ''
    const chartSymbol = symbol.replace('/', '-')

    const [selectedPeriod, setSelectedPeriod] = useState<TimePeriod>('1M')
    const [candles, setCandles] = useState<CandleData[]>([])
    const [quote, setQuote] = useState<Quote | null>(null)
    const [isLoadingChart, setIsLoadingChart] = useState(true)
    const [isLoadingQuote, setIsLoadingQuote] = useState(true)
    const [chartType, setChartType] = useState<'candlestick' | 'line' | 'area'>('area')

    // Order modal state
    const [showOrderModal, setShowOrderModal] = useState(false)
    const [orderSide, setOrderSide] = useState<OrderSide>('BUY')
    const [orderAmount, setOrderAmount] = useState('')
    const [orderType, setOrderType] = useState<'MARKET' | 'LIMIT'>('MARKET')
    const [limitPrice, setLimitPrice] = useState('')
    const [enableSLTP, setEnableSLTP] = useState(false)
    const [stopLoss, setStopLoss] = useState('')
    const [takeProfit, setTakeProfit] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)

    // Enhanced order options
    const [amountUnit, setAmountUnit] = useState<'USD' | 'SHARES'>('USD')
    const [leverage, setLeverage] = useState<number>(1)
    const LEVERAGE_OPTIONS = [1, 10, 50, 100]

    const instrument = instruments.find(i => i.symbol === symbol)

    // Ensure portfolio is loaded on mount
    useEffect(() => {
        if (!activePortfolio) {
            fetchPortfolios()
        }
    }, [activePortfolio, fetchPortfolios])

    // Fetch quote once on mount (with error handling)
    useEffect(() => {
        const loadQuote = async () => {
            if (!symbol) return
            setIsLoadingQuote(true)
            try {
                const data = await instrumentService.getQuote(symbol)
                setQuote(data)
            } catch (error) {
                console.log('Quote not available for', symbol)
                // Silently fail - we'll show candle data instead
            }
            setIsLoadingQuote(false)
        }
        loadQuote()
    }, [symbol])

    // Fetch candles when period changes
    useEffect(() => {
        const loadCandles = async () => {
            if (!symbol) return
            setIsLoadingChart(true)
            try {
                const period = TIME_PERIODS.find(p => p.label === selectedPeriod)!
                const to = Math.floor(Date.now() / 1000)
                const from = to - (period.days * 24 * 60 * 60)

                const data = await instrumentService.getCandles(symbol, period.resolution, from, to)
                setCandles(data)
            } catch (error) {
                console.log('Chart data not available for', symbol)
                setCandles([])
            }
            setIsLoadingChart(false)
        }
        loadCandles()
    }, [symbol, selectedPeriod])

    // Calculate stats from candles
    const high = candles.length > 0 ? Math.max(...candles.map(c => c.high)) : quote?.high ?? 0
    const low = candles.length > 0 ? Math.min(...candles.map(c => c.low)) : quote?.low ?? 0
    const open = candles.length > 0 ? candles[0].open : quote?.open ?? 0
    const volume = candles.length > 0 ? candles.reduce((sum, c) => sum + c.volume, 0) : 0

    const formatVolume = (vol: number) => {
        if (vol >= 1e9) return (vol / 1e9).toFixed(1) + 'B'
        if (vol >= 1e6) return (vol / 1e6).toFixed(1) + 'M'
        if (vol >= 1e3) return (vol / 1e3).toFixed(1) + 'K'
        return vol.toString()
    }

    const currentPrice = quote?.price ?? 0
    const balance = demoAccount?.balance ?? 0

    const openOrderModal = (side: OrderSide) => {
        setOrderSide(side)
        setOrderAmount('')
        setShowOrderModal(true)
    }

    const handlePlaceOrder = async () => {
        const amount = parseFloat(orderAmount)
        if (!amount || amount <= 0) {
            Alert.alert('Invalid Amount', 'Please enter a valid amount')
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

        // Calculate shares based on unit type
        const shares = amountUnit === 'USD' ? amount / currentPrice : amount
        const totalCost = shares * currentPrice
        const marginRequired = totalCost / leverage // Leverage reduces margin needed

        if (orderSide === 'BUY' && marginRequired > balance) {
            Alert.alert('Insufficient Balance', `You need $${marginRequired.toFixed(2)} margin (${leverage}x leverage)`)
            return
        }

        setIsSubmitting(true)
        try {
            const orderInput: CreateOrderInput = {
                accountId: demoAccount.accountId,
                portfolioId: activePortfolio.id,
                symbol: symbol,
                side: orderSide,
                type: orderType,
                quantity: shares,
                price: currentPrice, // Send actual price for execution
            }

            const order = await orderService.create(orderInput)

            // Refresh balance and positions
            await fetchDemo()
            await fetchPortfolios() // Refresh portfolio to see new position

            setShowOrderModal(false)
            const leverageText = leverage > 1 ? ` (${leverage}x leverage)` : ''
            Alert.alert(
                'Order Placed! 🎉',
                `${orderSide} order for ${shares.toFixed(4)} ${symbol}${leverageText}\n\nTotal: $${totalCost.toFixed(2)}`,
                [{ text: 'OK' }]
            )
        } catch (error: any) {
            const message = error?.response?.data?.message || error?.message || 'Failed to place order'
            Alert.alert('Order Failed', message)
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                {/* Header */}
                <XStack padding="$4" alignItems="center" gap="$3">
                    <Button
                        size="$3"
                        circular
                        backgroundColor={dimeTheme.colors.surface}
                        onPress={() => {
                            if (router.canGoBack()) {
                                router.back()
                            } else {
                                router.replace('/(tabs)/market')
                            }
                        }}
                    >
                        <Ionicons name="arrow-back" size={20} color={dimeTheme.colors.textPrimary} />
                    </Button>
                    <YStack flex={1}>
                        <XStack alignItems="center" gap="$2">
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$6" fontWeight="bold">
                                {symbol?.replace('/USD', '')}
                            </Text>
                            {instrument && (
                                <View style={styles.typeBadge}>
                                    <Text color={dimeTheme.colors.primary} fontSize={10} fontWeight="600">
                                        {instrument.type}
                                    </Text>
                                </View>
                            )}
                        </XStack>
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                            {instrument?.name ?? symbol}
                        </Text>
                    </YStack>
                    <Button size="$3" circular backgroundColor={dimeTheme.colors.surface}>
                        <Ionicons name="star-outline" size={20} color={dimeTheme.colors.primary} />
                    </Button>
                </XStack>

                <ScrollView showsVerticalScrollIndicator={false}>
                    {/* Price Section */}
                    <YStack paddingHorizontal="$4" marginBottom="$3">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$9" fontWeight="bold">
                            ${(quote?.price ?? 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                        </Text>
                        <XStack alignItems="center" gap="$2" marginTop="$1">
                            <PriceChip value={quote?.changePercent ?? 0} size="md" />
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                                {(quote?.change ?? 0) >= 0 ? '+' : ''}${(quote?.change ?? 0).toFixed(2)} today
                            </Text>
                        </XStack>
                    </YStack>

                    {/* Chart Type Selector */}
                    <XStack paddingHorizontal="$4" marginBottom="$2" gap="$2">
                        {(['candlestick', 'line', 'area'] as const).map(type => (
                            <Button
                                key={type}
                                size="$2"
                                backgroundColor={chartType === type ? dimeTheme.colors.primary + '30' : 'transparent'}
                                borderWidth={1}
                                borderColor={chartType === type ? dimeTheme.colors.primary : dimeTheme.colors.border}
                                onPress={() => setChartType(type)}
                            >
                                <Ionicons
                                    name={type === 'candlestick' ? 'bar-chart' : type === 'line' ? 'trending-up' : 'analytics'}
                                    size={14}
                                    color={chartType === type ? dimeTheme.colors.primary : dimeTheme.colors.textSecondary}
                                />
                            </Button>
                        ))}
                    </XStack>

                    {/* Time Period Selector */}
                    <ScrollView
                        horizontal
                        showsHorizontalScrollIndicator={false}
                        contentContainerStyle={styles.periodContainer}
                    >
                        {TIME_PERIODS.map(period => (
                            <Button
                                key={period.label}
                                size="$3"
                                backgroundColor={selectedPeriod === period.label ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                borderWidth={selectedPeriod === period.label ? 0 : 1}
                                borderColor={dimeTheme.colors.border}
                                pressStyle={{ opacity: 0.8 }}
                                marginRight="$2"
                                onPress={() => setSelectedPeriod(period.label)}
                            >
                                <Text
                                    color={selectedPeriod === period.label ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                    fontWeight="600"
                                >
                                    {period.label}
                                </Text>
                            </Button>
                        ))}
                    </ScrollView>

                    {/* TradingView Chart */}
                    <YStack marginVertical="$2">
                        {isLoadingChart ? (
                            <View style={styles.chartLoading}>
                                <Spinner size="large" color={dimeTheme.colors.primary} />
                                <Text color={dimeTheme.colors.textSecondary} marginTop="$2">
                                    Loading chart data...
                                </Text>
                            </View>
                        ) : candles.length > 0 ? (
                            <TradingViewChart
                                candles={candles}
                                symbol={symbol ?? ''}
                                height={320}
                                chartType={chartType}
                            />
                        ) : (
                            <View style={styles.chartLoading}>
                                <Ionicons name="bar-chart-outline" size={48} color={dimeTheme.colors.textTertiary} />
                                <Text color={dimeTheme.colors.textSecondary} marginTop="$2">
                                    No chart data available
                                </Text>
                            </View>
                        )}
                    </YStack>

                    {/* Stats Grid */}
                    <YStack paddingHorizontal="$4" marginTop="$3">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold" marginBottom="$3">
                            Statistics
                        </Text>
                        <View style={styles.statsGrid}>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Open</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    ${open.toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">High</Text>
                                <Text color={dimeTheme.colors.profit} fontWeight="600">
                                    ${high.toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Low</Text>
                                <Text color={dimeTheme.colors.loss} fontWeight="600">
                                    ${low.toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Volume</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    {formatVolume(volume)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Prev Close</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    ${(quote?.previousClose ?? 0).toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Exchange</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    {instrument?.exchange ?? '-'}
                                </Text>
                            </View>
                        </View>
                    </YStack>

                    {/* Buy/Sell Buttons */}
                    <XStack padding="$4" gap="$3" marginTop="$3">
                        <Button
                            flex={1}
                            size="$5"
                            backgroundColor={dimeTheme.colors.profit}
                            pressStyle={{ opacity: 0.9 }}
                            onPress={() => openOrderModal('BUY')}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons name="arrow-up-circle" size={20} color={dimeTheme.colors.background} />
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                    BUY
                                </Text>
                            </XStack>
                        </Button>
                        <Button
                            flex={1}
                            size="$5"
                            backgroundColor={dimeTheme.colors.loss}
                            pressStyle={{ opacity: 0.9 }}
                            onPress={() => openOrderModal('SELL')}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons name="arrow-down-circle" size={20} color={dimeTheme.colors.background} />
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                    SELL
                                </Text>
                            </XStack>
                        </Button>
                    </XStack>
                </ScrollView>
            </SafeAreaView>

            {/* Order Modal */}
            <Modal
                visible={showOrderModal}
                transparent
                animationType="slide"
                onRequestClose={() => setShowOrderModal(false)}
            >
                <View style={styles.modalOverlay}>
                    <View style={styles.modalContent}>
                        {/* Modal Header */}
                        <XStack justifyContent="space-between" alignItems="center" marginBottom="$4">
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$6" fontWeight="bold">
                                {orderSide} {symbol}
                            </Text>
                            <TouchableOpacity onPress={() => setShowOrderModal(false)}>
                                <Ionicons name="close" size={24} color={dimeTheme.colors.textSecondary} />
                            </TouchableOpacity>
                        </XStack>

                        {/* Price Info */}
                        <YStack marginBottom="$4">
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$3">Current Price</Text>
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$7" fontWeight="bold">
                                ${currentPrice.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                            </Text>
                        </YStack>

                        {/* Balance */}
                        <XStack justifyContent="space-between" marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary}>Available Balance</Text>
                            <Text color={dimeTheme.colors.primary} fontWeight="600">${balance.toFixed(2)}</Text>
                        </XStack>

                        {/* Order Type Selector */}
                        <YStack marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} marginBottom="$2">Order Type</Text>
                            <XStack gap="$2">
                                <TouchableOpacity
                                    style={[styles.orderTypeButton, orderType === 'MARKET' && styles.orderTypeActive]}
                                    onPress={() => setOrderType('MARKET')}
                                >
                                    <Text color={orderType === 'MARKET' ? dimeTheme.colors.primary : dimeTheme.colors.textSecondary} fontWeight="600">
                                        Market
                                    </Text>
                                </TouchableOpacity>
                                <TouchableOpacity
                                    style={[styles.orderTypeButton, orderType === 'LIMIT' && styles.orderTypeActive]}
                                    onPress={() => setOrderType('LIMIT')}
                                >
                                    <Text color={orderType === 'LIMIT' ? dimeTheme.colors.primary : dimeTheme.colors.textSecondary} fontWeight="600">
                                        Limit
                                    </Text>
                                </TouchableOpacity>
                            </XStack>
                        </YStack>

                        {/* Limit Price (if LIMIT order) */}
                        {orderType === 'LIMIT' && (
                            <YStack marginBottom="$3">
                                <Text color={dimeTheme.colors.textSecondary} marginBottom="$2">Limit Price</Text>
                                <View style={styles.amountInputContainer}>
                                    <Text color={dimeTheme.colors.textSecondary} fontSize="$5">$</Text>
                                    <TextInput
                                        style={styles.amountInput}
                                        value={limitPrice}
                                        onChangeText={setLimitPrice}
                                        placeholder={currentPrice.toFixed(2)}
                                        placeholderTextColor={dimeTheme.colors.textTertiary}
                                        keyboardType="decimal-pad"
                                    />
                                </View>
                            </YStack>
                        )}

                        {/* Unit Toggle (USD / Shares) */}
                        <YStack marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} marginBottom="$2">Amount in</Text>
                            <XStack gap="$2">
                                <TouchableOpacity
                                    style={[styles.orderTypeButton, amountUnit === 'USD' && styles.orderTypeActive]}
                                    onPress={() => setAmountUnit('USD')}
                                >
                                    <Text color={amountUnit === 'USD' ? dimeTheme.colors.primary : dimeTheme.colors.textSecondary} fontWeight="600">
                                        💵 USD
                                    </Text>
                                </TouchableOpacity>
                                <TouchableOpacity
                                    style={[styles.orderTypeButton, amountUnit === 'SHARES' && styles.orderTypeActive]}
                                    onPress={() => setAmountUnit('SHARES')}
                                >
                                    <Text color={amountUnit === 'SHARES' ? dimeTheme.colors.primary : dimeTheme.colors.textSecondary} fontWeight="600">
                                        📦 Shares
                                    </Text>
                                </TouchableOpacity>
                            </XStack>
                        </YStack>

                        {/* Leverage Selector */}
                        <YStack marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} marginBottom="$2">Leverage</Text>
                            <XStack gap="$2">
                                {LEVERAGE_OPTIONS.map((lev) => (
                                    <TouchableOpacity
                                        key={lev}
                                        style={[styles.leverageButton, leverage === lev && styles.leverageActive]}
                                        onPress={() => setLeverage(lev)}
                                    >
                                        <Text color={leverage === lev ? dimeTheme.colors.background : dimeTheme.colors.textSecondary} fontWeight="600" fontSize={12}>
                                            {lev}x
                                        </Text>
                                    </TouchableOpacity>
                                ))}
                            </XStack>
                        </YStack>

                        {/* Amount Input */}
                        <YStack marginBottom="$3">
                            <Text color={dimeTheme.colors.textSecondary} marginBottom="$2">
                                {amountUnit === 'USD' ? 'Amount (USD)' : 'Shares'}
                            </Text>
                            <View style={styles.amountInputContainer}>
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$5">
                                    {amountUnit === 'USD' ? '$' : '#'}
                                </Text>
                                <TextInput
                                    style={styles.amountInput}
                                    value={orderAmount}
                                    onChangeText={setOrderAmount}
                                    placeholder="0.00"
                                    placeholderTextColor={dimeTheme.colors.textTertiary}
                                    keyboardType="decimal-pad"
                                    autoFocus
                                />
                            </View>
                        </YStack>

                        {/* SL/TP Toggle */}
                        <XStack justifyContent="space-between" alignItems="center" marginBottom="$3">
                            <Text color={dimeTheme.colors.textPrimary} fontWeight="600">Stop Loss / Take Profit</Text>
                            <TouchableOpacity
                                style={[styles.toggle, enableSLTP && styles.toggleActive]}
                                onPress={() => setEnableSLTP(!enableSLTP)}
                            >
                                <View style={[styles.toggleKnob, enableSLTP && styles.toggleKnobActive]} />
                            </TouchableOpacity>
                        </XStack>

                        {/* SL/TP Inputs */}
                        {enableSLTP && (
                            <XStack gap="$3" marginBottom="$3">
                                <YStack flex={1}>
                                    <Text color={dimeTheme.colors.loss} fontSize="$2" marginBottom="$1">Stop Loss</Text>
                                    <View style={[styles.slTpInput, { borderColor: dimeTheme.colors.loss + '50' }]}>
                                        <Text color={dimeTheme.colors.loss}>$</Text>
                                        <TextInput
                                            style={styles.slTpInputText}
                                            value={stopLoss}
                                            onChangeText={setStopLoss}
                                            placeholder="0.00"
                                            placeholderTextColor={dimeTheme.colors.textTertiary}
                                            keyboardType="decimal-pad"
                                        />
                                    </View>
                                </YStack>
                                <YStack flex={1}>
                                    <Text color={dimeTheme.colors.profit} fontSize="$2" marginBottom="$1">Take Profit</Text>
                                    <View style={[styles.slTpInput, { borderColor: dimeTheme.colors.profit + '50' }]}>
                                        <Text color={dimeTheme.colors.profit}>$</Text>
                                        <TextInput
                                            style={styles.slTpInputText}
                                            value={takeProfit}
                                            onChangeText={setTakeProfit}
                                            placeholder="0.00"
                                            placeholderTextColor={dimeTheme.colors.textTertiary}
                                            keyboardType="decimal-pad"
                                        />
                                    </View>
                                </YStack>
                            </XStack>
                        )}

                        {/* Order Summary */}
                        {orderAmount && parseFloat(orderAmount) > 0 && (
                            <YStack backgroundColor={dimeTheme.colors.surface} padding="$3" borderRadius={12} marginBottom="$4">
                                <XStack justifyContent="space-between" marginBottom="$2">
                                    <Text color={dimeTheme.colors.textSecondary}>Shares</Text>
                                    <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                        {(parseFloat(orderAmount) / currentPrice).toFixed(4)}
                                    </Text>
                                </XStack>
                                <XStack justifyContent="space-between" marginBottom="$2">
                                    <Text color={dimeTheme.colors.textSecondary}>Price</Text>
                                    <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                        ${orderType === 'LIMIT' && limitPrice ? parseFloat(limitPrice).toFixed(2) : currentPrice.toFixed(2)}
                                    </Text>
                                </XStack>
                                {enableSLTP && stopLoss && (
                                    <XStack justifyContent="space-between" marginBottom="$2">
                                        <Text color={dimeTheme.colors.loss}>Stop Loss</Text>
                                        <Text color={dimeTheme.colors.loss} fontWeight="600">${stopLoss}</Text>
                                    </XStack>
                                )}
                                {enableSLTP && takeProfit && (
                                    <XStack justifyContent="space-between" marginBottom="$2">
                                        <Text color={dimeTheme.colors.profit}>Take Profit</Text>
                                        <Text color={dimeTheme.colors.profit} fontWeight="600">${takeProfit}</Text>
                                    </XStack>
                                )}
                                <View style={styles.divider} />
                                <XStack justifyContent="space-between">
                                    <Text color={dimeTheme.colors.textPrimary} fontWeight="bold">Total</Text>
                                    <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize="$5">
                                        ${parseFloat(orderAmount).toFixed(2)}
                                    </Text>
                                </XStack>
                            </YStack>
                        )}

                        {/* Place Order Button */}
                        <Button
                            size="$5"
                            backgroundColor={orderSide === 'BUY' ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                            pressStyle={{ opacity: 0.9 }}
                            onPress={handlePlaceOrder}
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? (
                                <ActivityIndicator color={dimeTheme.colors.background} />
                            ) : (
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$5">
                                    {orderSide} {symbol}
                                </Text>
                            )}
                        </Button>
                    </View>
                </View>
            </Modal>
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
    typeBadge: {
        backgroundColor: dimeTheme.colors.primary + '20',
        paddingHorizontal: 8,
        paddingVertical: 2,
        borderRadius: 4,
    },
    periodContainer: {
        paddingHorizontal: 16,
        marginBottom: 8,
    },
    chartLoading: {
        height: 320,
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: dimeTheme.colors.surface,
        marginHorizontal: 16,
        borderRadius: dimeTheme.radius.lg,
    },
    statsGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: dimeTheme.radius.lg,
        padding: 16,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    statItem: {
        width: '50%',
        paddingVertical: 8,
    },
    modalOverlay: {
        flex: 1,
        backgroundColor: 'rgba(0, 0, 0, 0.7)',
        justifyContent: 'flex-end',
    },
    modalContent: {
        backgroundColor: dimeTheme.colors.background,
        borderTopLeftRadius: 24,
        borderTopRightRadius: 24,
        padding: 24,
        paddingBottom: 40,
    },
    amountInputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 12,
        paddingHorizontal: 16,
        paddingVertical: 14,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    amountInput: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 24,
        fontWeight: '600',
        marginLeft: 8,
    },
    orderTypeButton: {
        flex: 1,
        paddingVertical: 10,
        paddingHorizontal: 16,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 8,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    orderTypeActive: {
        borderColor: dimeTheme.colors.primary,
        backgroundColor: dimeTheme.colors.primary + '15',
    },
    toggle: {
        width: 50,
        height: 28,
        borderRadius: 14,
        backgroundColor: dimeTheme.colors.surface,
        padding: 2,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    toggleActive: {
        backgroundColor: dimeTheme.colors.primary,
        borderColor: dimeTheme.colors.primary,
    },
    toggleKnob: {
        width: 22,
        height: 22,
        borderRadius: 11,
        backgroundColor: dimeTheme.colors.textTertiary,
    },
    toggleKnobActive: {
        backgroundColor: dimeTheme.colors.background,
        marginLeft: 'auto',
    },
    slTpInput: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 8,
        paddingHorizontal: 12,
        paddingVertical: 10,
        borderWidth: 1,
    },
    slTpInputText: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 16,
        fontWeight: '600',
        marginLeft: 6,
    },
    divider: {
        height: 1,
        backgroundColor: dimeTheme.colors.border,
        marginVertical: 12,
    },
    leverageButton: {
        flex: 1,
        paddingVertical: 8,
        paddingHorizontal: 12,
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: 8,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    leverageActive: {
        borderColor: dimeTheme.colors.primary,
        backgroundColor: dimeTheme.colors.primary,
    },
})
