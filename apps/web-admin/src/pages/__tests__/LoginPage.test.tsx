import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import LoginPage from '../LoginPage'
import { useAuthStore } from '@/stores/authStore'
import { authService } from '@/services/authService'

vi.mock('@/services/authService', () => ({
  authService: {
    login: vi.fn(),
    logout: vi.fn(),
    me: vi.fn(),
    refresh: vi.fn(),
  },
}))

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    useAuthStore.setState({ user: null, isAuthenticated: false, isLoading: false })
  })

  it('renders login form elements correctly', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    )

    expect(screen.getByText('OpenBench')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('admin@openbench.local')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('••••••••••••')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('handles successful login and navigation', async () => {
    const mockUser = { id: 'u1', email: 'admin@openbench.local', role: 'ADMIN' }
    vi.mocked(authService.login).mockResolvedValueOnce(mockUser)

    render(
      <MemoryRouter initialEntries={['/login']}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={<div>Dashboard Page</div>} />
        </Routes>
      </MemoryRouter>
    )

    fireEvent.change(screen.getByPlaceholderText('admin@openbench.local'), {
      target: { value: 'admin@openbench.local' },
    })
    fireEvent.change(screen.getByPlaceholderText('••••••••••••'), {
      target: { value: 'password123' },
    })

    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))

    await waitFor(() => {
      expect(screen.getByText('Dashboard Page')).toBeInTheDocument()
    })
    expect(useAuthStore.getState().isAuthenticated).toBe(true)
  })

  it('displays error message when login fails', async () => {
    const error = new Error('Invalid email or password')
    vi.mocked(authService.login).mockRejectedValueOnce(error)

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    )

    fireEvent.change(screen.getByPlaceholderText('admin@openbench.local'), {
      target: { value: 'wrong@openbench.local' },
    })
    fireEvent.change(screen.getByPlaceholderText('••••••••••••'), {
      target: { value: 'wrongpass' },
    })

    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))

    await waitFor(() => {
      expect(screen.getByText('An unexpected error occurred. Please try again.')).toBeInTheDocument()
    })
    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })

  it('redirects to dashboard automatically if user is already authenticated', async () => {
    useAuthStore.setState({
      user: { id: 'u1', email: 'admin@openbench.local', role: 'ADMIN' },
      isAuthenticated: true,
      isLoading: false,
    })

    render(
      <MemoryRouter initialEntries={['/login']}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={<div>Dashboard Page</div>} />
        </Routes>
      </MemoryRouter>
    )

    await waitFor(() => {
      expect(screen.getByText('Dashboard Page')).toBeInTheDocument()
    })
  })
})
