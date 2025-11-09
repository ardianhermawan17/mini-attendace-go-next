import { createSlice, PayloadAction, createEntityAdapter } from '@reduxjs/toolkit';
import { AttendanceRecord } from '@/types';
import { AttendanceState, OptimisticAttendance } from './attendance.types';

const attendanceAdapter = createEntityAdapter<OptimisticAttendance>({
  selectId: (record) => record.id,
  sortComparer: (a, b) => new Date(b.attendance_date).getTime() - new Date(a.attendance_date).getTime(),
});

const initialState: AttendanceState = {
  records: {},
  todayRecord: null,
  isCheckingIn: false,
  isCheckingOut: false,
  checkInError: null,
  checkOutError: null,
  lastCheckInAt: null,
  lastCheckOutAt: null,
};

const attendanceSlice = createSlice({
  name: 'attendance',
  initialState,
  reducers: {
    // Optimistic Check-In
    optimisticCheckInStart: (
      state,
      action: PayloadAction<{
        tempId: string;
        userId: string;
        date: string;
        checkInAt: string;
      }>
    ) => {
      const { tempId, userId, date, checkInAt } = action.payload;
      const optimisticRecord: OptimisticAttendance = {
        id: tempId,
        user_id: userId,
        attendance_date: date,
        check_in_time: checkInAt,
        check_out_time: null,
        status: 'hadir',
        isOptimistic: true,
        tempId,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      state.records[tempId] = optimisticRecord;
      state.todayRecord = optimisticRecord;
      state.isCheckingIn = true;
      state.checkInError = null;
    },

    optimisticCheckInSuccess: (
      state,
      action: PayloadAction<{ tempId: string; realRecord: AttendanceRecord }>
    ) => {
      const { tempId, realRecord } = action.payload;

      // Remove optimistic record
      delete state.records[tempId];

      // Add real record
      state.records[realRecord.id] = {
        ...realRecord,
        isOptimistic: false,
      };

      state.todayRecord = {
        ...realRecord,
        isOptimistic: false,
      };
      state.isCheckingIn = false;
      state.lastCheckInAt = realRecord.check_in_time;
    },

    optimisticCheckInFailure: (
      state,
      action: PayloadAction<{ tempId: string; error: string }>
    ) => {
      const { tempId, error } = action.payload;

      // Remove optimistic record
      delete state.records[tempId];

      state.todayRecord = null;
      state.isCheckingIn = false;
      state.checkInError = error;
    },

    // Optimistic Check-Out
    optimisticCheckOutStart: (
      state,
      action: PayloadAction<{
        recordId: string;
        checkOutAt: string;
      }>
    ) => {
      const { recordId, checkOutAt } = action.payload;

      if (state.records[recordId]) {
        state.records[recordId].check_out_time = checkOutAt;
        state.records[recordId].isOptimistic = true;
      }

      if (state.todayRecord?.id === recordId) {
        state.todayRecord.check_out_time = checkOutAt;
        state.todayRecord.isOptimistic = true;
      }

      state.isCheckingOut = true;
      state.checkOutError = null;
    },

    optimisticCheckOutSuccess: (
      state,
      action: PayloadAction<AttendanceRecord>
    ) => {
      const record = action.payload;

      state.records[record.id] = {
        ...record,
        isOptimistic: false,
      };

      state.todayRecord = {
        ...record,
        isOptimistic: false,
      };
      state.isCheckingOut = false;
      state.lastCheckOutAt = record.check_out_time;
    },

    optimisticCheckOutFailure: (
      state,
      action: PayloadAction<{ recordId: string; error: string }>
    ) => {
      const { recordId, error } = action.payload;

      // Rollback check-out time
      if (state.records[recordId]) {
        state.records[recordId].check_out_time = null;
        state.records[recordId].isOptimistic = false;
      }

      if (state.todayRecord?.id === recordId) {
        state.todayRecord.check_out_time = null;
        state.todayRecord.isOptimistic = false;
      }

      state.isCheckingOut = false;
      state.checkOutError = error;
    },

    // Set today's record
    setTodayRecord: (state, action: PayloadAction<AttendanceRecord | null>) => {
      if (action.payload) {
        state.todayRecord = {
          ...action.payload,
          isOptimistic: false,
        };
        state.records[action.payload.id] = state.todayRecord;
      } else {
        state.todayRecord = null;
      }
    },

    // Upsert records
    upsertRecords: (state, action: PayloadAction<AttendanceRecord[]>) => {
      action.payload.forEach((record) => {
        state.records[record.id] = {
          ...record,
          isOptimistic: false,
        };
      });
    },

    // Clear errors
    clearCheckInError: (state) => {
      state.checkInError = null;
    },

    clearCheckOutError: (state) => {
      state.checkOutError = null;
    },

    // Reset state
    resetAttendance: (state) => {
      state.records = {};
      state.todayRecord = null;
      state.isCheckingIn = false;
      state.isCheckingOut = false;
      state.checkInError = null;
      state.checkOutError = null;
    },
  },
});

export const {
  optimisticCheckInStart,
  optimisticCheckInSuccess,
  optimisticCheckInFailure,
  optimisticCheckOutStart,
  optimisticCheckOutSuccess,
  optimisticCheckOutFailure,
  setTodayRecord,
  upsertRecords,
  clearCheckInError,
  clearCheckOutError,
  resetAttendance,
} = attendanceSlice.actions;

export default attendanceSlice.reducer;
