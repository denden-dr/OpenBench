import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import AdminLayout from '@/layouts/AdminLayout'
import DashboardPage from '@/pages/DashboardPage'
import TicketsPage from '@/pages/TicketsPage'
import InventoryPage from '@/pages/InventoryPage'
import POSPage from '@/pages/POSPage'
import WarrantyPage from '@/pages/WarrantyPage'
import PlaceholderPage from '@/pages/PlaceholderPage'
import LoginPage from '@/pages/LoginPage'
import UnauthorizedPage from '@/pages/UnauthorizedPage'
import ProtectedRoute from '@/components/ProtectedRoute'
import ThemeProvider from '@/components/ThemeProvider'
import { useAuthStore } from '@/stores/authStore'
import { setupInterceptors } from '@/lib/api'

// Initialize Axios interceptor to trigger clearAuth on 401 refresh failure
setupInterceptors(() => useAuthStore.getState().clearAuth())

function App() {
  const checkAuth = useAuthStore((state) => state.checkAuth)

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  return (
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/unauthorized" element={<UnauthorizedPage />} />

          <Route
            path="/"
            element={
              <ProtectedRoute requiredRoles={['ADMIN']}>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<DashboardPage />} />
            <Route path="tickets" element={<TicketsPage />} />
            <Route path="inventory" element={<InventoryPage />} />
            <Route path="pos" element={<POSPage />} />
            <Route path="warranties" element={<WarrantyPage />} />
            <Route path="*" element={<PlaceholderPage title="Page Not Found" description="The page you requested does not exist." />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
