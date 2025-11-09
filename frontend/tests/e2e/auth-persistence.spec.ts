import { test, expect } from '@playwright/test';

/**
 * E2E Tests for Auth State Persistence
 * Tests that authentication state is properly persisted and restored
 */

test('should redirect to login when accessing home page without authentication', async ({
  page,
}) => {
  // Clear any existing data
  await page.context().clearCookies();
  await page.evaluate(() => localStorage.clear());

  // Navigate to home
  await page.goto('http://localhost:3001/');

  // Should be redirected to login
  await expect(page).toHaveURL('**/login');
  await expect(page.locator('input[name="email"]')).toBeVisible();
});

test('should persist auth state and redirect to dashboard after login', async ({ page }) => {
  // Clear any existing data
  await page.context().clearCookies();
  await page.evaluate(() => localStorage.clear());

  // Navigate to login
  await page.goto('http://localhost:3001/login');
  await expect(page.locator('input[name="email"]')).toBeVisible();

  // Login with test credentials
  await page.fill('input[name="email"]', 'admin@trustmedis.com');
  await page.fill('input[name="password"]', 'admin123');
  await page.click('button:has-text("Login")');

  // Wait for redirect to dashboard
  await expect(page).toHaveURL('**/dashboard', { timeout: 10000 });
  await expect(page.locator('text=Dashboard')).toBeVisible({ timeout: 5000 });

  // Verify tokens are in cookies
  const cookies = await page.context().cookies();
  const accessTokenCookie = cookies.find((c) => c.name === 'access_token');
  expect(accessTokenCookie).toBeDefined();
  expect(accessTokenCookie?.value).toBeTruthy();

  // Verify user data is in localStorage
  const userDataStr = await page.evaluate(() => localStorage.getItem('user'));
  expect(userDataStr).toBeTruthy();
  const userData = JSON.parse(userDataStr!);
  expect(userData.email).toBe('admin@trustmedis.com');
});

test('should stay logged in after page refresh', async ({ page }) => {
  // First, login
  await page.context().clearCookies();
  await page.evaluate(() => localStorage.clear());

  await page.goto('http://localhost:3001/login');
  await page.fill('input[name="email"]', 'admin@trustmedis.com');
  await page.fill('input[name="password"]', 'admin123');
  await page.click('button:has-text("Login")');

  // Wait for dashboard
  await expect(page).toHaveURL('**/dashboard', { timeout: 10000 });
  await expect(page.locator('text=Dashboard')).toBeVisible();

  // Now refresh the page
  await page.reload({ waitUntil: 'networkidle' });

  // Should still be on dashboard (not redirected to login)
  await expect(page).toHaveURL('**/dashboard');
  await expect(page.locator('text=Dashboard')).toBeVisible();
});
