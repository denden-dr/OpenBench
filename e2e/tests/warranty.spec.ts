import { test, expect } from '@playwright/test';

test.describe.serial('Warranty Claims Flow', () => {
  let ticketNumber = '';

  test.beforeEach(async ({ page, request }) => {
    // 1. Login as admin
    await page.goto('/login');
    await page.fill('input[name="email"]', 'admin@openbench.com');
    await page.fill('input[name="password"]', 'secretpassword123');
    await Promise.all([
      page.waitForURL((url) => url.pathname === '/', { timeout: 10000 }),
      page.click('button[type="submit"]'),
    ]);

    // 2. Setup Data: Create a Ticket and Complete it to generate a Warranty
    // We use the UI to get the cookies/session, but Playwright's `request` context
    // might not share the browser's cookie automatically if not configured.
    // However, since we need to seed data, we can just grab the token if it's in a cookie,
    // or we can just use the UI to create the ticket.

    // Let's use the UI to create a ticket and complete it to ensure a realistic flow.
    await page.goto('/tickets');
    await page.click('button:has-text("New Ticket")');
    await expect(page.locator('h2:has-text("Create New Service Ticket")')).toBeVisible();

    await page.fill('input[name="customer_name"]', 'Warranty Test User');
    await page.fill('input[name="customer_phone"]', '08123456789');
    await page.fill('input[name="device_brand"]', 'Apple');
    await page.fill('input[name="device_model"]', 'iPhone 13');
    await page.fill('input[name="device_passcode"]', '0000');
    await page.fill('textarea[name="issue_description"]', 'Screen replacement');
    await page.fill('input[name="repair_action"]', 'Replace OLED');
    await page.fill('input[name="cost"]', '2000000');
    await page.fill('input[name="warranty_days"]', '30');

    await page.click('button[type="submit"]:has-text("Save Ticket")');

    // Wait for it to appear in the table
    const tableBody = page.locator('tbody#tickets-table-body');
    await expect(tableBody).toContainText('Warranty Test User', { timeout: 10000 });

    // Extract the ticket number from the first row's specific column
    // The ticket number is usually in the first column (or part of the text).
    // Let's open the View drawer to get the exact ticket number.
    await page.locator('tbody#tickets-table-body tr').first().locator('button:has-text("View")').click();
    await expect(page.locator('form#update-ticket-form')).toBeVisible();
    
    // The ticket number is in the drawer title or readonly field.
    // In OpenBench, it's often in a readonly input or heading.
    // We can also extract it from the API response by intercepting the form submission.
  });

  test('Submit and evaluate a warranty claim', async ({ page }) => {
    // ---------------------------------------------------------
    // PART 1: Find the created ticket and mark as COMPLETED
    // ---------------------------------------------------------
    await page.goto('/tickets');
    
    const testRow = page.locator('tbody#tickets-table-body tr', { hasText: 'Warranty Test User' }).first();
    await testRow.locator('button:has-text("View")').click();
    await expect(page.locator('form#update-ticket-form')).toBeVisible();

    // Extract ticket number from the drawer title
    const drawerTitle = await page.locator('h2#slide-over-title').textContent();
    ticketNumber = drawerTitle?.trim() || '';
    
    // Change status to COMPLETED
    await page.selectOption('select[name="status"]', 'COMPLETED');
    await page.click('button:has-text("Save Changes")');

    // Ensure drawer closes
    await expect(page.locator('form#update-ticket-form')).not.toBeVisible();

    // ---------------------------------------------------------
    // PART 2: Submit a new warranty claim
    // ---------------------------------------------------------
    await page.goto('/warranties');
    await expect(page.locator('h1')).toContainText('Warranty Claims');

    // Test New Claim URL change and revert on Cancel
    await page.click('button:has-text("New Claim")');
    await expect(page.locator('h2:has-text("Submit New Claim")')).toBeVisible();
    await expect(page).toHaveURL(/\/warranties\/claims\/new$/);

    await page.locator('#slideover-container button:has-text("Cancel")').click();
    await expect(page.locator('h2:has-text("Submit New Claim")')).not.toBeVisible();
    await expect(page).toHaveURL(/\/warranties$/);

    // Open again to proceed with submission
    await page.click('button:has-text("New Claim")');
    await expect(page.locator('h2:has-text("Submit New Claim")')).toBeVisible();
    await expect(page).toHaveURL(/\/warranties\/claims\/new$/);

    await page.fill('#ticket_number', ticketNumber, { timeout: 5000 });
    await page.click('button:has-text("Verify Warranty")');

    // Step 2: Wait for Warranty Verified message and reason textarea
    await expect(page.locator('text=Warranty Verified')).toBeVisible({ timeout: 5000 });
    await page.fill('#reason', 'Screen flickers after 2 days');
    
    await page.click('button[type="submit"]:has-text("Submit Claim")');

    try {
      // Wait for submission and table update
      await expect(page.locator('h2:has-text("Submit New Claim")')).not.toBeVisible();
    } catch (e) {
      const html = await page.content();
      require('fs').writeFileSync('out_submit.html', html);
      throw e;
    }
    const claimTable = page.locator('tbody#warranties-table-body');
    await expect(claimTable).toContainText(ticketNumber, { timeout: 10000 });
    await expect(claimTable).toContainText('PENDING');

    // ---------------------------------------------------------
    // PART 3: Evaluate the warranty claim
    // ---------------------------------------------------------
    // Click View on the newly created claim
    const claimRow = page.locator('tbody#warranties-table-body tr', { hasText: ticketNumber }).first();
    await claimRow.locator('button:has-text("Verify")').click();
    
    const evalDrawerTitle = page.locator('h2:has-text("Claim:")');
    try {
      await expect(evalDrawerTitle).toBeVisible({ timeout: 5000 });
    } catch (e) {
      const html = await page.content();
      require('fs').writeFileSync('out_eval.html', html);
      throw e;
    }

    // Fill inspection notes
    await page.fill('textarea[name="notes"]', 'Verified flickering issue. Approve for rework.');
    
    // Click Approve
    await page.click('button:has-text("Approve & Create Ticket")');

    // Wait for drawer to close
    await expect(evalDrawerTitle).not.toBeVisible();

    // Verify status updated to ACCEPTED in the table
    const updatedClaimRow = page.locator('tbody#warranties-table-body tr', { hasText: ticketNumber }).first();
    await expect(updatedClaimRow).toContainText('ACCEPTED', { timeout: 10000 });

    // ---------------------------------------------------------
    // PART 4: Test UI Filters (Pending / Approved)
    // ---------------------------------------------------------
    // Verify "All" is active initially (bg-primary class)
    const allBtn = page.locator('button:has-text("All")');
    await expect(allBtn).toHaveClass(/bg-primary/);

    // Test Pending filter
    const pendingBtn = page.locator('button:has-text("Pending")');
    await pendingBtn.click();
    await expect(page).toHaveURL(/\?status=PENDING$/);
    await expect(pendingBtn).toHaveClass(/bg-primary/);
    await expect(allBtn).not.toHaveClass(/bg-primary/);

    // The ticket should disappear since it's ACCEPTED
    await expect(page.locator('tbody#warranties-table-body')).not.toContainText(ticketNumber, { timeout: 10000 });
    // Verify no nested page bug (h1 would be injected if whole page was returned)
    await expect(page.locator('tbody#warranties-table-body h1')).not.toBeVisible();

    // Test Approved filter
    const approvedBtn = page.locator('button:has-text("Approved")');
    await approvedBtn.click();
    await expect(page).toHaveURL(/\?status=ACCEPTED$/);
    await expect(approvedBtn).toHaveClass(/bg-primary/);
    await expect(pendingBtn).not.toHaveClass(/bg-primary/);

    // The ticket should reappear
    await expect(page.locator('tbody#warranties-table-body')).toContainText(ticketNumber, { timeout: 10000 });
    await expect(page.locator('tbody#warranties-table-body h1')).not.toBeVisible();

    // Test All filter again to revert URL
    await allBtn.click();
    await expect(page).toHaveURL(/\/warranties$/);
    await expect(allBtn).toHaveClass(/bg-primary/);
  });
});
