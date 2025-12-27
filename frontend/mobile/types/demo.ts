// Demo Account Types
export interface DemoAccountStats {
    accountId: string
    currency: string
    balance: number
    initialBalance: number
    totalDeposits: number
    totalPnL: number
    leverage: number
    pnlPercentage: number
    createdAt: string
}

export interface DemoDepositRequest {
    amount: number
}

export interface DemoDepositResponse {
    accountId: string
    newBalance: number
    totalDeposits: number
    message: string
}

export interface DemoResetRequest {
    initialBalance?: number
}

export interface DemoResetResponse {
    accountId: string
    newBalance: number
    initialBalance: number
    message: string
}

export interface CreateDemoAccountRequest {
    currency?: string
    leverage?: number
    initialBalance?: number
}
