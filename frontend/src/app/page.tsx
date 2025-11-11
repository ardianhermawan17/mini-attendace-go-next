'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAppSelector } from '@/store';
import Cookies from 'js-cookie';

export default function Home() {
  console.log('Home component: Rendering...');
  const router = useRouter();
  const { isAuthenticated, isHydrated } = useAppSelector((state) => state.auth);

  useEffect(() => {
    console.log('Home component: useEffect triggered');
    console.log('Home component: isHydrated:', isHydrated);
    console.log('Home component: isAuthenticated:', isAuthenticated);

    // Wait for hydration to complete before making routing decisions
    if (!isHydrated) {
      console.log('Home component: Waiting for hydration to complete...');
      return;
    }

    // Check if we have valid tokens
    const accessToken = Cookies.get('access_token');
    console.log('Home component: accessToken exists:', !!accessToken);
    console.log('Home component: Combined check (isAuthenticated && accessToken):', isAuthenticated && accessToken);

    if (isAuthenticated && accessToken) {
      // User is logged in - go to dashboard
      console.log('Home component: Redirecting to /dashboard');
      router.push('/dashboard');
    } else {
      // User is not logged in - go to login
      console.log('Home component: Redirecting to /login');
      router.push('/login');
    }
  }, [isAuthenticated, isHydrated, router]);

  return null;
}
