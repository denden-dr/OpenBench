import { test, expect } from '@playwright/test';

test.describe.serial('POS and Inventory Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin before each test
    await page.goto('/login');
    await page.fill('input[name="email"]', 'admin@openbench.com');
    await page.fill('input[name="password"]', 'secretpassword123');
    await Promise.all([
      page.waitForURL((url) => url.pathname === '/', { timeout: 10000 }),
      page.click('button[type="submit"]'),
    ]);
  });

  test('Add products to inventory, checkout, and verify history', async ({ page }) => {
    // ---------------------------------------------------------
    // PART 1: Inventory Management
    // ---------------------------------------------------------
    
    // 1. Navigate to inventory page
    await page.goto('/pos/inventory');
    await expect(page.locator('h1')).toContainText('POS & Inventory');

    // 2. Open the "Add Product" drawer
    await page.click('button:has-text("Add Product")');

    // Wait for the drawer title
    const drawerTitle = page.locator('h2:has-text("Add New Product")');
    await expect(drawerTitle).toBeVisible();

    // 3. Fill out the form for the first product
    const productName = `Test Charger ${Date.now()}`;
    await page.fill('input[name="name"]', productName);
    await page.fill('input[name="price"]', '250000');
    await page.fill('input[name="stock"]', '5');

    // 4. Submit the form and wait for navigation (since drawer does location.reload())
    await Promise.all([
      page.waitForNavigation(),
      page.click('button[type="submit"]:has-text("Save Product")')
    ]);

    // 5. Verify the new product appears in the table
    const tableBody = page.locator('tbody#inventory-table-body');
    await expect(tableBody).toContainText(productName, { timeout: 10000 });

    // ---------------------------------------------------------
    // PART 2: POS Checkout Flow
    // ---------------------------------------------------------

    // 1. Navigate to POS checkout page
    await page.goto('/pos');

    // 2. Search for the product we just added
    await page.fill('input[placeholder="Search accessories..."]', productName);

    // 3. Add the product to the cart by clicking its card
    // The name is inside a h4 tag within the product card
    const productCard = page.locator(`h4:has-text("${productName}")`).locator('..');
    await expect(productCard).toBeVisible();
    await productCard.click();

    // 4. Verify the cart updates
    const cartItems = page.locator('.flex-1.overflow-y-auto.p-4.space-y-4');
    await expect(cartItems).toContainText(productName);
    
    // Quantity span (using an exact match or structural selector might be better, 
    // but checking for text '1' is simple and works for this basic test)
    await expect(cartItems.locator('span.text-center')).toHaveText('1');

    // 5. Increase quantity
    await page.locator('button[\\@click^="increaseQty"]').click();
    await expect(cartItems.locator('span.text-center')).toHaveText('2'); // Quantity should now be 2

    // Verify subtotal (2 * 250000 = 500000)
    const checkoutSummary = page.locator('.p-4.border-t.border-slate-200\\/60.bg-slate-50');
    await expect(checkoutSummary).toContainText('Rp 500.000');

    // 6. Checkout using CASH
    await page.click('button:has-text("Cash")');

    // 7. Verify the success overlay
    const successMessage = page.locator('h3:has-text("Payment Successful!")');
    await expect(successMessage).toBeVisible({ timeout: 10000 });

    // Click "New Transaction" to reset
    await page.click('button:has-text("New Transaction")');

    // Wait for reload
    await page.waitForLoadState('networkidle');

    // Verify cart is empty
    await expect(page.locator('text=Cart is empty')).toBeVisible();

    // ---------------------------------------------------------
    // PART 3: Transaction History
    // ---------------------------------------------------------

    // 1. Navigate to history page
    await page.goto('/pos/history');

    // 2. Verify the recent transaction appears (Rp 500.000, CASH)
    const historyTableBody = page.locator('tbody#history-table-body');
    const firstHistoryRow = historyTableBody.locator('tr').first();
    
    // Check if the first row contains CASH and 500.000
    await expect(firstHistoryRow).toContainText('CASH');
    await expect(firstHistoryRow).toContainText('Rp 500.000');
  });
});
