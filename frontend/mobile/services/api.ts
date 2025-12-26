import axios, { type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'
import * as SecureStore from 'expo-secure-store'
import { API_CONFIG } from '@/constants'
import type { User, LoginCredentials, RegisterData } from '@/types/auth'
import type { Instrument, Quote } from '@/types/market'
import type { Portfolio, Account } from '@/types/portfolio'
import type { Order, CreateOrderInput, Trade } from '@/types/trade'

// Create axios instance
const api: AxiosInstance = axios.create({
    baseURL: API_CONFIG.baseUrl,
    timeout: API_CONFIG.timeout,
    headers: { 'Content-Type': 'application/json' },
})

// Request interceptor - add auth token
api.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
    const token = await SecureStore.getItemAsync('accessToken')
    if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Response interceptor - handle token refresh
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 401 && error.config && !error.config._retry) {
            error.config._retry = true
            const refreshToken = await SecureStore.getItemAsync('refreshToken')

            if (refreshToken) {
                try {
                    const { data } = await axios.post(`${API_CONFIG.baseUrl}/auth/refresh`, { refreshToken })
                    await SecureStore.setItemAsync('accessToken', data.accessToken)
                    error.config.headers.Authorization = `Bearer ${data.accessToken}`
                    return api(error.config)
                } catch {
                    await SecureStore.deleteItemAsync('accessToken')
                    await SecureStore.deleteItemAsync('refreshToken')
                }
            }
        }
        return Promise.reject(error)
    }
)

export default api

// Auth API - Backend returns wrapped response: { success, message, data }
interface ApiResponse<T> {
    success: boolean
    message?: string
    data: T
}

export const authService = {
    login: async (credentials: LoginCredentials) => {
        const { data: response } = await api.post<ApiResponse<{ user: User; accessToken: string; refreshToken: string }>>('/auth/login', credentials)
        if (response.data.accessToken) {
            await SecureStore.setItemAsync('accessToken', response.data.accessToken)
        }
        if (response.data.refreshToken) {
            await SecureStore.setItemAsync('refreshToken', response.data.refreshToken)
        }
        return response.data
    },

    register: async (userData: RegisterData) => {
        const { data: response } = await api.post<ApiResponse<{ id: string; email: string; fullName: string }>>('/auth/register', userData)
        return response.data
    },

    logout: async () => {
        try { await api.post('/auth/logout') } catch { /* ignore */ }
        await SecureStore.deleteItemAsync('accessToken')
        await SecureStore.deleteItemAsync('refreshToken')
    },

    getToken: () => SecureStore.getItemAsync('accessToken'),
}

// Instrument API
export const instrumentService = {
    getAll: async () => {
        const { data } = await api.get<{ instruments: Instrument[] }>('/instruments')
        return data.instruments
    },

    search: async (query: string) => {
        const { data } = await api.get<{ instruments: Instrument[] }>(`/instruments/search?q=${query}`)
        return data.instruments
    },

    getQuote: async (symbol: string) => {
        const { data } = await api.get<Quote>(`/instruments/${symbol}/quote`)
        return data
    },
}

// Portfolio API
export const portfolioService = {
    getAll: async () => {
        const { data } = await api.get<{ portfolios: Portfolio[] }>('/portfolios')
        return data.portfolios
    },

    getById: async (id: string) => {
        const { data } = await api.get<Portfolio>(`/portfolios/${id}`)
        return data
    },

    getSummary: async (id: string) => {
        const { data } = await api.get(`/portfolios/${id}/summary`)
        return data
    },
}

// Account API
export const accountService = {
    getAll: async () => {
        const { data } = await api.get<{ accounts: Account[] }>('/accounts')
        return data.accounts
    },

    deposit: async (id: string, amount: number) => {
        const { data } = await api.post(`/accounts/${id}/deposit`, { amount })
        return data
    },

    withdraw: async (id: string, amount: number) => {
        const { data } = await api.post(`/accounts/${id}/withdraw`, { amount })
        return data
    },
}

// Order API
export const orderService = {
    getAll: async () => {
        const { data } = await api.get<{ orders: Order[] }>('/orders')
        return data.orders
    },

    create: async (input: CreateOrderInput) => {
        const { data } = await api.post<Order>('/orders', input)
        return data
    },

    cancel: async (id: string) => {
        const { data } = await api.post<Order>(`/orders/${id}/cancel`)
        return data
    },
}

// Trade API
export const tradeService = {
    getAll: async () => {
        const { data } = await api.get<{ trades: Trade[] }>('/trades')
        return data.trades
    },

    getById: async (id: string) => {
        const { data } = await api.get<Trade>(`/trades/${id}`)
        return data
    },

    getSummary: async () => {
        const { data } = await api.get('/trades/summary')
        return data
    },
}

// Watchlist API
export const watchlistService = {
    getAll: async () => {
        const { data } = await api.get('/watchlists')
        return data.watchlists
    },

    create: async (name: string) => {
        const { data } = await api.post('/watchlists', { name })
        return data
    },

    addSymbol: async (id: string, symbol: string) => {
        const { data } = await api.post(`/watchlists/${id}/symbols`, { symbol })
        return data
    },

    removeSymbol: async (id: string, symbol: string) => {
        const { data } = await api.delete(`/watchlists/${id}/symbols/${symbol}`)
        return data
    },
}
