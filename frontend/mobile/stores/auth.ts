import { create } from 'zustand'
import type { User, AuthState } from '@/types/auth'
import { authService } from '@/services/api'

interface AuthActions {
    login: (email: string, password: string) => Promise<boolean>
    register: (email: string, password: string, fullName: string) => Promise<boolean>
    logout: () => Promise<void>
    checkAuth: () => Promise<void>
    clearError: () => void
}

export const useAuthStore = create<AuthState & AuthActions>((set, get) => ({
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null,

    login: async (email, password) => {
        set({ isLoading: true, error: null })
        try {
            const { user } = await authService.login({ email, password })
            set({ user, isAuthenticated: true, isLoading: false })
            return true
        } catch (error: unknown) {
            const message = (error as { response?: { data?: { error?: string } } }).response?.data?.error ?? 'Login failed'
            set({ error: message, isLoading: false })
            return false
        }
    },

    register: async (email, password, fullName) => {
        set({ isLoading: true, error: null })
        try {
            await authService.register({ email, password, fullName })
            return await get().login(email, password)
        } catch (error: unknown) {
            const message = (error as { response?: { data?: { error?: string } } }).response?.data?.error ?? 'Registration failed'
            set({ error: message, isLoading: false })
            return false
        }
    },

    logout: async () => {
        await authService.logout()
        set({ user: null, isAuthenticated: false })
    },

    checkAuth: async () => {
        set({ isLoading: true })
        const token = await authService.getToken()
        set({ isLoading: false, isAuthenticated: !!token })
    },

    clearError: () => set({ error: null }),
}))
