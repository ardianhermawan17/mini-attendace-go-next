import axios, { AxiosInstance, AxiosError } from 'axios';
import Cookies from 'js-cookie';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private client: AxiosInstance;
  private refreshTokenPromise: Promise<string> | null = null;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const token = Cookies.get('access_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any;

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            // Prevent multiple refresh token requests
            if (!this.refreshTokenPromise) {
              this.refreshTokenPromise = this.refreshAccessToken();
            }

            const newToken = await this.refreshTokenPromise;
            this.refreshTokenPromise = null;

            Cookies.set('access_token', newToken);
            originalRequest.headers.Authorization = `Bearer ${newToken}`;

            return this.client(originalRequest);
          } catch (refreshError) {
            // Refresh failed, redirect to login
            Cookies.remove('access_token');
            Cookies.remove('refresh_token');
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  private async refreshAccessToken(): Promise<string> {
    const refreshToken = Cookies.get('refresh_token');
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await axios.post(
      `${API_URL}/auth/refresh`,
      {},
      {
        headers: {
          Authorization: `Bearer ${refreshToken}`,
        },
      }
    );

    return response.data.access_token;
  }

  public getClient(): AxiosInstance {
    return this.client;
  }

  public setAuthToken(token: string): void {
    Cookies.set('access_token', token);
    this.client.defaults.headers.common.Authorization = `Bearer ${token}`;
  }

  public setRefreshToken(token: string): void {
    Cookies.set('refresh_token', token);
  }

  public clearTokens(): void {
    Cookies.remove('access_token');
    Cookies.remove('refresh_token');
    delete this.client.defaults.headers.common.Authorization;
  }

  public getAuthToken(): string | undefined {
    return Cookies.get('access_token');
  }

  public getRefreshToken(): string | undefined {
    return Cookies.get('refresh_token');
  }
}

export const apiClient = new ApiClient();
export default apiClient.getClient();
