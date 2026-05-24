import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname === '/api' || event.url.pathname.startsWith('/api/')) {
		const backendUrl = process.env.BACKEND_URL || 'http://localhost:3000';
		const targetUrl = new URL(event.url.pathname + event.url.search, backendUrl);
		
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
