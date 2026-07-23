import api from '@/lib/api'
import type {
  POSTransaction,
  POSCheckoutRequest,
  TransactionResponse,
  TransactionListResponse,
} from '@/types/pos'

export interface GetTransactionsParams {
  limit?: number
  cursor?: string
}

export const posService = {
  async checkout(data: POSCheckoutRequest): Promise<POSTransaction> {
    const response = await api.post<TransactionResponse>('/admin/pos/checkout', data)
    return response.data.data
  },

  async getTransactions(params?: GetTransactionsParams): Promise<TransactionListResponse> {
    const response = await api.get<TransactionListResponse>('/admin/pos/transactions', { params })
    return response.data
  },

  async getTransactionByID(id: string): Promise<POSTransaction> {
    const response = await api.get<TransactionResponse>(`/admin/pos/transactions/${id}`)
    return response.data.data
  },
}
