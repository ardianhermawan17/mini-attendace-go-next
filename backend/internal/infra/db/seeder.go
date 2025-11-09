package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SeederData struct {
	Users          []UserSeed
	AttendanceData []AttendanceRecordSeed
}

type UserSeed struct {
	ID         string
	Username   string
	Email      string
	Password   string
	FullName   string
	Department string
	Role       string
}

type AttendanceRecordSeed struct {
	UserEmail      string
	AttendanceDate string
	CheckInTime    *string
	CheckOutTime   *string
	Status         string
}

// SeedDatabase seeds the database with initial data
func SeedDatabase(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if data already seeded
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	if count > 0 {
		fmt.Println("Database already seeded, skipping...")
		return nil
	}

	data := getSeederData()

	// Seed users
	userMap := make(map[string]string) // email -> id
	for _, user := range data.Users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		userID := user.ID
		if userID == "" {
			userID = uuid.New().String()
		}

		err = pool.QueryRow(ctx,
			`INSERT INTO users (id, username, email, password_hash, full_name, department, role, status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (username) DO NOTHING
			 RETURNING id`,
			userID, user.Username, user.Email, string(hashedPassword), user.FullName, user.Department, user.Role, "active",
		).Scan(&userID)

		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Username, err)
		}

		userMap[user.Email] = userID
		fmt.Printf("✓ Seeded user: %s (Role: %s)\n", user.Email, user.Role)
	}

	// Seed attendance records
	for _, record := range data.AttendanceData {
		userID, exists := userMap[record.UserEmail]
		if !exists {
			return fmt.Errorf("user %s not found for attendance record", record.UserEmail)
		}

		var checkInTime, checkOutTime *time.Time
		if record.CheckInTime != nil {
			t, err := time.Parse("2006-01-02 15:04:05", *record.CheckInTime)
			if err != nil {
				return fmt.Errorf("failed to parse check-in time: %w", err)
			}
			checkInTime = &t
		}

		if record.CheckOutTime != nil {
			t, err := time.Parse("2006-01-02 15:04:05", *record.CheckOutTime)
			if err != nil {
				return fmt.Errorf("failed to parse check-out time: %w", err)
			}
			checkOutTime = &t
		}

		_, err := pool.Exec(ctx,
			`INSERT INTO attendance_records (id, user_id, attendance_date, check_in_time, check_out_time, status)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (user_id, attendance_date) DO NOTHING`,
			uuid.New(), userID, record.AttendanceDate, checkInTime, checkOutTime, record.Status,
		)

		if err != nil {
			return fmt.Errorf("failed to seed attendance record: %w", err)
		}
	}
	fmt.Printf("✓ Seeded %d attendance records\n", len(data.AttendanceData))

	fmt.Println("\n✓ Database seeding completed successfully!")
	return nil
}

func getSeederData() SeederData {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)
	threeDaysAgo := today.AddDate(0, 0, -3)

	checkInYesterday := yesterday.Format("2006-01-02") + " 08:30:00"
	checkOutYesterday := yesterday.Format("2006-01-02") + " 17:00:00"

	checkInTwoDaysAgo := twoDaysAgo.Format("2006-01-02") + " 09:15:00"
	checkOutTwoDaysAgo := twoDaysAgo.Format("2006-01-02") + " 17:30:00"

	checkInThreeDaysAgo := threeDaysAgo.Format("2006-01-02") + " 07:45:00"
	checkOutThreeDaysAgo := threeDaysAgo.Format("2006-01-02") + " 16:45:00"

	return SeederData{
		Users: []UserSeed{
			{
				ID:         "550e8400-e29b-41d4-a716-446655440001",
				Username:   "admin",
				Email:      "admin@trustmedis.com",
				Password:   "admin123",
				FullName:   "Admin User",
				Department: "IT",
				Role:       "admin",
			},
			{
				ID:         "550e8400-e29b-41d4-a716-446655440002",
				Username:   "john_doe",
				Email:      "john.doe@trustmedis.com",
				Password:   "password123",
				FullName:   "John Doe",
				Department: "HR",
				Role:       "employee",
			},
			{
				ID:         "550e8400-e29b-41d4-a716-446655440003",
				Username:   "jane_smith",
				Email:      "jane.smith@trustmedis.com",
				Password:   "password123",
				FullName:   "Jane Smith",
				Department: "Finance",
				Role:       "employee",
			},
			{
				ID:         "550e8400-e29b-41d4-a716-446655440004",
				Username:   "bob_manager",
				Email:      "bob.manager@trustmedis.com",
				Password:   "password123",
				FullName:   "Bob Manager",
				Department: "Operations",
				Role:       "manager",
			},
		},
		AttendanceData: []AttendanceRecordSeed{
			{
				UserEmail:      "john.doe@trustmedis.com",
				AttendanceDate: yesterday.Format("2006-01-02"),
				CheckInTime:    &checkInYesterday,
				CheckOutTime:   &checkOutYesterday,
				Status:         "present",
			},
			{
				UserEmail:      "jane.smith@trustmedis.com",
				AttendanceDate: yesterday.Format("2006-01-02"),
				CheckInTime:    &checkInYesterday,
				CheckOutTime:   &checkOutYesterday,
				Status:         "present",
			},
			{
				UserEmail:      "john.doe@trustmedis.com",
				AttendanceDate: twoDaysAgo.Format("2006-01-02"),
				CheckInTime:    &checkInTwoDaysAgo,
				CheckOutTime:   &checkOutTwoDaysAgo,
				Status:         "late",
			},
			{
				UserEmail:      "jane.smith@trustmedis.com",
				AttendanceDate: twoDaysAgo.Format("2006-01-02"),
				CheckInTime:    nil,
				CheckOutTime:   nil,
				Status:         "absent",
			},
			{
				UserEmail:      "bob_manager@trustmedis.com",
				AttendanceDate: threeDaysAgo.Format("2006-01-02"),
				CheckInTime:    &checkInThreeDaysAgo,
				CheckOutTime:   &checkOutThreeDaysAgo,
				Status:         "present",
			},
		},
	}
}
