export interface Product {
  id: string
  name: string
  price: number
  stock: number
  created_at: string
}

export type PaymentMethod = 'CASH' | 'QRIS'

export interface CartItem {
  product: Product
  quantity: number
}

export interface POSCheckoutItem {
  id: string
  product_id: string
  product_name: string
  quantity: number
  price: number
}

export interface POSTransaction {
  id: string
  payment_method: PaymentMethod
  total_amount: number
  created_at: string
  items: POSCheckoutItem[]
}

export interface POSCheckoutRequestItem {
  product_id: string
  quantity: number
}

export interface POSCheckoutRequest {
  payment_method: PaymentMethod
  items: POSCheckoutRequestItem[]
}

// --- API Response wrapper types ---

export interface ProductResponse {
  data: Product
}

export interface ProductListResponse {
  data: Product[]
  meta: PaginationMeta
}

export interface TransactionResponse {
  data: POSTransaction
}

export interface TransactionListResponse {
  data: POSTransaction[]
  meta: PaginationMeta
}

// --- Request types matching backend contract ---

export interface CreateProductRequest {
  name: string
  price: number
  stock: number
}

export interface UpdateProductRequest {
  name: string
  price: number
  stock: number
}

export interface AdjustStockRequest {
  quantity_change: number
}

// --- Pagination (shared) ---

export interface PaginationMeta {
  limit: number
  next_cursor?: string
}
