import { test, expect } from '@playwright/test';

// Serial: each test depends on a fresh warrant created in beforeEach.
// Running in parallel would cause ticketNumber to be shared across workers.
test.describe.serial('Admin Warranty Flow', () => {
  let ticketNumber = '';

  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('admin@openbench.local').fill('admin@openbench.com');
    await page.getByPlaceholder('••••••••••••').fill('secretpassword123');
    await page.getByRole('button', { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/(tickets|dashboard)/);

    // Create a ticket via API
    const createRes = await page.request.post('/api/v1/admin/services', {
      data: {
        customer_name: 'E2E Warranty Test',
        customer_phone: '081234567890',
        device_brand: 'Samsung',
        device_model: 'Galaxy A55',
        issue_description: 'Battery swelling',
        repair_action: 'Replace battery',
        cost: 350000,
        warranty_days: 30,
      },
    });
    expect(createRes.ok()).toBeTruthy();
    const ticket = (await createRes.json()).data;
    const ticketId = ticket.ticket_id;
    ticketNumber = ticket.ticket_number;

    // Advance status to COMPLETED (triggers warranty creation via event bus)
    for (const status of ['REPAIRING', 'FIXED', 'COMPLETED']) {
      const res = await page.request.patch(`/api/v1/admin/services/${ticketId}/status`, {
        data: { status },
      });
      expect(res.ok()).toBeTruthy();
    }

    // Poll until warranty is available (async event bus may have slight delay)
    await expect(async () => {
      const res = await page.request.get(`/api/v1/admin/warranties/by-ticket-number/${ticketNumber}`);
      expect(res.ok()).toBeTruthy();
      const body = await res.json();
      expect(body.data).toBeTruthy();
      expect(body.data.status).toBe('ACTIVE');
    }).toPass({ timeout: 10000 });
  });

  async function searchWarranty(page: import('@playwright/test').Page, num: string) {
    await page.goto('/warranties');
    await expect(page.getByRole('heading', { name: 'Warranty & Claim Ticketing' })).toBeVisible();

    // Fill input and submit form via button click (more reliable than Enter key)
    await page.getByPlaceholder('e.g. TKT-20260707-1234').fill(num);
    await page.getByRole('button', { name: 'Verify Warranty' }).click();
    await expect(page.getByText('Warranty Contract Details')).toBeVisible({ timeout: 15000 });
  }

  test('should navigate to warranty page and show initial state', async ({ page }) => {
    await page.goto('/warranties');
    await expect(page.getByRole('heading', { name: 'Warranty & Claim Ticketing' })).toBeVisible();
    await expect(page.getByText(/Awaiting Warranty Check/)).toBeVisible();
  });

  test('should search warranty and verify UI after search', async ({ page }) => {
    await searchWarranty(page, ticketNumber);
    await expect(page.getByText('ACTIVE', { exact: true })).toBeVisible();
  });

  test('should submit a claim and evaluate it via UI', async ({ page }) => {
    await searchWarranty(page, ticketNumber);

    // Submit claim
    await page.getByRole('button', { name: 'Submit Warranty Claim' }).click();
    await expect(page.getByRole('heading', { name: 'Create Warranty Claim' })).toBeVisible();
    await page.getByPlaceholder(/Describe the repeating/i).fill('Screen flickering');
    await page.getByRole('button', { name: 'Register Claim' }).click();

    await expect(page.getByRole('cell', { name: 'PENDING' })).toBeVisible({ timeout: 5000 });

    // Evaluate from Claims Queue tab
    await page.getByRole('button', { name: 'Evaluate' }).click();
    await expect(page.getByRole('heading', { name: 'Evaluate Warranty Claim' })).toBeVisible();
    await page.getByRole('button', { name: 'Submit Decision' }).click();
    await expect(page.getByRole('cell', { name: 'ACCEPTED' })).toBeVisible({ timeout: 5000 });
  });

  test('should void a warranty via UI', async ({ page }) => {
    await searchWarranty(page, ticketNumber);

    await page.getByRole('button', { name: 'Void Warranty' }).click();
    await expect(page.getByRole('heading', { name: 'Void Warranty Contract' })).toBeVisible();
    await page.getByPlaceholder(/Broken warranty seal/i).fill('Customer tampered');
    await page.getByRole('button', { name: 'Invalidate (Void)' }).click();

    await expect(page.getByText('VOIDED', { exact: true })).toBeVisible({ timeout: 10000 });
    await expect(page.getByText(/This warranty is void/)).toBeVisible();
  });
});
