import { env } from '$env/dynamic/public';
import { apiFetch as fetch } from './api';

const getApiUrl = () => {
  try {
    return env.PUBLIC_API_URL || '';
  } catch {
    return '';
  }
};

export interface SaleItem {
  productId: string;
  name: string;
  price: number;
  qty: number;
}

export interface Sale {
  id: string;
  invoice_number: string;
  items: SaleItem[];
  subtotal: number;
  discount: number;
  total: number;
  payment_method: 'cash' | 'qris';
  created_at: string;
}

export const saleService = {
  async getSales(): Promise<Sale[]> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/sales`, { credentials: 'include' });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to fetch sales');
    return body.data ?? [];
  },

  async createSale(sale: { items: { productId: string; qty: number }[]; discount: number; payment_method: 'cash' | 'qris' }): Promise<Sale> {
    const res = await fetch(`${getApiUrl()}/api/v1/admin/sales`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(sale)
    });
    const body = await res.json();
    if (!res.ok) throw new Error(body.message || 'Failed to log sale');
    return body.data;
  }
};
