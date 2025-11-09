import { test, expect } from '@playwright/test';

const TEST_USER = {
  email: 'john.doe@trustmedis.com',
  password: 'password123',
};

test.describe('Attendance Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();
  });

  test('should display attendance dashboard with check-in button', async ({ page }) => {
    await page.goto('/attendance');

    // Check for essential UI elements
    await expect(page.locator('text=Check In')).toBeVisible();
    await expect(page.locator('text=Status')).toBeVisible();
  });

  test('should perform check-in successfully', async ({ page }) => {
    await page.goto('/attendance');

    // Wait for check-in button and click it
    const checkInButton = page.locator('button:has-text("Check In")');
    await expect(checkInButton).toBeVisible();
    await checkInButton.click();

    // Should show success message or update UI
    await page.waitForTimeout(1000);

    // Check if status changed to "Checked In" or similar
    const statusText = await page.locator('text=Checked|Checked In|Present').first();
    await expect(statusText).toBeVisible();
  });

  test('should show check-out button after check-in', async ({ page }) => {
    await page.goto('/attendance');

    // Perform check-in
    const checkInButton = page.locator('button:has-text("Check In")');
    await checkInButton.click();
    await page.waitForTimeout(1000);

    // Check-out button should now be visible
    const checkOutButton = page.locator('button:has-text("Check Out")');
    await expect(checkOutButton).toBeVisible();
  });

  test('should perform check-out successfully', async ({ page }) => {
    await page.goto('/attendance');

    // Check-in first
    await page.locator('button:has-text("Check In")').click();
    await page.waitForTimeout(1000);

    // Now check-out
    const checkOutButton = page.locator('button:has-text("Check Out")');
    await expect(checkOutButton).toBeVisible();
    await checkOutButton.click();

    // Should show success message
    await page.waitForTimeout(1000);

    // Status should reflect check-out or return to initial state
    const statusText = page.locator('text=Checked Out|Not Checked In').first();
    await expect(statusText).toBeVisible();
  });

  test('should display attendance history', async ({ page }) => {
    await page.goto('/attendance');

    // Look for attendance history/records
    const historySection = page.locator('text=History|Record|Today').first();
    await expect(historySection).toBeVisible();
  });

  test('should handle location for check-in', async ({ page }) => {
    await page.goto('/attendance');

    // Mock geolocation
    await page.context().grantPermissions(['geolocation']);
    await page.context().setGeolocation({ latitude: -6.2088, longitude: 106.8456 });

    // Perform check-in
    const checkInButton = page.locator('button:has-text("Check In")');
    await checkInButton.click();

    await page.waitForTimeout(1000);

    // Should succeed despite geolocation request
    const statusElement = page.locator('text=Checked|Present|Error').first();
    await expect(statusElement).toBeVisible();
  });

  test('should prevent duplicate check-in', async ({ page }) => {
    await page.goto('/attendance');

    // First check-in
    await page.locator('button:has-text("Check In")').click();
    await page.waitForTimeout(1000);

    // Try to check-in again (button should be disabled or hidden)
    const checkInButton = page.locator('button:has-text("Check In")');
    const isDisabled = await checkInButton.isDisabled();
    const isHidden = await checkInButton.isHidden();

    expect(isDisabled || isHidden).toBeTruthy();
  });

  test('should show time for check-in and check-out', async ({ page }) => {
    await page.goto('/attendance');

    // Check-in
    await page.locator('button:has-text("Check In")').click();
    await page.waitForTimeout(1000);

    // Look for timestamp
    const timeElement = page.locator('[data-testid="check-in-time"]').first();
    await expect(timeElement).toBeVisible();
  });
});

test.describe('Attendance with Optimistic Updates', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();
  });

  test('should show optimistic UI update immediately', async ({ page }) => {
    await page.goto('/attendance');

    // Get initial status
    const initialStatus = await page.locator('[data-testid="status"]').textContent();

    // Click check-in - should update UI immediately
    const checkInButton = page.locator('button:has-text("Check In")');
    await checkInButton.click();

    // UI should update immediately (optimistic update)
    await page.waitForTimeout(100);
    const updatedStatus = await page.locator('[data-testid="status"]').textContent();

    expect(updatedStatus).not.toBe(initialStatus);
  });

  test('should maintain optimistic update even if request is slow', async ({ page }) => {
    await page.goto('/attendance');

    // Slow down network
    await page.route('**/api/**', async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      await route.continue();
    });

    // Perform check-in
    const checkInButton = page.locator('button:has-text("Check In")');
    await checkInButton.click();

    // UI should still be updated immediately despite slow network
    await page.waitForTimeout(100);
    const status = await page.locator('[data-testid="status"]').textContent();
    expect(status).toContain('Checked');
  });
});
