package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Port        int
	Version     string
}

type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	DBName         string
	SSLMode        string
	MaxConnections int
	MigrationsPath string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}

type KafkaConfig struct {
	Brokers []string
	Topics  KafkaTopics
}

type KafkaTopics struct {
	CheckIn  string
	CheckOut string
}

type JWTConfig struct {
	Secret            string
	ExpirationMinutes int
	RefreshMinutes    int
}

type LoggingConfig struct {
	Level      string
	Format     string
	LokiURL    string
	LokiLabels map[string]string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "mini-attendance"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnvInt("APP_PORT", 8080),
			Version:     getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvInt("DB_PORT", 5432),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", "postgres"),
			DBName:         getEnv("DB_NAME", "mini_attendance"),
			SSLMode:        getEnv("DB_SSLMODE", "disable"),
			MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 25),
			MigrationsPath: getEnv("MIGRATIONS_PATH", "file://migrations"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			TTL:      time.Duration(getEnvInt("REDIS_TTL", 3600)) * time.Second,
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topics: KafkaTopics{
				CheckIn:  getEnv("KAFKA_TOPIC_CHECKIN", "attendance.check-in"),
				CheckOut: getEnv("KAFKA_TOPIC_CHECKOUT", "attendance.check-out"),
			},
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationMinutes: getEnvInt("JWT_EXPIRATION_MINUTES", 60),
			RefreshMinutes:    getEnvInt("JWT_REFRESH_MINUTES", 1440),
		},
		Logging: LoggingConfig{
			Level:   getEnv("LOG_LEVEL", "info"),
			Format:  getEnv("LOG_FORMAT", "json"),
			LokiURL: getEnv("LOKI_URL", "http://localhost:3100"),
			LokiLabels: map[string]string{
				"app":  getEnv("APP_NAME", "mini-attendance"),
				"env":  getEnv("ENVIRONMENT", "development"),
				"host": getEnv("HOSTNAME", "localhost"),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
