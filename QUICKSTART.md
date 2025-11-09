# Quick Start Guide - Mini Attendance System

## 🚀 Get Started in 5 Minutes

### Prerequisites

- Docker & Docker Compose installed
- Git installed
- 4GB RAM available
- Ports 8080, 3000, 5432, 6379, 9092 available

### Step 1: Clone & Navigate

```bash
cd mini-attendace-fullstack-trustmedis
```

### Step 2: Start Services

```bash
docker-compose up -d
```

Wait for all services to be healthy (about 30 seconds):

```bash
docker-compose ps
```

### Step 3: Verify Installation

```bash
# Check API health
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "healthy",
#   "timestamp": "2024-01-15T08:30:00Z",
#   "services": {
#     "database": {"status": "healthy"},
#     "redis": {"status": "healthy"},
#     "kafka": {"status": "healthy"}
#   }
# }
```

### Step 4: Login & Test

```bash
# Login with sample credentials
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@trustmedis.com",
    "password": "password123"
  }'

# Save the access_token from response
export TOKEN="your_access_token_here"
```

### Step 5: Test Check-in

```bash
# Check-in
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "mobile",
    "metadata": {"location": "office"}
  }'

# Expected response:
# {
#   "id": "uuid",
#   "user_id": "uuid",
#   "attendance_date": "2024-01-15",
#   "check_in_time": "2024-01-15T08:30:00Z",
#   "status": "checked_in"
# }
```

### Step 6: Access Dashboards

- **API Documentation**: http://localhost:8080/swagger/index.html
- **Grafana Dashboards**: http://localhost:3000 (admin/admin)
- **PostgreSQL**: localhost:5432 (postgres/postgres)
- **Redis**: localhost:6379

## 📋 Sample Credentials

| Email | Password | Role |
|-------|----------|------|
| admin@trustmedis.com | admin123 | admin |
| john.doe@trustmedis.com | password123 | employee |
| jane.smith@trustmedis.com | password123 | employee |
| bob.wilson@trustmedis.com | password123 | employee |

## 🔧 Common Commands

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Stop Services

```bash
docker-compose down
```

### Reset Everything

```bash
docker-compose down -v
docker-compose up -d
```

### Access Database

```bash
# Connect to PostgreSQL
docker exec -it mini-attendance-postgres psql -U postgres -d mini_attendance

# List tables
\dt

# Query attendance records
SELECT * FROM attendance_records;
```

### Access Redis

```bash
# Connect to Redis
docker exec -it mini-attendance-redis redis-cli

# Check keys
KEYS *

# Get a value
GET lock:attendance:checkin:*
```

## 📚 API Examples

### Authentication

```bash
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@trustmedis.com",
    "password": "password123",
    "full_name": "New User"
  }'

# Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer $REFRESH_TOKEN"
```

### Attendance

```bash
# Get today's attendance
curl -X GET http://localhost:8080/api/v1/attendance/today \
  -H "Authorization: Bearer $TOKEN"

# Get attendance history
curl -X GET "http://localhost:8080/api/v1/attendance/history?start_date=2024-01-01&end_date=2024-01-31&page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"

# Check-out
curl -X POST http://localhost:8080/api/v1/attendance/check-out \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "mobile",
    "metadata": {"location": "office"}
  }'
```

### Reports

```bash
# Get absence report
curl -X GET "http://localhost:8080/api/v1/reports/absence?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer $TOKEN"

# Export as CSV
curl -X GET "http://localhost:8080/api/v1/reports/absence/export?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer $TOKEN" \
  -o report.csv
```

## 🐛 Troubleshooting

### Services not starting

```bash
# Check Docker daemon
docker ps

# Check logs
docker-compose logs

# Restart services
docker-compose restart
```

### Port already in use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database connection error

```bash
# Check PostgreSQL
docker-compose logs postgres

# Verify connection
docker exec mini-attendance-postgres pg_isready
```

### Redis connection error

```bash
# Check Redis
docker-compose logs redis

# Verify connection
docker exec mini-attendance-redis redis-cli ping
```

## 📖 Documentation

- **Full README**: [README.md](./README.md)
- **Backend README**: [backend/README.md](./backend/README.md)
- **Architecture**: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)

## 🎯 Next Steps

1. **Explore the API**: Use the Swagger UI at http://localhost:8080/swagger/index.html
2. **Check Logs**: View structured logs in Grafana at http://localhost:3000
3. **Test Workflows**: Try the complete check-in/check-out flow
4. **Review Code**: Explore the backend code structure
5. **Customize**: Modify configuration in `.env` file

## 💡 Tips

- **JWT Token**: Valid for 1 hour, use refresh token to get new one
- **Check-in Lock**: Prevents double check-in for 5 seconds
- **Database**: Automatically seeded with sample data
- **Logs**: All requests logged to Loki for debugging
- **Health Check**: Use `/health` endpoint to monitor system

## 🆘 Need Help?

1. Check logs: `docker-compose logs -f`
2. Review documentation: [README.md](./README.md)
3. Check architecture: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)
4. Verify configuration: `backend/.env`

## 🎉 You're Ready!

The system is now running and ready for testing. Start with the sample credentials and explore the API endpoints.

Happy coding! 🚀
