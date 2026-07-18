import { test, expect } from '@playwright/test';

test.describe.serial('Service Tickets Flow', () => {
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

  test('Create and update a service ticket', async ({ page }) => {
    // ---------------------------------------------------------
    // PART 1: Create a new service ticket
    // ---------------------------------------------------------
    
    // 1. Navigate to tickets page
    await page.goto('/tickets');
    
    // Check if the page loaded correctly
    await expect(page.locator('h1')).toContainText('Service Tickets');

    // 2. Open the "New Ticket" drawer
    await page.click('button:has-text("New Ticket")');

    // Wait for the drawer title to be visible to ensure HTMX loaded the form
    const drawerTitle = page.locator('h2:has-text("Create New Service Ticket")');
    await expect(drawerTitle).toBeVisible();

    // 3. Fill out the form
    // Customer Info
    await page.fill('input[name="customer_name"]', 'John Doe Test');
    await page.fill('input[name="customer_phone"]', '08123456789');
    
    // Device Info
    await page.fill('input[name="device_brand"]', 'Samsung');
    await page.fill('input[name="device_model"]', 'Galaxy S23');
    await page.fill('input[name="device_passcode"]', '1234');
    
    // Issue Info
    await page.fill('textarea[name="issue_description"]', 'Screen cracked from drop');
    await page.fill('input[name="repair_action"]', 'Replace LCD Display');
    
    // Cost & Warranty
    await page.fill('input[name="cost"]', '1500000');
    await page.fill('input[name="warranty_days"]', '30');

    // 4. Submit the form
    await page.click('button[type="submit"]:has-text("Save Ticket")');

    // 5. Verify the new ticket appears in the table
    // Playwright will automatically poll until the condition is met or it times out
    const tableBody = page.locator('tbody#tickets-table-body');
    await expect(tableBody).toContainText('John Doe Test', { timeout: 10000 });
    await expect(tableBody).toContainText('Galaxy S23');
    await expect(tableBody).toContainText('Samsung');
    
    // Verify the drawer is closed
    await expect(drawerTitle).not.toBeVisible();

    // ---------------------------------------------------------
    // PART 2: Update the created service ticket
    // ---------------------------------------------------------

    // We use the button in the table row that has our specific customer name
    const testRow = page.locator('tbody#tickets-table-body tr', { hasText: 'John Doe Test' }).first();
    const viewButton = testRow.locator('button:has-text("View")');
    await viewButton.click();

    // 3. Wait for the update drawer to open
    const updateForm = page.locator('form#update-ticket-form');
    await expect(updateForm).toBeVisible();

    // 4. Modify the status and cost
    // Select "REPAIRING" from the status dropdown
    await page.selectOption('select[name="status"]', 'REPAIRING');
    
    // Update the cost
    await page.fill('input[name="cost"]', '2000000');
    
    // Update internal notes
    await page.fill('textarea[name="notes"]', 'Screen replaced, testing functionality.');

    // 5. Submit the form
    await page.click('button[type="submit"]:has-text("Save Changes")');

    // 6. Verify the drawer closes and the table updates
    await expect(updateForm).not.toBeVisible();

    // The row should now display the updated status "REPAIRING"
    const updatedRow = page.locator('tbody#tickets-table-body tr', { hasText: 'John Doe Test' }).first();
    await expect(updatedRow).toContainText('REPAIRING', { timeout: 10000 });
  });
});
