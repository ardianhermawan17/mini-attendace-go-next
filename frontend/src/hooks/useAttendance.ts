import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store';
import {
  optimisticCheckInStart,
  optimisticCheckInSuccess,
  optimisticCheckInFailure,
  optimisticCheckOutStart,
  optimisticCheckOutSuccess,
  optimisticCheckOutFailure,
  setTodayRecord,
} from '@/store/slices/attendance/attendance.slice';
import {
  useCheckInMutation,
  useCheckOutMutation,
  useGetTodayAttendanceQuery,
} from '@/store/api/attendanceApi';
import { CheckInRequest, CheckOutRequest } from '@/types';

/**
 * Generate a temporary ID for optimistic updates
 */
function generateTempId(): string {
  return `temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

export const useAttendance = () => {
  const dispatch = useAppDispatch();
  const attendance = useAppSelector((state) => state.attendance);
  const userId = useAppSelector((state) => state.auth.user?.id);

  const [checkInMutation, { isLoading: isCheckingIn }] = useCheckInMutation();
  const [checkOutMutation, { isLoading: isCheckingOut }] = useCheckOutMutation();
  const { data: todayData, isLoading: isFetchingToday } = useGetTodayAttendanceQuery();

  const checkIn = useCallback(
    async (data: CheckInRequest) => {
      if (!userId) {
        throw new Error('User not authenticated');
      }

      const tempId = generateTempId();
      const now = new Date().toISOString();
      const today = new Date().toISOString().split('T')[0];

      // Dispatch optimistic update
      dispatch(
        optimisticCheckInStart({
          tempId,
          userId,
          date: today,
          checkInAt: now,
        })
      );

      try {
        const response = await checkInMutation(data).unwrap();

        // Replace optimistic record with real one
        dispatch(
          optimisticCheckInSuccess({
            tempId,
            realRecord: response,
          })
        );

        return response;
      } catch (error: any) {
        const errorMessage = error?.data?.error?.message || 'Check-in failed';

        // Rollback optimistic update
        dispatch(
          optimisticCheckInFailure({
            tempId,
            error: errorMessage,
          })
        );

        throw error;
      }
    },
    [userId, checkInMutation, dispatch]
  );

  const checkOut = useCallback(
    async (data: CheckOutRequest) => {
      if (!attendance.todayRecord?.id) {
        throw new Error('No active check-in found');
      }

      const recordId = attendance.todayRecord.id;
      const now = new Date().toISOString();

      // Dispatch optimistic update
      dispatch(
        optimisticCheckOutStart({
          recordId,
          checkOutAt: now,
        })
      );

      try {
        const response = await checkOutMutation(data).unwrap();

        // Update with real record
        dispatch(optimisticCheckOutSuccess(response));

        return response;
      } catch (error: any) {
        const errorMessage = error?.data?.error?.message || 'Check-out failed';

        // Rollback optimistic update
        dispatch(
          optimisticCheckOutFailure({
            recordId,
            error: errorMessage,
          })
        );

        throw error;
      }
    },
    [attendance.todayRecord?.id, checkOutMutation, dispatch]
  );

  // Sync today's record from API
  const syncTodayRecord = useCallback(() => {
    if (todayData) {
      dispatch(setTodayRecord(todayData));
    }
  }, [todayData, dispatch]);

  return {
    ...attendance,
    checkIn,
    checkOut,
    syncTodayRecord,
    isCheckingIn,
    isCheckingOut,
    isFetchingToday,
    todayRecord: attendance.todayRecord || todayData || null,
  };
};
