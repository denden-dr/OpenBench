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

// Only start a local dev server if we aren't testing an external/composed environment
if (!process.env.BASE_URL) {
	config.webServer = {
		command: playMode === 'dev' ? 'npm run dev' : 'npm run dev:mock',
		url: 'http://localhost:5173',
		reuseExistingServer: !process.env.CI,
		timeout: 120000,
	};
}

export default config;
