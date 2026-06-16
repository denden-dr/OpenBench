import { test, expect } from '@playwright/test';

test.describe('Admin Dashboard Flow', () => {
	test.beforeEach(async ({ page }) => {
		page.on('console', msg => {
			console.log(`[BROWSER CONSOLE] ${msg.type()}: ${msg.text()}`);
		});
		page.on('pageerror', err => {
			console.error(`[BROWSER ERROR] ${err.message}`);
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
		await page.waitForSelector('main[data-hydrated="true"]');

		// Login
		await page.fill('#email', 'admin@openbench.dev');
		await page.fill('#password', 'SecureAdminPassword123!');
		await page.locator('button[type="submit"]').click();

		// Wait for dashboard to load
		await expect(page).toHaveURL('/admin');
		await expect(page.locator('h1').filter({ hasText: 'Welcome to the Workbench' })).toBeVisible();
	});

	test('should navigate to tickets and display mock data correctly', async ({ page }) => {
		await page.click('a[href="/admin/tickets"]');
		await expect(page).toHaveURL('/admin/tickets');

		// Wait for tickets to load (mock API latency)
		await expect(page.getByText('OB-202606-0001')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('Denden Hidayat')).toBeVisible();

		// Switch to Archive tab
		await page.click('button:has-text("Archive")');
		
		// The loading skeleton might appear, but wait for the resolved text
		await expect(page.getByText('OB-202606-0003')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('Alice Cooper')).toBeVisible();
	});

	test('should navigate to inventory and display mock data correctly', async ({ page }) => {
		await page.click('a[href="/admin/inventory"]');
		await expect(page).toHaveURL('/admin/inventory');

		// Wait for inventory to load
		await expect(page.getByText('Charger 25W Fast Charging Type-C')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('LCD Screen Module Samsung S23 Ultra')).toBeVisible();
	});

	test('should navigate to sales and perform checkout correctly', async ({ page }) => {
		await page.click('a[href="/admin/sales"]');
		await expect(page).toHaveURL('/admin/sales');

		// Wait for catalog to load
		const catalogItem = page.getByRole('button', { name: 'Charger 25W Fast Charging Type-C' });
		await expect(catalogItem).toBeVisible({ timeout: 10000 });

		// Add item to cart
		await catalogItem.click();

		// Check if cart has items
		await expect(page.locator('h3:has-text("Shopping Cart")')).toBeVisible();
		await expect(page.getByText('1 Items')).toBeVisible();

		// Type cash paid
		await page.fill('input[placeholder="0"] >> nth=1', '300000');

		// Click process checkout
		await page.click('button:has-text("PROCESS CHECKOUT")');

		// Verify receipt modal and message
		await expect(page.getByText('Transaction completed successfully!')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('OB-INV-')).toBeVisible();
	});

	test('should navigate to warranties and display mock data correctly', async ({ page }) => {
		await page.click('a[href="/admin/warranties"]');
		await expect(page).toHaveURL('/admin/warranties');

		// Wait for warranties to load
		await expect(page.getByText('Alice Cooper')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('Xiaomi Redmi Note 12')).toBeVisible();
	});
});
