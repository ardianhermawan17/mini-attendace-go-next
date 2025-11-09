'use client';

import { ReactNode, useEffect } from 'react';
import { useAppDispatch } from '@/store';
import { setUser, setTokens } from '@/store/slices/auth/auth.slice';
import Cookies from 'js-cookie';
import { User } from '@/types';

interface AuthHydrationProviderProps {
  children: ReactNode;
}

/**
 * Provider component that hydrates auth state on app initialization
 * This runs once when the app loads to restore user from localStorage/cookies
 */
export const AuthHydrationProvider: React.FC<AuthHydrationProviderProps> = ({ children }) => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    // Hydrate auth state from persisted storage on app initialization
    const accessToken = Cookies.get('access_token');
    const refreshToken = Cookies.get('refresh_token');
    const userJson = typeof window !== 'undefined' ? localStorage.getItem('user') : null;

    if (accessToken && userJson) {
      try {
        const user = JSON.parse(userJson) as User;

        // Restore auth state from persisted data
        dispatch(setUser(user));
        dispatch(
          setTokens({
            accessToken,
            refreshToken: refreshToken || '',
            expiresIn: 3600, // Default 1 hour
          })
        );
      } catch (error) {
        console.error('Failed to hydrate auth state:', error);
        // Clear invalid data
        Cookies.remove('access_token');
        Cookies.remove('refresh_token');
        if (typeof window !== 'undefined') {
          localStorage.removeItem('user');
        }
      }
    }
  }, [dispatch]);

  return <>{children}</>;
};
