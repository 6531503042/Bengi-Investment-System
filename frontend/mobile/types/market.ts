// Market/Quote types
export interface Quote {
    symbol: string
    price: number
    change: number
    changePercent: number
    high: number
    low: number
    open: number
    previousClose: number
    volume: number
    timestamp: number
}

export interface Instrument {
    symbol: string
    name: string
    type: 'stock' | 'crypto' | 'forex'
    exchange: string
    currency: string
}

export interface MarketState {
    instruments: Instrument[]
    quotes: Record<string, Quote>
    watchedSymbols: string[]
    isLoading: boolean
    wsConnected: boolean
}
