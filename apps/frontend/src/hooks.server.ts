import type { Handle } from '@sveltejs/kit';

const mutationMethods = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

function parseAllowedOrigins(value: string | undefined): string[] {
	if (!value) return [];

	return value
		.split(',')
		.map((origin) => origin.trim())
		.filter(Boolean)
		.map((origin) => new URL(origin).origin);
}

function getAllowedCsrfOrigins(requestOrigin: string): string[] {
	const configuredOrigins = parseAllowedOrigins(process.env.CSRF_ALLOWED_ORIGINS);
	if (configuredOrigins.length > 0) {
		return configuredOrigins;
	}

	const adapterOrigin = parseAllowedOrigins(process.env.ORIGIN);
	if (adapterOrigin.length > 0) {
		return adapterOrigin;
	}

	return [requestOrigin];
}

export const handle: Handle = async ({ event, resolve }) => {
	if (process.env.MOCK_API === 'true') {
		const { handleMockRequest } = await import('$lib/mocks/handlers');
		const response = await handleMockRequest(event);
		if (response) {
			return response;
		}
		return resolve(event);
	}

	if (event.url.pathname === '/api' || event.url.pathname.startsWith('/api/')) {
		const isProduction = process.env.NODE_ENV === 'production';
		const backendUrl = process.env.BACKEND_URL;
		
		if (isProduction && !backendUrl) {
			return new Response(JSON.stringify({ success: false, error: 'Server configuration error: BACKEND_URL missing' }), {
				status: 500,
				headers: { 'Content-Type': 'application/json' }
			});
		}

		const actualBackendUrl = backendUrl || 'http://localhost:3000';

		// CSRF prevention for proxy routes. In containers/proxies, event.url.origin
		// can differ from the public browser origin, so prefer explicit config.
		if (mutationMethods.has(event.request.method)) {
			const origin = event.request.headers.get('origin');
			const allowedOrigins = getAllowedCsrfOrigins(event.url.origin);
			if (origin && !allowedOrigins.includes(new URL(origin).origin)) {
				return new Response(JSON.stringify({ success: false, error: 'CSRF Forbidden' }), {
					status: 403,
					headers: { 'Content-Type': 'application/json' }
				});
			}
		}

		const targetUrl = new URL(event.url.pathname + event.url.search, actualBackendUrl);
		
		const headers = new Headers(event.request.headers);
		headers.delete('host');
		headers.delete('connection');
		
		let body: ArrayBuffer | null = null;
		if (event.request.method !== 'GET' && event.request.method !== 'HEAD') {
			body = await event.request.arrayBuffer();
		}
		
		try {
			const response = await fetch(targetUrl.toString(), {
				method: event.request.method,
				headers,
				body,
				duplex: 'half'
			} as RequestInit);
			
			return response;
		} catch (error) {
			console.error(`[Proxy Error] Failed to fetch from backend: ${error}`);
			return new Response(JSON.stringify({ success: false, error: 'Backend proxy error' }), {
				status: 502,
				headers: { 'Content-Type': 'application/json' }
			});
		}
	}

	return resolve(event);
};
