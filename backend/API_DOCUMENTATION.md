# API Documentation

## Overview

This document describes all API endpoints available in the TrustMedis Attendance Management System. The API uses JWT-based authentication and follows RESTful principles.

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** Include the JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

---

## Table of Contents

1. [Authentication](#authentication)
2. [Attendance Management](#attendance-management)
3. [Reports](#reports)
4. [Leave Management](#leave-management)
5. [Overtime Management](#overtime-management)
6. [Absence Management](#absence-management)
7. [User Management](#user-management)
8. [Error Responses](#error-responses)

---

## Authentication

### POST /auth/login

Login with username and password to receive JWT tokens.

**Request:**

```json
{
	"username": "string",
	"password": "string"
}
```

**Response (200):**

```json
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	"expires_in": 3600,
	"token_type": "Bearer",
	"user": {
		"id": "550e8400-e29b-41d4-a716-446655440001",
		"username": "admin",
		"email": "admin@trustmedis.com",
		"full_name": "Admin User",
		"department": "IT",
		"role": "admin"
	}
}
```

**Error Responses:**

- `400 Bad Request`: Missing username or password
- `401 Unauthorized`: Invalid credentials

---

### POST /auth/refresh

Refresh the access token using a refresh token.

**Request:**

```
Header: Authorization: Bearer <refresh_token>
```

**Response (200):**

```json
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	"expires_in": 3600,
	"token_type": "Bearer"
}
```

**Error Responses:**

- `401 Unauthorized`: Invalid or expired refresh token

---

### POST /auth/logout

Logout and invalidate the current session.

**Request:**

```
Header: Authorization: Bearer <access_token>
```

**Response (200):**

```json
{
	"message": "Logged out successfully"
}
```

**Error Responses:**

- `401 Unauthorized`: No token provided

---

## Attendance Management

### POST /attendance/check-in

Record employee check-in with geolocation data.

**Request:**

```json
{
	"latitude": -6.1753,
	"longitude": 106.8249,
	"device": "mobile"
}
```

**Response (201):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440101",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"attendance_date": "2024-01-15",
	"check_in_time": "2024-01-15T08:30:00Z",
	"check_out_time": null,
	"status": "present"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid coordinates or missing required fields
- `401 Unauthorized`: Not authenticated
- `409 Conflict`: Already checked in today

---

### POST /attendance/check-out

Record employee check-out.

**Request:**

```json
{
	"latitude": -6.1753,
	"longitude": 106.8249,
	"device": "mobile"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440101",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"attendance_date": "2024-01-15",
	"check_in_time": "2024-01-15T08:30:00Z",
	"check_out_time": "2024-01-15T17:00:00Z",
	"status": "present"
}
```

**Error Responses:**

- `400 Bad Request`: Not checked in yet
- `401 Unauthorized`: Not authenticated

---

### GET /attendance/today

Get current day's attendance record for logged-in user.

**Query Parameters:**

- None

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440101",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"attendance_date": "2024-01-15",
	"check_in_time": "2024-01-15T08:30:00Z",
	"check_out_time": null,
	"status": "present"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `404 Not Found`: No record for today

---

### GET /attendance/history

Get attendance history for logged-in user.

**Query Parameters:**

- `from` (required): Start date (YYYY-MM-DD)
- `to` (required): End date (YYYY-MM-DD)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Records per page (default: 20, max: 100)

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440101",
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"attendance_date": "2024-01-15",
			"check_in_time": "2024-01-15T08:30:00Z",
			"check_out_time": "2024-01-15T17:00:00Z",
			"status": "present"
		}
	],
	"pagination": {
		"total": 22,
		"page": 1,
		"limit": 20,
		"total_pages": 2
	}
}
```

**Error Responses:**

- `400 Bad Request`: Invalid date format or invalid date range
- `401 Unauthorized`: Not authenticated

---

### GET /attendance/users/:userId

Get attendance records for a specific user. **Requires Manager or Admin role.**

**Path Parameters:**

- `userId`: UUID of the target user

**Query Parameters:**

- `from` (required): Start date (YYYY-MM-DD)
- `to` (required): End date (YYYY-MM-DD)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Records per page (default: 20, max: 100)

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440101",
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"attendance_date": "2024-01-15",
			"check_in_time": "2024-01-15T08:30:00Z",
			"check_out_time": "2024-01-15T17:00:00Z",
			"status": "present"
		}
	],
	"pagination": {
		"total": 22,
		"page": 1,
		"limit": 20,
		"total_pages": 2
	}
}
```

**Error Responses:**

- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: User not found

---

## Reports

### GET /reports/monthly

Get monthly attendance summary. **Requires Admin role.**

**Query Parameters:**

- `year` (required): Year (YYYY)
- `month` (required): Month (1-12)

**Response (200):**

```json
{
	"year": 2024,
	"month": 1,
	"total_employees": 10,
	"present": 230,
	"absent": 5,
	"late": 15,
	"summary_by_status": {
		"present": {
			"count": 230,
			"percentage": 88.46
		},
		"absent": {
			"count": 5,
			"percentage": 1.92
		},
		"late": {
			"count": 15,
			"percentage": 5.77
		},
		"leave": {
			"count": 10,
			"percentage": 3.85
		}
	}
}
```

**Error Responses:**

- `400 Bad Request`: Invalid month or year
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can access

---

### GET /reports/users/:userId/monthly

Get monthly attendance summary for a specific user. **Requires Manager or Admin role.**

**Path Parameters:**

- `userId`: UUID of the target user

**Query Parameters:**

- `year` (required): Year (YYYY)
- `month` (required): Month (1-12)

**Response (200):**

```json
{
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"user_name": "John Doe",
	"year": 2024,
	"month": 1,
	"total_workdays": 22,
	"present": 20,
	"absent": 1,
	"late": 1,
	"leave": 0,
	"overtime_hours": 8.5,
	"avg_check_in_time": "08:25:00",
	"avg_check_out_time": "17:15:00"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: User not found

---

### GET /reports/summary

Get overall attendance summary across all employees. **Requires Admin role.**

**Query Parameters:**

- `from` (required): Start date (YYYY-MM-DD)
- `to` (required): End date (YYYY-MM-DD)

**Response (200):**

```json
{
	"period": {
		"from": "2024-01-01",
		"to": "2024-01-31"
	},
	"total_employees": 25,
	"total_records": 550,
	"attendance_stats": {
		"present": 480,
		"absent": 20,
		"late": 35,
		"leave": 15
	},
	"by_department": [
		{
			"department": "IT",
			"total_employees": 8,
			"present": 152,
			"absent": 3,
			"late": 10,
			"leave": 5
		}
	]
}
```

**Error Responses:**

- `400 Bad Request`: Invalid date format
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can access

---

### GET /reports/department/:department

Get department-level attendance report. **Requires Admin role.**

**Path Parameters:**

- `department`: Department name or code

**Query Parameters:**

- `from` (required): Start date (YYYY-MM-DD)
- `to` (required): End date (YYYY-MM-DD)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Records per page (default: 50, max: 200)

**Response (200):**

```json
{
	"department": "IT",
	"period": {
		"from": "2024-01-01",
		"to": "2024-01-31"
	},
	"employees": [
		{
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"full_name": "John Doe",
			"present": 20,
			"absent": 1,
			"late": 1,
			"leave": 0
		}
	],
	"pagination": {
		"total": 8,
		"page": 1,
		"limit": 50,
		"total_pages": 1
	}
}
```

**Error Responses:**

- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can access
- `404 Not Found`: Department not found

---

### GET /reports/export

Export attendance data in various formats. **Requires Admin role.**

**Query Parameters:**

- `format` (required): Export format (csv, xlsx, pdf)
- `from` (required): Start date (YYYY-MM-DD)
- `to` (required): End date (YYYY-MM-DD)
- `include_departments` (optional): Comma-separated department list

**Response (200):**

- Content-Type: application/csv, application/vnd.ms-excel, or application/pdf
- Returns binary file data

**Error Responses:**

- `400 Bad Request`: Invalid format or parameters
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can access

---

## Leave Management

### POST /leave/request

Submit a leave request.

**Request:**

```json
{
	"start_date": "2024-02-01",
	"end_date": "2024-02-05",
	"leave_type": "annual",
	"reason": "Vacation"
}
```

**Response (201):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440201",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"start_date": "2024-02-01",
	"end_date": "2024-02-05",
	"leave_type": "annual",
	"reason": "Vacation",
	"status": "pending",
	"created_at": "2024-01-20T10:30:00Z",
	"approved_by": null
}
```

**Error Responses:**

- `400 Bad Request`: Invalid dates or leave type
- `401 Unauthorized`: Not authenticated
- `409 Conflict`: Overlapping leave request exists

---

### GET /leave/my-requests

Get current user's leave requests.

**Query Parameters:**

- `status` (optional): Filter by status (pending, approved, rejected)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Records per page (default: 20)

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440201",
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"start_date": "2024-02-01",
			"end_date": "2024-02-05",
			"leave_type": "annual",
			"reason": "Vacation",
			"status": "pending",
			"created_at": "2024-01-20T10:30:00Z",
			"approved_by": null
		}
	],
	"pagination": {
		"total": 5,
		"page": 1,
		"limit": 20,
		"total_pages": 1
	}
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated

---

### GET /leave/requests/:requestId

Get details of a specific leave request.

**Path Parameters:**

- `requestId`: UUID of the leave request

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440201",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"user": {
		"full_name": "John Doe",
		"email": "john.doe@trustmedis.com"
	},
	"start_date": "2024-02-01",
	"end_date": "2024-02-05",
	"leave_type": "annual",
	"reason": "Vacation",
	"status": "pending",
	"created_at": "2024-01-20T10:30:00Z",
	"approved_by": null,
	"rejection_reason": null
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `404 Not Found`: Request not found

---

### PUT /leave/requests/:requestId/approve

Approve a leave request. **Requires Manager or Admin role.**

**Path Parameters:**

- `requestId`: UUID of the leave request

**Request:**

```json
{
	"notes": "Approved"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440201",
	"status": "approved",
	"approved_at": "2024-01-20T11:00:00Z",
	"approved_by": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only manager/admin can approve
- `404 Not Found`: Request not found
- `409 Conflict`: Request already processed

---

### PUT /leave/requests/:requestId/reject

Reject a leave request. **Requires Manager or Admin role.**

**Path Parameters:**

- `requestId`: UUID of the leave request

**Request:**

```json
{
	"reason": "Maximum leave quota exceeded"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440201",
	"status": "rejected",
	"rejection_reason": "Maximum leave quota exceeded",
	"rejected_at": "2024-01-20T11:00:00Z",
	"rejected_by": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only manager/admin can reject
- `404 Not Found`: Request not found
- `409 Conflict`: Request already processed

---

## Overtime Management

### POST /overtime/request

Request overtime approval.

**Request:**

```json
{
	"date": "2024-01-20",
	"hours": 2.5,
	"reason": "Project deadline"
}
```

**Response (201):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440301",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"date": "2024-01-20",
	"hours": 2.5,
	"reason": "Project deadline",
	"status": "pending",
	"created_at": "2024-01-20T10:30:00Z"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid data
- `401 Unauthorized`: Not authenticated

---

### GET /overtime/my-records

Get current user's overtime records.

**Query Parameters:**

- `status` (optional): Filter by status (pending, approved, rejected)
- `from` (optional): Start date
- `to` (optional): End date

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440301",
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"date": "2024-01-20",
			"hours": 2.5,
			"reason": "Project deadline",
			"status": "approved",
			"created_at": "2024-01-20T10:30:00Z"
		}
	]
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated

---

### PUT /overtime/:overtimeId/approve

Approve overtime request. **Requires Manager or Admin role.**

**Path Parameters:**

- `overtimeId`: UUID of the overtime record

**Request:**

```json
{
	"notes": "Approved"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440301",
	"status": "approved",
	"approved_at": "2024-01-20T11:00:00Z"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only manager/admin can approve
- `404 Not Found`: Record not found

---

## Absence Management

### POST /absence/record

Record an employee absence. **Requires Manager or Admin role.**

**Request:**

```json
{
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"date": "2024-01-20",
	"type": "unplanned",
	"reason": "Emergency"
}
```

**Response (201):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440401",
	"user_id": "550e8400-e29b-41d4-a716-446655440002",
	"date": "2024-01-20",
	"type": "unplanned",
	"reason": "Emergency",
	"created_at": "2024-01-20T10:30:00Z"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only manager/admin can record

---

### GET /absence/user/:userId

Get absence records for a user. **Requires Manager or Admin role.**

**Path Parameters:**

- `userId`: UUID of the target user

**Query Parameters:**

- `from` (optional): Start date
- `to` (optional): End date
- `type` (optional): Filter by type (planned, unplanned)

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440401",
			"user_id": "550e8400-e29b-41d4-a716-446655440002",
			"date": "2024-01-20",
			"type": "unplanned",
			"reason": "Emergency",
			"created_at": "2024-01-20T10:30:00Z"
		}
	]
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only manager/admin can view
- `404 Not Found`: User not found

---

## User Management

### GET /users/profile

Get current user's profile.

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440002",
	"username": "john_doe",
	"email": "john.doe@trustmedis.com",
	"full_name": "John Doe",
	"department": "HR",
	"role": "employee",
	"status": "active",
	"created_at": "2024-01-01T00:00:00Z",
	"last_login": "2024-01-15T08:30:00Z"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated

---

### PUT /users/profile

Update current user's profile.

**Request:**

```json
{
	"full_name": "John Doe Updated",
	"email": "john.doe.new@trustmedis.com"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440002",
	"username": "john_doe",
	"email": "john.doe.new@trustmedis.com",
	"full_name": "John Doe Updated",
	"department": "HR",
	"role": "employee",
	"status": "active",
	"updated_at": "2024-01-20T11:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid data
- `401 Unauthorized`: Not authenticated

---

### PUT /users/password

Change current user's password.

**Request:**

```json
{
	"current_password": "oldpassword",
	"new_password": "newpassword",
	"confirm_password": "newpassword"
}
```

**Response (200):**

```json
{
	"message": "Password changed successfully"
}
```

**Error Responses:**

- `400 Bad Request`: Passwords don't match
- `401 Unauthorized`: Current password is incorrect

---

### GET /users

Get list of all users. **Requires Admin role.**

**Query Parameters:**

- `department` (optional): Filter by department
- `role` (optional): Filter by role
- `status` (optional): Filter by status (active, inactive)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Records per page (default: 20, max: 100)

**Response (200):**

```json
{
	"data": [
		{
			"id": "550e8400-e29b-41d4-a716-446655440002",
			"username": "john_doe",
			"email": "john.doe@trustmedis.com",
			"full_name": "John Doe",
			"department": "HR",
			"role": "employee",
			"status": "active",
			"created_at": "2024-01-01T00:00:00Z"
		}
	],
	"pagination": {
		"total": 25,
		"page": 1,
		"limit": 20,
		"total_pages": 2
	}
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can access

---

### POST /users

Create a new user. **Requires Admin role.**

**Request:**

```json
{
	"username": "new_user",
	"email": "new.user@trustmedis.com",
	"password": "password123",
	"full_name": "New User",
	"department": "Finance",
	"role": "employee"
}
```

**Response (201):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440050",
	"username": "new_user",
	"email": "new.user@trustmedis.com",
	"full_name": "New User",
	"department": "Finance",
	"role": "employee",
	"status": "active",
	"created_at": "2024-01-20T12:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid data or username already exists
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can create users

---

### GET /users/:userId

Get user details. **Requires Admin role or user viewing themselves.**

**Path Parameters:**

- `userId`: UUID of the target user

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440002",
	"username": "john_doe",
	"email": "john.doe@trustmedis.com",
	"full_name": "John Doe",
	"department": "HR",
	"role": "employee",
	"status": "active",
	"created_at": "2024-01-01T00:00:00Z",
	"last_login": "2024-01-15T08:30:00Z"
}
```

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: User not found

---

### PUT /users/:userId

Update user details. **Requires Admin role.**

**Path Parameters:**

- `userId`: UUID of the target user

**Request:**

```json
{
	"full_name": "John Doe",
	"department": "Finance",
	"role": "manager",
	"status": "active"
}
```

**Response (200):**

```json
{
	"id": "550e8400-e29b-41d4-a716-446655440002",
	"username": "john_doe",
	"email": "john.doe@trustmedis.com",
	"full_name": "John Doe",
	"department": "Finance",
	"role": "manager",
	"status": "active",
	"updated_at": "2024-01-20T12:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can update users
- `404 Not Found`: User not found

---

### DELETE /users/:userId

Delete a user. **Requires Admin role.**

**Path Parameters:**

- `userId`: UUID of the target user

**Response (204):**
No content

**Error Responses:**

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Only admin can delete users
- `404 Not Found`: User not found

---

## Error Responses

All error responses follow this format:

```json
{
	"error": {
		"code": "ERROR_CODE",
		"message": "Human-readable error message",
		"details": {
			"field": "Additional context"
		}
	},
	"timestamp": "2024-01-20T12:00:00Z",
	"path": "/api/v1/endpoint"
}
```

### Common HTTP Status Codes

| Status | Meaning                                                 |
| ------ | ------------------------------------------------------- |
| 200    | OK - Request successful                                 |
| 201    | Created - Resource created successfully                 |
| 204    | No Content - Successful request with no response body   |
| 400    | Bad Request - Invalid request data or parameters        |
| 401    | Unauthorized - Authentication required or invalid token |
| 403    | Forbidden - User lacks permission for this resource     |
| 404    | Not Found - Resource not found                          |
| 409    | Conflict - Request conflicts with current state         |
| 500    | Internal Server Error - Server error occurred           |

### Common Error Codes

- `INVALID_REQUEST`: Request data is malformed or missing required fields
- `AUTHENTICATION_FAILED`: Authentication failed (invalid credentials)
- `INVALID_TOKEN`: JWT token is invalid or expired
- `UNAUTHORIZED`: User is not authenticated
- `FORBIDDEN`: User lacks required permissions
- `NOT_FOUND`: Requested resource not found
- `CONFLICT`: Request conflicts with existing data
- `VALIDATION_ERROR`: Data validation failed
- `INTERNAL_ERROR`: Internal server error

---

## Rate Limiting

API endpoints are rate-limited to prevent abuse:

- **Default limit**: 100 requests per minute per user
- **Auth endpoints**: 5 requests per minute per IP address

Rate limit information is included in response headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1705699200
```

---

## Webhook Events

The system publishes events for critical operations. Subscribe to these events:

### Events Published

- `attendance.check_in`: User checked in
- `attendance.check_out`: User checked out
- `leave.request_created`: Leave request submitted
- `leave.request_approved`: Leave request approved
- `leave.request_rejected`: Leave request rejected
- `overtime.request_created`: Overtime request submitted
- `overtime.request_approved`: Overtime request approved

**Event Payload Format:**

```json
{
	"event": "attendance.check_in",
	"timestamp": "2024-01-20T12:00:00Z",
	"data": {
		"user_id": "550e8400-e29b-41d4-a716-446655440002",
		"attendance_id": "550e8400-e29b-41d4-a716-446655440101",
		"check_in_time": "2024-01-20T08:30:00Z"
	}
}
```

---

## Example cURL Commands

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Check In

```bash
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "latitude": -6.1753,
    "longitude": 106.8249,
    "device": "mobile"
  }'
```

### Get Attendance History

```bash
curl -X GET "http://localhost:8080/api/v1/attendance/history?from=2024-01-01&to=2024-01-31" \
  -H "Authorization: Bearer <access_token>"
```

### Get Monthly Report

```bash
curl -X GET "http://localhost:8080/api/v1/reports/monthly?year=2024&month=1" \
  -H "Authorization: Bearer <access_token>"
```

### Request Leave

```bash
curl -X POST http://localhost:8080/api/v1/leave/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "start_date": "2024-02-01",
    "end_date": "2024-02-05",
    "leave_type": "annual",
    "reason": "Vacation"
  }'
```

---

## API Versioning

Current API Version: **v1**

The API follows semantic versioning. Major version changes (v1 → v2) indicate breaking changes. Minor and patch updates are backward compatible.

---

## Support

For API support and issues, contact: `api-support@trustmedis.com`

For detailed implementation examples, refer to the test files in `internal/api/http/handlers/`.
