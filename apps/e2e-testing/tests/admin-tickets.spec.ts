import { test, expect } from '@playwright/test';

test.describe('Admin Ticket Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Perform login first with correct seeded email
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('admin@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('secretpassword123');
    await page.getByRole('button', { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/(tickets|dashboard)/);
  });

  test('should navigate to tickets page and list tickets', async ({ page }) => {
    await page.goto('/tickets');
    await expect(page.getByRole('heading', { name: /repair tickets/i })).toBeVisible();
    await expect(page.getByRole('table')).toBeVisible();
  });

  test('should create a new ticket successfully', async ({ page }) => {
    await page.goto('/tickets');

    // Click New Ticket button
    await page.getByRole('button', { name: /new ticket/i }).click();

    // Fill Form
    await page.getByPlaceholder('e.g. John Doe').fill('E2E Customer');
    await page.getByPlaceholder('e.g. 08123456789').fill('089988776655');
    await page.getByPlaceholder('e.g. Apple').fill('Google');
    await page.getByPlaceholder('e.g. iPhone 13 Pro').fill('Pixel 8');
    await page.getByPlaceholder('e.g. Cracked LCD, Touch unresponsive').fill('Screen glitch');
    await page.getByPlaceholder('e.g. Replacement LCD OLED Screen').fill('Replace Screen Assembly');

    // Submit
    await page.getByRole('button', { name: 'Create Ticket', exact: true }).click();

    // Verify ticket appears in table
    await expect(page.getByText('E2E Customer')).toBeVisible();
    await expect(page.getByText('Pixel 8')).toBeVisible();
  });
});
