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
    
    // 2. Click the Logout button in the sidebar
    // Wait for the redirect to /login
    await Promise.all([
      page.waitForURL((url) => url.pathname.includes('/login'), { timeout: 10000 }),
      page.click('button:has-text("Logout")'),
    ]);

    // 3. Verify we are on the login page by checking for the login button
    await expect(page.locator('button[type="submit"]')).toContainText('Sign In');
  });
});
