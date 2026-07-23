import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import { ThemeToggle } from '@/components/ThemeToggle'
import { Wrench } from 'lucide-react'
import LoginForm from '@/features/auth/components/LoginForm'

export default function LoginPage() {
  const { isAuthenticated, isLoading } = useAuthStore()
  const navigate = useNavigate()
  const location = useLocation()

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname || '/dashboard'

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      navigate(from, { replace: true })
    }
  }, [isAuthenticated, isLoading, navigate, from])

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex items-center justify-center p-4 relative overflow-hidden selection:bg-primary selection:text-white transition-colors duration-200">
      {/* Theme Toggle Top Right */}
      <div className="absolute top-6 right-6 z-20">
        <ThemeToggle />
      </div>

      {/* Dynamic Background Accents */}
      <div className="absolute -top-40 -left-40 w-96 h-96 bg-primary/20 dark:bg-primary/30 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute -bottom-40 -right-40 w-96 h-96 bg-tertiary/20 dark:bg-tertiary/20 rounded-full blur-3xl pointer-events-none" />

      {/* Main Glassmorphism Card */}
      <div className="w-full max-w-md relative z-10 bg-white/80 dark:bg-slate-900/70 backdrop-blur-2xl border border-slate-200/80 dark:border-white/10 shadow-2xl rounded-3xl p-8 space-y-8 animate-in fade-in zoom-in-95 duration-300">
        
        {/* Header Branding */}
        <div className="text-center space-y-3">
          <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-gradient-to-br from-primary to-secondary text-white shadow-lg shadow-primary/25 border border-white/20 mx-auto">
            <Wrench className="w-7 h-7" />
          </div>
          <div>
            <h1 className="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight font-heading flex items-center justify-center gap-1.5">
              OpenBench <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-primary/10 dark:bg-primary/20 text-primary dark:text-accent border border-primary/20 dark:border-accent/20">Admin</span>
            </h1>
            <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">Sign in to access your repair bench dashboard</p>
          </div>
        </div>

        {/* LoginForm component */}
        <LoginForm from={from} />

        {/* Footer info */}
        <p className="text-center text-xs text-slate-400 dark:text-slate-500">
          OpenBench Repair Management System &bull; Enterprise Admin
        </p>
      </div>
    </div>
  )
}
