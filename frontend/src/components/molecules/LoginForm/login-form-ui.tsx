import React from 'react';
import { VStack, Box, Alert, AlertIcon } from '@chakra-ui/react';
import { InputUI } from '@/components/atoms/Input/input-ui';
import { ButtonUI } from '@/components/atoms/Button/button-ui';
import { HeadingUI } from '@/components/atoms/Text/text-ui';

export interface LoginFormUIProps {
  username: string;
  password: string;
  usernameError?: string;
  passwordError?: string;
  isLoading?: boolean;
  error?: string | null;
  onUsernameChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onPasswordChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
}

/**
 * Login Form UI Component
 * Presentational component that only accepts props
 */
export const LoginFormUI: React.FC<LoginFormUIProps> = ({
  username,
  password,
  usernameError,
  passwordError,
  isLoading = false,
  error,
  onUsernameChange,
  onPasswordChange,
  onSubmit,
}) => {
  return (
    <Box as="form" onSubmit={onSubmit} w="full" maxW="md">
      <VStack spacing={6} align="stretch">
        <HeadingUI level={1} size="xl" textAlign="center">
          Login
        </HeadingUI>

        {error && (
          <Alert status="error" borderRadius="md">
            <AlertIcon />
            {error}
          </Alert>
        )}

        <InputUI
          label="Username"
          name="username"
          placeholder="Enter your username"
          value={username}
          onChange={onUsernameChange}
          error={usernameError}
          disabled={isLoading}
          autoFocus
        />

        <InputUI
          label="Password"
          name="password"
          type="password"
          placeholder="Enter your password"
          value={password}
          onChange={onPasswordChange}
          error={passwordError}
          disabled={isLoading}
        />

        <ButtonUI
          type="submit"
          colorScheme="blue"
          size="lg"
          isLoading={isLoading}
          isDisabled={isLoading}
          w="full"
        >
          {isLoading ? 'Logging in...' : 'Login'}
        </ButtonUI>
      </VStack>
    </Box>
  );
};

export default LoginFormUI;
