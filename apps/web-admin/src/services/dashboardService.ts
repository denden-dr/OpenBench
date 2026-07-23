import api from '@/lib/api'
import type { TicketStatus } from '@/types/ticket'

export interface DashboardMetrics {
  active_tickets: number
  pending_diagnoses: number
  sales_today: number
  active_warranties: number
}

export interface RecentTicket {
  ticket_id: string
  ticket_number: string
  status: TicketStatus
  customer_name: string
  device_brand: string
  device_model: string
  cost: number
  created_at: string
}

export interface DashboardData {
  metrics: DashboardMetrics
  recent_tickets: RecentTicket[]
}

export interface DashboardResponse {
  data: DashboardData
}

export const dashboardService = {
  async getDashboard(): Promise<DashboardData> {
    const response = await api.get<DashboardResponse>('/admin/dashboard')
    return response.data.data
  }
}
