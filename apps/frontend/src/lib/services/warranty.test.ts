import { describe, it, expect, vi, afterEach } from 'vitest';
import { warrantyService } from './warranty';

describe('warrantyService Unit Tests', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  const createMockResponse = (ok: boolean, data: any = null, message: string = '') => {
    return {
      ok,
      json: async () => ({ code: ok ? 200 : 400, message, data })
    };
  };

  it('should get all warranties', async () => {
    const mockWarranties = [{ id: 'war-1', status: 'active' }];
    const fetchSpy = vi.fn().mockResolvedValue(createMockResponse(true, mockWarranties));
    vi.stubGlobal('fetch', fetchSpy);

    const result = await warrantyService.getWarranties();
    
    expect(fetchSpy).toHaveBeenCalledWith(expect.stringContaining('/api/v1/admin/warranties'), { credentials: 'include' });
    expect(result).toEqual(mockWarranties);
  });

  it('should handle getWarranties error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(createMockResponse(false, null, 'Error fetching warranties')));
    await expect(warrantyService.getWarranties()).rejects.toThrow('Error fetching warranties');
  });
});
