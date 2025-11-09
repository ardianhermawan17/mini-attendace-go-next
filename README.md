# Mini Attendance System - Fullstack

A comprehensive attendance management system with event-driven architecture. This project includes both backend (Go/Gin) and frontend (Next.js) components.

## Project Structure

```
mini-attendace-fullstack-trustmedis/
├── backend/                 # Go backend API
│   ├── cmd/                # Application binaries
│   ├── internal/           # Internal packages
│   ├── go.mod             # Go module definition
│   ├── Dockerfile         # Docker image
│   ├── Makefile           # Build commands
│   └── README.md          # Backend documentation
├── frontend/              # Next.js frontend (future)
├── docker-compose.yml     # Docker Compose configuration
└── README.md             # This file
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Git

### Running the Application

1. **Clone the repository**
   ```bash
   cd mini-attendace-fullstack-trustmedis
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Verify services are running**
   ```bash
   docker-compose ps
   ```

4. **Access the application**
   - **API**: http://localhost:8080
   - **Health Check**: http://localhost:8080/health
   - **Grafana**: http://localhost:3000 (admin/admin)
   - **PostgreSQL**: localhost:5432
   - **Redis**: localhost:6379
   - **Kafka**: localhost:9092

### Sample Credentials

After seeding, use these credentials to login:

- **Admin**: admin@trustmedis.com / admin123
- **Employee 1**: john.doe@trustmedis.com / password123
- **Employee 2**: jane.smith@trustmedis.com / password123
- **Employee 3**: bob.wilson@trustmedis.com / password123

## Architecture Overview

### Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Backend | Go (Gin) | 1.21 |
| Frontend | Next.js | 14+ |
| Database | PostgreSQL | 16 |
| Cache | Redis | 7 |
| Message Broker | Apache Kafka | 7.5 |
| Logging | Grafana Loki | 2.9 |
| Visualization | Grafana | 10.2 |
| Container | Docker | Latest |

### System Design

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (Next.js)                       │
│                    http://localhost:3000                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Backend API (Go/Gin)                        │
│                  http://localhost:8080                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ • Authentication (JWT)                               │   │
│  │ • Check-in/Check-out (Distributed Locks)            │   │
│  │ • Attendance Reports                                 │   │
│  │ • Health Checks                                      │   │
│  └──────────────────────────────────────────────────────┘   │
└────────┬──────────────────────────────────────────────┬─────┘
         │                                              │
         ▼                                              ▼
    ┌─────────────┐                            ┌──────────────┐
    │ PostgreSQL  │                            │    Redis     │
    │  Database   │                            │    Cache     │
    └─────────────┘                            └──────────────┘
         │
         ▼
    ┌─────────────┐
    │   Kafka     │
    │   Broker    │
    └─────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────────┐
    │         Report Service (Consumer - Future)              │
    │  • Processes attendance events                          │
    │  • Generates absence reports                            │
    │  • Updates cache                                        │
    └─────────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────────┐
    │    Observability Stack                                  │
    │  • Loki (Log Aggregation)                              │
    │  • Grafana (Visualization)                             │
    └─────────────────────────────────────────────────────────┘
```

## Features

### Core Features

✅ **User Authentication**
- JWT-based authentication
- Access and refresh tokens
- User registration and login

✅ **Attendance Management**
- Check-in/Check-out functionality
- Prevent double check-in with distributed locks
- Attendance history tracking

✅ **Reporting**
- Absence reports with status tracking
- CSV export functionality
- Filtering and pagination

✅ **Event-Driven Architecture**
- Kafka-based event publishing
- Asynchronous report generation
- Event idempotency tracking

✅ **Observability**
- Structured logging with Loki
- Grafana dashboards
- Health check endpoints
- System monitoring

✅ **Security**
- Password hashing with bcrypt
- JWT token validation
- CORS support
- Database transaction safety

## API Documentation

### Authentication Endpoints

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@trustmedis.com",
    "password": "password123"
  }'

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@trustmedis.com",
    "password": "password123",
    "full_name": "New User"
  }'

# Refresh Token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer {refresh_token}"
```

### Attendance Endpoints

```bash
# Check-in
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "mobile",
    "metadata": {"location": "office"}
  }'

# Check-out
curl -X POST http://localhost:8080/api/v1/attendance/check-out \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "mobile",
    "metadata": {"location": "office"}
  }'

# Get Today's Attendance
curl -X GET http://localhost:8080/api/v1/attendance/today \
  -H "Authorization: Bearer {access_token}"

# Get Attendance History
curl -X GET "http://localhost:8080/api/v1/attendance/history?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer {access_token}"
```

### Report Endpoints

```bash
# Get Absence Report
curl -X GET "http://localhost:8080/api/v1/reports/absence?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer {access_token}"

# Export Report as CSV
curl -X GET "http://localhost:8080/api/v1/reports/absence/export?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer {access_token}" \
  -o report.csv
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Database Schema

### Key Tables

- **users**: User/employee master data
- **work_schedules**: Expected work hours per user/day
- **holidays**: Public holidays
- **attendance_records**: Check-in/check-out events
- **absence_reports**: Denormalized report data
- **processed_events**: Event idempotency tracking
- **audit_events**: Audit trail

## Deployment

### Local Development

See [Backend README](./backend/README.md) for local development setup.

### Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Cloud Deployment

The system is designed for cloud deployment with:
- Environment-based configuration
- Health check endpoints
- Graceful shutdown
- Structured logging
- Stateless API design

## Monitoring

### Grafana Dashboard

Access at http://localhost:3000:
- Username: admin
- Password: admin

### Available Metrics

- Application logs (via Loki)
- Database performance
- Redis usage
- Kafka lag
- API response times

## Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run locally
go run cmd/attendance-api/main.go

# Run tests
go test ./...

# Build
make build
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## Troubleshooting

### Services not starting

```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart

# Full reset
docker-compose down -v
docker-compose up -d
```

### Database connection issues

```bash
# Check PostgreSQL
docker-compose logs postgres

# Connect to database
psql -h localhost -U postgres -d mini_attendance
```

### API not responding

```bash
# Check health
curl http://localhost:8080/health

# View logs
docker-compose logs backend
```

## Performance Optimization

- Connection pooling (25 DB connections)
- Redis caching for reports
- Kafka partitioning by user_id
- Database indexes on frequently queried columns
- Distributed locks with short TTL

## Security Considerations

- JWT tokens expire after 1 hour
- Passwords hashed with bcrypt
- CORS enabled for frontend
- Database transactions ensure consistency
- Non-root Docker user
- Environment-based secrets

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## Support

For issues and questions, please contact the development team.

## License

This project is proprietary and confidential.
