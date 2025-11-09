-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'employee', -- 'admin', 'manager', 'employee'
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'inactive', 'suspended'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create attendance records table
CREATE TABLE IF NOT EXISTS attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in_time TIMESTAMP,
    check_in_location VARCHAR(255),
    check_in_device_info VARCHAR(255),
    check_out_time TIMESTAMP,
    check_out_location VARCHAR(255),
    check_out_device_info VARCHAR(255),
    attendance_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'present', -- 'present', 'absent', 'late', 'early_leave'
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, attendance_date)
);

-- Create overtime records table
CREATE TABLE IF NOT EXISTS overtime_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    overtime_hours DECIMAL(5,2) NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create leave records table
CREATE TABLE IF NOT EXISTS leave_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    leave_type VARCHAR(50) NOT NULL, -- 'sick', 'annual', 'unpaid', 'special'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create absence records table
CREATE TABLE IF NOT EXISTS absence_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    absence_date DATE NOT NULL,
    absence_type VARCHAR(50) NOT NULL, -- 'unexcused', 'excused', 'sick'
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, absence_date)
);

-- Create audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_attendance_user_id ON attendance_records(user_id);
CREATE INDEX idx_attendance_date ON attendance_records(attendance_date);
CREATE INDEX idx_attendance_user_date ON attendance_records(user_id, attendance_date);
CREATE INDEX idx_overtime_user_id ON overtime_records(user_id);
CREATE INDEX idx_overtime_date ON overtime_records(attendance_date);
CREATE INDEX idx_leave_user_id ON leave_records(user_id);
CREATE INDEX idx_leave_dates ON leave_records(start_date, end_date);
CREATE INDEX idx_absence_user_id ON absence_records(user_id);
CREATE INDEX idx_absence_date ON absence_records(absence_date);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_timestamp BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER attendance_records_update_timestamp BEFORE UPDATE ON attendance_records
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER overtime_records_update_timestamp BEFORE UPDATE ON overtime_records
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER leave_records_update_timestamp BEFORE UPDATE ON leave_records
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER absence_records_update_timestamp BEFORE UPDATE ON absence_records
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
