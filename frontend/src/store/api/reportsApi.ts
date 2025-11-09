import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { MonthlyReportSummary, UserMonthlyReport, AttendanceSummary } from '@/types';
import Cookies from 'js-cookie';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const reportsApi = createApi({
  reducerPath: 'reportsApi',
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
  tagTypes: ['Reports'],
  endpoints: (builder) => ({
    getMonthlyReport: builder.query<MonthlyReportSummary, { year: number; month: number }>({
      query: ({ year, month }) => `/reports/monthly?year=${year}&month=${month}`,
      providesTags: ['Reports'],
    }),

    getUserMonthlyReport: builder.query<
      UserMonthlyReport,
      { userId: string; year: number; month: number }
    >({
      query: ({ userId, year, month }) =>
        `/reports/users/${userId}/monthly?year=${year}&month=${month}`,
      providesTags: ['Reports'],
    }),

    getAttendanceSummary: builder.query<AttendanceSummary, { from: string; to: string }>({
      query: ({ from, to }) => `/reports/summary?from=${from}&to=${to}`,
      providesTags: ['Reports'],
    }),

    getDepartmentReport: builder.query<
      { department: string; period: { from: string; to: string }; employees: any[] },
      { department: string; from: string; to: string; page?: number; limit?: number }
    >({
      query: ({ department, from, to, page = 1, limit = 50 }) =>
        `/reports/department/${department}?from=${from}&to=${to}&page=${page}&limit=${limit}`,
      providesTags: ['Reports'],
    }),

    exportReport: builder.query<
      Blob,
      { format: 'csv' | 'xlsx' | 'pdf'; from: string; to: string; departments?: string }
    >({
      query: ({ format, from, to, departments }) => {
        let url = `/reports/export?format=${format}&from=${from}&to=${to}`;
        if (departments) {
          url += `&include_departments=${departments}`;
        }
        return url;
      },
      providesTags: ['Reports'],
    }),
  }),
});

export const {
  useGetMonthlyReportQuery,
  useGetUserMonthlyReportQuery,
  useGetAttendanceSummaryQuery,
  useGetDepartmentReportQuery,
  useExportReportQuery,
} = reportsApi;
