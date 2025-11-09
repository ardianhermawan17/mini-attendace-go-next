import React from 'react';
import { LoginFormUI } from '@/components/molecules/LoginForm/login-form-ui';
import { useLogin } from './use-login';

/**
 * Login Form Container
 * Combines the UI component with the logic hook
 */
export const LoginForm: React.FC = () => {
  const { formData, formErrors, isLoading, error, handleInputChange, handleSubmit } =
    useLogin();

  return (
    <LoginFormUI
      username={formData.username}
      password={formData.password}
      usernameError={formErrors.username}
      passwordError={formErrors.password}
      isLoading={isLoading}
      error={error}
      onUsernameChange={(e) =>
        handleInputChange({ ...e, target: { ...e.target, name: 'username' } })
      }
      onPasswordChange={(e) =>
        handleInputChange({ ...e, target: { ...e.target, name: 'password' } })
      }
      onSubmit={handleSubmit}
    />
  );
};

export default LoginForm;
