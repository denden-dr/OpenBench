import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
	test.beforeEach(async ({ page }) => {
		page.on('console', msg => {
			console.log(`[BROWSER CONSOLE] ${msg.type()}: ${msg.text()}`);
		});
		page.on('pageerror', err => {
			console.error(`[BROWSER ERROR] ${err.message}`);
			if (err.stack) console.error(err.stack);
		});

		// Abort external fonts to avoid waiting for network timeouts
		await page.route('**/*', route => {
			const url = route.request().url();
			if (url.includes('fonts.googleapis.com') || url.includes('fonts.gstatic.com')) {
				route.abort();
			} else {
				route.continue();
			}
		});

		// Start at the sign-in page
		await page.goto('/auth/signin');
		// Wait for Svelte hydration to be fully complete
		await page.waitForSelector('main[data-hydrated="true"]');
	});

	test('should display the sign-in page elements correctly', async ({ page }) => {
		// Check title
		await expect(page).toHaveTitle('Sign In - OpenBench');

		// Check heading
		const heading = page.locator('h1');
		await expect(heading).toContainText('OPEN');
		await expect(heading).toContainText('BENCH');

		// Check fields
		const emailInput = page.locator('#email');
		const passwordInput = page.locator('#password');
		const submitButton = page.locator('button[type="submit"]');

		await expect(emailInput).toBeVisible();
		await expect(passwordInput).toBeVisible();
		await expect(submitButton).toContainText('SIGN IN');
	});

	test('should show error when fields are empty', async ({ page }) => {
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		// Svelte validation check
		const errorAlert = page.locator('[role="alert"]');
		await expect(errorAlert).toBeVisible();
		await expect(errorAlert).toContainText('Email address is required.');
	});

	test('should show error when password is empty', async ({ page }) => {
		await page.fill('#email', 'admin@openbench.dev');
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		const errorAlert = page.locator('[role="alert"]');
		await expect(errorAlert).toBeVisible();
		await expect(errorAlert).toContainText('Password is required.');
	});

	test('should show error on invalid credentials', async ({ page }) => {
		await page.fill('#email', 'admin@openbench.dev');
		await page.fill('#password', 'WrongPassword123!');
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		const errorAlert = page.locator('[role="alert"]');
		await expect(errorAlert).toBeVisible();
		// In mock mode, only admin@openbench.dev / SecureAdminPassword123! is accepted
		await expect(errorAlert).toContainText('Authentication Failed');
	});

	test('should sign in successfully and then log out', async ({ page }) => {
		// 1. Fill in valid mock credentials
		await page.fill('#email', 'admin@openbench.dev');
		await page.fill('#password', 'SecureAdminPassword123!');

		// 2. Submit form
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		// 3. Verify redirection to /admin
		await expect(page).toHaveURL('/admin');

		// 4. Verify admin dashboard welcome header
		const welcomeHeader = page.locator('h1').filter({ hasText: 'Welcome to the Workbench' });
		await expect(welcomeHeader).toContainText('Welcome to the Workbench, Admin!');

		// 5. Log out
		const logoutButton = page.locator('button:has-text("LOGOUT")');
		await expect(logoutButton).toBeVisible();
		await logoutButton.click();

		// 6. Verify redirected back to signin page
		await expect(page).toHaveURL('/auth/signin');
	});
});
