import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('Admin can log in successfully', async ({ page }) => {
    await page.goto('/login');

    // Assuming the input has name="email" or id="email"
    // Since we seeded the DB, we can use the admin credentials
    await page.fill('input[name="email"]', 'admin@openbench.com');
    await page.fill('input[name="password"]', 'secretpassword123');

    // Click the submit button and wait for the HTMX request to complete and redirect to root (dashboard)
    await Promise.all([
      page.waitForURL((url) => url.pathname === '/', { timeout: 10000 }),
      page.click('button[type="submit"]'),
    ]);

    // Ensure there is some admin-specific text visible
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('Admin can log out successfully', async ({ page }) => {
    // 1. Log in first
    await page.goto('/login');
    await page.fill('input[name="email"]', 'admin@openbench.com');
    await page.fill('input[name="password"]', 'secretpassword123');
    await Promise.all([
      page.waitForURL((url) => url.pathname === '/', { timeout: 10000 }),
      page.click('button[type="submit"]'),
    ]);
    
    // 2. Click the Logout button in the sidebar to open confirmation dialog
    await page.click('button:has-text("Logout")');
    
    // 3. Wait for the confirmation dialog and confirm logout
    await expect(page.locator('text=Confirm Logout')).toBeVisible({ timeout: 5000 });
    await page.locator('button[hx-post="/api/v1/auth/logout"]').click();

    // 4. Wait for redirect to login page
    await page.waitForURL((url) => url.pathname.includes('/login'), { timeout: 10000 });

    // 5. Verify we are on the login page by checking for the login button
    await expect(page.locator('button[type="submit"]')).toContainText('Sign In');
  });
});
