# Mini Attendance System - Backend

A comprehensive attendance management system built with Go (Gin), PostgreSQL, Kafka, and Redis. This system provides check-in/check-out functionality with event-driven architecture for generating attendance reports.

## Features

- **JWT Authentication**: Secure user authentication with access and refresh tokens
- **Check-in/Check-out**: Prevent double check-in with distributed locks
- **Event-Driven Architecture**: Kafka-based event publishing for asynchronous report generation
- **Distributed Locking**: Redis-based locks to prevent race conditions
- **Caching**: Redis caching for improved performance
- **Structured Logging**: Integrated with Grafana Loki for log aggregation
- **Health Checks**: Comprehensive health check endpoint monitoring all services
- **API Documentation**: Swagger/OpenAPI documentation
- **Database Migrations**: Automated schema management
- **Docker Support**: Full Docker and Docker Compose setup

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Message Broker**: Apache Kafka
- **Logging**: Grafana Loki
- **Visualization**: Grafana
- **Container**: Docker & Docker Compose

## Project Structure

```
backend/
├── cmd/
│   ├── attendance-api/          # Main API binary
│   │   └── main.go
│   └── attendance-worker/       # Consumer service (future)
│       └── main.go
├── internal/
│   ├── config/                  # Configuration management
│   ├── api/
│   │   └── http/
│   │       ├── handler/         # HTTP handlers
│   │       └── middleware/      # HTTP middleware
│   ├── domain/                  # Domain models
│   ├── application/             # Application services
│   ├── infra/
│   │   ├── db/                  # Database layer
│   │   ├── kafka/               # Kafka producer
│   │   ├── redis/               # Redis client
│   │   └── observability/       # Logging
│   └── services/
│       └── auth/                # Authentication service
├── migrations/                  # SQL migration files
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies lock
├── Dockerfile                   # Docker image definition
└── README.md                    # This file
```

## Prerequisites

- Docker & Docker Compose (for containerized setup)
- Go 1.21+ (for local development)
- PostgreSQL 16+ (for local development)
- Redis 7+ (for local development)
- Kafka (for local development)

## Installation & Setup

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   cd mini-attendace-fullstack-trustmedis
   ```

2. **Copy environment file**
   ```bash
   cp backend/.env.example backend/.env
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Verify services are running**
   ```bash
   docker-compose ps
   ```

5. **Access the application**
   - API: http://localhost:8080
   - Grafana: http://localhost:3000 (admin/admin)
   - Health Check: http://localhost:8080/health

### Option 2: Local Development Setup

1. **Install dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   # Create database
   createdb mini_attendance
   ```

3. **Set up Redis**
   ```bash
   # Start Redis server
   redis-server
   ```

4. **Set up Kafka**
   ```bash
   # Start Zookeeper and Kafka (or use Docker)
   docker run -d --name zookeeper -p 2181:2181 confluentinc/cp-zookeeper:7.5.0
   docker run -d --name kafka -p 9092:9092 confluentinc/cp-kafka:7.5.0
   ```

5. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

6. **Run the application**
   ```bash
   go run cmd/attendance-api/main.go
   ```

## Configuration

### Environment Variables

See `.env.example` for all available configuration options:

```env
# Application
APP_NAME=mini-attendance
ENVIRONMENT=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mini_attendance

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka
KAFKA_BROKERS=localhost:9092

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_MINUTES=60

# Logging
LOG_LEVEL=info
LOKI_URL=http://localhost:3100
```

## Database Migrations

Migrations run automatically on application startup. The system tracks migration state in the `schema_migrations` table.

### Manual Migration

To manually run migrations:

```bash
go run cmd/attendance-api/main.go
```

### Database Schema

The system creates the following tables:

- **users**: User/employee master data
- **work_schedules**: Expected work hours per user/day
- **holidays**: Public holidays
- **attendance_records**: Check-in/check-out events (main service writes)
- **absence_reports**: Denormalized report data (consumer writes)
- **processed_events**: Event idempotency tracking
- **audit_events**: Audit trail for compliance

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh access token

### Attendance

- `POST /api/v1/attendance/check-in` - Record check-in
- `POST /api/v1/attendance/check-out` - Record check-out
- `GET /api/v1/attendance/today` - Get today's attendance
- `GET /api/v1/attendance/history` - Get attendance history

### Reports

- `GET /api/v1/reports/absence` - Get absence report
- `GET /api/v1/reports/absence/export` - Export report as CSV

### Health

- `GET /health` - System health check

## API Documentation

Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## Seeding Data

The application automatically seeds sample data on first run:

- **Admin User**: admin@trustmedis.com / admin123
- **Employee Users**: john.doe@trustmedis.com, jane.smith@trustmedis.com, bob.wilson@trustmedis.com
- **Sample Attendance**: Historical check-in/check-out records

## Monitoring & Logging

### Grafana Dashboard

Access Grafana at http://localhost:3000:
- Username: admin
- Password: admin

### Loki Logs

Logs are aggregated in Loki and can be queried through Grafana:

```
{app="mini-attendance"}
```

### Health Check

Monitor system health:

```bash
curl http://localhost:8080/health
```

Response includes status of:
- Database connection
- Redis connection
- Kafka connectivity

## Event-Driven Architecture

### Check-in Event Flow

1. User calls `POST /api/v1/attendance/check-in`
2. API acquires distributed lock (Redis)
3. Validates no active check-in exists
4. Records check-in in database
5. Publishes `attendance.check-in` event to Kafka
6. Releases lock
7. Returns response

### Check-out Event Flow

1. User calls `POST /api/v1/attendance/check-out`
2. API finds active check-in record
3. Records check-out time
4. Publishes `attendance.check-out` event to Kafka
5. Returns response

### Report Generation (Consumer)

The report service (future implementation) will:
1. Consume events from Kafka
2. Track processed events for idempotency
3. Generate/update absence reports
4. Cache results in Redis

## Security Considerations

- JWT tokens expire after 1 hour
- Refresh tokens expire after 24 hours
- Distributed locks prevent race conditions
- Database transactions ensure consistency
- Passwords hashed with bcrypt
- CORS enabled for frontend integration
- Non-root Docker user for container security

## Performance Optimization

- Connection pooling for database (25 connections)
- Redis caching for frequently accessed data
- Kafka partitioning by user_id for ordering
- Database indexes on frequently queried columns
- Distributed locks with short TTL (5 seconds)

## Deployment

### Cloud Deployment

The system is designed for cloud deployment:

1. **Environment Configuration**: All settings via environment variables
2. **Health Checks**: Kubernetes-ready health endpoints
3. **Graceful Shutdown**: 30-second shutdown timeout
4. **Logging**: Structured JSON logs for log aggregation
5. **Scalability**: Stateless API design

### Kubernetes Deployment

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mini-attendance-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mini-attendance-api
  template:
    metadata:
      labels:
        app: mini-attendance-api
    spec:
      containers:
      - name: api
        image: mini-attendance:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db_host
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker-compose logs postgres

# Verify connection
psql -h localhost -U postgres -d mini_attendance
```

### Redis Connection Issues

```bash
# Check Redis is running
docker-compose logs redis

# Test connection
redis-cli ping
```

### Kafka Issues

```bash
# Check Kafka is running
docker-compose logs kafka

# List topics
docker exec mini-attendance-kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Application Logs

```bash
# View application logs
docker-compose logs backend

# Follow logs in real-time
docker-compose logs -f backend
```

## Development

### Running Tests

```bash
go test ./...
```

### Building Locally

```bash
go build -o attendance-api ./cmd/attendance-api/main.go
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Run vet
go vet ./...
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

This project is proprietary and confidential.

## Support

For issues and questions, please contact the development team.
