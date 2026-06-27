/**
 * Real authentication service for OpenBench frontend interacting with Go Fiber API.
 */

import { env } from '$env/dynamic/public';

export const isMockEnabled = () => {
  // Try Vite's built-in env mode variables first (populated based on --mode)
  try {
    if (import.meta.env.PUBLIC_MOCK_API !== undefined) {
      return import.meta.env.PUBLIC_MOCK_API === 'true';
    }
  } catch {
    // ignore
  }

  try {
    if (env.PUBLIC_MOCK_API !== undefined && env.PUBLIC_MOCK_API !== '') {
      return env.PUBLIC_MOCK_API === 'true';
    }
  } catch {
    // ignore env access errors
  }

  if (typeof window !== 'undefined') {
    const override = localStorage.getItem('MOCK_API');
    if (override !== null) {
      return override === 'true';
    }
  }

  try {
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
  username?: string;
  full_name?: string;
  phone_number?: string;
}



export const authService = {
  /**
   * Authenticates user against backend API.
   */
  async signIn(email: string, password: string): Promise<UserSession> {
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
      throw new Error(resBody.message || 'Invalid email or password.');
    }

    const session: UserSession = {
      email: resBody.data.email,
      role: resBody.data.role,
      userId: resBody.data.user_id || '',
      username: resBody.data.username,
      full_name: resBody.data.full_name,
      phone_number: resBody.data.phone_number
    };

    if (typeof window !== 'undefined') {
      sessionStorage.setItem('openbench_session', JSON.stringify(session));
    }

    return session;
  },

  /**
   * Registers a new user and authenticates them immediately.
   */
  async signUp(email: string, password: string): Promise<UserSession> {
    const response = await fetch(`${getApiUrl()}/api/v1/auth/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({ email, password })
    });

    const resBody = await response.json();

    if (!response.ok) {
      throw new Error(resBody.message || 'Failed to sign up.');
    }

    const session: UserSession = {
      email: resBody.data.email,
      role: resBody.data.role,
      userId: resBody.data.user_id || '',
      username: resBody.data.username,
      full_name: resBody.data.full_name,
      phone_number: resBody.data.phone_number
    };

    if (typeof window !== 'undefined') {
      sessionStorage.setItem('openbench_session', JSON.stringify(session));
    }

    return session;
  },

  /**
   * Updates user profile attributes.
   */
  async updateProfile(profile: { username: string; full_name: string; phone_number?: string }): Promise<UserSession> {
    const response = await fetch(`${getApiUrl()}/api/v1/auth/profile`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(profile)
    });

    const resBody = await response.json();

    if (!response.ok) {
      throw new Error(resBody.message || 'Failed to update profile.');
    }

    const session: UserSession = {
      email: resBody.data.email,
      role: resBody.data.role,
      userId: resBody.data.user_id || '',
      username: resBody.data.username,
      full_name: resBody.data.full_name,
      phone_number: resBody.data.phone_number
    };

    if (typeof window !== 'undefined') {
      sessionStorage.setItem('openbench_session', JSON.stringify(session));
    }

    return session;
  },

  /**
   * Calls sign out endpoint and clears local session cache.
   */
  async signOut(): Promise<void> {
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
    const response = await fetch(`${getApiUrl()}/api/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include'
    });

    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.message || 'Failed to refresh token.');
    }
  },

  /**
   * Retrieves active session user details.
   * Returns null if unauthenticated.
   */
  async checkSession(): Promise<UserSession | null> {
    try {
      let response = await fetch(`${getApiUrl()}/api/v1/auth/me`, {
        method: 'GET',
        credentials: 'include'
      });

      if (!response.ok) {
        if (response.status === 401) {
          try {
            await this.refresh();
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
          userId: resBody.data.user_id,
          username: resBody.data.username,
          full_name: resBody.data.full_name,
          phone_number: resBody.data.phone_number
        };

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
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('openbench_session');
    }
  },

  getSession(): UserSession | null {
    if (typeof window !== 'undefined') {
      try {
        const stored = sessionStorage.getItem('openbench_session');
        if (stored) return JSON.parse(stored);
      } catch (e) {}
    }
    return null;
  },

  isAdminAuthenticated(): boolean {
    const session = this.getSession();
    return session !== null && session.role === 'admin';
  }
};
