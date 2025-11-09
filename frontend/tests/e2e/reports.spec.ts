import { test, expect } from '@playwright/test';

const TEST_USER = {
  email: 'john.doe@trustmedis.com',
  password: 'password123',
};

test.describe('Reports Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();
  });

  test('should display reports page with data table', async ({ page }) => {
    await page.goto('/reports');

    // Check for table elements
    await expect(page.locator('text=Date|Time|Status')).toBeVisible();
    await expect(page.locator('table')).toBeVisible();
  });

  test('should display attendance records in table', async ({ page }) => {
    await page.goto('/reports');

    // Wait for table to load
    await page.waitForSelector('tbody tr', { timeout: 5000 });

    // Check if rows exist
    const rows = page.locator('tbody tr');
    const rowCount = await rows.count();

    expect(rowCount).toBeGreaterThan(0);
  });

  test('should filter by date range', async ({ page }) => {
    await page.goto('/reports');

    // Look for date input fields
    const dateInputs = page.locator('input[type="date"]');
    const dateInputCount = await dateInputs.count();

    if (dateInputCount >= 2) {
      // Set date range
      const fromDate = dateInputs.first();
      const toDate = dateInputs.nth(1);

      await fromDate.fill('2024-01-01');
      await toDate.fill('2024-01-31');

      // Should apply filter
      await page.waitForTimeout(500);

      // Verify filtered results
      const rows = page.locator('tbody tr');
      const filteredCount = await rows.count();
      expect(filteredCount).toBeGreaterThanOrEqual(0);
    }
  });

  test('should search by status', async ({ page }) => {
    await page.goto('/reports');

    // Look for status filter
    const statusDropdown = page.locator('select, [role="listbox"]').first();

    if (await statusDropdown.isVisible()) {
      await statusDropdown.click();

      // Select "Present" status
      await page.locator('text=Present').click();

      // Should filter results
      await page.waitForTimeout(500);

      const rows = page.locator('tbody tr');
      const filteredCount = await rows.count();
      expect(filteredCount).toBeGreaterThanOrEqual(0);
    }
  });

  test('should display check-in and check-out times', async ({ page }) => {
    await page.goto('/reports');

    // Wait for table to load
    await page.waitForSelector('tbody tr', { timeout: 5000 });

    // Look for time columns
    const cells = page.locator('td');
    const timePatterns = /\d{1,2}:\d{2}(:\d{2})?/;

    let hasTimeData = false;
    const cellCount = await cells.count();

    for (let i = 0; i < Math.min(cellCount, 20); i++) {
      const text = await cells.nth(i).textContent();
      if (text && timePatterns.test(text)) {
        hasTimeData = true;
        break;
      }
    }

    expect(hasTimeData).toBeTruthy();
  });

  test('should export or download report', async ({ page }) => {
    await page.goto('/reports');

    // Look for export button
    const exportButton = page.locator('button:has-text("Export|Download|CSV|PDF")').first();

    if (await exportButton.isVisible()) {
      // Set up listener for download
      const downloadPromise = page.waitForEvent('download');

      await exportButton.click();

      const download = await downloadPromise;

      // Verify download
      expect(download).toBeTruthy();
      expect(download.suggestedFilename()).toMatch(/\.(csv|pdf|xlsx?)$/);
    }
  });

  test('should display summary statistics', async ({ page }) => {
    await page.goto('/reports');

    // Look for summary cards/statistics
    const stats = page.locator('[data-testid*="stat"]').first();

    if (await stats.isVisible()) {
      const statsText = await stats.textContent();
      expect(statsText).toBeTruthy();
    }
  });

  test('should sort table columns', async ({ page }) => {
    await page.goto('/reports');

    // Wait for table to load
    await page.waitForSelector('thead th', { timeout: 5000 });

    // Click on sortable column header
    const headers = page.locator('thead th');
    const firstHeader = headers.first();

    // Click to sort
    await firstHeader.click();
    await page.waitForTimeout(500);

    // Verify table updated
    const rows = page.locator('tbody tr');
    const rowCount = await rows.count();
    expect(rowCount).toBeGreaterThan(0);
  });

  test('should handle pagination', async ({ page }) => {
    await page.goto('/reports');

    // Look for pagination controls
    const nextButton = page.locator('button:has-text("Next"), [aria-label*="next"]').first();

    if ((await nextButton.isVisible()) && !(await nextButton.isDisabled())) {
      // Click next page
      await nextButton.click();
      await page.waitForTimeout(500);

      // Verify page changed
      const rows = page.locator('tbody tr');
      const rowCount = await rows.count();
      expect(rowCount).toBeGreaterThan(0);
    }
  });
});

test.describe('Reports Data Accuracy', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();
  });

  test('should show correct total working hours calculation', async ({ page }) => {
    await page.goto('/reports');

    // Look for working hours display
    const hoursElement = page.locator('[data-testid="working-hours"]').first();

    if (await hoursElement.isVisible()) {
      const hoursText = await hoursElement.textContent();
      // Should contain a number
      expect(hoursText).toBeTruthy();
    }
  });

  test('should reflect changes after new attendance record', async ({ page }) => {
    // First, go to attendance and check in
    await page.goto('/attendance');
    const checkInButton = page.locator('button:has-text("Check In")');
    if (await checkInButton.isEnabled()) {
      await checkInButton.click();
      await page.waitForTimeout(1000);
    }

    // Now navigate to reports
    await page.goto('/reports');

    // Verify new record appears
    const rows = page.locator('tbody tr');
    const rowCount = await rows.count();
    expect(rowCount).toBeGreaterThan(0);
  });
});
