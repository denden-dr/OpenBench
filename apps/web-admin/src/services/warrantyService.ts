import api from '@/lib/api'
import type {
  Warranty,
  WarrantyClaim,
  WarrantyResponse,
  ClaimResponse,
  ClaimsListResponse,
  CreateClaimRequest,
  UpdateClaimRequest,
  EvaluateClaimRequest,
  UpdateWarrantyStatusRequest,
} from '@/types/warranty'

export interface GetClaimsParams {
  status?: string
  search?: string
  limit?: number
  cursor?: string
}

export const warrantyService = {
  async getWarrantyByTicketNumber(ticketNumber: string): Promise<Warranty> {
    const response = await api.get<WarrantyResponse>(`/admin/warranties/by-ticket-number/${ticketNumber}`)
    return response.data.data
  },

  async updateWarrantyStatus(warrantyId: string, data: UpdateWarrantyStatusRequest): Promise<Warranty> {
    const response = await api.patch<WarrantyResponse>(`/admin/warranties/${warrantyId}/status`, data)
    return response.data.data
  },

  async createClaim(data: CreateClaimRequest): Promise<WarrantyClaim> {
    const response = await api.post<ClaimResponse>('/admin/claims', data)
    return response.data.data
  },

  async getClaims(params?: GetClaimsParams): Promise<ClaimsListResponse> {
    const response = await api.get<ClaimsListResponse>('/admin/claims', { params })
    return response.data
  },

  async getClaimByID(claimId: string): Promise<WarrantyClaim> {
    const response = await api.get<ClaimResponse>(`/admin/claims/${claimId}`)
    return response.data.data
  },

  async updateClaim(claimId: string, data: UpdateClaimRequest): Promise<WarrantyClaim> {
    const response = await api.put<ClaimResponse>(`/admin/claims/${claimId}`, data)
    return response.data.data
  },

  async evaluateClaim(claimId: string, data: EvaluateClaimRequest): Promise<WarrantyClaim> {
    const response = await api.post<ClaimResponse>(`/admin/claims/${claimId}/evaluate`, data)
    return response.data.data
  },
}
