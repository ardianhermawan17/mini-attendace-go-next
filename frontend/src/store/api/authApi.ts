import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { AuthResponse, LoginRequest, User } from '@/types';
import Cookies from 'js-cookie';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const authApi = createApi({
  reducerPath: 'authApi',
  baseQuery: fetchBaseQuery({
    baseUrl: API_URL,
    prepareHeaders: (headers) => {
      const token = Cookies.get('access_token');
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  endpoints: (builder) => ({
    login: builder.mutation<AuthResponse, LoginRequest>({
      query: (credentials) => ({
        url: '/auth/login',
        method: 'POST',
        body: credentials,
      }),
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          Cookies.set('access_token', data.access_token);
          Cookies.set('refresh_token', data.refresh_token);
        } catch (error) {
          console.error('Login failed:', error);
        }
      },
    }),

    logout: builder.mutation<{ message: string }, void>({
      query: () => ({
        url: '/auth/logout',
        method: 'POST',
      }),
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          await queryFulfilled;
          Cookies.remove('access_token');
          Cookies.remove('refresh_token');
        } catch (error) {
          console.error('Logout failed:', error);
        }
      },
    }),

    refreshToken: builder.mutation<{ access_token: string; expires_in: number }, void>({
      query: () => {
        const refreshToken = Cookies.get('refresh_token');
        return {
          url: '/auth/refresh',
          method: 'POST',
          headers: {
            Authorization: `Bearer ${refreshToken}`,
          },
        };
      },
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          Cookies.set('access_token', data.access_token);
        } catch (error) {
          console.error('Token refresh failed:', error);
          Cookies.remove('access_token');
          Cookies.remove('refresh_token');
        }
      },
    }),

    getProfile: builder.query<User, void>({
      query: () => '/users/profile',
    }),

    updateProfile: builder.mutation<User, Partial<User>>({
      query: (updates) => ({
        url: '/users/profile',
        method: 'PUT',
        body: updates,
      }),
    }),

    changePassword: builder.mutation<
      { message: string },
      { current_password: string; new_password: string; confirm_password: string }
    >({
      query: (passwords) => ({
        url: '/users/password',
        method: 'PUT',
        body: passwords,
      }),
    }),
  }),
});

export const {
  useLoginMutation,
  useLogoutMutation,
  useRefreshTokenMutation,
  useGetProfileQuery,
  useUpdateProfileMutation,
  useChangePasswordMutation,
} = authApi;
