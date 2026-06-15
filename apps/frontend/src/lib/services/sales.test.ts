import { describe, it, expect, vi, afterEach } from 'vitest';
import { saleService } from './sales';

describe('saleService Unit Tests', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  const createMockResponse = (ok: boolean, data: any = null, message: string = '') => {
    return {
      ok,
      json: async () => ({ code: ok ? 200 : 400, message, data })
    };
  };

  it('should get all sales', async () => {
    const mockSales = [{ id: 'sale-1', invoice_number: 'INV-001' }];
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockSales));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await saleService.getSales();
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/sales'), { credentials: 'include' });
    expect(result).toEqual(mockSales);
  });

  it('should handle getSales error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(createMockResponse(false, null, 'Error fetching sales')));
    await expect(saleService.getSales()).rejects.toThrow('Error fetching sales');
  });

  it('should create a sale', async () => {
    const newSale = { items: [], subtotal: 100, discount: 0, total: 100, payment_method: 'cash' } as any;
    const createdSale = { ...newSale, id: 'sale-1', invoice_number: 'INV-001' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, createdSale));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await saleService.createSale(newSale);
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/sales'), expect.objectContaining({
      method: 'POST',
      body: JSON.stringify(newSale)
    }));
    expect(result).toEqual(createdSale);
  });
});
