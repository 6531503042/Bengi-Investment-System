import { useState, useEffect, useMemo } from 'react'
import { StyleSheet, ScrollView, StatusBar, RefreshControl, TouchableOpacity, ActivityIndicator } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useRouter } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { useDemoStore } from '@/stores/demo'
import { usePortfolioStore, type PortfolioPosition } from '@/stores/portfolio'
import { dimeTheme } from '@/constants/theme'
import { PortfolioCard } from '@/components/portfolio/PortfolioCard'
import { HoldingItem } from '@/components/portfolio/HoldingItem'
import { OptionItem } from '@/components/portfolio/OptionItem'

// Mock data for when no real positions exist (demo mode)
const MOCK_HOLDINGS = [
    { symbol: 'NVDA', name: 'NVIDIA Corporation', logoUrl: 'https://img.logo.dev/nvidia.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', quantity: 15, avgCost: 120.50, currentPrice: 134.82 },
    { symbol: 'TSLA', name: 'Tesla Inc.', logoUrl: 'https://img.logo.dev/tesla.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', quantity: 8, avgCost: 245.00, currentPrice: 421.06 },
    { symbol: 'AAPL', name: 'Apple Inc.', logoUrl: 'https://img.logo.dev/apple.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', quantity: 25, avgCost: 168.00, currentPrice: 254.49 },
    { symbol: 'AMZN', name: 'Amazon.com', logoUrl: 'https://img.logo.dev/amazon.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', quantity: 12, avgCost: 142.50, currentPrice: 227.05 },
    { symbol: 'PLTR', name: 'Palantir Technologies', logoUrl: 'https://img.logo.dev/palantir.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', quantity: 100, avgCost: 18.50, currentPrice: 75.14 },
]

// Mock options (options not in backend yet)
const MOCK_OPTIONS = [
    { symbol: 'TSLA', name: 'Tesla Inc.', logoUrl: 'https://img.logo.dev/tesla.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', type: 'Call' as const, strike: 450.00, expiry: '2025-02-21', contracts: 2, premium: 12.50, currentPrice: 28.75, delta: 0.45, theta: -0.15, iv: 0.52 },
    { symbol: 'NVDA', name: 'NVIDIA Corporation', logoUrl: 'https://img.logo.dev/nvidia.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', type: 'Call' as const, strike: 150.00, expiry: '2025-03-14', contracts: 3, premium: 8.20, currentPrice: 12.35, delta: 0.38, theta: -0.08, iv: 0.48 },
    { symbol: 'AAPL', name: 'Apple Inc.', logoUrl: 'https://img.logo.dev/apple.com?token=pk_X-1ZO13GSgeOoUrIuJ6GMQ', type: 'Put' as const, strike: 240.00, expiry: '2025-01-15', contracts: 5, premium: 3.50, currentPrice: 2.15, delta: -0.32, theta: -0.22, iv: 0.35 },
]

// Portfolio tabs
const PORTFOLIO_TABS = [
    { id: 'all', name: '← All', color: dimeTheme.colors.textSecondary },
    { id: 'stocks', name: '📈 Stocks', color: '#4CAF50' },
    { id: 'options', name: '🎯 Options', color: '#9C27B0' },
]

type FilterType = 'all' | 'stocks' | 'options'
type SortOption = 'value' | 'pnl' | 'allocation'

export default function PortfolioScreen() {
    const router = useRouter()
    const { account, fetchDemo } = useDemoStore()
    const {
        positions: realPositions,
        summary,
        isLoading: portfolioLoading,
        fetchPortfolios,
        refreshPositions
    } = usePortfolioStore()

    const [refreshing, setRefreshing] = useState(false)
    const [activeTab, setActiveTab] = useState<FilterType>('all')
    const [sortBy, setSortBy] = useState<SortOption>('value')

    // Fetch data on mount
    useEffect(() => {
        fetchDemo()
        fetchPortfolios()
    }, [])

    const onRefresh = async () => {
        setRefreshing(true)
        await Promise.all([fetchDemo(), refreshPositions()])
        setRefreshing(false)
    }

    // Use real positions if available, fallback to mock
    const hasRealPositions = realPositions.length > 0
    const displayHoldings = hasRealPositions
        ? realPositions.map(p => ({
            symbol: p.symbol,
            name: p.name,
            logoUrl: p.logoUrl,
            quantity: p.quantity,
            avgCost: p.avgCost,
            currentPrice: p.currentPrice,
        }))
        : MOCK_HOLDINGS

    // Calculate portfolio stats from holdings
    const calculateStats = () => {
        let stocksValue = 0
        let stocksCost = 0
        let optionsValue = 0
        let optionsCost = 0

        // Use displayHoldings (real data or mock)
        displayHoldings.forEach(h => {
            stocksValue += h.quantity * h.currentPrice
            stocksCost += h.quantity * h.avgCost
        })

        MOCK_OPTIONS.forEach(o => {
            optionsValue += o.contracts * 100 * o.currentPrice
            optionsCost += o.contracts * 100 * o.premium
        })

        const totalValue = stocksValue + optionsValue
        const totalCost = stocksCost + optionsCost
        const totalPnl = totalValue - totalCost
        const totalPnlPercent = totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0

        return {
            totalValue,
            totalCost,
            totalPnl,
            totalPnlPercent,
            stocksValue,
            optionsValue,
            initialValue: account?.initialBalance ?? 10000,
            dailyChange: 2.15, // Mock positive daily change
        }
    }

    const stats = calculateStats()

    // Calculate allocation for each holding
    const holdingsWithAllocation = displayHoldings.map(h => ({
        ...h,
        allocation: (h.quantity * h.currentPrice) / stats.totalValue * 100
    }))

    // Sort holdings
    const sortedHoldings = [...holdingsWithAllocation].sort((a, b) => {
        switch (sortBy) {
            case 'value':
                return (b.quantity * b.currentPrice) - (a.quantity * a.currentPrice)
            case 'pnl':
                const pnlA = ((a.currentPrice - a.avgCost) / a.avgCost) * 100
                const pnlB = ((b.currentPrice - b.avgCost) / b.avgCost) * 100
                return pnlB - pnlA
            case 'allocation':
                return b.allocation - a.allocation
            default:
                return 0
        }
    })

    // Filter based on active tab
    const showStocks = activeTab === 'all' || activeTab === 'stocks'
    const showOptions = activeTab === 'all' || activeTab === 'options'

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                {/* Header */}
                <XStack paddingHorizontal="$4" paddingVertical="$3" alignItems="center" justifyContent="space-between">
                    <XStack alignItems="center" gap="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontSize={24} fontWeight="bold">
                            Portfolio
                        </Text>
                        <Ionicons name="eye-outline" size={20} color={dimeTheme.colors.textSecondary} />
                    </XStack>
                    <TouchableOpacity>
                        <Ionicons name="ellipsis-vertical" size={20} color={dimeTheme.colors.textSecondary} />
                    </TouchableOpacity>
                </XStack>

                {/* Filter Tabs */}
                <ScrollView
                    horizontal
                    showsHorizontalScrollIndicator={false}
                    contentContainerStyle={styles.tabsContainer}
                >
                    {PORTFOLIO_TABS.map(tab => (
                        <TouchableOpacity
                            key={tab.id}
                            style={[
                                styles.tabButton,
                                activeTab === tab.id && styles.tabButtonActive,
                                activeTab === tab.id && { borderColor: tab.color }
                            ]}
                            onPress={() => setActiveTab(tab.id as FilterType)}
                        >
                            <Text
                                color={activeTab === tab.id ? tab.color : dimeTheme.colors.textSecondary}
                                fontWeight={activeTab === tab.id ? '600' : '400'}
                                fontSize={13}
                            >
                                {tab.name}
                            </Text>
                        </TouchableOpacity>
                    ))}
                </ScrollView>

                <ScrollView
                    showsVerticalScrollIndicator={false}
                    refreshControl={
                        <RefreshControl
                            refreshing={refreshing}
                            onRefresh={onRefresh}
                            tintColor={dimeTheme.colors.primary}
                        />
                    }
                >
                    {/* Portfolio Card */}
                    <PortfolioCard
                        totalValue={activeTab === 'options' ? stats.optionsValue : activeTab === 'stocks' ? stats.stocksValue : stats.totalValue}
                        cashBalance={account?.balance ?? 10000}
                        investedValue={stats.stocksValue + stats.optionsValue}
                        dailyChange={stats.dailyChange * stats.totalValue / 100}
                        dailyChangePercent={stats.dailyChange}
                        totalPnl={stats.totalPnl}
                        totalPnlPercent={stats.totalPnlPercent}
                        onDetailsPress={() => console.log('Details pressed')}
                    />

                    {/* Holdings Header */}
                    <XStack
                        paddingHorizontal="$4"
                        paddingTop="$4"
                        paddingBottom="$2"
                        justifyContent="space-between"
                        alignItems="center"
                    >
                        <XStack alignItems="center" gap="$2">
                            <Text color={dimeTheme.colors.textSecondary} fontSize={12}>
                                Sort by:
                            </Text>
                            <TouchableOpacity
                                style={styles.sortButton}
                                onPress={() => {
                                    const options: SortOption[] = ['value', 'pnl', 'allocation']
                                    const current = options.indexOf(sortBy)
                                    setSortBy(options[(current + 1) % options.length])
                                }}
                            >
                                <Text color={dimeTheme.colors.textPrimary} fontSize={12} fontWeight="600">
                                    {sortBy === 'value' ? 'Value' : sortBy === 'pnl' ? 'P&L %' : 'Allocation'}
                                </Text>
                                <Ionicons name="chevron-down" size={14} color={dimeTheme.colors.textPrimary} />
                            </TouchableOpacity>
                        </XStack>
                        <TouchableOpacity style={styles.filterButton}>
                            <Text color={dimeTheme.colors.primary} fontSize={12} fontWeight="600">
                                P&L Filter
                            </Text>
                            <Ionicons name="filter" size={14} color={dimeTheme.colors.primary} />
                        </TouchableOpacity>
                    </XStack>

                    {/* Stocks Section */}
                    {showStocks && sortedHoldings.length > 0 && (
                        <>
                            <XStack paddingHorizontal="$4" paddingBottom="$2" justifyContent="space-between">
                                <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                                    {sortedHoldings.length} Stocks
                                </Text>
                                <XStack gap="$4">
                                    <Text color={dimeTheme.colors.textTertiary} fontSize={11}>Value (USD)</Text>
                                    <Text color={dimeTheme.colors.textTertiary} fontSize={11}>P&L %</Text>
                                </XStack>
                            </XStack>
                            {sortedHoldings.map((holding) => (
                                <HoldingItem
                                    key={holding.symbol}
                                    symbol={holding.symbol}
                                    name={holding.name}
                                    logoUrl={holding.logoUrl}
                                    quantity={holding.quantity}
                                    avgCost={holding.avgCost}
                                    currentPrice={holding.currentPrice}
                                    allocation={holding.allocation}
                                    onPress={() => router.push(`/(tabs)/market/${encodeURIComponent(holding.symbol)}`)}
                                />
                            ))}
                        </>
                    )}

                    {/* Options Section */}
                    {showOptions && MOCK_OPTIONS.length > 0 && (
                        <>
                            <XStack paddingHorizontal="$4" paddingTop="$3" paddingBottom="$2">
                                <Text color={dimeTheme.colors.textTertiary} fontSize={11}>
                                    {MOCK_OPTIONS.length} Options
                                </Text>
                            </XStack>
                            {MOCK_OPTIONS.map((option, index) => (
                                <OptionItem
                                    key={`${option.symbol}-${option.type}-${index}`}
                                    symbol={option.symbol}
                                    name={option.name}
                                    logoUrl={option.logoUrl}
                                    type={option.type}
                                    strike={option.strike}
                                    expiry={option.expiry}
                                    contracts={option.contracts}
                                    premium={option.premium}
                                    currentPrice={option.currentPrice}
                                    delta={option.delta}
                                    theta={option.theta}
                                    iv={option.iv}
                                />
                            ))}
                        </>
                    )}

                    {/* Bottom padding */}
                    <View style={{ height: 120 }} />
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
    tabsContainer: {
        paddingHorizontal: 16,
        paddingBottom: 12,
        gap: 8,
    },
    tabButton: {
        paddingHorizontal: 14,
        paddingVertical: 8,
        borderRadius: 20,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        marginRight: 8,
        minHeight: 36,
        justifyContent: 'center',
        alignItems: 'center',
    },
    tabButtonActive: {
        backgroundColor: 'transparent',
        borderWidth: 1.5,
    },
    sortButton: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 4,
        backgroundColor: dimeTheme.colors.surface,
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 6,
    },
    filterButton: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 4,
    },
})
