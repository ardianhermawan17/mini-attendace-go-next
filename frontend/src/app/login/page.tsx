'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Center, Box } from '@chakra-ui/react';
import { LoginForm } from '@/features/auth';
import { useAppSelector } from '@/store';

export default function LoginPage() {
  const router = useRouter();
  const isAuthenticated = useAppSelector((state) => state.auth.isAuthenticated);

  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isAuthenticated, router]);

  return (
    <Center minH="100vh" bg="gray.50">
      <Box w="full" px={4}>
        <LoginForm />
      </Box>
    </Center>
  );
}
