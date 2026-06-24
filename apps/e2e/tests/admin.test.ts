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

	test('should manage tickets and warranties lifecycle dynamically', async ({ page }) => {
		await page.click('a[href="/admin/tickets"]');
		await expect(page).toHaveURL('/admin/tickets');

		// Click NEW TICKET button
		await page.click('a[href="/admin/tickets/new"]');
		await expect(page).toHaveURL('/admin/tickets/new');

		// Fill out the new ticket form
		const testCustomer = `E2E Cust ${Date.now()}`;
		await page.fill('#cust-name', testCustomer);
		await page.fill('#cust-phone', '081234567890');
		await page.fill('#dev-brand', 'Google');
		await page.fill('#dev-model', 'Pixel 8 Pro');
		await page.fill('#damage', 'Layar retak rambut');
		await page.fill('#est-cost', '1500000');
		await page.selectOption('#warranty-select', '30'); // 30 Days

		// Submit ticket
		await page.click('button[type="submit"]:has-text("CREATE TICKET")');

		// Redirects to /admin/tickets
		await expect(page).toHaveURL('/admin/tickets');

		// Search for the newly created ticket
		const searchInput = page.locator('input[placeholder*="Search ticket number"]');
		await searchInput.fill(testCustomer);
		
		const card = page.locator('div.hover\\:shadow-neubrutalism-lg', { hasText: testCustomer }).first();
		await expect(card).toBeVisible({ timeout: 10000 });
		await expect(card.getByText('Google Pixel 8 Pro')).toBeVisible();

		// Navigate to details page (avoiding /new link matches)
		await card.locator('a[href*="/admin/tickets/"]:not([href*="/new"])').click();
		await expect(page).toHaveURL(/\/admin\/tickets\/[a-f0-9-]+/);

		// Enter Emergency Edit mode
		await page.click('button:has-text("EMERGENCY EDIT")');
		await page.click('button:has-text("CONFIRM EDIT")');

		// Get current URL to extract ticket ID and verify route redirection
		const detailUrl = page.url();
		const tId = detailUrl.split('/').pop();
		await expect(page).toHaveURL(`/admin/tickets/${tId}/emergency`);

		// Update ticket status to completed, location to picked_up, and pay
		await page.selectOption('#status-select', 'completed');
		await page.selectOption('#pos-select', 'picked_up');
		await page.selectOption('#pay-status', 'paid');
		await page.selectOption('#pay-method', 'cash');

		// Save changes
		await page.click('button[type="submit"]:has-text("SAVE EMERGENCY OVERRIDES")');

		// Redirection back to detail page
		await expect(page).toHaveURL(`/admin/tickets/${tId}`);

		// Go to warranties page and verify it appears there
		await page.goto('/admin/warranties');
		await expect(page).toHaveURL('/admin/warranties');
		await expect(page.getByText(testCustomer)).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('Google Pixel 8 Pro').first()).toBeVisible();

		// --- EMERGENCY EDIT E2E TEST ---
		// Navigate directly back to the ticket's details page
		await page.goto(`/admin/tickets/${tId}`);
		await expect(page).toHaveURL(`/admin/tickets/${tId}`);

		// Verify standard EDIT TICKET is not visible (since it's picked up)
		await expect(page.locator('button:has-text("EDIT TICKET")')).not.toBeVisible();

		// Click EMERGENCY EDIT
		await page.click('button:has-text("EMERGENCY EDIT")');

		// Click CONFIRM EDIT in modal
		await page.click('button:has-text("CONFIRM EDIT")');

		// Verify redirection to emergency page
		await expect(page).toHaveURL(`/admin/tickets/${tId}/emergency`);

		// Select warehouse location
		await page.selectOption('#pos-select', 'warehouse');

		// Save emergency override
		await page.click('button[type="submit"]:has-text("SAVE EMERGENCY OVERRIDES")');

		// Redirection back to detail page
		await expect(page).toHaveURL(`/admin/tickets/${tId}`);

		// Go to warranties page and verify it's gone
		await page.goto('/admin/warranties');
		await expect(page).toHaveURL('/admin/warranties');
		await expect(page.getByText(testCustomer)).not.toBeVisible({ timeout: 10000 });
	});

	test('should navigate to inventory and display data correctly', async ({ page }) => {
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

		// Negative Scenario: Insufficient Stock Check
		// Set qty to 999 which exceeds stock
		const increaseQtyBtn = page.locator('button').filter({ has: page.locator('svg.lucide-plus') });
		// Just click a few times and then we can simulate an error or type the quantity.
		// Since we don't have a direct input for qty, we'll just bypass and use the mock network or assume
		// the backend returns an error if we manage to submit it. But clicking 999 times is bad.
		// Wait, instead of 999, maybe the mock stock is 100? Let's just do a normal checkout for the positive flow
		// and mock the API or rely on the actual mock db logic which we updated to throw error if stock < qty.
		// For a robust test, let's just make sure the positive flow works first.

		// Type cash paid using robust testid
		await page.getByTestId('cash-paid-input').fill('300000');

		// Click process checkout
		await page.click('button:has-text("PROCESS CHECKOUT")');

		// Verify receipt modal and message
		await expect(page.getByText('Transaction completed successfully!')).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('INV-')).toBeVisible(); // We changed receipt text prefix from OB-INV- to INV-
	});

	test('should show error when checkout fails due to insufficient stock', async ({ page }) => {
		await page.click('a[href="/admin/sales"]');
		await expect(page).toHaveURL('/admin/sales');

		// Add item with low stock (has 3 in seed.ts)
		const catalogItem = page.getByRole('button', { name: 'Tempered Glass Ultra Clear iPhone 14 Pro' });
		await catalogItem.click();

		// Increase quantity to 4 (exceeding stock of 3)
		const increaseQtyBtn = page.locator('button').filter({ has: page.locator('svg.lucide-plus') });
		await increaseQtyBtn.click();
		await increaseQtyBtn.click();
		await increaseQtyBtn.click();

		await expect(page.getByText('Stock limited! Only 3 units available.')).toBeVisible({ timeout: 10000 });
	});
});
