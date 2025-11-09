import { useState, useCallback } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { LoginRequest } from '@/types';

export const useLogin = () => {
  const { login, isLoginLoading, error } = useAuth();
  const [formData, setFormData] = useState<LoginRequest>({
    username: '',
    password: '',
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const validateForm = useCallback((): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.username.trim()) {
      errors.username = 'Username is required';
    }

    if (!formData.password) {
      errors.password = 'Password is required';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  }, [formData]);

  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const { name, value } = e.target;
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
      // Clear error for this field when user starts typing
      if (formErrors[name]) {
        setFormErrors((prev) => ({
          ...prev,
          [name]: '',
        }));
      }
    },
    [formErrors]
  );

  const handleSubmit = useCallback(
    async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();

      if (!validateForm()) {
        return;
      }

      try {
        await login(formData);
      } catch (err) {
        // Error is handled by the hook
        console.error('Login error:', err);
      }
    },
    [formData, validateForm, login]
  );

  return {
    formData,
    formErrors,
    isLoading: isLoginLoading,
    error,
    handleInputChange,
    handleSubmit,
  };
};
