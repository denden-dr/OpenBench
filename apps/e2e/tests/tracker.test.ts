import { test, expect, type Page } from '@playwright/test';

// Helper function to create a ticket dynamically and return its details
async function createTicketAndGetInfo(page: Page): Promise<{ ticketId: string; ticketNumber: string }> {
	// Sign in as admin
	await page.goto('/auth/signin');
	await page.waitForSelector('main[data-hydrated="true"]');
	await page.fill('#email', 'admin@openbench.dev');
	await page.fill('#password', 'SecureAdminPassword123!');
	await page.locator('button[type="submit"]').click();
	await expect(page).toHaveURL('/admin');

	// Navigate to tickets page
	await page.click('a[href="/admin/tickets"]');
	await expect(page).toHaveURL('/admin/tickets');

	// Click NEW TICKET button
	await page.click('a[href="/admin/tickets/new"]');
	await expect(page).toHaveURL('/admin/tickets/new');

	// Fill out the new ticket form
	const testCustomer = `Tracker Cust ${Date.now()}`;
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

	// Navigate to details page (avoiding /new link matches)
	await card.locator('a[href*="/admin/tickets/"]:not([href*="/new"])').click();
	await expect(page).toHaveURL(/\/admin\/tickets\/[a-f0-9-]+/);

	// Extract the UUID from the URL
	const url = page.url();
	const parts = url.split('/');
	const ticketId = parts[parts.length - 1];

	// Extract the human-readable ticket number (OB-...) after it loads
	const ticketNumberLocator = page.locator('span.font-mono.text-sm.font-extrabold').first();
	await expect(ticketNumberLocator).toContainText('OB-');
	const ticketNumber = (await ticketNumberLocator.textContent())?.trim() || '';

	return { ticketId, ticketNumber };
}

test.describe('Public Tracker Flow', () => {
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
	});

	test('should display tracking page elements and enforce validations', async ({ page }) => {
		// Start at the public tracker page
		await page.goto('/tracker');
		await page.waitForSelector('main[data-hydrated="true"]');
		await expect(page).toHaveTitle('Track Ticket Status - OpenBench Tracker');

		// Input and submit button
		const input = page.locator('input[placeholder*="Enter Tracking Number"]');
		const button = page.locator('button[type="submit"]');

		await expect(input).toBeVisible();
		await expect(button).toBeVisible();

		// Empty input validation
		await button.click();
		await expect(page.getByText('Please enter a Tracking Number first.')).toBeVisible();

		// Invalid format validation
		await input.fill('invalid-tracking-format');
		await button.click();
		await expect(page.getByText('Invalid tracking number format. Format should be OB-YYYYMM-XXXX-XXXX (e.g., OB-202606-0001-A9X2).')).toBeVisible();
	});

	test('should successfully retrieve and display ticket status for valid ticket number', async ({ page }) => {
		// 1. Create a ticket dynamically and get its details
		const { ticketId, ticketNumber } = await createTicketAndGetInfo(page);

		// 2. Go to the public tracker page
		await page.goto('/tracker');
		await page.waitForSelector('main[data-hydrated="true"]');

		const input = page.locator('input[placeholder*="Enter Tracking Number"]');
		const button = page.locator('button[type="submit"]');

		// Search using the dynamic ticket number
		await input.fill(ticketNumber);
		await button.click();

		// Check progress card elements
		await expect(page.getByText('TICKET FOUND')).toBeVisible({ timeout: 10000 });
		await expect(page.getByRole('heading', { name: 'Google Pixel 8 Pro' })).toBeVisible();
		await expect(page.getByText('received', { exact: true })).toBeVisible();

		// 3. Go back to admin details page to mark it as picked up
		await page.goto(`/admin/tickets/${ticketId}`);
		await expect(page).toHaveURL(`/admin/tickets/${ticketId}`);

		// Enter Emergency Edit mode
		await page.click('button:has-text("EMERGENCY EDIT")');
		await page.click('button:has-text("CONFIRM EDIT")');

		// Wait for redirection to emergency page
		await expect(page).toHaveURL(`/admin/tickets/${ticketId}/emergency`);

		// Update ticket status to completed, location to picked_up, and pay
		await page.selectOption('#status-select', 'completed');
		await page.selectOption('#pos-select', 'picked_up');
		await page.selectOption('#pay-status', 'paid');
		await page.selectOption('#pay-method', 'cash');

		// Save changes
		await page.click('button[type="submit"]:has-text("SAVE EMERGENCY OVERRIDES")');

		// Wait for redirect back to detail page
		await expect(page).toHaveURL(`/admin/tickets/${ticketId}`);

		// 4. Go back to public tracker and search again to verify active warranty
		await page.goto('/tracker');
		await page.waitForSelector('main[data-hydrated="true"]');

		const input2 = page.locator('input[placeholder*="Enter Tracking Number"]');
		const button2 = page.locator('button[type="submit"]');

		await input2.fill(ticketNumber);
		await button2.click();

		// Check warranty info is active for picked_up status
		await expect(page.getByRole('heading', { name: 'Google Pixel 8 Pro' })).toBeVisible({ timeout: 10000 });
		await expect(page.getByText('Device Warranty Active')).toBeVisible();
		await expect(page.getByText('Your repair warranty is valid until')).toBeVisible();
	});
});
