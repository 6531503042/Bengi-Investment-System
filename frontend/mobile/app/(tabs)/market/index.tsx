import { useState, useEffect, useCallback } from 'react'
import { StyleSheet, FlatList, StatusBar, TextInput, Image, RefreshControl, TouchableOpacity } from 'react-native'
import { YStack, XStack, Text, View, Spinner } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useRouter } from 'expo-router'
import { useMarketStore } from '@/stores/market'
import { dimeTheme } from '@/constants/theme'
import { Ionicons } from '@expo/vector-icons'
import { LinearGradient } from 'expo-linear-gradient'
import type { Instrument, InstrumentType } from '@/types/market'

const CATEGORIES: Array<{ key: 'All' | InstrumentType; label: string; icon: string; color: string }> = [
    { key: 'All', label: 'All', icon: 'apps', color: dimeTheme.colors.primary },
    { key: 'Stock', label: 'Stocks', icon: 'trending-up', color: '#4CAF50' },
    { key: 'ETF', label: 'ETFs', icon: 'pie-chart', color: '#9C27B0' },
    { key: 'Crypto', label: 'Crypto', icon: 'logo-bitcoin', color: '#FF9800' },
]

export default function MarketScreen() {
    const router = useRouter()
    const { instruments, fetchInstruments, isLoading } = useMarketStore()
    const [searchQuery, setSearchQuery] = useState('')
    const [activeCategory, setActiveCategory] = useState<'All' | InstrumentType>('All')
    const [refreshing, setRefreshing] = useState(false)

    useEffect(() => {
        fetchInstruments()
    }, [])

    const onRefresh = useCallback(async () => {
        setRefreshing(true)
        await fetchInstruments()
        setRefreshing(false)
    }, [])

    // Filter instruments
    const filteredInstruments = instruments.filter(instrument => {
        const matchesSearch =
            instrument.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
            instrument.name.toLowerCase().includes(searchQuery.toLowerCase())
        const matchesCategory = activeCategory === 'All' || instrument.type === activeCategory
        return matchesSearch && matchesCategory
    })

    const getTypeColor = (type: InstrumentType) => {
        switch (type) {
            case 'Stock': return '#4CAF50'
            case 'ETF': return '#9C27B0'
            case 'Crypto': return '#FF9800'
            default: return dimeTheme.colors.primary
        }
    }

    const [failedLogos, setFailedLogos] = useState<Set<string>>(new Set())

    const renderInstrument = ({ item }: { item: Instrument }) => {
        const logoFailed = failedLogos.has(item.symbol) || !item.logoUrl

        return (
            <TouchableOpacity
                style={styles.instrumentCard}
                activeOpacity={0.7}
                onPress={() => router.push(`/(tabs)/market/${encodeURIComponent(item.symbol)}`)}
            >
                {/* Logo */}
                <View style={styles.logoContainer}>
                    {!logoFailed ? (
                        <Image
                            source={{ uri: item.logoUrl }}
                            style={styles.logo}
                            onError={() => setFailedLogos(prev => new Set(prev).add(item.symbol))}
                        />
                    ) : (
                        <View style={[styles.logoPlaceholder, { backgroundColor: getTypeColor(item.type) + '20' }]}>
                            <Text color={getTypeColor(item.type)} fontSize={16} fontWeight="bold">
                                {item.symbol.charAt(0)}
                            </Text>
                        </View>
                    )}
                </View>

                {/* Info */}
                <View style={styles.instrumentInfo}>
                    <XStack alignItems="center" gap="$1">
                        <Text color={dimeTheme.colors.textPrimary} fontWeight="700" fontSize={15}>
                            {item.symbol.replace('/USD', '')}
                        </Text>
                        <View style={[styles.typeBadge, { backgroundColor: getTypeColor(item.type) + '20' }]}>
                            <Text color={getTypeColor(item.type)} fontSize={9} fontWeight="600">
                                {item.type}
                            </Text>
                        </View>
                    </XStack>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={12} numberOfLines={1}>
                        {item.name}
                    </Text>
                </View>

                {/* Arrow */}
                <Ionicons name="chevron-forward" size={18} color={dimeTheme.colors.textTertiary} />
            </TouchableOpacity>
        )
    }

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                {/* Header */}
                <View style={styles.header}>
                    <Text color={dimeTheme.colors.textPrimary} fontSize={28} fontWeight="bold">
                        Market
                    </Text>
                    <Text color={dimeTheme.colors.textSecondary} fontSize={13}>
                        {instruments.length} instruments
                    </Text>
                </View>

                {/* Search */}
                <View style={styles.searchContainer}>
                    <Ionicons name="search" size={18} color={dimeTheme.colors.textTertiary} />
                    <TextInput
                        style={styles.searchInput}
                        placeholder="Search stocks, ETFs, crypto..."
                        placeholderTextColor={dimeTheme.colors.textTertiary}
                        value={searchQuery}
                        onChangeText={setSearchQuery}
                    />
                    {searchQuery.length > 0 && (
                        <TouchableOpacity onPress={() => setSearchQuery('')}>
                            <Ionicons name="close-circle" size={18} color={dimeTheme.colors.textTertiary} />
                        </TouchableOpacity>
                    )}
                </View>

                {/* Categories */}
                <View style={styles.categoriesContainer}>
                    {CATEGORIES.map(cat => (
                        <TouchableOpacity
                            key={cat.key}
                            style={[
                                styles.categoryChip,
                                activeCategory === cat.key && styles.categoryChipActive,
                                activeCategory === cat.key && { borderColor: cat.color }
                            ]}
                            onPress={() => setActiveCategory(cat.key)}
                        >
                            <Ionicons
                                name={cat.icon as any}
                                size={14}
                                color={activeCategory === cat.key ? cat.color : dimeTheme.colors.textSecondary}
                            />
                            <Text
                                color={activeCategory === cat.key ? cat.color : dimeTheme.colors.textSecondary}
                                fontSize={12}
                                fontWeight={activeCategory === cat.key ? '600' : '400'}
                            >
                                {cat.label}
                            </Text>
                        </TouchableOpacity>
                    ))}
                </View>

                {/* List */}
                {isLoading && instruments.length === 0 ? (
                    <View style={styles.loadingContainer}>
                        <Spinner size="large" color={dimeTheme.colors.primary} />
                    </View>
                ) : (
                    <FlatList
                        data={filteredInstruments}
                        renderItem={renderInstrument}
                        keyExtractor={item => item.id}
                        contentContainerStyle={styles.listContent}
                        showsVerticalScrollIndicator={false}
                        refreshControl={
                            <RefreshControl
                                refreshing={refreshing}
                                onRefresh={onRefresh}
                                tintColor={dimeTheme.colors.primary}
                            />
                        }
                        ListEmptyComponent={
                            <View style={styles.emptyContainer}>
                                <Ionicons name="search-outline" size={48} color={dimeTheme.colors.textTertiary} />
                                <Text color={dimeTheme.colors.textSecondary} marginTop="$2">
                                    No instruments found
                                </Text>
                            </View>
                        }
                    />
                )}
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
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'baseline',
        paddingHorizontal: 16,
        paddingTop: 8,
        paddingBottom: 12,
    },
    searchContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        marginHorizontal: 16,
        paddingHorizontal: 14,
        paddingVertical: 10,
        borderRadius: 12,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        gap: 10,
    },
    searchInput: {
        flex: 1,
        color: dimeTheme.colors.textPrimary,
        fontSize: 15,
    },
    categoriesContainer: {
        flexDirection: 'row',
        paddingHorizontal: 16,
        paddingVertical: 12,
        gap: 8,
    },
    categoryChip: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 20,
        backgroundColor: dimeTheme.colors.surface,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
        gap: 4,
    },
    categoryChipActive: {
        backgroundColor: 'transparent',
        borderWidth: 1.5,
    },
    listContent: {
        paddingHorizontal: 16,
        paddingBottom: 100,
    },
    instrumentCard: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: dimeTheme.colors.surface,
        padding: 12,
        borderRadius: 12,
        marginBottom: 8,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    logoContainer: {
        width: 40,
        height: 40,
        borderRadius: 20,
        overflow: 'hidden',
        marginRight: 12,
    },
    logo: {
        width: 40,
        height: 40,
        borderRadius: 20,
    },
    logoPlaceholder: {
        width: 40,
        height: 40,
        borderRadius: 20,
        alignItems: 'center',
        justifyContent: 'center',
    },
    instrumentInfo: {
        flex: 1,
    },
    typeBadge: {
        paddingHorizontal: 6,
        paddingVertical: 2,
        borderRadius: 4,
    },
    loadingContainer: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
    },
    emptyContainer: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        paddingTop: 100,
    },
})
