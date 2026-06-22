import { authService } from './auth';

/**
 * apiFetch is a wrapper around the native fetch API.
 * It intercepts 401 Unauthorized responses to automatically attempt
 * refreshing the session token. If the refresh succeeds, it retries
 * the original request. If the refresh fails, it clears local state
 * and redirects the user to the sign-in page.
 */
export async function apiFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  let response = init !== undefined ? await fetch(input, init) : await fetch(input);

  // Intercept 401 Unauthorized responses
  if (response.status === 401) {
    try {
      // Attempt to refresh the session token
      await authService.refresh();
      
      // If refresh is successful, retry the original request
      // Note: Request objects can only be consumed once if they have a body.
      // However, our services typically pass URL strings and init objects,
      // which allows fetch to be cleanly re-invoked.
      response = init !== undefined ? await fetch(input, init) : await fetch(input);
    } catch (refreshErr) {
      console.warn('Session refresh failed during apiFetch intercept:', refreshErr);
      
      // Token is fully expired or revoked.
      authService.clearLocalSession();
      
      // Perform a hard client-side redirect to sign in, ensuring state reset.
      if (typeof window !== 'undefined') {
        window.location.href = '/sign-in';
      }
    }
  }

  return response;
}
