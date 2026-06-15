import type { UserSession } from './types';
import { mockDbService } from './db';

export const mockAuthService = {
  /**
   * Validates credentials against users inside the mock database.
   */
  async signIn(email: string, password: string): Promise<UserSession> {
    const user = await mockDbService.getUserByEmail(email);

    if (user && user.passwordHash === password) {
      const session: UserSession = {
        email: user.email,
        role: user.role,
        user_id: user.id
      };

      mockDbService.saveActiveSession(session);
      return session;
    }

    throw new Error('Invalid email or password. Please use admin@openbench.dev and SecureAdminPassword123!.');
  },

  /**
   * Destroys current active session.
   */
  async signOut(): Promise<void> {
    mockDbService.clearActiveSession();
  },

  /**
   * Verifies active session token.
   */
  getSession(): UserSession | null {
    return mockDbService.getActiveSession() as UserSession | null;
  },

  /**
   * Synchronously checks if admin session is active.
   */
  isAdminAuthenticated(): boolean {
    const session = this.getSession();
    return session !== null && session.role === 'admin';
  }
};
