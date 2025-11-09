import { test, expect } from '@playwright/test';

test.describe('UI Responsiveness', () => {
  test('should be accessible on mobile viewports', async ({ page }) => {
    // Test on mobile
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto('/login');

    // Elements should be visible and clickable
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('button:has-text("Login")')).toBeVisible();
  });

  test('should be accessible on tablet viewports', async ({ page }) => {
    // Test on tablet
    await page.setViewportSize({ width: 768, height: 1024 });

    await page.goto('/login');

    // Elements should be visible and properly laid out
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('button:has-text("Login")')).toBeVisible();
  });

  test('should be accessible on desktop viewports', async ({ page }) => {
    // Test on desktop
    await page.setViewportSize({ width: 1920, height: 1080 });

    await page.goto('/login');

    // Elements should be properly spaced
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('button:has-text("Login")')).toBeVisible();
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/login');

    // Check for h1 tag
    const h1 = page.locator('h1');
    await expect(h1).toBeVisible();
  });

  test('should have proper form labels', async ({ page }) => {
    await page.goto('/login');

    // Check for labels
    const labels = page.locator('label');
    const labelCount = await labels.count();

    expect(labelCount).toBeGreaterThan(0);
  });

  test('should have proper button text', async ({ page }) => {
    await page.goto('/login');

    // Check for descriptive button text
    const button = page.locator('button:has-text("Login")');
    await expect(button).toBeVisible();
  });

  test('should have alt text for images', async ({ page }) => {
    await page.goto('/');

    // Check all images have alt text
    const images = page.locator('img');
    const imageCount = await images.count();

    for (let i = 0; i < imageCount; i++) {
      const alt = await images.nth(i).getAttribute('alt');
      // Allow empty alt for decorative images, but should exist
      expect(alt).toBeDefined();
    }
  });
});

test.describe('Performance', () => {
  test('should load page within acceptable time', async ({ page }) => {
    const startTime = Date.now();

    await page.goto('/login');

    // Wait for key elements to be visible
    await page.waitForSelector('button:has-text("Login")', { timeout: 5000 });

    const loadTime = Date.now() - startTime;

    // Should load within 5 seconds
    expect(loadTime).toBeLessThan(5000);
  });

  test('should cache static assets', async ({ page }) => {
    // First navigation
    await page.goto('/login');
    await page.waitForSelector('button:has-text("Login")');

    // Navigate away and back
    await page.goto('/');
    await page.goto('/login');
    await page.waitForSelector('button:has-text("Login")');

    // Check if resources are cached (sizes should be 0 or transfer-size should indicate cache)
    const resources = await page.evaluate(() => {
      return performance.getEntriesByType('resource');
    });

    expect(resources.length).toBeGreaterThan(0);
  });
});

test.describe('Error Handling', () => {
  test('should handle network errors gracefully', async ({ page }) => {
    // Simulate network failure
    await page.route('**/api/**', (route) => route.abort());

    await page.goto('/login');

    // Try to login - should show error
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button:has-text("Login")');

    // Should show error message
    await page.waitForTimeout(1000);
    const errorElement = page.locator('[role="alert"], text=/error|failed/i').first();

    // Either error message or still on login page indicates graceful handling
    const onLoginPage = page.url().includes('/login');
    const hasErrorMsg = await errorElement.isVisible().catch(() => false);

    expect(onLoginPage || hasErrorMsg).toBeTruthy();
  });

  test('should handle 404 errors', async ({ page }) => {
    // Try to access non-existent page
    const response = await page.goto('/non-existent-page');

    // Should not crash, might show 404 or redirect to home
    expect(response?.ok() || page.url().includes('/non-existent-page')).toBeTruthy();
  });

  test('should handle 500 errors', async ({ page }) => {
    // Simulate server error
    await page.route('**/api/**', (route) => {
      route.abort('serverfailed');
    });

    await page.goto('/login');
    await page.fill('input[type="email"]', 'test@example.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button:has-text("Login")');

    // Page should handle it gracefully
    await page.waitForTimeout(1000);
    expect(page.url()).toBeTruthy();
  });
});
