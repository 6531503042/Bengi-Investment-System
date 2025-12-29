import { create } from 'zustand'
import { portfolioService, instrumentService } from '@/services/api'
import type { Portfolio } from '@/types/portfolio'

// Position with instrument info for display
export interface PortfolioPosition {
    id: string
    portfolioId: string
    instrumentId: string
    symbol: string
    name: string
    logoUrl?: string
    quantity: number
    avgCost: number
    totalCost: number
    currentPrice: number
    marketValue: number
    unrealizedPnL: number
    unrealizedPnLPct: number
}

interface PortfolioSummary {
    totalValue: number
    totalCost: number
    totalPnL: number
    totalPnLPct: number
    cashBalance: number
    investedValue: number
    dailyChange: number
    dailyChangePercent: number
}

interface PortfolioState {
    portfolios: Portfolio[]
    activePortfolio: Portfolio | null
    positions: PortfolioPosition[]
    summary: PortfolioSummary
    isLoading: boolean
    error: string | null

    // Actions
    fetchPortfolios: () => Promise<void>
    fetchPortfolioSummary: (id: string) => Promise<void>
    setActivePortfolio: (portfolio: Portfolio) => void
    refreshPositions: () => Promise<void>
}

// Initial empty summary
const emptySummary: PortfolioSummary = {
    totalValue: 0,
    totalCost: 0,
    totalPnL: 0,
    totalPnLPct: 0,
    cashBalance: 0,
    investedValue: 0,
    dailyChange: 0,
    dailyChangePercent: 0,
}

export const usePortfolioStore = create<PortfolioState>((set, get) => ({
    portfolios: [],
    activePortfolio: null,
    positions: [],
    summary: emptySummary,
    isLoading: false,
    error: null,

    fetchPortfolios: async () => {
        set({ isLoading: true, error: null })
        try {
            const portfolios = await portfolioService.getAll() ?? []
            set({ portfolios, isLoading: false })

            // Auto-select first portfolio if none active
            if (portfolios.length > 0 && !get().activePortfolio) {
                const defaultPortfolio = portfolios.find(p => p.isDefault) || portfolios[0]
                set({ activePortfolio: defaultPortfolio })
                // Fetch summary for the active portfolio
                await get().fetchPortfolioSummary(defaultPortfolio.id)
            }
        } catch (error: any) {
            console.error('Failed to fetch portfolios:', error)
            // Don't show error for empty portfolio (demo mode)
            set({
                portfolios: [],
                isLoading: false,
                error: null // Suppress error for demo mode
            })
        }
    },

    fetchPortfolioSummary: async (id: string) => {
        set({ isLoading: true, error: null })
        try {
            const summaryData = await portfolioService.getSummary(id)

            // Enrich positions with instrument data
            const enrichedPositions: PortfolioPosition[] = await Promise.all(
                (summaryData.positions || []).map(async (pos) => {
                    let instrumentData = { name: pos.symbol, logoUrl: undefined }
                    try {
                        const instrument = await instrumentService.getBySymbol(pos.symbol)
                        instrumentData = { name: instrument.name, logoUrl: instrument.logoUrl }
                    } catch {
                        // Use symbol as name if instrument not found
                    }

                    return {
                        id: pos.id,
                        portfolioId: pos.portfolioId,
                        instrumentId: pos.instrumentId,
                        symbol: pos.symbol,
                        name: instrumentData.name,
                        logoUrl: instrumentData.logoUrl,
                        quantity: pos.quantity,
                        avgCost: pos.avgCost,
                        totalCost: pos.totalCost,
                        currentPrice: pos.currentPrice || 0,
                        marketValue: pos.marketValue || pos.quantity * (pos.currentPrice || pos.avgCost),
                        unrealizedPnL: pos.unrealizedPnL || 0,
                        unrealizedPnLPct: pos.unrealizedPnLPct || 0,
                    }
                })
            )

            // Calculate summary
            const investedValue = enrichedPositions.reduce((sum, p) => sum + p.marketValue, 0)

            const summary: PortfolioSummary = {
                totalValue: summaryData.totalValue || investedValue,
                totalCost: summaryData.totalCost || enrichedPositions.reduce((sum, p) => sum + p.totalCost, 0),
                totalPnL: summaryData.totalPnL || 0,
                totalPnLPct: summaryData.totalPnLPct || 0,
                cashBalance: (summaryData.totalValue || 0) - investedValue,
                investedValue,
                dailyChange: summaryData.totalPnL * 0.1, // Mock daily (10% of total P&L)
                dailyChangePercent: summaryData.totalPnLPct * 0.1,
            }

            set({
                positions: enrichedPositions,
                summary,
                isLoading: false
            })
        } catch (error: any) {
            console.error('Failed to fetch portfolio summary:', error)
            set({
                isLoading: false,
                error: error.response?.data?.message || 'Failed to fetch portfolio summary'
            })
        }
    },

    setActivePortfolio: (portfolio: Portfolio) => {
        set({ activePortfolio: portfolio })
        get().fetchPortfolioSummary(portfolio.id)
    },

    refreshPositions: async () => {
        const { activePortfolio } = get()
        if (activePortfolio) {
            await get().fetchPortfolioSummary(activePortfolio.id)
        }
    },
}))
