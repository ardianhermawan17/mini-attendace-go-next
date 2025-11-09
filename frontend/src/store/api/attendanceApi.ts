import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import {
  AttendanceRecord,
  CheckInRequest,
  CheckOutRequest,
  AttendanceHistoryResponse,
} from '@/types';
import Cookies from 'js-cookie';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const attendanceApi = createApi({
  reducerPath: 'attendanceApi',
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
  tagTypes: ['Attendance', 'AttendanceToday'],
  endpoints: (builder) => ({
    checkIn: builder.mutation<AttendanceRecord, CheckInRequest>({
      query: (data) => ({
        url: '/attendance/check-in',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Attendance', 'AttendanceToday'],
    }),

    checkOut: builder.mutation<AttendanceRecord, CheckOutRequest>({
      query: (data) => ({
        url: '/attendance/check-out',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Attendance', 'AttendanceToday'],
    }),

    getTodayAttendance: builder.query<AttendanceRecord, void>({
      query: () => '/attendance/today',
      providesTags: ['AttendanceToday'],
    }),

    getAttendanceHistory: builder.query<
      AttendanceHistoryResponse,
      { from: string; to: string; page?: number; limit?: number }
    >({
      query: ({ from, to, page = 1, limit = 20 }) =>
        `/attendance/history?from=${from}&to=${to}&page=${page}&limit=${limit}`,
      providesTags: ['Attendance'],
    }),

    getUserAttendance: builder.query<
      AttendanceHistoryResponse,
      { userId: string; from: string; to: string; page?: number; limit?: number }
    >({
      query: ({ userId, from, to, page = 1, limit = 20 }) =>
        `/attendance/users/${userId}?from=${from}&to=${to}&page=${page}&limit=${limit}`,
      providesTags: ['Attendance'],
    }),
  }),
});

export const {
  useCheckInMutation,
  useCheckOutMutation,
  useGetTodayAttendanceQuery,
  useGetAttendanceHistoryQuery,
  useGetUserAttendanceQuery,
} = attendanceApi;
