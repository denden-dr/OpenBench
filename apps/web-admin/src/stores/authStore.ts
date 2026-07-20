import { create } from 'zustand'
import type { User, LoginCredentials } from '@/types/auth'
import { authService } from '@/services/authService'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  logout: () => Promise<void>
  checkAuth: () => Promise<void>
  setAuth: (user: User | null) => void
  clearAuth: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: true,

  setAuth: (user) => set({ user, isAuthenticated: !!user, isLoading: false }),

  clearAuth: () => set({ user: null, isAuthenticated: false, isLoading: false }),

  login: async (credentials) => {
    set({ isLoading: true })
    try {
      const user = await authService.login(credentials)
      set({ user, isAuthenticated: true, isLoading: false })
    } catch (error) {
      set({ user: null, isAuthenticated: false, isLoading: false })
      throw error
    }
  },

  logout: async () => {
    try {
      await authService.logout()
    } catch {
      // Ignore logout errors
    } finally {
      set({ user: null, isAuthenticated: false, isLoading: false })
    }
  },

  checkAuth: async () => {
    set({ isLoading: true })
    try {
      const user = await authService.me()
      set({ user, isAuthenticated: true, isLoading: false })
    } catch {
      set({ user: null, isAuthenticated: false, isLoading: false })
    }
  },
}))
