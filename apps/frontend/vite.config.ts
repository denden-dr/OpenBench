import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const backendUrl = process.env.BACKEND_URL || 'http://localhost:3000';
const useMock = process.env.MOCK_API === 'true';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: useMock ? undefined : {
			'/api': {
				target: backendUrl,
				changeOrigin: true
			}
		}
	}
});
