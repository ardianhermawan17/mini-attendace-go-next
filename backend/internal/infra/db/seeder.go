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
	WorkSchedules  []WorkScheduleSeed
	Holidays       []HolidaySeed
	AttendanceData []AttendanceRecordSeed
}

type UserSeed struct {
	Email    string
	Password string
	FullName string
	Role     string
}

type WorkScheduleSeed struct {
	Email      string
	DayOfWeek  int
	StartTime  string
	EndTime    string
}

type HolidaySeed struct {
	Date        string
	Name        string
	Description string
}

type AttendanceRecordSeed struct {
	Email          string
	AttendanceDate string
	CheckInTime    *string
	CheckOutTime   *string
}

func SeedDatabase(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if data already seeded
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
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

		var userID string
		err = pool.QueryRow(ctx,
			`INSERT INTO users (email, password_hash, full_name, role, is_active)
			 VALUES ($1, $2, $3, $4, TRUE)
			 RETURNING id`,
			user.Email, string(hashedPassword), user.FullName, user.Role,
		).Scan(&userID)

		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Email, err)
		}

		userMap[user.Email] = userID
		fmt.Printf("✓ Seeded user: %s\n", user.Email)
	}

	// Seed work schedules
	for _, schedule := range data.WorkSchedules {
		userID, exists := userMap[schedule.Email]
		if !exists {
			return fmt.Errorf("user %s not found for work schedule", schedule.Email)
		}

		_, err := pool.Exec(ctx,
			`INSERT INTO work_schedules (user_id, day_of_week, start_time, end_time, is_active)
			 VALUES ($1, $2, $3, $4, TRUE)`,
			userID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime,
		)

		if err != nil {
			return fmt.Errorf("failed to seed work schedule: %w", err)
		}
	}
	fmt.Printf("✓ Seeded %d work schedules\n", len(data.WorkSchedules))

	// Seed holidays
	for _, holiday := range data.Holidays {
		_, err := pool.Exec(ctx,
			`INSERT INTO holidays (holiday_date, holiday_name, description)
			 VALUES ($1, $2, $3)`,
			holiday.Date, holiday.Name, holiday.Description,
		)

		if err != nil {
			return fmt.Errorf("failed to seed holiday: %w", err)
		}
	}
	fmt.Printf("✓ Seeded %d holidays\n", len(data.Holidays))

	// Seed attendance records
	for _, record := range data.AttendanceData {
		userID, exists := userMap[record.Email]
		if !exists {
			return fmt.Errorf("user %s not found for attendance record", record.Email)
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
			`INSERT INTO attendance_records (id, user_id, attendance_date, check_in_at, check_out_at, check_in_source, check_out_source)
			 VALUES ($1, $2, $3, $4, $5, 'mobile', 'mobile')`,
			uuid.New(), userID, record.AttendanceDate, checkInTime, checkOutTime,
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

	checkInYesterday := yesterday.Format("2006-01-02") + " 08:30:00"
	checkOutYesterday := yesterday.Format("2006-01-02") + " 17:00:00"

	checkInTwoDaysAgo := twoDaysAgo.Format("2006-01-02") + " 09:15:00"
	checkOutTwoDaysAgo := twoDaysAgo.Format("2006-01-02") + " 17:30:00"

	return SeederData{
		Users: []UserSeed{
			{
				Email:    "admin@trustmedis.com",
				Password: "admin123",
				FullName: "Admin User",
				Role:     "admin",
			},
			{
				Email:    "john.doe@trustmedis.com",
				Password: "password123",
				FullName: "John Doe",
				Role:     "employee",
			},
			{
				Email:    "jane.smith@trustmedis.com",
				Password: "password123",
				FullName: "Jane Smith",
				Role:     "employee",
			},
			{
				Email:    "bob.wilson@trustmedis.com",
				Password: "password123",
				FullName: "Bob Wilson",
				Role:     "employee",
			},
		},
		WorkSchedules: []WorkScheduleSeed{
			// Monday to Friday: 08:00 - 17:00
			{Email: "john.doe@trustmedis.com", DayOfWeek: 1, StartTime: "08:00", EndTime: "17:00"},
			{Email: "john.doe@trustmedis.com", DayOfWeek: 2, StartTime: "08:00", EndTime: "17:00"},
			{Email: "john.doe@trustmedis.com", DayOfWeek: 3, StartTime: "08:00", EndTime: "17:00"},
			{Email: "john.doe@trustmedis.com", DayOfWeek: 4, StartTime: "08:00", EndTime: "17:00"},
			{Email: "john.doe@trustmedis.com", DayOfWeek: 5, StartTime: "08:00", EndTime: "17:00"},

			{Email: "jane.smith@trustmedis.com", DayOfWeek: 1, StartTime: "08:00", EndTime: "17:00"},
			{Email: "jane.smith@trustmedis.com", DayOfWeek: 2, StartTime: "08:00", EndTime: "17:00"},
			{Email: "jane.smith@trustmedis.com", DayOfWeek: 3, StartTime: "08:00", EndTime: "17:00"},
			{Email: "jane.smith@trustmedis.com", DayOfWeek: 4, StartTime: "08:00", EndTime: "17:00"},
			{Email: "jane.smith@trustmedis.com", DayOfWeek: 5, StartTime: "08:00", EndTime: "17:00"},

			{Email: "bob.wilson@trustmedis.com", DayOfWeek: 1, StartTime: "08:00", EndTime: "17:00"},
			{Email: "bob.wilson@trustmedis.com", DayOfWeek: 2, StartTime: "08:00", EndTime: "17:00"},
			{Email: "bob.wilson@trustmedis.com", DayOfWeek: 3, StartTime: "08:00", EndTime: "17:00"},
			{Email: "bob.wilson@trustmedis.com", DayOfWeek: 4, StartTime: "08:00", EndTime: "17:00"},
			{Email: "bob.wilson@trustmedis.com", DayOfWeek: 5, StartTime: "08:00", EndTime: "17:00"},
		},
		Holidays: []HolidaySeed{
			{
				Date:        "2024-01-01",
				Name:        "New Year's Day",
				Description: "Public holiday",
			},
			{
				Date:        "2024-12-25",
				Name:        "Christmas Day",
				Description: "Public holiday",
			},
		},
		AttendanceData: []AttendanceRecordSeed{
			{
				Email:          "john.doe@trustmedis.com",
				AttendanceDate: yesterday.Format("2006-01-02"),
				CheckInTime:    &checkInYesterday,
				CheckOutTime:   &checkOutYesterday,
			},
			{
				Email:          "jane.smith@trustmedis.com",
				AttendanceDate: yesterday.Format("2006-01-02"),
				CheckInTime:    &checkInYesterday,
				CheckOutTime:   &checkOutYesterday,
			},
			{
				Email:          "john.doe@trustmedis.com",
				AttendanceDate: twoDaysAgo.Format("2006-01-02"),
				CheckInTime:    &checkInTwoDaysAgo,
				CheckOutTime:   &checkOutTwoDaysAgo,
			},
		},
	}
}
