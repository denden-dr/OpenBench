import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) => filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		})
	],
	resolve: {
		conditions: ['browser']
	},
	optimizeDeps: {
		include: ['lucide-svelte']
	},
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/health': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		globals: true,
		setupFiles: ['src/tests/setup.ts']
	}
});
