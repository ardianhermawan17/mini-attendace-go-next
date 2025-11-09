import { AttendanceRecord } from '@/types';

export interface OptimisticAttendance extends AttendanceRecord {
  isOptimistic?: boolean;
  tempId?: string;
}

export interface AttendanceState {
  records: Record<string, OptimisticAttendance>;
  todayRecord: OptimisticAttendance | null;
  isCheckingIn: boolean;
  isCheckingOut: boolean;
  checkInError: string | null;
  checkOutError: string | null;
  lastCheckInAt: string | null;
  lastCheckOutAt: string | null;
}
