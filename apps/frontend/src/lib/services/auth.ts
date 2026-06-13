/**
 * Real authentication service for OpenBench frontend interacting with Go Fiber API.
 * Supports fallback to mockAuthService when mock mode is enabled.
 */

import { env } from '$env/dynamic/public';
import { mockAuthService } from './mockAuth';

export const isMockEnabled = () => {
  if (typeof window !== 'undefined') {
    const override = localStorage.getItem('MOCK_API');
    if (override !== null) {
      return override === 'true';
    }
  }
  // Try reading dynamic public runtime env or build-time injected environment variable
  try {
    if (env.PUBLIC_MOCK_API !== undefined) {
      return env.PUBLIC_MOCK_API === 'true';
    }
    return import.meta.env.VITE_MOCK_API === 'true';
  } catch {
    return false;
  }
};

const getApiUrl = () => {
  try {
    return env.PUBLIC_API_URL || '';
  } catch {
    return '';
  }
};

export interface UserSession {
  email: string;
  role: 'admin' | 'user';
  userId: string;
}

// In-memory or sessionStorage cached session for synchronous check
let cachedSession: UserSession | null = null;

if (typeof window !== 'undefined') {
  try {
    const stored = sessionStorage.getItem('openbench_session');
    if (stored) {
      cachedSession = JSON.parse(stored);
    }
  } catch (e) {
    // Ignore parse errors
  }
}

export const authService = {
  /**
   * Authenticates user against backend API or mock service.
   */
  async signIn(email: string, password: string): Promise<UserSession> {
    if (isMockEnabled()) {
      const mockSession = await mockAuthService.signIn(email, password);
      const session: UserSession = {
        email: mockSession.email,
        role: mockSession.role,
        userId: 'mock-admin-id'
      };
      cachedSession = session;
      return session;
    }

    const response = await fetch(`${getApiUrl()}/api/v1/auth/signin`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({ email, password })
    });

    const resBody = await response.json();

    if (!response.ok) {
      throw new Error(resBody.error || 'Invalid email or password.');
    }

    const session: UserSession = {
      email: resBody.data.email,
      role: resBody.data.role,
      userId: resBody.data.user_id || ''
    };

    cachedSession = session;
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('openbench_session', JSON.stringify(session));
    }

    return session;
  },

  /**
   * Calls sign out endpoint and clears local session cache.
   */
  async signOut(): Promise<void> {
    if (isMockEnabled()) {
      await mockAuthService.signOut();
      this.clearLocalSession();
      return;
    }

    try {
      await fetch(`${getApiUrl()}/api/v1/auth/signout`, {
        method: 'POST',
        credentials: 'include'
      });
    } catch (err) {
      console.error('Failed to sign out from backend:', err);
    } finally {
      this.clearLocalSession();
    }
  },

  /**
   * Refreshes access token using refresh token cookie.
   */
  async refresh(): Promise<void> {
    if (isMockEnabled()) {
      return;
    }

    const response = await fetch(`${getApiUrl()}/api/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include'
    });

    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || 'Failed to refresh token.');
    }
  },

  /**
   * Retrieves active session user details.
   * Returns null if unauthenticated.
   */
  async checkSession(): Promise<UserSession | null> {
    if (isMockEnabled()) {
      const mockSession = mockAuthService.getSession();
      if (mockSession) {
        const session: UserSession = {
          email: mockSession.email,
          role: mockSession.role,
          userId: 'mock-admin-id'
        };
        cachedSession = session;
        return session;
      }
      this.clearLocalSession();
      return null;
    }

    try {
      let response = await fetch(`${getApiUrl()}/api/v1/auth/me`, {
        method: 'GET',
        credentials: 'include'
      });

      if (!response.ok) {
        // Try refreshing if token is expired (401)
        if (response.status === 401) {
          try {
            await this.refresh();
            // Retry fetch /me after successful refresh
            response = await fetch(`${getApiUrl()}/api/v1/auth/me`, {
              method: 'GET',
              credentials: 'include'
            });
          } catch (refreshErr) {
            console.warn('Session refresh failed:', refreshErr);
            this.clearLocalSession();
            return null;
          }
        } else {
          this.clearLocalSession();
          return null;
        }
      }

      if (response.ok) {
        const resBody = await response.json();
        const session: UserSession = {
          email: resBody.data.email,
          role: resBody.data.role,
          userId: resBody.data.user_id
        };

        cachedSession = session;
        if (typeof window !== 'undefined') {
          sessionStorage.setItem('openbench_session', JSON.stringify(session));
        }
        return session;
      }

      this.clearLocalSession();
      return null;
    } catch (err) {
      console.error('Error checking session:', err);
      this.clearLocalSession();
      return null;
    }
  },

  clearLocalSession(): void {
    cachedSession = null;
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('openbench_session');
    }
  },

  getSession(): UserSession | null {
    if (isMockEnabled()) {
      const mockSession = mockAuthService.getSession();
      if (mockSession) {
        return {
          email: mockSession.email,
          role: mockSession.role,
          userId: 'mock-admin-id'
        };
      }
      return null;
    }
    return cachedSession;
  },

  isAdminAuthenticated(): boolean {
    if (isMockEnabled()) {
      return mockAuthService.isAdminAuthenticated();
    }
    return cachedSession !== null && cachedSession.role === 'admin';
  }
};
