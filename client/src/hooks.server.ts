import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const handle: Handle = async ({ event, resolve }) => {
	// Intercept requests starting with /api
	if (event.url.pathname.startsWith('/api')) {
		const backendUrl = env.BACKEND_URL || 'http://localhost:3000';
		const targetUrl = `${backendUrl}${event.url.pathname}${event.url.search}`;

		console.log(`[Proxy] Forwarding ${event.url.pathname} to ${targetUrl}`);

		try {
			// Forward the request to the backend
			const response = await fetch(targetUrl, {
				method: event.request.method,
				headers: event.request.headers,
				body: event.request.method !== 'GET' && event.request.method !== 'HEAD' 
					? await event.request.blob() 
					: undefined,
                // @ts-expect-error - SvelteKit fetch can handle this but TS might complain about duplex for streams
                duplex: 'half'
			});

			return response;
		} catch (error) {
			console.error(`[Proxy] Error forwarding request to ${targetUrl}:`, error);
			return new Response(JSON.stringify({ 
                message: 'Service Unavailable', 
                status: 503,
                error: String(error)
            }), {
				status: 503,
				headers: { 'Content-Type': 'application/json' }
			});
		}
	}

	// For non-api requests, continue as normal
	return resolve(event);
};
