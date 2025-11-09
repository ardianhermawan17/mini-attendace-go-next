# Mini Attendance System - Architecture Documentation

## System Overview

The Mini Attendance System is a comprehensive attendance management platform built with a modern, event-driven architecture. It provides real-time check-in/check-out functionality with asynchronous report generation.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend Layer                            │
│                      (Next.js - Future)                          │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway Layer                             │
│                   (Gin Web Framework)                            │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ • Authentication (JWT)                                   │   │
│  │ • Request Validation                                     │   │
│  │ • Rate Limiting                                          │   │
│  │ • CORS Handling                                          │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────┬────────────────────────────────────────────────��─┬─────┘
         │                                                  │
         ▼                                                  ▼
    ┌─────────────┐                            ┌──────────────────┐
    │ PostgreSQL  │                            │ Redis Cache      │
    │  Database   │                            │ (Distributed     │
    │             │                            │  Locks, Cache)   │
    └─────────────┘                            └──────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────────────┐
    │              Event Broker (Kafka)                            │
    │  Topics:                                                     │
    │  • attendance.check-in                                       │
    │  • attendance.check-out                                      │
    └─────────────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────────────┐
    │         Report Service (Consumer - Future)                   │
    │  • Processes events asynchronously                           │
    │  • Generates absence reports                                 │
    │  • Updates cache                                             │
    │  • Ensures idempotency                                       │
    └─────────────────────────────────────────────────────────────┘
         │
         ▼
    ┌─────────────────────────────────────────────────────────────┐
    │    Observability Stack                                       │
    │  • Loki (Log Aggregation)                                   │
    │  • Grafana (Visualization)                                  │
    │  • Structured Logging                                       │
    └─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. API Layer (cmd/attendance-api)

**Responsibilities:**
- HTTP request handling
- JWT authentication
- Input validation
- Response formatting
- Error handling

**Key Endpoints:**
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/attendance/check-in` - Record check-in
- `POST /api/v1/attendance/check-out` - Record check-out
- `GET /api/v1/reports/absence` - Retrieve reports
- `GET /health` - System health check

### 2. Domain Layer (internal/domain)

**Responsibilities:**
- Business logic encapsulation
- Domain models and entities
- Domain events
- Repository interfaces

**Key Models:**
- `User` - Employee/user entity
- `Attendance` - Check-in/check-out record
- `AbsenceReport` - Generated report
- `WorkSchedule` - Expected work hours

### 3. Application Layer (internal/application)

**Responsibilities:**
- Use case orchestration
- Command/Query handling
- DTO mapping
- Transaction management

**Key Services:**
- `CheckInCommand` - Handle check-in logic
- `CheckOutCommand` - Handle check-out logic
- `AbsenceReportQuery` - Retrieve reports

### 4. Infrastructure Layer (internal/infra)

#### Database (PostgreSQL)

**Tables:**
- `users` - User master data
- `work_schedules` - Expected work hours
- `holidays` - Public holidays
- `attendance_records` - Check-in/check-out events
- `absence_reports` - Denormalized reports
- `processed_events` - Event tracking
- `audit_events` - Audit trail

**Indexes:**
- `idx_users_email` - Fast user lookup
- `idx_attendance_records_user_date` - Efficient date range queries
- `idx_absence_reports_status` - Report filtering

#### Cache (Redis)

**Usage:**
- Distributed locks for check-in prevention
- Report caching
- Session management
- Rate limiting counters

**Key Patterns:**
```
lock:attendance:checkin:{user_id}:{date} - Check-in lock
report:cache:{user_id}:{date_range} - Report cache
session:{token} - Session data
```

#### Message Broker (Kafka)

**Topics:**
- `attendance.check-in` - Check-in events
- `attendance.check-out` - Check-out events

**Event Schema:**
```json
{
  "event_id": "uuid",
  "event_type": "attendance.check_in",
  "aggregate_id": "attendance_id",
  "aggregate_type": "attendance",
  "timestamp": 1234567890,
  "data": {
    "user_id": "uuid",
    "attendance_date": "2024-01-15",
    "check_in_time": "2024-01-15T08:30:00Z"
  },
  "metadata": {
    "user_id": "uuid",
    "source": "mobile"
  }
}
```

#### Logging (Loki)

**Log Levels:**
- `DEBUG` - Detailed debugging information
- `INFO` - General information
- `WARN` - Warning messages
- `ERROR` - Error messages

**Log Format:**
```json
{
  "timestamp": "2024-01-15T08:30:00Z",
  "level": "INFO",
  "message": "User checked in successfully",
  "user_id": "uuid",
  "attendance_id": "uuid",
  "labels": {
    "app": "mini-attendance",
    "env": "production",
    "host": "api-pod-1"
  }
}
```

## Data Flow

### Check-in Flow

```
1. User sends POST /api/v1/attendance/check-in
   ↓
2. API validates JWT token
   ↓
3. API acquires distributed lock (Redis)
   ↓
4. API checks for existing active check-in
   ↓
5. API begins database transaction
   ↓
6. API inserts/updates attendance_records
   ↓
7. API commits transaction
   ↓
8. API publishes CHECK_IN event to Kafka
   ↓
9. API releases lock
   ↓
10. API returns 201 Created response
   ↓
11. Report Service (Consumer) receives event
   ↓
12. Consumer checks idempotency (processed_events)
   ↓
13. Consumer updates absence_reports
   ↓
14. Consumer marks event as processed
```

### Check-out Flow

```
1. User sends POST /api/v1/attendance/check-out
   ↓
2. API validates JWT token
   ↓
3. API begins database transaction
   ↓
4. API finds active check-in record (FOR UPDATE)
   ↓
5. API updates check_out_at timestamp
   ↓
6. API commits transaction
   ↓
7. API publishes CHECK_OUT event to Kafka
   ↓
8. API returns 200 OK response
   ↓
9. Report Service (Consumer) receives event
   ↓
10. Consumer calculates work hours
   ↓
11. Consumer determines status (hadir/terlambat/pulang_cepat)
   ↓
12. Consumer updates absence_reports
```

## Security Architecture

### Authentication

- **JWT Tokens**: Stateless authentication
- **Access Token**: 1-hour expiration
- **Refresh Token**: 24-hour expiration
- **Token Validation**: On every protected endpoint

### Authorization

- **Role-Based Access Control (RBAC)**
  - `admin` - Full system access
  - `employee` - Personal attendance access
  - `manager` - Team attendance access (future)

### Data Protection

- **Password Hashing**: bcrypt with salt
- **Database Transactions**: ACID compliance
- **Distributed Locks**: Prevent race conditions
- **Audit Trail**: All changes logged

### Network Security

- **CORS**: Configured for frontend domain
- **HTTPS**: Enforced in production
- **TLS**: Encrypted connections
- **Rate Limiting**: Prevent abuse (future)

## Scalability Considerations

### Horizontal Scaling

**Stateless API Design:**
- No session state in memory
- All state in database/cache
- Multiple API instances possible

**Database Scaling:**
- Connection pooling (25 connections)
- Read replicas for reports
- Partitioning by date for large datasets

**Cache Scaling:**
- Redis cluster for high availability
- Cache invalidation strategy
- TTL-based expiration

**Message Queue Scaling:**
- Kafka partitioning by user_id
- Consumer groups for parallel processing
- Offset management for reliability

### Performance Optimization

**Database:**
- Indexes on frequently queried columns
- Query optimization
- Connection pooling
- Prepared statements

**Cache:**
- Report caching (1-6 hours)
- Lock caching (5 seconds)
- Cache warming on updates

**API:**
- Response compression
- Pagination for large datasets
- Lazy loading of related data

## Reliability & Resilience

### Error Handling

- **Graceful Degradation**: Service continues with reduced functionality
- **Retry Logic**: Exponential backoff for transient failures
- **Circuit Breaker**: Prevent cascading failures
- **Fallback Mechanisms**: Default responses

### Data Consistency

- **Database Transactions**: ACID properties
- **Event Idempotency**: Duplicate event handling
- **Distributed Locks**: Prevent concurrent modifications
- **Audit Trail**: Track all changes

### Monitoring & Alerting

- **Health Checks**: Regular service status
- **Metrics**: Performance monitoring
- **Logs**: Structured logging for debugging
- **Alerts**: Proactive issue detection

## Deployment Architecture

### Local Development

```
Docker Compose:
├── PostgreSQL (5432)
├── Redis (6379)
├── Kafka (9092)
├── Zookeeper (2181)
├── Loki (3100)
├── Grafana (3000)
└── Backend API (8080)
```

### Cloud Deployment

```
Kubernetes:
├── API Deployment (3 replicas)
├── PostgreSQL StatefulSet
├── Redis StatefulSet
├── Kafka StatefulSet
├── Loki Deployment
├── Grafana Deployment
├── ConfigMaps (configuration)
├── Secrets (credentials)
└── Services (networking)
```

## Technology Decisions

### Why Go?

- **Performance**: Fast execution, low memory footprint
- **Concurrency**: Goroutines for handling multiple requests
- **Simplicity**: Clean syntax, easy to maintain
- **Deployment**: Single binary, easy containerization

### Why PostgreSQL?

- **ACID Compliance**: Data consistency
- **JSON Support**: Flexible metadata storage
- **Indexing**: Efficient queries
- **Scalability**: Proven at scale

### Why Kafka?

- **Event Streaming**: Real-time data processing
- **Durability**: Persistent message storage
- **Scalability**: Horizontal scaling
- **Ordering**: Per-partition message ordering

### Why Redis?

- **Performance**: In-memory data store
- **Distributed Locks**: Atomic operations
- **Caching**: Fast data retrieval
- **TTL Support**: Automatic expiration

## Future Enhancements

1. **Report Service Implementation**
   - Full consumer implementation
   - Advanced report generation
   - Scheduled report generation

2. **Frontend Application**
   - Next.js web application
   - Mobile-responsive design
   - Real-time updates

3. **Advanced Features**
   - Geolocation tracking
   - Biometric authentication
   - Overtime management
   - Leave management

4. **Observability**
   - Distributed tracing (Jaeger)
   - Metrics collection (Prometheus)
   - Custom dashboards

5. **Performance**
   - Database query optimization
   - Caching strategies
   - API rate limiting

## Conclusion

The Mini Attendance System is designed with scalability, reliability, and maintainability in mind. The event-driven architecture allows for asynchronous processing and future expansion, while the comprehensive monitoring ensures system health and performance.
