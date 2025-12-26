// Portfolio types
export interface Position {
    id: string
    symbol: string
    quantity: number
    averageCost: number
    currentPrice: number
    marketValue: number
    unrealizedPL: number
    unrealizedPLPercent: number
}

export interface Portfolio {
    id: string
    name: string
    userId: string
    positions: Position[]
    totalValue: number
    totalCost: number
    totalPL: number
    totalPLPercent: number
}

export interface Account {
    id: string
    userId: string
    currency: string
    balance: number
    availableBalance: number
}

export interface PortfolioState {
    portfolios: Portfolio[]
    accounts: Account[]
    selectedPortfolioId: string | null
    isLoading: boolean
}
