package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create migrations table if not exists
	createMigrationsTable := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		version BIGINT NOT NULL UNIQUE,
		dirty BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := pool.Exec(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run all migration files
	migrations := []struct {
		version int
		name    string
		up      string
	}{
		{
			version: 1,
			name:    "create_users_table",
			up: `
			CREATE TABLE IF NOT EXISTS users (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				email VARCHAR(255) NOT NULL UNIQUE,
				password_hash VARCHAR(255) NOT NULL,
				full_name VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL DEFAULT 'employee',
				is_active BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
			CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
			`,
		},
		{
			version: 2,
			name:    "create_work_schedules_table",
			up: `
			CREATE TABLE IF NOT EXISTS work_schedules (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				day_of_week INT NOT NULL,
				start_time TIME NOT NULL,
				end_time TIME NOT NULL,
				is_active BOOLEAN DEFAULT TRUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(user_id, day_of_week)
			);
			CREATE INDEX IF NOT EXISTS idx_work_schedules_user_id ON work_schedules(user_id);
			`,
		},
		{
			version: 3,
			name:    "create_holidays_table",
			up: `
			CREATE TABLE IF NOT EXISTS holidays (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				holiday_date DATE NOT NULL UNIQUE,
				holiday_name VARCHAR(255) NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_holidays_date ON holidays(holiday_date);
			`,
		},
		{
			version: 4,
			name:    "create_attendance_records_table",
			up: `
			CREATE TABLE IF NOT EXISTS attendance_records (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				attendance_date DATE NOT NULL,
				check_in_at TIMESTAMP,
				check_out_at TIMESTAMP,
				check_in_source VARCHAR(50),
				check_out_source VARCHAR(50),
				check_in_meta JSONB,
				check_out_meta JSONB,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(user_id, attendance_date)
			);
			CREATE INDEX IF NOT EXISTS idx_attendance_records_user_id ON attendance_records(user_id);
			CREATE INDEX IF NOT EXISTS idx_attendance_records_date ON attendance_records(attendance_date);
			CREATE INDEX IF NOT EXISTS idx_attendance_records_user_date ON attendance_records(user_id, attendance_date);
			CREATE INDEX IF NOT EXISTS idx_attendance_records_check_in ON attendance_records(check_in_at);
			CREATE INDEX IF NOT EXISTS idx_attendance_records_check_out ON attendance_records(check_out_at);
			`,
		},
		{
			version: 5,
			name:    "create_absence_reports_table",
			up: `
			CREATE TABLE IF NOT EXISTS absence_reports (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				report_date DATE NOT NULL,
				status VARCHAR(50) NOT NULL,
				check_in_time TIMESTAMP,
				check_out_time TIMESTAMP,
				work_hours DECIMAL(5, 2),
				notes TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(user_id, report_date)
			);
			CREATE INDEX IF NOT EXISTS idx_absence_reports_user_id ON absence_reports(user_id);
			CREATE INDEX IF NOT EXISTS idx_absence_reports_date ON absence_reports(report_date);
			CREATE INDEX IF NOT EXISTS idx_absence_reports_status ON absence_reports(status);
			CREATE INDEX IF NOT EXISTS idx_absence_reports_user_date ON absence_reports(user_id, report_date);
			`,
		},
		{
			version: 6,
			name:    "create_processed_events_table",
			up: `
			CREATE TABLE IF NOT EXISTS processed_events (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				event_id UUID NOT NULL UNIQUE,
				event_type VARCHAR(100) NOT NULL,
				processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_processed_events_event_id ON processed_events(event_id);
			CREATE INDEX IF NOT EXISTS idx_processed_events_type ON processed_events(event_type);
			`,
		},
		{
			version: 7,
			name:    "create_audit_events_table",
			up: `
			CREATE TABLE IF NOT EXISTS audit_events (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				event_id UUID NOT NULL,
				event_type VARCHAR(100) NOT NULL,
				user_id UUID REFERENCES users(id) ON DELETE SET NULL,
				action VARCHAR(100) NOT NULL,
				resource_type VARCHAR(100),
				resource_id UUID,
				old_values JSONB,
				new_values JSONB,
				ip_address VARCHAR(45),
				user_agent TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_audit_events_event_id ON audit_events(event_id);
			CREATE INDEX IF NOT EXISTS idx_audit_events_user_id ON audit_events(user_id);
			CREATE INDEX IF NOT EXISTS idx_audit_events_type ON audit_events(event_type);
			CREATE INDEX IF NOT EXISTS idx_audit_events_created_at ON audit_events(created_at);
			`,
		},
	}

	for _, migration := range migrations {
		// Check if migration already ran
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", migration.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if !exists {
			// Run migration
			if _, err := pool.Exec(ctx, migration.up); err != nil {
				return fmt.Errorf("failed to run migration %d (%s): %w", migration.version, migration.name, err)
			}

			// Mark migration as done
			if _, err := pool.Exec(ctx, "INSERT INTO schema_migrations (version, dirty) VALUES ($1, FALSE)", migration.version); err != nil {
				return fmt.Errorf("failed to mark migration as done: %w", err)
			}
		}
	}

	return nil
}
