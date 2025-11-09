'use client';

import React from 'react';
import { ChakraProvider } from '@chakra-ui/react';
import { CacheProvider } from '@chakra-ui/next-js';

interface ChakraProviderProps {
  children: React.ReactNode;
}

export const ChakraUIProvider: React.FC<ChakraProviderProps> = ({ children }) => {
  return (
    <CacheProvider>
      <ChakraProvider>{children}</ChakraProvider>
    </CacheProvider>
  );
};

export default ChakraUIProvider;
