/**
 * Test suite for useAuthHydration hook
 * Tests auth state persistence and recovery from localStorage/cookies
 */

import { ReactNode } from 'react';
import { renderHook, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '@/store/slices/auth/auth.slice';
import { useAuthHydration } from '@/hooks/useAuthHydration';

// Mock js-cookie
jest.mock('js-cookie', () => ({
  get: jest.fn(),
  remove: jest.fn(),
  set: jest.fn(),
}));

describe('useAuthHydration Hook', () => {
  let store: any;

  beforeEach(() => {
    store = configureStore({
      reducer: {
        auth: authReducer,
      },
    });
    localStorage.clear();
    jest.clearAllMocks();
  });

  const renderHookWithProvider = (hook: any) => {
    const Wrapper = ({ children }: { children: ReactNode }) => (
      <Provider store={store}>{children}</Provider>
    );
    return renderHook(hook, { wrapper: Wrapper });
  };

  test('should return false initially, then true after hydration', async () => {
    const { result } = renderHookWithProvider(() => useAuthHydration());

    expect(result.current).toBe(false);

    await waitFor(() => {
      expect(result.current).toBe(true);
    });
  });

  test('should restore auth state when tokens exist in cookies and user in localStorage', async () => {
    const mockUser = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      role: 'user',
      department: 'IT',
    };

    // Setup mock cookies
    const Cookies = require('js-cookie');
    Cookies.get.mockImplementation((key: string) => {
      if (key === 'access_token') return 'test-access-token';
      if (key === 'refresh_token') return 'test-refresh-token';
      return undefined;
    });

    localStorage.setItem('user', JSON.stringify(mockUser));

    // Render hook
    const { result } = renderHookWithProvider(() => useAuthHydration());

    // Wait for hydration
    await waitFor(() => {
      expect(result.current).toBe(true);
    });

    // Check Redux state
    const state = store.getState();
    expect(state.auth.isAuthenticated).toBe(true);
    expect(state.auth.user).toEqual(mockUser);
    expect(state.auth.accessToken).toBe('test-access-token');
    expect(state.auth.refreshToken).toBe('test-refresh-token');
  });

  test('should NOT restore auth when only user exists but no tokens', async () => {
    const mockUser = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      role: 'user',
      department: 'IT',
    };

    // Setup - tokens missing
    const Cookies = require('js-cookie');
    Cookies.get.mockReturnValue(undefined);
    localStorage.setItem('user', JSON.stringify(mockUser));

    // Render hook
    const { result } = renderHookWithProvider(() => useAuthHydration());

    // Wait for hydration
    await waitFor(() => {
      expect(result.current).toBe(true);
    });

    // Check Redux state - should NOT be authenticated
    const state = store.getState();
    expect(state.auth.isAuthenticated).toBe(false);
    expect(state.auth.user).toBeNull();
  });

  test('should NOT restore auth when only tokens exist but no user', async () => {
    // Setup - user missing
    const Cookies = require('js-cookie');
    Cookies.get.mockImplementation((key: string) => {
      if (key === 'access_token') return 'test-access-token';
      if (key === 'refresh_token') return 'test-refresh-token';
      return undefined;
    });
    localStorage.clear();

    // Render hook
    const { result } = renderHookWithProvider(() => useAuthHydration());

    // Wait for hydration
    await waitFor(() => {
      expect(result.current).toBe(true);
    });

    // Check Redux state - should NOT be authenticated
    const state = store.getState();
    expect(state.auth.isAuthenticated).toBe(false);
    expect(state.auth.user).toBeNull();
  });

  test('should clear corrupted data from localStorage', async () => {
    // Setup - invalid JSON
    const Cookies = require('js-cookie');
    Cookies.get.mockImplementation((key: string) => {
      if (key === 'access_token') return 'test-access-token';
      if (key === 'refresh_token') return 'test-refresh-token';
      return undefined;
    });
    localStorage.setItem('user', 'invalid-json-{{{');

    // Render hook
    const { result } = renderHookWithProvider(() => useAuthHydration());

    // Wait for hydration
    await waitFor(() => {
      expect(result.current).toBe(true);
    });

    // Check that invalid data was cleared
    expect(localStorage.getItem('user')).toBeNull();
    expect(Cookies.remove).toHaveBeenCalledWith('access_token');
    expect(Cookies.remove).toHaveBeenCalledWith('refresh_token');
  });

  test('should complete hydration even when no data is found', async () => {
    // Setup - nothing in storage
    const Cookies = require('js-cookie');
    Cookies.get.mockReturnValue(undefined);
    localStorage.clear();

    // Render hook
    const { result } = renderHookWithProvider(() => useAuthHydration());

    expect(result.current).toBe(false);

    // Should still complete hydration
    await waitFor(() => {
      expect(result.current).toBe(true);
    });

    const state = store.getState();
    expect(state.auth.isAuthenticated).toBe(false);
  });
});
