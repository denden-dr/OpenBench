import { describe, it, expect, beforeEach } from 'vitest'
import { useThemeStore } from '../themeStore'

describe('themeStore', () => {
  beforeEach(() => {
    useThemeStore.setState({ theme: 'system' })
  })

  it('initializes with default system theme', () => {
    expect(useThemeStore.getState().theme).toBe('system')
  })

  it('updates theme correctly', () => {
    useThemeStore.getState().setTheme('dark')
    expect(useThemeStore.getState().theme).toBe('dark')

    useThemeStore.getState().setTheme('light')
    expect(useThemeStore.getState().theme).toBe('light')
  })
})
