import { useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { SafeAreaView } from 'react-native-safe-area-context'
import { usePortfolioStore } from '@/stores/portfolio'
import { useMarketStore } from '@/stores/market'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { HoldingCard } from '@/components/portfolio/HoldingCard'
import type { Portfolio, Position } from '@/types/portfolio'

// Mock holdings data for demo
const MOCK_HOLDINGS = [
    { symbol: 'AAPL', name: 'Apple Inc.', quantity: 10, avgCost: 145.00, currentPrice: 178.50 },
    { symbol: 'GOOGL', name: 'Alphabet Inc.', quantity: 5, avgCost: 125.00, currentPrice: 142.30 },
    { symbol: 'TSLA', name: 'Tesla Inc.', quantity: 8, avgCost: 220.00, currentPrice: 248.75 },
]

export default function PortfolioScreen() {
    const { portfolios = [], fetchPortfolios } = usePortfolioStore()
    const { quotes } = useMarketStore()

    useEffect(() => {
        fetchPortfolios()
    }, [])

    const safePortfolios = portfolios ?? []
    const totalValue = MOCK_HOLDINGS.reduce((sum, h) => sum + (h.quantity * h.currentPrice), 0)
    const totalCost = MOCK_HOLDINGS.reduce((sum, h) => sum + (h.quantity * h.avgCost), 0)
    const totalPL = totalValue - totalCost
    const totalPLPercent = totalCost > 0 ? (totalPL / totalCost) * 100 : 0

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                <ScrollView showsVerticalScrollIndicator={false}>
                    {/* Header */}
                    <YStack padding="$4" paddingBottom="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                            Portfolio
                        </Text>
                    </YStack>

                    {/* Portfolio Value Card */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <LinearGradient
                            colors={[dimeTheme.colors.surface, dimeTheme.colors.backgroundSecondary]}
                            start={{ x: 0, y: 0 }}
                            end={{ x: 1, y: 1 }}
                            style={styles.valueCard}
                        >
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$3" marginBottom="$2">
                                Total Value
                            </Text>
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$10" fontWeight="bold">
                                ${totalValue.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                            </Text>
                            <XStack marginTop="$3" alignItems="center" gap="$3">
                                <PriceChip value={totalPLPercent} size="md" />
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                                    {totalPL >= 0 ? '+' : ''}${totalPL.toFixed(2)} all time
                                </Text>
                            </XStack>
                        </LinearGradient>
                    </YStack>

                    {/* Holdings Summary */}
                    <XStack paddingHorizontal="$4" gap="$3" marginBottom="$4">
                        <View style={styles.summaryCard}>
                            <Text color={dimeTheme.colors.textTertiary} fontSize="$2">
                                Invested
                            </Text>
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$4" fontWeight="bold">
                                ${totalCost.toFixed(2)}
                            </Text>
                        </View>
                        <View style={styles.summaryCard}>
                            <Text color={dimeTheme.colors.textTertiary} fontSize="$2">
                                Profit/Loss
                            </Text>
                            <Text
                                color={totalPL >= 0 ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                fontSize="$4"
                                fontWeight="bold"
                            >
                                {totalPL >= 0 ? '+' : ''}${totalPL.toFixed(2)}
                            </Text>
                        </View>
                    </XStack>

                    {/* Holdings List */}
                    <YStack paddingHorizontal="$4" paddingBottom="$4">
                        <XStack justifyContent="space-between" alignItems="center" marginBottom="$3">
                            <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                                Holdings
                            </Text>
                            <Text color={dimeTheme.colors.primary} fontSize="$3">
                                {MOCK_HOLDINGS.length} positions
                            </Text>
                        </XStack>

                        {MOCK_HOLDINGS.length > 0 ? (
                            MOCK_HOLDINGS.map(holding => (
                                <HoldingCard
                                    key={holding.symbol}
                                    symbol={holding.symbol}
                                    name={holding.name}
                                    quantity={holding.quantity}
                                    avgCost={holding.avgCost}
                                    currentPrice={quotes[holding.symbol]?.price ?? holding.currentPrice}
                                />
                            ))
                        ) : (
                            <View style={styles.emptyState}>
                                <Text color={dimeTheme.colors.textTertiary} textAlign="center" fontSize="$4">
                                    No positions yet
                                </Text>
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$2" textAlign="center" marginTop="$2">
                                    Start trading to build your portfolio!
                                </Text>
                            </View>
                        )}
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
    valueCard: {
        padding: 20,
        borderRadius: dimeTheme.radius.xl,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    summaryCard: {
        flex: 1,
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    emptyState: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 32,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
