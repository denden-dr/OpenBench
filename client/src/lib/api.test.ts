import { describe, it, expect, vi } from 'vitest';
import { fetchHealth } from './api';

describe('api client', () => {
    it('fetchHealth should return data when successful', async () => {
        const mockResponse = {
            message: 'OK',
            status: 200,
            data: { version: '0.1.0' }
        };

        global.fetch = vi.fn().mockResolvedValue({
            ok: true,
            json: () => Promise.resolve(mockResponse)
        });

        const result = await fetchHealth();
        expect(result).toEqual(mockResponse);
        expect(global.fetch).toHaveBeenCalledWith('/api/health');
    });

    it('fetchHealth should throw error when not ok', async () => {
        global.fetch = vi.fn().mockResolvedValue({
            ok: false,
            statusText: 'Internal Server Error'
        });

        await expect(fetchHealth()).rejects.toThrow('API error: Internal Server Error');
    });
});
