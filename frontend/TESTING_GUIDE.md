# Frontend Testing & Storybook Setup

This document describes the comprehensive testing infrastructure set up for the Mini Attendance Frontend application.

## Overview

The frontend includes three levels of testing and documentation:

1. **Unit Tests** - Jest (existing setup)
2. **E2E Tests** - Playwright
3. **Component Documentation** - Storybook

## Testing Stack

### Dependencies

```json
{
  "devDependencies": {
    "@playwright/test": "^1.40.0",
    "storybook": "^7.6.0",
    "@storybook/nextjs": "^7.6.0",
    "@storybook/react": "^7.6.0",
    "@storybook/addon-essentials": "^7.6.0",
    "@storybook/addon-interactions": "^7.6.0",
    "@storybook/blocks": "^7.6.0",
    "@storybook/test": "^7.6.0"
  }
}
```

## 1. E2E Testing with Playwright

### Configuration

- **Config File**: `playwright.config.ts`
- **Test Directory**: `tests/e2e/`
- **Base URL**: `http://localhost:3001`
- **Browsers**: Chromium, Firefox, WebKit
- **Mobile**: Pixel 5, iPhone 12

### Running E2E Tests

```bash
# Run all E2E tests
npm run e2e

# Run with UI mode (interactive debugging)
npm run e2e:ui

# Run in debug mode (pause on breakpoints)
npm run e2e:debug

# Run specific test file
npx playwright test tests/e2e/auth.spec.ts

# Run with specific browser
npx playwright test --project=chromium

# Generate HTML report after test
npx playwright show-report
```

### Test Files

#### `tests/e2e/auth.spec.ts`

Tests for authentication flows:

- Login page display
- Invalid credentials rejection
- Successful login with token storage
- JWT token persistence
- Protected route access control
- Logout functionality
- User data clearing

#### `tests/e2e/attendance.spec.ts`

Tests for attendance operations:

- Dashboard display
- Check-in functionality
- Check-out functionality
- Attendance history display
- Location handling
- Duplicate check-in prevention
- Time recording accuracy
- Optimistic UI updates

#### `tests/e2e/reports.spec.ts`

Tests for reporting features:

- Reports page display
- Data table rendering
- Date range filtering
- Status filtering
- Check-in/out time display
- Export functionality
- Summary statistics
- Column sorting
- Pagination

#### `tests/e2e/general.spec.ts`

Tests for general features:

- **Responsiveness**: Mobile (375x667), Tablet (768x1024), Desktop (1920x1080)
- **Accessibility**: Headings, labels, button text, image alt text
- **Performance**: Page load time, asset caching
- **Error Handling**: Network errors, 404 errors, 500 errors

### Test Coverage

- **Auth Flow**: 7 test cases
- **Attendance Operations**: 9 test cases
- **Reports**: 8 test cases
- **General Features**: 10 test cases
- **Total**: 34+ E2E test scenarios

## 2. Storybook Component Documentation

### Configuration

- **Config Directory**: `.storybook/`
- **Main Config**: `main.ts`
- **Preview Config**: `preview.tsx`
- **Stories Directory**: `src/components/**/*.stories.tsx`
- **Port**: 6006

### Running Storybook

```bash
# Start Storybook development server
npm run storybook

# Build static Storybook site
npm run storybook:build

# View at http://localhost:6006
```

### Story Files

#### Atoms

**Button** (`src/components/atoms/Button/button-ui.stories.tsx`)

- Solid variant
- Outline variant
- Ghost variant
- Link variant
- Size variations (small, large)
- Loading state
- Disabled state
- Color schemes (success, danger, warning)

**Input** (`src/components/atoms/Input/input-ui.stories.tsx`)

- Default input
- With label
- Email input
- Password input
- Error state
- Disabled state
- Size variations
- Prefilled values
- Phone number format

**Card** (`src/components/atoms/Card/card-ui.stories.tsx`)

- Default variant
- Elevated variant
- Outlined variant
- Filled variant
- Clickable card with hover effects

**Text** (`src/components/atoms/Text/text-ui.stories.tsx`)

- Text sizes (small, medium, large)
- Text variants (label, caption)
- Heading levels (h1, h2, h3)
- All variants showcase

#### Molecules

**LoginForm** (`src/components/molecules/LoginForm/login-form-ui.stories.tsx`)

- Default state
- With server error
- With validation errors
- Loading state
- Prefilled data
- All errors combined

**CheckInCard** (`src/components/molecules/CheckInCard/check-in-card-ui.stories.tsx`)

- Not checked in state
- Checked in state
- Checked out state
- Checking in (loading)
- Checking out (loading)
- Getting location (loading)
- Location error
- Check-in error
- Check-out error

### Story Count

- **Atoms**: 4 components × 10+ variations = 40+ stories
- **Molecules**: 2 components × 10+ variations = 20+ stories
- **Total**: 60+ component variations documented

## Testing Workflow

### Pre-commit Testing

```bash
# Run unit tests
npm run test

# Check types
npm run type-check

# Lint
npm run lint

# Build
npm run build
```

### Local Development

```bash
# Terminal 1: Start dev server
npm run dev

# Terminal 2: Run Storybook
npm run storybook

# Terminal 3: Run E2E tests in UI mode
npm run e2e:ui
```

### CI/CD Integration

In your CI pipeline:

```bash
# Install dependencies
npm ci

# Type checking
npm run type-check

# Linting
npm run lint

# Unit tests
npm run test -- --coverage

# Build check
npm run build

# E2E tests (requires running backend)
npm run e2e -- --project=chromium
```

## Best Practices

### E2E Testing

1. **Selectors**: Use semantic locators:

   ```typescript
   // Good
   await page.locator('button:has-text("Login")').click();

   // Avoid
   await page.locator('.btn-primary').click();
   ```

2. **Assertions**: Use explicit assertions:

   ```typescript
   await expect(page.locator('text=Logged In')).toBeVisible();
   ```

3. **Wait Strategies**: Let Playwright auto-wait:

   ```typescript
   // Playwright auto-waits for element
   await page.locator('button').click();

   // Use explicit waits for navigation
   await page.waitForNavigation();
   ```

### Storybook Stories

1. **One concept per story**: Each story showcases a single state/variant
2. **Use descriptive names**: Stories like `WithError`, `Loading`, `Disabled`
3. **Include props**: Pass relevant props for documentation
4. **Add interactions**: Show button clicks, form input, etc.

## Troubleshooting

### E2E Tests Failing

1. **Port already in use**:

   ```bash
   # Check what's using port 3001
   netstat -ano | findstr :3001
   ```

2. **Geolocation permission denied**: Playwright mocks geolocation, should work in tests

3. **Element not found**: Check selectors with:
   ```bash
   npm run e2e -- tests/e2e/auth.spec.ts --debug
   ```

### Storybook Not Starting

1. **Clear cache**:

   ```bash
   rm -rf node_modules/.cache
   npm run storybook
   ```

2. **Port already in use**: Use `--port` flag:
   ```bash
   npm run storybook -- --port 6007
   ```

## Performance Optimization

### E2E Testing

- Tests run in parallel by default
- Use `fullyParallel: false` for sequential execution
- Screenshots and videos only on failure to save space

### Storybook

- Uses Webpack for fast development
- Addon essentials provide essential features
- Static build for production deployment

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Run E2E tests
  run: npm run e2e

- name: Upload test results
  if: always()
  uses: actions/upload-artifact@v2
  with:
    name: playwright-report
    path: playwright-report/
```

## Future Enhancements

1. **Visual Regression Testing**: Add Percy or similar
2. **Performance Monitoring**: Lighthouse CI integration
3. **Test Analytics**: Track test execution trends
4. **MSW Integration**: Mock API responses for better test isolation
5. **Component Testing**: Add more interaction-based tests

## Resources

- [Playwright Documentation](https://playwright.dev)
- [Storybook Documentation](https://storybook.js.org)
- [Testing Library Best Practices](https://testing-library.com/docs/queries/about)
- [Atomic Design Methodology](https://bradfrost.com/blog/post/atomic-web-design/)
