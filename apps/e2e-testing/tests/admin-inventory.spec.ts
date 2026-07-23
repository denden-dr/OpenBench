import { test, expect } from '@playwright/test';

test.describe('Admin Inventory Management Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('admin@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('secretpassword123');
    await page.getByRole('button', { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/(tickets|dashboard)/);
  });

  test('should navigate to inventory page and show product list', async ({ page }) => {
    await page.getByRole('link', { name: /product inventory/i }).click();
    await expect(page).toHaveURL(/\/inventory/);
    await expect(page.getByRole('heading', { name: 'Product Inventory' })).toBeVisible();
    await expect(page.getByText('Tempered Glass iPhone 15 Pro Max')).toBeVisible({ timeout: 10000 });
    await expect(page.getByRole('table')).toBeVisible();
  });

  test('should create a new product', async ({ page }) => {
    await page.goto('/inventory');
    await expect(page.getByRole('heading', { name: 'Product Inventory' })).toBeVisible();

    await page.getByRole('button', { name: 'Add Product' }).click();
    await expect(page.getByRole('heading', { name: 'Add New Product' })).toBeVisible();

    const addDialog = page.getByRole('dialog')
    await addDialog.getByPlaceholder('e.g. Tempered Glass iPhone 15').fill('E2E Test Charger 65W');
    await addDialog.locator('input[type="number"]').nth(0).fill('500000');
    await addDialog.locator('input[type="number"]').nth(1).fill('10');

    await page.getByRole('button', { name: 'Save Product' }).click();
    await expect(page.getByText('E2E Test Charger 65W')).toBeVisible({ timeout: 5000 });
  });

  test('should adjust product stock', async ({ page }) => {
    await page.goto('/inventory');
    await expect(page.getByRole('heading', { name: 'Product Inventory' })).toBeVisible();
    await expect(page.getByText('MicroUSB Cable 1m')).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row', { name: /MicroUSB Cable 1m/ });
    await row.getByRole('button', { name: 'Adjust Stock' }).click();

    await expect(page.getByRole('heading', { name: 'Adjust Stock Count' })).toBeVisible();

    const adjustDialog = page.getByRole('dialog')
    await adjustDialog.getByPlaceholder('e.g. 5 or -3').fill('5');
    await page.getByRole('button', { name: 'Apply Adjustment' }).click();

    await expect(row.getByText('In Stock (30)')).toBeVisible({ timeout: 5000 });
  });

  test('should edit a product', async ({ page }) => {
    await page.goto('/inventory');
    await expect(page.getByRole('heading', { name: 'Product Inventory' })).toBeVisible();
    await expect(page.getByText('Silicon Case iPhone 15')).toBeVisible({ timeout: 10000 });

    const row = page.getByRole('row', { name: /Silicon Case iPhone 15/ });
    await row.getByRole('button', { name: 'Edit Product' }).click();

    await expect(page.getByRole('heading', { name: 'Edit Product Details' })).toBeVisible();

    const editDialog = page.getByRole('dialog')
    await editDialog.locator('input').first().clear();
    await editDialog.locator('input').first().fill('Silicon Case iPhone 15 - Updated');

    await page.getByRole('button', { name: 'Save Changes' }).click();
    await expect(page.getByText('Silicon Case iPhone 15 - Updated')).toBeVisible({ timeout: 5000 });
  });

  test('should delete a product', async ({ page }) => {
    await page.goto('/inventory');
    await expect(page.getByRole('heading', { name: 'Product Inventory' })).toBeVisible();
    await expect(page.getByText('Lightning to USB-C Cable 2m')).toBeVisible({ timeout: 10000 });

    page.on('dialog', async (dialog) => {
      expect(dialog.message()).toContain('Are you sure you want to delete this product?');
      await dialog.accept();
    });

    const row = page.getByRole('row', { name: /Lightning to USB-C Cable 2m/ });
    await row.getByRole('button', { name: 'Delete Product' }).click();

    await expect(page.getByText('Lightning to USB-C Cable 2m')).not.toBeVisible({ timeout: 5000 });
  });
});
