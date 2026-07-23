import { describe, it, expect, vi, beforeEach } from 'vitest'
import { posService } from '../posService'
import api from '@/lib/api'
import type { POSTransaction, POSCheckoutRequest } from '@/types/pos'

vi.mock('@/lib/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

describe('posService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('checkout', () => {
    it('sends POST request to process checkout', async () => {
      const checkoutReq: POSCheckoutRequest = {
        payment_method: 'CASH',
        items: [{ product_id: 'p-1', quantity: 2 }],
      }
      const mockTx: POSTransaction = {
        id: 'tx-1',
        payment_method: 'CASH',
        total_amount: 150000,
        created_at: '2026-07-22T10:00:00Z',
        items: [
          { id: 'item-1', product_id: 'p-1', product_name: 'Product', quantity: 2, price: 75000 },
        ],
      }
      vi.mocked(api.post).mockResolvedValueOnce({ data: { data: mockTx } })

      const result = await posService.checkout(checkoutReq)

      expect(api.post).toHaveBeenCalledWith('/admin/pos/checkout', checkoutReq)
      expect(result).toEqual(mockTx)
    })
  })

  describe('getTransactions', () => {
    it('fetches transaction list with params', async () => {
      const mockTxs: POSTransaction[] = [
        { id: 'tx-1', payment_method: 'QRIS', total_amount: 150000, created_at: '2026-07-22T10:00:00Z', items: [] },
      ]
      const mockResponse = { data: { data: mockTxs, meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await posService.getTransactions({ limit: 10 })

      expect(api.get).toHaveBeenCalledWith('/admin/pos/transactions', { params: { limit: 10 } })
      expect(result).toEqual(mockResponse.data)
    })

    it('fetches transaction list without params', async () => {
      const mockResponse = { data: { data: [], meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await posService.getTransactions()

      expect(api.get).toHaveBeenCalledWith('/admin/pos/transactions', { params: undefined })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('getTransactionByID', () => {
    it('fetches single transaction by ID', async () => {
      const mockTx: POSTransaction = {
        id: 'tx-1',
        payment_method: 'CASH',
        total_amount: 100000,
        created_at: '2026-07-22T10:00:00Z',
        items: [{ id: 'item-1', product_id: 'p-1', product_name: 'Product', quantity: 1, price: 100000 }],
      }
      vi.mocked(api.get).mockResolvedValueOnce({ data: { data: mockTx } })

      const result = await posService.getTransactionByID('tx-1')

      expect(api.get).toHaveBeenCalledWith('/admin/pos/transactions/tx-1')
      expect(result).toEqual(mockTx)
    })
  })
})
