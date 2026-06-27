import { test, expect } from '@playwright/test';

test.describe('Customer Registration Flow', () => {
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

		// Start at the sign-up page
		await page.goto('/auth/signup');
		// Wait for Svelte hydration to be fully complete
		await page.waitForSelector('main[data-hydrated="true"]');
	});

	test('should display the sign-up page elements correctly', async ({ page }) => {
		// Check title
		await expect(page).toHaveTitle('Sign Up - OpenBench');

		// Check heading
		const heading = page.locator('h2');
		await expect(heading).toContainText('SIGN UP');

		// Check inputs
		const emailInput = page.locator('#email');
		const passwordInput = page.locator('#password');
		const submitButton = page.locator('button[type="submit"]');

		await expect(emailInput).toBeVisible();
		await expect(passwordInput).toBeVisible();
		await expect(submitButton).toContainText('CREATE ACCOUNT');
	});

	test('should show error when email format is invalid', async ({ page }) => {
		await page.fill('#email', 'invalid-email');
		await page.fill('#password', 'ValidPassword123!');
		
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		const errorAlert = page.locator('[role="alert"]');
		await expect(errorAlert).toBeVisible();
		await expect(errorAlert).toContainText('Please enter a valid email address.');
	});

	test('should show error when password is less than 6 characters', async ({ page }) => {
		await page.fill('#email', 'valid-email@openbench.dev');
		await page.fill('#password', 'short');
		
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		const errorAlert = page.locator('[role="alert"]');
		await expect(errorAlert).toBeVisible();
		await expect(errorAlert).toContainText('Password must be at least 6 characters long.');
	});

	test('should sign up successfully and redirect to setup onboarding', async ({ page }) => {
		// Fill in valid email and password
		await page.fill('#email', 'new-customer@openbench.dev');
		await page.fill('#password', 'SecurePassword123!');

		// Submit form
		const submitButton = page.locator('button[type="submit"]');
		await submitButton.click();

		// Verify redirect to portal setup onboarding page
		await expect(page).toHaveURL(/\/portal\/setup/);
	});
});
