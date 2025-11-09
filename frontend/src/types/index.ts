// User Types
export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  department: string;
  role: 'admin' | 'manager' | 'employee';
  status: 'active' | 'inactive' | 'suspended';
  created_at: string;
  updated_at: string;
  last_login?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  user: User;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

// Attendance Types
export type AttendanceStatus = 'hadir' | 'terlambat' | 'pulang_cepat' | 'tidak_hadir';

export interface AttendanceRecord {
  id: string;
  user_id: string;
  attendance_date: string;
  check_in_time: string | null;
  check_in_location?: string;
  check_in_device_info?: string;
  check_out_time: string | null;
  check_out_location?: string;
  check_out_device_info?: string;
  status: AttendanceStatus;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CheckInRequest {
  latitude: number;
  longitude: number;
  device: string;
}

export interface CheckOutRequest {
  latitude: number;
  longitude: number;
  device: string;
}

export interface AttendanceHistoryResponse {
  data: AttendanceRecord[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

// Report Types
export interface MonthlyReportSummary {
  year: number;
  month: number;
  total_employees: number;
  present: number;
  absent: number;
  late: number;
  summary_by_status: {
    present: { count: number; percentage: number };
    absent: { count: number; percentage: number };
    late: { count: number; percentage: number };
    leave: { count: number; percentage: number };
  };
}

export interface UserMonthlyReport {
  user_id: string;
  user_name: string;
  year: number;
  month: number;
  total_workdays: number;
  present: number;
  absent: number;
  late: number;
  leave: number;
  overtime_hours: number;
  avg_check_in_time: string;
  avg_check_out_time: string;
}

export interface AttendanceSummary {
  period: {
    from: string;
    to: string;
  };
  total_employees: number;
  total_records: number;
  attendance_stats: {
    present: number;
    absent: number;
    late: number;
    leave: number;
  };
  by_department: DepartmentStats[];
}

export interface DepartmentStats {
  department: string;
  total_employees: number;
  present: number;
  absent: number;
  late: number;
  leave: number;
}

// Leave Types
export type LeaveType = 'sick' | 'annual' | 'unpaid' | 'special';
export type LeaveStatus = 'pending' | 'approved' | 'rejected';

export interface LeaveRecord {
  id: string;
  user_id: string;
  leave_type: LeaveType;
  start_date: string;
  end_date: string;
  reason: string;
  status: LeaveStatus;
  approved_by?: string;
  approved_at?: string;
  rejection_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface LeaveRequest {
  start_date: string;
  end_date: string;
  leave_type: LeaveType;
  reason: string;
}

export interface LeaveHistoryResponse {
  data: LeaveRecord[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

// Overtime Types
export type OvertimeStatus = 'pending' | 'approved' | 'rejected';

export interface OvertimeRecord {
  id: string;
  user_id: string;
  attendance_date: string;
  overtime_hours: number;
  reason: string;
  status: OvertimeStatus;
  approved_by?: string;
  approved_at?: string;
  created_at: string;
  updated_at: string;
}

export interface OvertimeRequest {
  date: string;
  hours: number;
  reason: string;
}

// Absence Types
export type AbsenceType = 'unexcused' | 'excused' | 'sick';

export interface AbsenceRecord {
  id: string;
  user_id: string;
  absence_date: string;
  absence_type: AbsenceType;
  reason?: string;
  created_at: string;
  updated_at: string;
}

// Error Response
export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
  timestamp: string;
  path: string;
}

// Pagination
export interface PaginationParams {
  page?: number;
  limit?: number;
}

// API Response Wrapper
export interface ApiResponse<T> {
  data?: T;
  error?: ErrorResponse;
  message?: string;
}
