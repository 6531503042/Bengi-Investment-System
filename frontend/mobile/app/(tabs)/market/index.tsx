import { useState, useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar, TextInput } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useRouter } from 'expo-router'
import { useMarketStore } from '@/stores/market'
import { dimeTheme } from '@/constants/theme'
import { WatchlistItem } from '@/components/market/WatchlistItem'
import { Ionicons } from '@expo/vector-icons'

// Default watchlist symbols
const DEFAULT_SYMBOLS = ['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'TSLA', 'NVDA', 'META']

export default function MarketScreen() {
    const router = useRouter()
    const [searchQuery, setSearchQuery] = useState('')
    const [activeCategory, setActiveCategory] = useState('All')
    const { quotes, watchedSymbols, watchSymbol } = useMarketStore()

    // Initialize watchlist with default symbols
    useEffect(() => {
        DEFAULT_SYMBOLS.forEach(symbol => {
            if (!watchedSymbols.includes(symbol)) {
                watchSymbol(symbol)
            }
        })
    }, [])

    const categories = ['All', 'Stocks', 'ETFs', 'Crypto']

    // Filter symbols based on search
    const filteredSymbols = watchedSymbols.filter(symbol =>
        symbol.toLowerCase().includes(searchQuery.toLowerCase())
    )

    // Mock company names
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
                <YStack padding="$4" paddingBottom="$2">
                    <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                        Market
                    </Text>
                </YStack>

                {/* Search Bar */}
                <YStack paddingHorizontal="$4" marginBottom="$3">
                    <View style={styles.searchContainer}>
                        <Ionicons name="search" size={20} color={dimeTheme.colors.textSecondary} />
                        <TextInput
                            style={styles.searchInput}
                            placeholder="Search stocks, ETFs..."
                            placeholderTextColor={dimeTheme.colors.textTertiary}
                            value={searchQuery}
                            onChangeText={setSearchQuery}
                        />
                    </View>
                </YStack>

                {/* Categories */}
                <ScrollView
                    horizontal
                    showsHorizontalScrollIndicator={false}
                    contentContainerStyle={styles.categoriesContainer}
                >
                    {categories.map(category => (
                        <Button
                            key={category}
                            size="$3"
                            backgroundColor={activeCategory === category ? dimeTheme.colors.primary : dimeTheme.colors.surface}
                            pressStyle={{ opacity: 0.8 }}
                            marginRight="$2"
                            onPress={() => setActiveCategory(category)}
                        >
                            <Text
                                color={activeCategory === category ? dimeTheme.colors.background : dimeTheme.colors.textPrimary}
                                fontWeight="600"
                            >
                                {category}
                            </Text>
                        </Button>
                    ))}
                </ScrollView>

                {/* Watchlist */}
                <ScrollView
                    style={styles.watchlistContainer}
                    showsVerticalScrollIndicator={false}
                    contentContainerStyle={styles.watchlistContent}
                >
                    <XStack justifyContent="space-between" alignItems="center" marginBottom="$3">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                            Watchlist
                        </Text>
                        <Text color={dimeTheme.colors.primary} fontSize="$3">
                            {filteredSymbols.length} items
                        </Text>
                    </XStack>

                    {filteredSymbols.length > 0 ? (
                        filteredSymbols.map(symbol => {
                            const quote = quotes[symbol]
                            return (
                                <WatchlistItem
                                    key={symbol}
                                    symbol={symbol}
                                    name={companyNames[symbol] ?? symbol}
                                    price={quote?.price ?? 0}
                                    change={quote?.change ?? 0}
                                    changePercent={quote?.changePercent ?? 0}
                                    onPress={() => router.push(`/(tabs)/market/${symbol}`)}
                                />
                            )
                        })
                    ) : (
                        <View style={styles.emptyState}>
                            <Text color={dimeTheme.colors.textTertiary} textAlign="center">
                                No stocks found
                            </Text>
                            <Text color={dimeTheme.colors.textSecondary} fontSize="$2" textAlign="center" marginTop="$1">
                                Try a different search term
                            </Text>
                        </View>
                    )}
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
    searchContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        borderRadius: dimeTheme.radius.lg,
        paddingHorizontal: 16,
        paddingVertical: 12,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        gap: 12,
    },
    searchInput: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 16,
    },
    categoriesContainer: {
        paddingHorizontal: 16,
        marginBottom: 16,
    },
    watchlistContainer: {
        flex: 1,
    },
    watchlistContent: {
        paddingHorizontal: 16,
        paddingBottom: 24,
    },
    emptyState: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 32,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
