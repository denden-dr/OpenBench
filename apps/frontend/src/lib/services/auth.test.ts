import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { authService } from './auth';

describe('authService Unit Tests', () => {
  beforeEach(() => {
    if (typeof window !== 'undefined') {
      sessionStorage.clear();
      authService.clearLocalSession();
    }
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should authenticate successfully with correct credentials', async () => {
    const mockResponse = {
      ok: true,
      json: async () => ({
        code: 200,
        message: 'Successfully signed in',
        data: { role: 'admin', user_id: 'user-id-123', email: 'admin@openbench.dev' }
      })
    };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(mockResponse));

    const session = await authService.signIn('admin@openbench.dev', 'adminpassword');
    expect(session).toBeDefined();
    expect(session.email).toBe('admin@openbench.dev');
    expect(session.role).toBe('admin');
    expect(session.userId).toBe('user-id-123');
    
    // Check if saved to storage
    const stored = sessionStorage.getItem('openbench_session');
    expect(stored).not.toBeNull();
    expect(JSON.parse(stored!)).toEqual(session);
    expect(authService.isAdminAuthenticated()).toBe(true);
  });

  it('should throw an error on failed sign in', async () => {
    const mockResponse = {
      ok: false,
      json: async () => ({ code: 400, message: 'Invalid email or password', data: null })
    };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(mockResponse));

    await expect(authService.signIn('wrong@openbench.dev', 'adminpassword'))
      .rejects
      .toThrow('Invalid email or password');
    
    expect(sessionStorage.getItem('openbench_session')).toBeNull();
    expect(authService.isAdminAuthenticated()).toBe(false);
  });

  it('should clear session and call signout endpoint on signOut', async () => {
    // First, populate session
    const mockSignInResponse = {
      ok: true,
      json: async () => ({
        code: 200,
        message: 'Successfully signed in',
        data: { role: 'admin', user_id: 'user-id-123', email: 'admin@openbench.dev' }
      })
    };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(mockSignInResponse));
    await authService.signIn('admin@openbench.dev', 'adminpassword');
    expect(authService.isAdminAuthenticated()).toBe(true);

    // Mock signout endpoint response
    const mockSignOutResponse = {
      ok: true,
      json: async () => ({ message: 'Successfully signed out' })
    };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(mockSignOutResponse));

    await authService.signOut();
    expect(authService.isAdminAuthenticated()).toBe(false);
    expect(sessionStorage.getItem('openbench_session')).toBeNull();
  });

  it('should verify session and handle refresh logic in checkSession', async () => {
    // Mock /me failing first with 401, then /refresh succeeding, then /me succeeding
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ code: 401, message: 'Authentication required', data: null })
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({ message: 'Tokens rotated successfully' })
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          code: 200,
          message: 'Success',
          data: { role: 'admin', user_id: 'user-id-123', email: 'admin@openbench.dev' }
        })
      });
    vi.stubGlobal('fetch', fetchMock);

    const session = await authService.checkSession();
    expect(session).not.toBeNull();
    expect(session!.role).toBe('admin');
    expect(session!.userId).toBe('user-id-123');
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });
});
