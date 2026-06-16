/// <reference types="node" />
import { defineConfig, devices } from '@playwright/test';

const baseURL = process.env.BASE_URL || 'http://localhost:5173';
const playMode = process.env.PLAYWRIGHT_MODE || 'mock';

const config = defineConfig({
	testDir: './tests',
	timeout: 60000,
	fullyParallel: false,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 2 : 0,
	workers: 1,
	reporter: 'html',
	use: {
		baseURL,
		trace: 'on-first-retry',
	},
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] },
		},
	],
});

// Start a local dev server or wait for composed environment to be ready
if (process.env.BASE_URL) {
	config.webServer = {
		command: 'echo "Waiting for external environment..."',
		url: baseURL,
		reuseExistingServer: true,
		timeout: 60000,
	};
} else {
	config.webServer = {
		command: playMode === 'dev' ? 'npm --prefix ../frontend run dev' : 'npm --prefix ../frontend run dev:mock',
		url: 'http://localhost:5173',
		reuseExistingServer: !process.env.CI,
		timeout: 120000,
	};
}

export default config;
