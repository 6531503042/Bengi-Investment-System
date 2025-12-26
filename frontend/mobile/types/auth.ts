// User types
export interface User {
    id: string
    email: string
    fullName: string
    role?: 'user' | 'admin'
    createdAt?: string
}

export interface AuthState {
    user: User | null
    isAuthenticated: boolean
    isLoading: boolean
    error: string | null
}

export interface LoginCredentials {
    email: string
    password: string
}

export interface RegisterData {
    email: string
    password: string
    fullName: string
}

