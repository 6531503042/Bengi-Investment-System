import { create } from 'zustand'
import { demoService } from '@/services/api'
import type { DemoAccountStats, DemoDepositResponse, DemoResetResponse } from '@/types/demo'

interface DemoState {
    account: DemoAccountStats | null
    isLoading: boolean
    error: string | null
}

interface DemoActions {
    fetchDemo: () => Promise<void>
    deposit: (amount: number) => Promise<DemoDepositResponse | null>
    reset: (initialBalance?: number) => Promise<DemoResetResponse | null>
    refresh: () => Promise<void>
    clear: () => void
}

export const useDemoStore = create<DemoState & DemoActions>((set, get) => ({
    account: null,
    isLoading: false,
    error: null,

    fetchDemo: async () => {
        set({ isLoading: true, error: null })
        try {
            const account = await demoService.getOrCreate()
            set({ account, isLoading: false })
        } catch (error) {
            set({
                error: error instanceof Error ? error.message : 'Failed to fetch demo account',
                isLoading: false
            })
        }
    },

    deposit: async (amount: number) => {
        const { account } = get()
        if (!account) return null

        set({ isLoading: true, error: null })
        try {
            const response = await demoService.deposit(account.accountId, amount)
            // Update local state with new balance
            set({
                account: {
                    ...account,
                    balance: response.newBalance,
                    totalDeposits: response.totalDeposits,
                },
                isLoading: false
            })
            return response
        } catch (error) {
            set({
                error: error instanceof Error ? error.message : 'Failed to deposit',
                isLoading: false
            })
            return null
        }
    },

    reset: async (initialBalance?: number) => {
        const { account } = get()
        if (!account) return null

        set({ isLoading: true, error: null })
        try {
            const response = await demoService.reset(account.accountId, initialBalance)
            // Update local state with reset values
            set({
                account: {
                    ...account,
                    balance: response.newBalance,
                    initialBalance: response.initialBalance,
                    totalDeposits: response.initialBalance,
                    totalPnL: 0,
                    pnlPercentage: 0,
                },
                isLoading: false
            })
            return response
        } catch (error) {
            set({
                error: error instanceof Error ? error.message : 'Failed to reset account',
                isLoading: false
            })
            return null
        }
    },

    refresh: async () => {
        const { account } = get()
        if (!account) {
            await get().fetchDemo()
            return
        }
        try {
            const stats = await demoService.getStats(account.accountId)
            set({ account: stats })
        } catch (error) {
            // Fall back to getOrCreate if stats fails
            await get().fetchDemo()
        }
    },

    clear: () => {
        set({ account: null, error: null, isLoading: false })
    },
}))
