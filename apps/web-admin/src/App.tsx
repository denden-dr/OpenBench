import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import AdminLayout from '@/layouts/AdminLayout'
import DashboardPage from '@/pages/DashboardPage'
import TicketsPage from '@/pages/TicketsPage'
import PlaceholderPage from '@/pages/PlaceholderPage'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<AdminLayout />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="tickets" element={<TicketsPage />} />
          <Route path="inventory" element={<PlaceholderPage title="Product Inventory" />} />
          <Route path="pos" element={<PlaceholderPage title="Point of Sale (POS)" />} />
          <Route path="warranties" element={<PlaceholderPage title="Warranty Claims" />} />
          <Route path="*" element={<PlaceholderPage title="Page Not Found" description="The page you requested does not exist." />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
