-- Drop triggers
DROP TRIGGER IF EXISTS absence_records_update_timestamp ON absence_records;
DROP TRIGGER IF EXISTS leave_records_update_timestamp ON leave_records;
DROP TRIGGER IF EXISTS overtime_records_update_timestamp ON overtime_records;
DROP TRIGGER IF EXISTS attendance_records_update_timestamp ON attendance_records;
DROP TRIGGER IF EXISTS users_update_timestamp ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS absence_records;
DROP TABLE IF EXISTS leave_records;
DROP TABLE IF EXISTS overtime_records;
DROP TABLE IF EXISTS attendance_records;
DROP TABLE IF EXISTS users;
