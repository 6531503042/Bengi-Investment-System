import { useState, useEffect, useMemo } from 'react'
import { StyleSheet, ScrollView, StatusBar, Dimensions } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useLocalSearchParams, useRouter } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { useMarketStore } from '@/stores/market'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { KLineChart, type KLineDataPoint } from '@/components/charts/KLineChart'

const { width: screenWidth } = Dimensions.get('window')

// Time period options
const TIME_PERIODS = ['1D', '1W', '1M', '3M', '1Y', 'ALL'] as const
type TimePeriod = typeof TIME_PERIODS[number]

export default function SymbolDetailScreen() {
    const { symbol } = useLocalSearchParams<{ symbol: string }>()
    const router = useRouter()
    const { quotes } = useMarketStore()
    const [selectedPeriod, setSelectedPeriod] = useState<TimePeriod>('1D')

    const quote = quotes[symbol ?? '']

    // Generate mock chart data
    const chartData = useMemo<KLineDataPoint[]>(() => {
        const data: KLineDataPoint[] = []
        const basePrice = quote?.price ?? 150
        const now = Date.now()
        const dataPoints = selectedPeriod === '1D' ? 78 :
            selectedPeriod === '1W' ? 35 :
                selectedPeriod === '1M' ? 22 :
                    selectedPeriod === '3M' ? 65 :
                        selectedPeriod === '1Y' ? 252 : 500

        for (let i = 0; i < dataPoints; i++) {
            const timestamp = now - (dataPoints - i) * 60000 * 5
            const volatility = 0.02
            const trend = Math.sin(i / 10) * 0.01
            const open = basePrice * (1 + (Math.random() - 0.5) * volatility + trend)
            const close = open * (1 + (Math.random() - 0.5) * volatility)
            const high = Math.max(open, close) * (1 + Math.random() * volatility * 0.5)
            const low = Math.min(open, close) * (1 - Math.random() * volatility * 0.5)
            const vol = Math.random() * 10000000 + 5000000

            data.push({
                time: timestamp,
                open: +open.toFixed(2),
                high: +high.toFixed(2),
                low: +low.toFixed(2),
                close: +close.toFixed(2),
                volume: Math.round(vol),
            })
        }
        return data
    }, [quote?.price, selectedPeriod])

    // Company info
    const companyNames: Record<string, string> = {
        AAPL: 'Apple Inc.',
        GOOGL: 'Alphabet Inc.',
        MSFT: 'Microsoft Corp.',
        AMZN: 'Amazon.com Inc.',
        TSLA: 'Tesla Inc.',
        NVDA: 'NVIDIA Corp.',
        META: 'Meta Platforms Inc.',
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
                        onPress={() => router.back()}
                    >
                        <Ionicons name="arrow-back" size={20} color={dimeTheme.colors.textPrimary} />
                    </Button>
                    <YStack flex={1}>
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$6" fontWeight="bold">
                            {symbol}
                        </Text>
                        <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                            {companyNames[symbol ?? ''] ?? symbol}
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
                            ${(quote?.price ?? 0).toFixed(2)}
                        </Text>
                        <XStack alignItems="center" gap="$2" marginTop="$1">
                            <PriceChip value={quote?.changePercent ?? 0} size="md" />
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                                {(quote?.change ?? 0) >= 0 ? '+' : ''}${(quote?.change ?? 0).toFixed(2)} today
                            </Text>
                        </XStack>
                    </YStack>

                    {/* Time Period Selector */}
                    <ScrollView
                        horizontal
                        showsHorizontalScrollIndicator={false}
                        contentContainerStyle={styles.periodContainer}
                    >
                        {TIME_PERIODS.map(period => (
                            <Button
                                key={period}
                                size="$3"
                                backgroundColor={selectedPeriod === period ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                                pressStyle={{ opacity: 0.8 }}
                                marginRight="$2"
                                onPress={() => setSelectedPeriod(period)}
                            >
                                <Text
                                    color={selectedPeriod === period ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                    fontWeight="600"
                                >
                                    {period}
                                </Text>
                            </Button>
                        ))}
                    </ScrollView>

                    {/* Chart */}
                    <YStack marginVertical="$2">
                        <KLineChart
                            data={chartData}
                            height={300}
                            showVolume={true}
                        />
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
                                    ${(quote?.price ?? 0 * 0.995).toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">High</Text>
                                <Text color={dimeTheme.colors.profit} fontWeight="600">
                                    ${((quote?.price ?? 0) * 1.02).toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Low</Text>
                                <Text color={dimeTheme.colors.loss} fontWeight="600">
                                    ${((quote?.price ?? 0) * 0.98).toFixed(2)}
                                </Text>
                            </View>
                            <View style={styles.statItem}>
                                <Text color={dimeTheme.colors.textTertiary} fontSize="$2">Volume</Text>
                                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                    12.5M
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
                            pressStyle={{ backgroundColor: dimeTheme.colors.primaryDark }}
                            onPress={() => router.push({ pathname: '/trade', params: { symbol, side: 'buy' } })}
                        >
                            <Text color={dimeTheme.colors.background} fontWeight="bold" fontSize="$4">
                                Buy
                            </Text>
                        </Button>
                        <Button
                            flex={1}
                            size="$5"
                            backgroundColor={dimeTheme.colors.loss}
                            pressStyle={{ opacity: 0.8 }}
                            onPress={() => router.push({ pathname: '/trade', params: { symbol, side: 'sell' } })}
                        >
                            <Text color={dimeTheme.colors.textPrimary} fontWeight="bold" fontSize="$4">
                                Sell
                            </Text>
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
    periodContainer: {
        paddingHorizontal: 16,
        marginBottom: 8,
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
