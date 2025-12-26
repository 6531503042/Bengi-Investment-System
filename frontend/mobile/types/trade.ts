// Trade/Order types
export type OrderSide = 'BUY' | 'SELL'
export type OrderType = 'MARKET' | 'LIMIT' | 'STOP'
export type OrderStatus = 'PENDING' | 'FILLED' | 'PARTIALLY_FILLED' | 'CANCELLED'

export interface Order {
    id: string
    userId: string
    portfolioId: string
    symbol: string
    side: OrderSide
    type: OrderType
    quantity: number
    price?: number
    status: OrderStatus
    filledQty: number
    avgPrice: number
    createdAt: string
    updatedAt: string
}

export interface CreateOrderInput {
    portfolioId: string
    symbol: string
    side: OrderSide
    type: OrderType
    quantity: number
    price?: number
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
    executedAt: string
}
