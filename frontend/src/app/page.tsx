'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAppSelector } from '@/store';
import Cookies from 'js-cookie';

export default function Home() {
  const router = useRouter();
  const { isAuthenticated, isHydrated } = useAppSelector((state) => state.auth);

  useEffect(() => {
    // Wait for auth state to be hydrated before redirecting
    if (!isHydrated) {
      console.log('Waiting for auth hydration...');
      return;
    }

    // Check if we have valid tokens
    const accessToken = Cookies.get('access_token');

    console.log('isAuthenticated:', isAuthenticated);
    console.log('accessToken:', accessToken);
    console.log('akses page index', isAuthenticated && accessToken);

    if (isAuthenticated && accessToken) {
      // User is logged in - go to dashboard
      console.log('Redirecting to /dashboard');
      router.push('/dashboard');
    } else {
      // User is not logged in - go to login
      console.log('Redirecting to /login');
      router.push('/login');
    }
  }, [isAuthenticated, isHydrated, router]);

  return null;
}
