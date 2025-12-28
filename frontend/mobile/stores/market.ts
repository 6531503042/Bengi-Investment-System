import { create } from 'zustand'
import type { Quote, Instrument, MarketState } from '@/types/market'
import { instrumentService } from '@/services/api'
import { DEFAULT_WATCHLIST } from '@/constants'

interface PaginationState {
    page: number
    hasMore: boolean
    total: number
}

interface MarketActions {
    fetchInstruments: (reset?: boolean) => Promise<void>
    loadMore: () => Promise<void>
    fetchQuote: (symbol: string) => Promise<Quote | null>
    updateQuote: (quote: Quote) => void
    watchSymbol: (symbol: string) => void
    unwatchSymbol: (symbol: string) => void
    setWsConnected: (connected: boolean) => void
    searchInstruments: (query: string, type?: string) => Promise<void>
}

interface ExtendedMarketState extends MarketState {
    pagination: PaginationState
    isLoadingMore: boolean
    searchQuery: string
    searchResults: Instrument[]
}

const ITEMS_PER_PAGE = 100

export const useMarketStore = create<ExtendedMarketState & MarketActions>((set, get) => ({
    instruments: [],
    quotes: {},
    watchedSymbols: [...DEFAULT_WATCHLIST],
    isLoading: false,
    wsConnected: false,
    pagination: {
        page: 1,
        hasMore: true,
        total: 0,
    },
    isLoadingMore: false,
    searchQuery: '',
    searchResults: [],

    fetchInstruments: async (reset = true) => {
        const page = reset ? 1 : get().pagination.page
        set({ isLoading: true })
        try {
            const result = await instrumentService.getAll(page, ITEMS_PER_PAGE)
            set({
                instruments: reset ? (result.instruments ?? []) : [...get().instruments, ...(result.instruments ?? [])],
                pagination: {
                    page,
                    hasMore: (result.instruments?.length ?? 0) >= ITEMS_PER_PAGE,
                    total: result.total,
                },
                isLoading: false,
            })
        } catch (error) {
            console.error('Failed to fetch instruments:', error)
            set({ isLoading: false })
        }
    },

    loadMore: async () => {
        const { pagination, isLoading, isLoadingMore } = get()
        if (isLoading || isLoadingMore || !pagination.hasMore) return

        const nextPage = pagination.page + 1
        set({ isLoadingMore: true })

        try {
            const result = await instrumentService.getAll(nextPage, ITEMS_PER_PAGE)
            set((state) => ({
                instruments: [...state.instruments, ...(result.instruments ?? [])],
                pagination: {
                    page: nextPage,
                    hasMore: (result.instruments?.length ?? 0) >= ITEMS_PER_PAGE,
                    total: result.total,
                },
                isLoadingMore: false,
            }))
        } catch (error) {
            console.error('Failed to load more instruments:', error)
            set({ isLoadingMore: false })
        }
    },

    searchInstruments: async (query: string, type?: string) => {
        if (!query.trim()) {
            set({ searchQuery: '', searchResults: [] })
            return
        }

        set({ searchQuery: query, isLoading: true })
        try {
            const result = await instrumentService.search(query, type)
            set({
                searchResults: result.instruments ?? [],
                isLoading: false
            })
        } catch (error) {
            console.error('Search failed:', error)
            set({ isLoading: false, searchResults: [] })
        }
    },

    fetchQuote: async (symbol) => {
        try {
            const quote = await instrumentService.getQuote(symbol)
            set((state) => ({ quotes: { ...state.quotes, [symbol]: quote } }))
            return quote
        } catch (error) {
            console.error(`Failed to fetch quote for ${symbol}:`, error)
            return null
        }
    },

    updateQuote: (quote) => {
        set((state) => ({ quotes: { ...state.quotes, [quote.symbol]: quote } }))
    },

    watchSymbol: (symbol) => {
        const { watchedSymbols } = get()
        if (!watchedSymbols.includes(symbol)) {
            set({ watchedSymbols: [...watchedSymbols, symbol] })
        }
    },

    unwatchSymbol: (symbol) => {
        set((state) => ({
            watchedSymbols: state.watchedSymbols.filter((s) => s !== symbol),
        }))
    },

    setWsConnected: (connected) => set({ wsConnected: connected }),
}))
