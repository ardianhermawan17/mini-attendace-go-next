import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store';
import { setUser, setTokens, logout, setError } from '@/store/slices/auth/auth.slice';
import { useLoginMutation, useLogoutMutation } from '@/store/api/authApi';
import { LoginRequest } from '@/types';

export const useAuth = () => {
  const dispatch = useAppDispatch();
  const [loginMutation, { isLoading: isLoginLoading }] = useLoginMutation();
  const [logoutMutation, { isLoading: isLogoutLoading }] = useLogoutMutation();

  const auth = useAppSelector((state) => state.auth);

  const login = useCallback(
    async (credentials: LoginRequest) => {
      try {
        const response = await loginMutation(credentials).unwrap();
        dispatch(setUser(response.user));
        dispatch(
          setTokens({
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
            expiresIn: response.expires_in,
          })
        );
        return response;
      } catch (error: any) {
        const errorMessage = error?.data?.error?.message || 'Login failed';
        dispatch(setError(errorMessage));
        throw error;
      }
    },
    [loginMutation, dispatch]
  );

  const handleLogout = useCallback(async () => {
    try {
      await logoutMutation().unwrap();
      dispatch(logout());
    } catch (error: any) {
      const errorMessage = error?.data?.error?.message || 'Logout failed';
      dispatch(setError(errorMessage));
      throw error;
    }
  }, [logoutMutation, dispatch]);

  return {
    ...auth,
    login,
    logout: handleLogout,
    isLoginLoading,
    isLogoutLoading,
  };
};
