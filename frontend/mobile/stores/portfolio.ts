import { create } from 'zustand'
import type { Portfolio, Account, PortfolioState } from '@/types/portfolio'
import { portfolioService, accountService } from '@/services/api'

interface PortfolioActions {
    fetchPortfolios: () => Promise<void>
    fetchAccounts: () => Promise<void>
    selectPortfolio: (id: string) => void
    getSelectedPortfolio: () => Portfolio | null
    deposit: (accountId: string, amount: number) => Promise<boolean>
    withdraw: (accountId: string, amount: number) => Promise<boolean>
}

export const usePortfolioStore = create<PortfolioState & PortfolioActions>((set, get) => ({
    portfolios: [],
    accounts: [],
    selectedPortfolioId: null,
    isLoading: false,

    fetchPortfolios: async () => {
        set({ isLoading: true })
        try {
            const portfolios = await portfolioService.getAll()
            set({
                portfolios: portfolios ?? [],
                isLoading: false,
                selectedPortfolioId: portfolios?.[0]?.id ?? null,
            })
        } catch (error) {
            console.error('Failed to fetch portfolios:', error)
            set({ portfolios: [], isLoading: false })
        }
    },

    fetchAccounts: async () => {
        try {
            const accounts = await accountService.getAll()
            set({ accounts: accounts ?? [] })
        } catch (error) {
            console.error('Failed to fetch accounts:', error)
            set({ accounts: [] })
        }
    },

    selectPortfolio: (id) => set({ selectedPortfolioId: id }),

    getSelectedPortfolio: () => {
        const { portfolios, selectedPortfolioId } = get()
        return portfolios.find((p) => p.id === selectedPortfolioId) ?? null
    },

    deposit: async (accountId, amount) => {
        try {
            await accountService.deposit(accountId, amount)
            await get().fetchAccounts()
            return true
        } catch (error) {
            console.error('Deposit failed:', error)
            return false
        }
    },

    withdraw: async (accountId, amount) => {
        try {
            await accountService.withdraw(accountId, amount)
            await get().fetchAccounts()
            return true
        } catch (error) {
            console.error('Withdraw failed:', error)
            return false
        }
    },
}))
