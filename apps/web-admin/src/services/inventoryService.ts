import api from '@/lib/api'
import type {
  Product,
  ProductResponse,
  ProductListResponse,
  CreateProductRequest,
  UpdateProductRequest,
  AdjustStockRequest,
} from '@/types/pos'

export interface GetProductsParams {
  search?: string
  limit?: number
  cursor?: string
}

export const inventoryService = {
  async getProducts(params?: GetProductsParams): Promise<ProductListResponse> {
    const response = await api.get<ProductListResponse>('/admin/products', { params })
    return response.data
  },

  async getProductByID(id: string): Promise<Product> {
    const response = await api.get<ProductResponse>(`/admin/products/${id}`)
    return response.data.data
  },

  async createProduct(data: CreateProductRequest): Promise<Product> {
    const response = await api.post<ProductResponse>('/admin/products', data)
    return response.data.data
  },

  async updateProduct(id: string, data: UpdateProductRequest): Promise<Product> {
    const response = await api.put<ProductResponse>(`/admin/products/${id}`, data)
    return response.data.data
  },

  async adjustStock(id: string, data: AdjustStockRequest): Promise<Product> {
    const response = await api.patch<ProductResponse>(`/admin/products/${id}/stock`, data)
    return response.data.data
  },

  async deleteProduct(id: string): Promise<void> {
    await api.delete(`/admin/products/${id}`)
  },
}
