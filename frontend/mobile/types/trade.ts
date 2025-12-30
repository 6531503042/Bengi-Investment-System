// Trade/Order types
export type OrderSide = 'BUY' | 'SELL'
export type OrderType = 'MARKET' | 'LIMIT' | 'STOP'
export type OrderStatus = 'PENDING' | 'FILLED' | 'PARTIALLY_FILLED' | 'CANCELLED'
export type TimeInForce = 'GTC' | 'DAY' | 'IOC' | 'FOK'

export interface Order {
    id: string
    userId: string
    accountId: string
    portfolioId: string
    instrumentId: string
    symbol: string
    side: OrderSide
    type: OrderType
    status: OrderStatus
    timeInForce: TimeInForce
    quantity: number
    filledQty: number
    price?: number
    stopPrice?: number
    avgFillPrice?: number
    commission: number
    createdAt: string
    filledAt?: string
    cancelledAt?: string
}

// Input for creating new order - matches backend DTO
export interface CreateOrderInput {
    accountId: string
    portfolioId: string
    symbol: string
    side: OrderSide
    type: OrderType
    quantity: number
    price?: number        // Required for LIMIT orders
    stopPrice?: number    // Required for STOP orders
    timeInForce?: TimeInForce
}

export interface Trade {
    id: string
    orderId: string
    userId: string
    symbol: string
    side: OrderSide
    quantity: number
    price: number
    commission: number
    total: number
    netAmount: number
    executedAt: string
}
