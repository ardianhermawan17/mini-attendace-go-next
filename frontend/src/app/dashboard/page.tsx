'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Container, VStack, HStack, Box, Heading, Text } from '@chakra-ui/react';
import { CheckInCard } from '@/features/attendance';
import { useAppSelector } from '@/store';
import { format } from 'date-fns';

export default function DashboardPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAppSelector((state) => state.auth);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) {
    return null;
  }

  return (
    <Container maxW="container.lg" py={8}>
      <VStack spacing={8} align="stretch">
        {/* Header */}
        <Box>
          <Heading as="h1" size="2xl" mb={2}>
            Welcome, {user?.full_name}
          </Heading>
          <Text color="gray.600">
            {format(new Date(), 'EEEE, MMMM d, yyyy')}
          </Text>
        </Box>

        {/* Main Content */}
        <HStack spacing={8} align="start">
          <Box flex={1}>
            <CheckInCard />
          </Box>

          {/* Sidebar Info */}
          <Box w="300px">
            <Box bg="white" p={6} borderRadius="lg" boxShadow="md">
              <Heading as="h3" size="md" mb={4}>
                Your Information
              </Heading>
              <VStack align="start" spacing={3}>
                <Box>
                  <Text fontSize="sm" color="gray.600">
                    Department
                  </Text>
                  <Text fontWeight="medium">{user?.department}</Text>
                </Box>
                <Box>
                  <Text fontSize="sm" color="gray.600">
                    Role
                  </Text>
                  <Text fontWeight="medium" textTransform="capitalize">
                    {user?.role}
                  </Text>
                </Box>
                <Box>
                  <Text fontSize="sm" color="gray.600">
                    Email
                  </Text>
                  <Text fontWeight="medium">{user?.email}</Text>
                </Box>
              </VStack>
            </Box>
          </Box>
        </HStack>
      </VStack>
    </Container>
  );
}
