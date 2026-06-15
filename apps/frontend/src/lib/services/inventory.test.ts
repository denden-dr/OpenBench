import { describe, it, expect, vi, afterEach } from 'vitest';
import { inventoryService } from './inventory';

describe('inventoryService Unit Tests', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  const createMockResponse = (ok: boolean, data: any = null, message: string = '') => {
    return {
      ok,
      json: async () => ({ code: ok ? 200 : 400, message, data })
    };
  };

  it('should get all inventory', async () => {
    const mockProducts = [{ id: 'prod-1', name: 'Product 1' }];
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockProducts));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await inventoryService.getInventory();
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/inventory'), { credentials: 'include' });
    expect(result).toEqual(mockProducts);
  });

  it('should handle getInventory error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(createMockResponse(false, null, 'Error fetching')));
    await expect(inventoryService.getInventory()).rejects.toThrow('Error fetching');
  });

  it('should get a product by id', async () => {
    const mockProduct = { id: 'prod-1', name: 'Product 1' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockProduct));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await inventoryService.getProduct('prod-1');
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/inventory/prod-1'), { credentials: 'include' });
    expect(result).toEqual(mockProduct);
  });

  it('should create a product', async () => {
    const newProduct = { name: 'Product 2', category: 'retail', stock: 10, price: 100, cost_price: 50, min_stock: 5 } as any;
    const createdProduct = { ...newProduct, id: 'prod-2' };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, createdProduct));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await inventoryService.createProduct(newProduct);
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/inventory'), expect.objectContaining({
      method: 'POST',
      body: JSON.stringify(newProduct)
    }));
    expect(result).toEqual(createdProduct);
  });

  it('should update a product', async () => {
    const updates = { stock: 15 };
    const updatedProduct = { id: 'prod-1', name: 'Product 1', stock: 15 };
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, updatedProduct));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await inventoryService.updateProduct('prod-1', updates);
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/inventory/prod-1'), expect.objectContaining({
      method: 'PATCH',
      body: JSON.stringify(updates)
    }));
    expect(result).toEqual(updatedProduct);
  });

  it('should delete a product', async () => {
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true));
    vi.stubGlobal('fetch', fetchSpy);

    await inventoryService.deleteProduct('prod-1');
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/inventory/prod-1'), expect.objectContaining({
      method: 'DELETE'
    }));
  });
});
