import api from '@/lib/api'
import type { 
  TicketDetail, 
  TicketListResponse, 
  TicketDetailResponse, 
  CreateTicketRequest, 
  UpdateTicketStatusRequest, 
  UpdateTicketDetailsRequest,
  TicketStatus
} from '@/types/ticket'

export interface GetTicketsParams {
  status?: string
  search?: string
  limit?: number
  cursor?: string
}

export const ticketService = {
  async getTickets(params?: GetTicketsParams): Promise<TicketListResponse> {
    const response = await api.get<TicketListResponse>('/admin/services', { params })
    return response.data
  },

  async getTicketByID(id: string): Promise<TicketDetail> {
    const response = await api.get<TicketDetailResponse>(`/admin/services/${id}`)
    return response.data.data
  },

  async createTicket(data: CreateTicketRequest): Promise<TicketDetail> {
    const response = await api.post<TicketDetailResponse>('/admin/services', data)
    return response.data.data
  },

  async updateTicketStatus(id: string, status: TicketStatus): Promise<{ ticket_id: string; status: TicketStatus; updated_at: string }> {
    const payload: UpdateTicketStatusRequest = { status }
    const response = await api.patch<{ data: { ticket_id: string; status: TicketStatus; updated_at: string } }>(`/admin/services/${id}/status`, payload)
    return response.data.data
  },

  async updateTicketDetails(id: string, data: UpdateTicketDetailsRequest): Promise<TicketDetail> {
    const response = await api.put<TicketDetailResponse>(`/admin/services/${id}`, data)
    return response.data.data
  },

  async searchTickets(data: {
    search?: string
    start_date?: string
    end_date?: string
    exact_date?: string
    is_active?: boolean
    limit?: number
    cursor?: string
  }): Promise<TicketListResponse> {
    const response = await api.post<TicketListResponse>('/admin/services/search', data)
    return response.data
  }
}
