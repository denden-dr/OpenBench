import { env } from '$env/dynamic/public';

const getApiUrl = () => {
  try {
    return env.PUBLIC_API_URL || '';
  } catch {
    return '';
  }
};

export interface Product {
  id: string;
  name: string;
  category: 'retail' | 'spare_part';
  stock: number;
  price: number;
  cost_price: number;
  min_stock: number;
}

export const inventoryService = {
  async getInventory(): Promise<Product[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/inventory`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch inventory');
    return body.data;
  },

  async getProduct(id: string): Promise<Product> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/inventory/${id}`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch product');
    return body.data;
  },

  async createProduct(product: Omit<Product, 'id'>): Promise<Product> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/inventory`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(product)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to create product');
    return body.data;
  },

  async updateProduct(id: string, updates: Partial<Product>): Promise<Product> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/inventory/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(updates)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to update product');
    return body.data;
  },

  async deleteProduct(id: string): Promise<void> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/inventory/${id}`, {
      method: 'DELETE',
      credentials: 'include'
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to delete product');
  }
};
