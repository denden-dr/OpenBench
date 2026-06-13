import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { mockAuthService } from './mockAuth';

describe('mockAuthService Unit Tests', () => {
  beforeEach(() => {
    // Reset sessionStorage before each test
    if (typeof window !== 'undefined') {
      sessionStorage.clear();
    }
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should authenticate successfully with correct admin credentials', async () => {
    const session = await mockAuthService.signIn('admin@openbench.dev', 'SecureAdminPassword123!');
    expect(session).toBeDefined();
    expect(session.email).toBe('admin@openbench.dev');
    expect(session.role).toBe('admin');
    expect(session.token).toBe('mock-jwt-token-admin-12345');
    
    // Check if saved to storage
    const stored = sessionStorage.getItem('openbench_session');
    expect(stored).not.toBeNull();
    expect(JSON.parse(stored!)).toEqual(session);
  });

  it('should throw an error with incorrect email', async () => {
    await expect(mockAuthService.signIn('wrong@openbench.dev', 'SecureAdminPassword123!'))
      .rejects
      .toThrow('Invalid email or password');
    
    expect(sessionStorage.getItem('openbench_session')).toBeNull();
  });

  it('should throw an error with incorrect password', async () => {
    await expect(mockAuthService.signIn('admin@openbench.dev', 'wrongpassword'))
      .rejects
      .toThrow('Invalid email or password');
    
    expect(sessionStorage.getItem('openbench_session')).toBeNull();
  });

  it('should return correct auth status when check is executed', async () => {
    expect(mockAuthService.isAdminAuthenticated()).toBe(false);
    
    await mockAuthService.signIn('admin@openbench.dev', 'SecureAdminPassword123!');
    expect(mockAuthService.isAdminAuthenticated()).toBe(true);
  });

  it('should invalidate session and clear storage on sign out', async () => {
    await mockAuthService.signIn('admin@openbench.dev', 'SecureAdminPassword123!');
    expect(mockAuthService.isAdminAuthenticated()).toBe(true);

    await mockAuthService.signOut();
    expect(mockAuthService.isAdminAuthenticated()).toBe(false);
    expect(sessionStorage.getItem('openbench_session')).toBeNull();
  });
});
