import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuthStore } from '../authStore'
import { authService } from '@/services/authService'
import type { User } from '@/types/auth'

vi.mock('@/services/authService', () => ({
  authService: {
    login: vi.fn(),
    logout: vi.fn(),
    me: vi.fn(),
    refresh: vi.fn(),
  },
}))

describe('useAuthStore', () => {
  const mockUser: User = {
    id: 'u-123',
    email: 'admin@openbench.local',
    role: 'ADMIN',
  }

  beforeEach(() => {
    vi.clearAllMocks()
    useAuthStore.setState({ user: null, isAuthenticated: false, isLoading: false })
  })

  it('login successfully updates auth state', async () => {
    vi.mocked(authService.login).mockResolvedValueOnce(mockUser)

    await useAuthStore.getState().login({ email: 'admin@openbench.local', password: 'password123' })

    const state = useAuthStore.getState()
    expect(state.user).toEqual(mockUser)
    expect(state.isAuthenticated).toBe(true)
    expect(state.isLoading).toBe(false)
    expect(authService.login).toHaveBeenCalledWith({
      email: 'admin@openbench.local',
      password: 'password123',
    })
  })

  it('login failure resets state and throws error', async () => {
    vi.mocked(authService.login).mockRejectedValueOnce(new Error('Invalid credentials'))

    await expect(
      useAuthStore.getState().login({ email: 'admin@openbench.local', password: 'wrong' })
    ).rejects.toThrow('Invalid credentials')

    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(state.isLoading).toBe(false)
  })

  it('logout clears auth state', async () => {
    useAuthStore.setState({ user: mockUser, isAuthenticated: true, isLoading: false })
    vi.mocked(authService.logout).mockResolvedValueOnce()

    await useAuthStore.getState().logout()

    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(authService.logout).toHaveBeenCalled()
  })

  it('checkAuth populates user when me() succeeds', async () => {
    vi.mocked(authService.me).mockResolvedValueOnce(mockUser)

    await useAuthStore.getState().checkAuth()

    const state = useAuthStore.getState()
    expect(state.user).toEqual(mockUser)
    expect(state.isAuthenticated).toBe(true)
    expect(state.isLoading).toBe(false)
  })

  it('checkAuth clears state when me() fails', async () => {
    vi.mocked(authService.me).mockRejectedValueOnce(new Error('401 Unauthorized'))

    await useAuthStore.getState().checkAuth()

    const state = useAuthStore.getState()
    expect(state.user).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(state.isLoading).toBe(false)
  })
})
