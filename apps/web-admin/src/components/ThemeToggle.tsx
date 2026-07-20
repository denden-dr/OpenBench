import { Sun, Moon, Monitor } from 'lucide-react'
import { useThemeStore, type Theme } from '@/stores/themeStore'

interface ThemeToggleProps {
  className?: string
}

export function ThemeToggle({ className = '' }: ThemeToggleProps) {
  const { theme, setTheme } = useThemeStore()

  const handleCycleTheme = () => {
    const nextTheme: Record<Theme, Theme> = {
      system: 'light',
      light: 'dark',
      dark: 'system',
    }
    setTheme(nextTheme[theme])
  }

  const getLabel = () => {
    switch (theme) {
      case 'light':
        return 'Light Mode'
      case 'dark':
        return 'Dark Mode'
      case 'system':
      default:
        return 'System Theme'
    }
  }

  return (
    <button
      type="button"
      onClick={handleCycleTheme}
      title={`Theme: ${getLabel()} (Click to toggle)`}
      aria-label={`Current theme: ${getLabel()}. Click to switch.`}
      className={`p-2 rounded-xl border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-white transition-all shadow-xs flex items-center gap-1.5 text-xs font-medium cursor-pointer ${className}`}
    >
      {theme === 'light' && <Sun className="w-4 h-4 text-amber-500" />}
      {theme === 'dark' && <Moon className="w-4 h-4 text-indigo-400" />}
      {theme === 'system' && <Monitor className="w-4 h-4 text-slate-400" />}
      <span className="hidden sm:inline capitalize">{theme}</span>
    </button>
  )
}
