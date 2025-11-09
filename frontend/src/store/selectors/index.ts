import { createSelector } from '@reduxjs/toolkit';
import { RootState } from '../index';

// ============================================
// Attendance Selectors
// ============================================

/**
 * Select all attendance records as an array
 */
export const selectAllAttendanceRecords = (state: RootState) =>
  Object.values(state.attendance.records);

/**
 * Select today's attendance record
 */
export const selectTodayRecord = (state: RootState) => state.attendance.todayRecord;

/**
 * Select attendance record by user ID and date
 * Memoized for performance
 */
export const selectAttendanceByUserAndDate = createSelector(
  [
    (state: RootState) => state.attendance.records,
    (_: RootState, userId: string) => userId,
    (_: RootState, __: string, date: string) => date,
  ],
  (records, userId, date) => {
    return Object.values(records).find(
      (record) => record.user_id === userId && record.attendance_date === date
    );
  }
);

/**
 * Select attendance records for a date range
 * Memoized for performance
 */
export const selectAttendanceListForRange = createSelector(
  [
    (state: RootState) => state.attendance.records,
    (_: RootState, from: string, to: string) => ({ from, to }),
  ],
  (records, { from, to }) => {
    const fromDate = new Date(from);
    const toDate = new Date(to);

    return Object.values(records)
      .filter((record) => {
        const recordDate = new Date(record.attendance_date);
        return recordDate >= fromDate && recordDate <= toDate;
      })
      .sort(
        (a, b) => new Date(b.attendance_date).getTime() - new Date(a.attendance_date).getTime()
      );
  }
);

/**
 * Count present days in a date range
 */
export const selectPresentCountForRange = createSelector(
  [selectAttendanceListForRange],
  (records) => {
    return records.filter((record) => record.status === 'hadir').length;
  }
);

/**
 * Select attendance records grouped by status
 */
export const selectAttendanceByStatus = createSelector([selectAllAttendanceRecords], (records) => {
  return records.reduce(
    (acc, record) => {
      if (!acc[record.status]) {
        acc[record.status] = [];
      }
      acc[record.status].push(record);
      return acc;
    },
    {} as Record<string, typeof records>
  );
});

/**
 * Select optimistic (pending) records
 */
export const selectOptimisticRecords = createSelector([selectAllAttendanceRecords], (records) => {
  return records.filter((record) => 'isOptimistic' in record && record.isOptimistic);
});

/**
 * Check if user has checked in today
 */
export const selectHasCheckedInToday = (state: RootState) => {
  return !!state.attendance.todayRecord?.check_in_time;
};

/**
 * Check if user has checked out today
 */
export const selectHasCheckedOutToday = (state: RootState) => {
  return !!state.attendance.todayRecord?.check_out_time;
};

// ============================================
// Auth Selectors
// ============================================

/**
 * Select current user
 */
export const selectCurrentUser = (state: RootState) => state.auth.user;

/**
 * Select authentication status
 */
export const selectIsAuthenticated = (state: RootState) => state.auth.isAuthenticated;

/**
 * Select user role
 */
export const selectUserRole = (state: RootState) => state.auth.user?.role;

/**
 * Check if user is admin
 */
export const selectIsAdmin = createSelector([selectUserRole], (role) => role === 'admin');

/**
 * Check if user is manager or admin
 */
export const selectIsManagerOrAdmin = createSelector(
  [selectUserRole],
  (role) => role === 'admin' || role === 'manager'
);

// ============================================
// UI Selectors
// ============================================

/**
 * Select if any operation is loading
 */
export const selectIsLoading = (state: RootState) =>
  state.attendance.isCheckingIn || state.attendance.isCheckingOut || state.auth.isLoading;

/**
 * Select if there are any errors
 */
export const selectHasErrors = (state: RootState) =>
  !!state.attendance.checkInError || !!state.attendance.checkOutError || !!state.auth.error;

/**
 * Select all errors
 */
export const selectAllErrors = createSelector(
  [
    (state: RootState) => state.attendance.checkInError,
    (state: RootState) => state.attendance.checkOutError,
    (state: RootState) => state.auth.error,
  ],
  (checkInError, checkOutError, authError) => {
    return {
      checkIn: checkInError,
      checkOut: checkOutError,
      auth: authError,
    };
  }
);

/**
 * Select sidebar state
 */
export const selectSidebarOpen = (state: RootState) => state.ui.sidebarOpen;

// ============================================
// Combined Selectors
// ============================================

/**
 * Select attendance summary for current user
 */
export const selectUserAttendanceSummary = createSelector(
  [selectAllAttendanceRecords, selectCurrentUser],
  (records, user) => {
    if (!user) return null;

    const userRecords = records.filter((record) => record.user_id === user.id);

    const summary = {
      total: userRecords.length,
      present: userRecords.filter((r) => r.status === 'hadir').length,
      late: userRecords.filter((r) => r.status === 'terlambat').length,
      early_leave: userRecords.filter((r) => r.status === 'pulang_cepat').length,
      absent: userRecords.filter((r) => r.status === 'tidak_hadir').length,
    };

    return summary;
  }
);

/**
 * Select if user can check in/out
 */
export const selectCanCheckIn = createSelector(
  [selectHasCheckedInToday, (state: RootState) => state.attendance.isCheckingIn],
  (hasCheckedIn, isCheckingIn) => !hasCheckedIn && !isCheckingIn
);

export const selectCanCheckOut = createSelector(
  [
    selectHasCheckedInToday,
    selectHasCheckedOutToday,
    (state: RootState) => state.attendance.isCheckingOut,
  ],
  (hasCheckedIn, hasCheckedOut, isCheckingOut) => hasCheckedIn && !hasCheckedOut && !isCheckingOut
);
