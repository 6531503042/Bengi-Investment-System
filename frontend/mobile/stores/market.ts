import { create } from 'zustand'
import type { Quote, Instrument, MarketState } from '@/types/market'
import { instrumentService } from '@/services/api'
import { DEFAULT_WATCHLIST } from '@/constants'

interface MarketActions {
    fetchInstruments: () => Promise<void>
    fetchQuote: (symbol: string) => Promise<Quote | null>
    updateQuote: (quote: Quote) => void
    watchSymbol: (symbol: string) => void
    unwatchSymbol: (symbol: string) => void
    setWsConnected: (connected: boolean) => void
}

export const useMarketStore = create<MarketState & MarketActions>((set, get) => ({
    instruments: [],
    quotes: {},
    watchedSymbols: [...DEFAULT_WATCHLIST],
    isLoading: false,
    wsConnected: false,

    fetchInstruments: async () => {
        set({ isLoading: true })
        try {
            const result = await instrumentService.getAll()
            set({ instruments: result.instruments ?? [], isLoading: false })
        } catch (error) {
            console.error('Failed to fetch instruments:', error)
            set({ isLoading: false })
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
