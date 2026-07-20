import api from '@/lib/api'
import type { User, LoginCredentials, LoginResponse, UserResponse } from '@/types/auth'

export const authService = {
  async login(credentials: LoginCredentials): Promise<User> {
    const response = await api.post<LoginResponse>('/auth/login', credentials)
    return response.data.data.user
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout')
  },

  async me(): Promise<User> {
    const response = await api.get<UserResponse>('/auth/me')
    return response.data.data
  },

  async refresh(): Promise<void> {
    await api.post('/auth/refresh')
  },
}
