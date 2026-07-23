import { test, expect } from '@playwright/test';

test.describe('Admin Authentication Flow', () => {
  test('should display login page elements', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByRole('heading', { name: /OpenBench/i })).toBeVisible();
    await expect(page.getByPlaceholder('admin@openbench.local')).toBeVisible();
    await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible();
  });

  test('should fail login with invalid credentials', async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('invalid@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('wrongpassword');
    await page.getByRole('button', { name: /sign in/i }).click();

    // Check for error alert message
    await expect(page.getByText(/kredensial|invalid|error|gagal|failed/i)).toBeVisible();
  });

  test('should login successfully with admin credentials and redirect to tickets/dashboard', async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('admin@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('secretpassword123');
    await page.getByRole('button', { name: /sign in/i }).click();

    await expect(page).toHaveURL(/\/(tickets|dashboard)/);
  });
});
