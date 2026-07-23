import { test, expect } from '@playwright/test';

test.describe('Admin POS Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('admin@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('secretpassword123');
    await page.getByRole('button', { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/(tickets|dashboard)/);
  });

  test('should navigate to POS page and show product catalog', async ({ page }) => {
    await page.getByRole('link', { name: /point of sale/i }).click();
    await expect(page).toHaveURL(/\/pos/);
    await expect(page.getByRole('heading', { name: 'Point of Sale (POS)' })).toBeVisible();
    await expect(page.getByText('Tempered Glass iPhone 15 Pro Max')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText(/Your shopping cart is empty/)).toBeVisible();
  });

  test('should add items to cart and process checkout with Cash', async ({ page }) => {
    await page.goto('/pos');
    await expect(page.getByRole('heading', { name: 'Point of Sale (POS)' })).toBeVisible();
    await expect(page.getByText('Tempered Glass iPhone 15 Pro Max')).toBeVisible({ timeout: 10000 });

    await page.getByText('Tempered Glass iPhone 15 Pro Max').first().click();
    await expect(page.getByText('1 items')).toBeVisible();

    await page.getByRole('button', { name: 'Cash' }).click();

    await page.getByRole('button', { name: 'Process Checkout' }).click();
    await expect(page.getByText(/Your shopping cart is empty/)).toBeVisible({ timeout: 5000 });
  });

  test('should process checkout with QRIS payment', async ({ page }) => {
    await page.goto('/pos');
    await expect(page.getByText('USB-C Charger Adapter 20W')).toBeVisible({ timeout: 10000 });

    await page.getByText('USB-C Charger Adapter 20W').first().click();
    await expect(page.getByText('1 items')).toBeVisible();

    await page.getByRole('button', { name: 'QRIS Code' }).click();
    await page.getByRole('button', { name: 'Process Checkout' }).click();

    await expect(page.getByText(/Your shopping cart is empty/)).toBeVisible({ timeout: 5000 });
  });

  test('should view transaction history after checkout', async ({ page }) => {
    await page.goto('/pos');
    await expect(page.getByText('MicroUSB Cable 1m')).toBeVisible({ timeout: 10000 });

    await page.getByText('MicroUSB Cable 1m').first().click();
    await expect(page.getByText('1 items')).toBeVisible();
    await page.getByRole('button', { name: 'Process Checkout' }).click();

    await page.getByRole('tab', { name: 'Transaction History' }).click();

    await expect(page.getByRole('table')).toBeVisible({ timeout: 5000 });

    const detailsButton = page.getByRole('button', { name: 'Details' }).first();
    await detailsButton.click();

    await expect(page.getByRole('heading', { name: 'Transaction Details' })).toBeVisible();
    await expect(page.getByText(/Total Paid/)).toBeVisible();

    await page.getByRole('button', { name: 'Close Invoice' }).click();
  });

  test('should not allow checkout with empty cart', async ({ page }) => {
    await page.goto('/pos');
    await expect(page.getByRole('heading', { name: 'Point of Sale (POS)' })).toBeVisible();
    await expect(page.getByText('Tempered Glass iPhone 15 Pro Max')).toBeVisible({ timeout: 10000 });

    const checkoutBtn = page.getByRole('button', { name: 'Process Checkout' });
    await expect(checkoutBtn).toBeDisabled();
  });
});
