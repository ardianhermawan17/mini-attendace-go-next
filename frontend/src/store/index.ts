import { configureStore } from '@reduxjs/toolkit';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import authReducer from './slices/auth/auth.slice';
import attendanceReducer from './slices/attendance/attendance.slice';
import uiReducer from './slices/ui/ui.slice';
import { attendanceApi } from './api/attendanceApi';
import { authApi } from './api/authApi';
import { reportsApi } from './api/reportsApi';
import idempotencyMiddleware from './middlewares/idempotencyMiddleware';
import loggerMiddleware from './middlewares/logger';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    attendance: attendanceReducer,
    ui: uiReducer,
    [attendanceApi.reducerPath]: attendanceApi.reducer,
    [authApi.reducerPath]: authApi.reducer,
    [reportsApi.reducerPath]: reportsApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActionMatchers: [
          // Ignore all RTK Query async thunk actions
          /^(attendanceApi|authApi|reportsApi)\/.*\/(pending|fulfilled|rejected)$/,
        ],
        ignoredPaths: [attendanceApi.reducerPath, authApi.reducerPath, reportsApi.reducerPath],
      },
    })
      .concat(attendanceApi.middleware)
      .concat(authApi.middleware)
      .concat(reportsApi.middleware)
      .concat(idempotencyMiddleware)
      .concat(loggerMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Export typed hooks
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
