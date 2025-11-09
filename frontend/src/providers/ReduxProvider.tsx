'use client';

import React, { ReactNode } from 'react';
import { Provider } from 'react-redux';
import { store } from '@/store';
import { AuthHydrationProvider } from './AuthHydrationProvider';

interface ReduxProviderProps {
  children: ReactNode;
}

export const ReduxProvider: React.FC<ReduxProviderProps> = ({ children }) => {
  return (
    <Provider store={store}>
      <AuthHydrationProvider>{children}</AuthHydrationProvider>
    </Provider>
  );
};

export default ReduxProvider;
