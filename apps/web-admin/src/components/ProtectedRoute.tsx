import type { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import { Loader2 } from 'lucide-react'

interface ProtectedRouteProps {
  children: ReactNode
  requiredRoles?: string[]
}

export default function ProtectedRoute({ children, requiredRoles }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuthStore()
  const location = useLocation()

  if (isLoading) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <Loader2 className="w-8 h-8 text-primary animate-spin" />
          <p className="text-sm font-medium text-slate-500">Checking authentication...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  if (requiredRoles && requiredRoles.length > 0 && user && !requiredRoles.includes(user.role)) {
    return <Navigate to="/unauthorized" replace />
  }

  return <>{children}</>
}
