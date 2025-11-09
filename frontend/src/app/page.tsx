'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAppSelector } from '@/store';
import Cookies from 'js-cookie';

export default function Home() {
  const router = useRouter();
  const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);

  useEffect(() => {
    // Check if we have valid tokens
    const accessToken = Cookies.get('access_token');

    if (isAuthenticated && accessToken) {
      // User is logged in - go to dashboard
      router.push('/dashboard');
    } else {
      // User is not logged in - go to login
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  return null;
}
