/**
 * Simulated authentication service for OpenBench frontend.
 * Allows testing of login, logout, and route guarding.
 */

// Toggle mock mode. In development, it defaults to true unless explicitly disabled.
export const isMockEnabled = () => {
  if (typeof window !== 'undefined') {
    const override = localStorage.getItem('MOCK_API');
    if (override !== null) {
      return override === 'true';
    }
  }
  return true; // Default to true for phase 1 frontend development
};

export interface UserSession {
  email: string;
  role: 'admin' | 'user';
  token: string;
}

const SIMULATED_LATENCY_MS = 600;

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const mockAuthService = {
  /**
   * Validates credentials against hardcoded admin user.
   */
  async signIn(email: string, password: string): Promise<UserSession> {
    await delay(SIMULATED_LATENCY_MS);

    // Hardcoded credentials for local/mock development
    const mockEmail = 'admin@openbench.dev';
    const mockPassword = 'SecureAdminPassword123!';

    if (email.trim().toLowerCase() === mockEmail && password === mockPassword) {
      const session: UserSession = {
        email: mockEmail,
        role: 'admin',
        token: 'mock-jwt-token-admin-12345'
      };

      if (typeof window !== 'undefined') {
        sessionStorage.setItem('openbench_session', JSON.stringify(session));
      }

      return session;
    }

    throw new Error('Invalid email or password. Please use admin@openbench.dev and SecureAdminPassword123!.');
  },

  /**
   * Destroys current active session.
   */
  async signOut(): Promise<void> {
    await delay(300);
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('openbench_session');
    }
  },

  /**
   * Verifies active session token synchronously.
   */
  getSession(): UserSession | null {
    if (typeof window === 'undefined') return null;
    const data = sessionStorage.getItem('openbench_session');
    if (!data) return null;
    try {
      return JSON.parse(data) as UserSession;
    } catch {
      return null;
    }
  },

  /**
   * Synchronously checks if admin session is active.
   */
  isAdminAuthenticated(): boolean {
    const session = this.getSession();
    return session !== null && session.role === 'admin';
  }
};
