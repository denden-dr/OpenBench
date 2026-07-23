import { describe, it, expect, vi, beforeEach } from 'vitest'
import { inventoryService } from '../inventoryService'
import api from '@/lib/api'
import type { Product } from '@/types/pos'

vi.mock('@/lib/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('inventoryService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getProducts', () => {
    it('fetches product list with search params', async () => {
      const mockProducts: Product[] = [
        { id: 'p-1', name: 'Tempered Glass', price: 75000, stock: 18, created_at: '2026-07-15T09:00:00Z' },
      ]
      const mockResponse = { data: { data: mockProducts, meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await inventoryService.getProducts({ search: 'Glass', limit: 10 })

      expect(api.get).toHaveBeenCalledWith('/admin/products', {
        params: { search: 'Glass', limit: 10 },
      })
      expect(result).toEqual(mockResponse.data)
    })

    it('fetches product list without params', async () => {
      const mockResponse = { data: { data: [], meta: { limit: 10 } } }
      vi.mocked(api.get).mockResolvedValueOnce(mockResponse)

      const result = await inventoryService.getProducts()

      expect(api.get).toHaveBeenCalledWith('/admin/products', { params: undefined })
      expect(result).toEqual(mockResponse.data)
    })
  })

  describe('getProductByID', () => {
    it('fetches single product by ID', async () => {
      const mockProduct: Product = { id: 'p-1', name: 'USB Cable', price: 35000, stock: 25, created_at: '2026-07-15T09:00:00Z' }
      vi.mocked(api.get).mockResolvedValueOnce({ data: { data: mockProduct } })

      const result = await inventoryService.getProductByID('p-1')

      expect(api.get).toHaveBeenCalledWith('/admin/products/p-1')
      expect(result).toEqual(mockProduct)
    })
  })

  describe('createProduct', () => {
    it('sends POST request to create a product', async () => {
      const createReq = { name: 'New Product', price: 50000, stock: 10 }
      const created: Product = { id: 'p-new', name: 'New Product', price: 50000, stock: 10, created_at: '2026-07-22T10:00:00Z' }
      vi.mocked(api.post).mockResolvedValueOnce({ data: { data: created } })

      const result = await inventoryService.createProduct(createReq)

      expect(api.post).toHaveBeenCalledWith('/admin/products', createReq)
      expect(result).toEqual(created)
    })
  })

  describe('updateProduct', () => {
    it('sends PUT request to update a product', async () => {
      const updateReq = { name: 'Updated Product', price: 60000, stock: 15 }
      const updated: Product = { id: 'p-1', ...updateReq, created_at: '2026-07-15T09:00:00Z' }
      vi.mocked(api.put).mockResolvedValueOnce({ data: { data: updated } })

      const result = await inventoryService.updateProduct('p-1', updateReq)

      expect(api.put).toHaveBeenCalledWith('/admin/products/p-1', updateReq)
      expect(result).toEqual(updated)
    })
  })

  describe('adjustStock', () => {
    it('sends PATCH request to adjust stock', async () => {
      const adjustReq = { quantity_change: 5 }
      const adjusted: Product = { id: 'p-1', name: 'Product', price: 10000, stock: 15, created_at: '2026-07-15T09:00:00Z' }
      vi.mocked(api.patch).mockResolvedValueOnce({ data: { data: adjusted } })

      const result = await inventoryService.adjustStock('p-1', adjustReq)

      expect(api.patch).toHaveBeenCalledWith('/admin/products/p-1/stock', adjustReq)
      expect(result).toEqual(adjusted)
    })
  })

  describe('deleteProduct', () => {
    it('sends DELETE request to remove a product', async () => {
      vi.mocked(api.delete).mockResolvedValueOnce({})

      await inventoryService.deleteProduct('p-1')

      expect(api.delete).toHaveBeenCalledWith('/admin/products/p-1')
    })
  })
})
