import { useEffect, useRef } from 'react';
import { useAppDispatch } from '@/store';
import { setUser, setTokens } from '@/store/slices/auth/auth.slice';
import Cookies from 'js-cookie';
import { User } from '@/types';

/**
 * Hook to hydrate auth state from localStorage/cookies on app load
 * This ensures user stays logged in after page refresh
 * Returns true when hydration has been attempted
 */
export const useAuthHydration = () => {
  const dispatch = useAppDispatch();
  const hydrationAttempted = useRef(false);

  useEffect(() => {
    // Only run once per component mount
    if (hydrationAttempted.current) return;
    hydrationAttempted.current = true;

    // Check if tokens exist in cookies
    const accessToken = Cookies.get('access_token');
    const refreshToken = Cookies.get('refresh_token');
    const userJson = localStorage.getItem('user');

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
        localStorage.removeItem('user');
      }
    }
  }, [dispatch]);

  // Return true if hydration was attempted and completed
  return hydrationAttempted.current;
};
