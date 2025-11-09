import { test, expect } from '@playwright/test';

const TEST_USER = {
  email: 'john.doe@trustmedis.com',
  password: 'password123',
};

test.describe('Authentication Flow', () => {
  test('should display login page', async ({ page }) => {
    await page.goto('/login');

    // Check for login form elements
    await expect(page.locator('h1')).toContainText(/login|sign in/i);
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"]')).toBeVisible();
    await expect(page.locator('button:has-text("Login")')).toBeVisible();
  });

  test('should reject invalid credentials', async ({ page }) => {
    await page.goto('/login');

    // Fill in invalid credentials
    await page.fill('input[type="email"]', 'invalid@example.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button:has-text("Login")');

    // Should show error message or stay on login page
    await page.waitForTimeout(1000);
    const url = page.url();
    expect(url).toContain('/login');
  });

  test('should successfully login with valid credentials', async ({ page }) => {
    await page.goto('/login');

    // Fill in valid credentials
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');

    // Should redirect to dashboard
    await page.waitForNavigation();
    const url = page.url();
    expect(url).toContain('/dashboard');
  });

  test('should show JWT token in localStorage after login', async ({ page }) => {
    await page.goto('/login');

    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');

    await page.waitForNavigation();

    // Check localStorage for JWT token
    const token = await page.evaluate(() => {
      return localStorage.getItem('token');
    });

    expect(token).toBeTruthy();
    expect(typeof token).toBe('string');
  });

  test('should redirect to login when accessing protected routes without token', async ({
    context,
  }) => {
    // Create new context without any stored tokens
    const newPage = await context.newPage();

    await newPage.goto('/dashboard');

    // Should redirect to login
    await newPage.waitForNavigation();
    const url = newPage.url();
    expect(url).toContain('/login');

    await newPage.close();
  });

  test('should persist login state after page reload', async ({ page }) => {
    // Login first
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();

    // Store the token
    const token = await page.evaluate(() => localStorage.getItem('token'));

    // Reload page
    await page.reload();

    // Should still be on dashboard
    const url = page.url();
    expect(url).toContain('/dashboard');

    // Token should still exist
    const newToken = await page.evaluate(() => localStorage.getItem('token'));
    expect(newToken).toBe(token);
  });
});

test.describe('Logout Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button:has-text("Login")');
    await page.waitForNavigation();
  });

  test('should logout successfully', async ({ page }) => {
    // Click logout button (usually in header/menu)
    await page.click('button:has-text("Logout")');

    // Should redirect to login
    await page.waitForNavigation();
    const url = page.url();
    expect(url).toContain('/login');

    // Token should be cleared
    const token = await page.evaluate(() => localStorage.getItem('token'));
    expect(token).toBeNull();
  });

  test('should clear user data after logout', async ({ page }) => {
    // Store some user data in localStorage
    const userDataBefore = await page.evaluate(() => {
      return JSON.parse(localStorage.getItem('user') || '{}');
    });
    expect(Object.keys(userDataBefore).length).toBeGreaterThan(0);

    // Logout
    await page.click('button:has-text("Logout")');
    await page.waitForNavigation();

    // User data should be cleared
    const userDataAfter = await page.evaluate(() => {
      return JSON.parse(localStorage.getItem('user') || '{}');
    });
    expect(Object.keys(userDataAfter).length).toBe(0);
  });
});
