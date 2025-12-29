import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, Dimensions } from 'react-native'
import { YStack, XStack, Text, View, Button, Spinner } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams, useRouter } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { useMarketStore } from '@/stores/market'
import { instrumentService } from '@/services/api'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { TradingViewChart } from '@/components/chart/TradingViewChart'
import type { CandleData, Instrument, Quote } from '@/types/market'

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

    // Decode URL-encoded symbol and create chart symbol for Yahoo Finance
    const symbol = rawSymbol ? decodeURIComponent(rawSymbol) : ''
    // Convert crypto symbols for Yahoo Finance (BTC/USD -> BTC-USD)
    const chartSymbol = symbol.replace('/', '-')

    const [selectedPeriod, setSelectedPeriod] = useState<TimePeriod>('1M')
    const [candles, setCandles] = useState<CandleData[]>([])
    const [quote, setQuote] = useState<Quote | null>(null)
    const [isLoadingChart, setIsLoadingChart] = useState(true)
    const [isLoadingQuote, setIsLoadingQuote] = useState(true)
    const [chartType, setChartType] = useState<'candlestick' | 'line' | 'area'>('area')

    const instrument = instruments.find(i => i.symbol === symbol)

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
                            onPress={() => router.push({ pathname: '/(tabs)/trade', params: { symbol, side: 'long' } })}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons name="trending-up" size={18} color={dimeTheme.colors.background} />
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                    Long
                                </Text>
                            </XStack>
                        </Button>
                        <Button
                            flex={1}
                            size="$5"
                            backgroundColor={dimeTheme.colors.loss}
                            pressStyle={{ opacity: 0.9 }}
                            onPress={() => router.push({ pathname: '/(tabs)/trade', params: { symbol, side: 'short' } })}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons name="trending-down" size={18} color={dimeTheme.colors.background} />
                                <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                    Short
                                </Text>
                            </XStack>
                        </Button>
                    </XStack>
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
})
