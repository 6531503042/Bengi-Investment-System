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
    timestamp: number | string
}

export type InstrumentType = 'Stock' | 'ETF' | 'Crypto' | 'Future' | 'Option'

export interface Instrument {
    id: string
    symbol: string
    name: string
    type: InstrumentType
    exchange: string
    currency: string
    status: string
    description?: string
    logoUrl?: string
    createdAt?: string
    updatedAt?: string
}

export interface MarketState {
    instruments: Instrument[]
    quotes: Record<string, Quote>
    watchedSymbols: string[]
    isLoading: boolean
    wsConnected: boolean
}

export interface CandleData {
    time: number
    open: number
    high: number
    low: number
    close: number
    volume: number
}

