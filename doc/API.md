# Backend Portfolio API Documentation

## Overview

Backend Portfolio is a RESTful API built with Go and Gin framework for managing a personal portfolio website. It provides endpoints for managing portfolio items, skills, qualifications, authentication, and visitor tracking.

**Framework:** Gin 1.9.1  
**Database:** PostgreSQL with GORM ORM  
**Go Version:** 1.21+  
**API Version:** 2.0.0

---

## Table of Contents

1. [Authentication](#authentication)
2. [Base URL & Headers](#base-url--headers)
3. [Public Routes](#public-routes)
4. [Protected Routes (Admin Only)](#protected-routes-admin-only)
5. [Error Handling](#error-handling)
6. [Data Models](#data-models)
7. [Example Usage](#example-usage)

---

## Authentication

This API uses **JWT (JSON Web Tokens)** for authentication with the following mechanisms:

### Cookie-Based Authentication
- Access tokens and refresh tokens are stored in **HTTP-only cookies**
- CSRF protection is enabled for protected routes
- Cookies are automatically set on login and cleared on logout

### Token Expiration
- **Access Token**: Short-lived (default: 15 minutes)
- **Refresh Token**: Long-lived (default: 7 days)

### CSRF Protection
- Protected routes require an `X-XSRF-TOKEN` header
- CSRF token is automatically set in cookies upon login

---

## Base URL & Headers

### Base URL
```
http://localhost:8080
```
Or via Docker:
```
docker compose up -d
# Server runs on configured PORT (default: 8080)
```

### Required Headers for Protected Routes
```
Content-Type: application/json
X-XSRF-TOKEN: <csrf-token-from-cookie>
Authorization: Bearer <access-token> (optional, cookie takes precedence)
```

### CORS Configuration
- Allowed Origins: Configurable via `ALLOWED_ORIGINS` env variable
- Allowed Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
- Credentials: Allowed
- Max Age: 12 hours

---

## Public Routes

### 1. Health Check
**Endpoint:** `GET /health`

**Description:** Verifies database connection and server status.

**Response:**
```json
{
  "status": "OK",
  "message": "Server is running"
}
```

**Status Codes:** 
- `200 OK` - Server healthy
- `500 Internal Server Error` - Database connection failed

---

### 2. API Info
**Endpoint:** `GET /`

**Description:** Returns basic API information.

**Response:**
```json
{
  "message": "Backend Portfolio API is running!",
  "version": "2.0.0"
}
```

---

### 3. Get About Information
**Endpoint:** `GET /api/about`

**Description:** Retrieves the portfolio owner's about information.

**Response:**
```json
{
  "id": 1,
  "bio": "Full-stack developer with 5+ years experience",
  "summary": "Passionate about building scalable web applications",
  "email": "user@example.com",
  "phone": "+1234567890",
  "location": "San Francisco, CA",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Success
- `404 Not Found` - About information not found

---

### 4. Get All Portfolio Items
**Endpoint:** `GET /api/portfolio`

**Description:** Retrieves all portfolio items with associated media files.

**Query Parameters:** None

**Response:**
```json
[
  {
    "id": 1,
    "title": "E-commerce Platform",
    "description": "Full-stack e-commerce application",
    "category": "Web Development",
    "tags": "Go,React,PostgreSQL",
    "project_url": "https://github.com/user/ecommerce",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "media_files": [
      {
        "id": 1,
        "portfolio_id": 1,
        "type": "image",
        "url": "/uploads/portfolio/1_1234567890_screenshot.png",
        "order_index": 0
      }
    ]
  }
]
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

### 5. Get Single Portfolio Item
**Endpoint:** `GET /api/portfolio/:id`

**Description:** Retrieves a specific portfolio item by ID.

**Path Parameters:**
- `id` (uint) - Portfolio item ID

**Response:** Single portfolio item object (see Get All Portfolio Items)

**Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Portfolio not found

---

### 6. Get All Skills
**Endpoint:** `GET /api/skills`

**Description:** Retrieves all skills without grouping.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Go",
    "level": "Expert",
    "score": 95,
    "category": "Backend",
    "icon": "go-icon",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "name": "React",
    "level": "Advanced",
    "score": 85,
    "category": "Frontend",
    "icon": "react-icon",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

### 7. Get Single Skill
**Endpoint:** `GET /api/skills/:id`

**Description:** Retrieves a specific skill by ID.

**Path Parameters:**
- `id` (uint) - Skill ID

**Response:** Single skill object (see Get All Skills)

**Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Skill not found

---

### 8. Get Skills by Category
**Endpoint:** `GET /api/skills/category`

**Description:** Retrieves all skills grouped by category.

**Response:**
```json
{
  "Backend": [
    {
      "id": 1,
      "name": "Go",
      "level": "Expert",
      "score": 95,
      "category": "Backend",
      "icon": "go-icon"
    }
  ],
  "Frontend": [
    {
      "id": 2,
      "name": "React",
      "level": "Advanced",
      "score": 85,
      "category": "Frontend",
      "icon": "react-icon"
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

### 9. Get All Qualifications
**Endpoint:** `GET /api/qualifications`

**Description:** Retrieves all qualifications (degrees, certifications, etc.).

**Response:**
```json
[
  {
    "id": 1,
    "type": "Degree",
    "title": "Bachelor of Science in Computer Science",
    "issuer": "University of California",
    "issue_date": "2020-05-15",
    "expiry_date": null,
    "credential_url": "https://example.com/cert",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

### 10. Get Single Qualification
**Endpoint:** `GET /api/qualifications/:id`

**Description:** Retrieves a specific qualification by ID.

**Path Parameters:**
- `id` (uint) - Qualification ID

**Response:** Single qualification object (see Get All Qualifications)

**Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Qualification not found

---

### 11. Track Visitor
**Endpoint:** `POST /api/visitors/track`

**Description:** Records a visitor hit (called by React frontend on every page load).

**Request Body:**
```json
{
  "path": "/about",
  "referer": "https://google.com"
}
```

**Notes:**
- Both `path` and `referer` are optional
- If `path` is not provided, defaults to `/`
- Referer from body takes precedence over HTTP Referer header
- IP address is extracted from client IP
- GeoIP lookup is performed (if configured) to get country/city
- User agent is captured from request headers

**Response:**
```json
{
  "message": "Visitor recorded"
}
```

**Status Codes:**
- `200 OK` - Visitor recorded
- `500 Internal Server Error` - Database error

---

### 12. Get Visitor Stats (Public)
**Endpoint:** `GET /api/visitors/stats`

**Description:** Retrieves public visitor statistics (no IP addresses or user agents).

**Response:**
```json
{
  "today": 42,
  "this_week": 215,
  "this_month": 980,
  "this_year": 5432,
  "total": 12845,
  "unique_today": 38,
  "unique_this_week": 180,
  "unique_this_month": 850,
  "unique_this_year": 4200,
  "unique_total": 10000,
  "daily_chart": [
    { "date": "2024-01-15", "count": 42 },
    { "date": "2024-01-14", "count": 51 }
  ],
  "monthly_chart": [
    { "month": "2024-01", "count": 980 }
  ],
  "top_pages": [
    { "path": "/", "count": 500 },
    { "path": "/portfolio", "count": 450 }
  ]
}
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

## Protected Routes (Admin Only)

### Authentication Requirement
All protected routes require:
1. Valid JWT token in cookies (automatic after login)
2. CSRF token in `X-XSRF-TOKEN` header

### 1. Login
**Endpoint:** `POST /api/auth/login` or `POST /api/login` (backward compatibility)

**Description:** Authenticates a user and sets HTTP-only cookies with tokens.

**Request Body:**
```json
{
  "username": "admin",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Status Codes:**
- `200 OK` - Login successful
- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Invalid credentials

**Side Effects:**
- Sets `access_token` HTTP-only cookie
- Sets `refresh_token` HTTP-only cookie
- Sets `csrf_token` cookie

---

### 2. Logout
**Endpoint:** `POST /api/auth/logout`

**Description:** Clears authentication and CSRF cookies.

**Request Body:** Empty or no body required

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

**Status Codes:**
- `200 OK` - Logout successful

---

### 3. Refresh Token
**Endpoint:** `POST /api/auth/refresh`

**Description:** Validates refresh token and issues new token pair (token rotation).

**Request Body:** Empty (uses refresh token from cookie)

**Response:**
```json
{
  "message": "Token refreshed successfully"
}
```

**Status Codes:**
- `200 OK` - Tokens refreshed
- `401 Unauthorized` - Invalid or expired refresh token
- `500 Internal Server Error` - Token generation failed

**Side Effects:**
- New access token in `access_token` cookie
- New refresh token in `refresh_token` cookie
- New CSRF token in `csrf_token` cookie

---

### 4. Get Current User
**Endpoint:** `GET /api/auth/me`

**Description:** Returns information about the currently authenticated user.

**Request Body:** None

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Status Codes:**
- `200 OK` - Success
- `401 Unauthorized` - Not authenticated
- `404 Not Found` - User not found

---

### 5. Create User
**Endpoint:** `POST /api/admin/create-user`

**Description:** Creates a new user account (admin only).

**Request Body:**
```json
{
  "username": "newuser",
  "password": "securepassword123",
  "role": "admin"
}
```

**Validation:**
- `username`: Required
- `password`: Required, minimum 8 characters
- `role`: Optional, defaults to "admin"

**Response:**
```json
{
  "message": "User created successfully",
  "user": {
    "id": 2,
    "username": "newuser",
    "role": "admin"
  }
}
```

**Status Codes:**
- `201 Created` - User created
- `400 Bad Request` - Invalid request
- `500 Internal Server Error` - Creation failed

---

### 6. Reset Admin Password
**Endpoint:** `POST /api/admin/reset-admin`

**Description:** Resets the admin password using a secret key.

**Request Headers:**
```
X-Reset-Secret: <reset-secret-from-env>
```

**Request Body:**
```json
{
  "new_password": "newsecurepassword123"
}
```

**Validation:**
- `new_password`: Required, minimum 8 characters
- `X-Reset-Secret` header must match `RESET_SECRET` environment variable

**Response:**
```json
{
  "message": "Admin password reset successfully"
}
```

**Status Codes:**
- `200 OK` - Password reset
- `400 Bad Request` - Invalid request
- `403 Forbidden` - Reset not configured or invalid secret
- `500 Internal Server Error` - Reset failed

---

### 7. Create or Update About
**Endpoint:** `POST /api/admin/about`

**Description:** Creates a new about entry or updates the existing one.

**Request Body:**
```json
{
  "bio": "Full-stack developer with 5+ years experience",
  "summary": "Passionate about building scalable web applications",
  "email": "user@example.com",
  "phone": "+1234567890",
  "location": "San Francisco, CA"
}
```

**Response:**
```json
{
  "id": 1,
  "bio": "Full-stack developer with 5+ years experience",
  "summary": "Passionate about building scalable web applications",
  "email": "user@example.com",
  "phone": "+1234567890",
  "location": "San Francisco, CA",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - New entry created
- `200 OK` - Existing entry updated
- `400 Bad Request` - Invalid request
- `500 Internal Server Error` - Operation failed

---

### 8. Update About
**Endpoint:** `PUT /api/admin/about/:id`

**Description:** Updates an existing about entry by ID.

**Path Parameters:**
- `id` (uint) - About entry ID

**Request Body:** (same as Create or Update About)

**Response:** (same as Create or Update About)

**Status Codes:**
- `200 OK` - Updated
- `400 Bad Request` - Invalid request
- `404 Not Found` - About entry not found
- `500 Internal Server Error` - Update failed

---

### 9. Create Portfolio Item
**Endpoint:** `POST /api/admin/portfolio`

**Content-Type:** `multipart/form-data`

**Description:** Creates a new portfolio item with optional media files.

**Form Parameters:**
```
title (required): string - Portfolio title
description: string - Portfolio description
category: string - Project category
tags: string - Comma-separated tags
project_url: string - URL to project
media_files: file[] - Optional image/video files (max 32 MB total)
```

**Response:**
```json
{
  "id": 1,
  "title": "E-commerce Platform",
  "description": "Full-stack e-commerce application",
  "category": "Web Development",
  "tags": "Go,React,PostgreSQL",
  "project_url": "https://github.com/user/ecommerce",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "media_files": [
    {
      "id": 1,
      "portfolio_id": 1,
      "type": "image",
      "url": "/uploads/portfolio/1_1234567890_screenshot.png",
      "order_index": 0
    }
  ]
}
```

**Status Codes:**
- `201 Created` - Portfolio created
- `400 Bad Request` - Missing title or invalid form
- `500 Internal Server Error` - Creation failed

---

### 10. Update Portfolio Item
**Endpoint:** `PUT /api/admin/portfolio/:id`

**Content-Type:** `multipart/form-data` or `application/json`

**Description:** Updates a portfolio item (supports both JSON and multipart).

**Path Parameters:**
- `id` (uint) - Portfolio item ID

**Form/JSON Parameters:** (same as Create Portfolio Item, all optional)

**Response:** (same as Create Portfolio Item)

**Status Codes:**
- `200 OK` - Updated
- `400 Bad Request` - Invalid request
- `404 Not Found` - Portfolio not found
- `500 Internal Server Error` - Update failed

**Notes:**
- When using multipart, any existing media files are kept
- New media files are added (not replaced)
- Use Delete Portfolio Media endpoint to remove files

---

### 11. Delete Portfolio Item
**Endpoint:** `DELETE /api/admin/portfolio/:id`

**Description:** Removes a portfolio item and all associated media files from storage.

**Path Parameters:**
- `id` (uint) - Portfolio item ID

**Response:**
```json
{
  "message": "Portfolio item deleted successfully"
}
```

**Status Codes:**
- `200 OK` - Deleted
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Portfolio not found
- `500 Internal Server Error` - Deletion failed

---

### 12. Delete Portfolio Media File
**Endpoint:** `DELETE /api/admin/portfolio-media/:portfolio_id/:media_id`

**Description:** Removes a specific media file from a portfolio.

**Path Parameters:**
- `portfolio_id` (uint) - Portfolio item ID
- `media_id` (uint) - Media file ID

**Response:**
```json
{
  "message": "Media file deleted successfully"
}
```

**Status Codes:**
- `200 OK` - Deleted
- `400 Bad Request` - Invalid ID
- `404 Not Found` - Media not found
- `500 Internal Server Error` - Deletion failed

---

### 13. Create Skill
**Endpoint:** `POST /api/admin/skills`

**Content-Type:** `application/json`

**Description:** Creates a new skill with validation.

**Request Body:**
```json
{
  "name": "Go",
  "level": "Expert",
  "score": 95,
  "category": "Backend",
  "icon": "go-icon"
}
```

**Validation:**
- `name`: Required
- `level`: Required, must be one of: Beginner, Intermediate, Advanced, Expert
- `score`: Required, must be between 0 and 100
- `category`: Required
- `icon`: Optional

**Response:**
```json
{
  "id": 1,
  "name": "Go",
  "level": "Expert",
  "score": 95,
  "category": "Backend",
  "icon": "go-icon",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Skill created
- `400 Bad Request` - Validation failed
- `500 Internal Server Error` - Creation failed

---

### 14. Update Skill
**Endpoint:** `PUT /api/admin/skills/:id`

**Description:** Updates an existing skill (all fields optional).

**Path Parameters:**
- `id` (uint) - Skill ID

**Request Body:** (all fields from Create Skill are optional)

**Response:** (same as Create Skill)

**Status Codes:**
- `200 OK` - Updated
- `400 Bad Request` - Validation failed
- `404 Not Found` - Skill not found
- `500 Internal Server Error` - Update failed

---

### 15. Delete Skill
**Endpoint:** `DELETE /api/admin/skills/:id`

**Description:** Removes a skill by ID.

**Path Parameters:**
- `id` (uint) - Skill ID

**Response:**
```json
{
  "message": "Skill deleted successfully"
}
```

**Status Codes:**
- `200 OK` - Deleted
- `400 Bad Request` - Invalid ID
- `500 Internal Server Error` - Deletion failed

---

### 16. Create Qualification
**Endpoint:** `POST /api/admin/qualifications`

**Content-Type:** `application/json`

**Description:** Creates a new qualification.

**Request Body:**
```json
{
  "type": "Degree",
  "title": "Bachelor of Science in Computer Science",
  "issuer": "University of California",
  "issue_date": "2020-05-15",
  "expiry_date": null,
  "credential_url": "https://example.com/cert"
}
```

**Response:**
```json
{
  "id": 1,
  "type": "Degree",
  "title": "Bachelor of Science in Computer Science",
  "issuer": "University of California",
  "issue_date": "2020-05-15",
  "expiry_date": null,
  "credential_url": "https://example.com/cert",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Qualification created
- `400 Bad Request` - Invalid request
- `500 Internal Server Error` - Creation failed

---

### 17. Update Qualification
**Endpoint:** `PUT /api/admin/qualifications/:id`

**Description:** Updates an existing qualification.

**Path Parameters:**
- `id` (uint) - Qualification ID

**Request Body:** (all fields from Create Qualification, optional)

**Response:** (same as Create Qualification)

**Status Codes:**
- `200 OK` - Updated
- `400 Bad Request` - Invalid request
- `404 Not Found` - Qualification not found
- `500 Internal Server Error` - Update failed

---

### 18. Delete Qualification
**Endpoint:** `DELETE /api/admin/qualifications/:id`

**Description:** Removes a qualification by ID.

**Path Parameters:**
- `id` (uint) - Qualification ID

**Response:**
```json
{
  "message": "Qualification deleted successfully"
}
```

**Status Codes:**
- `200 OK` - Deleted
- `400 Bad Request` - Invalid ID
- `500 Internal Server Error` - Deletion failed

---

### 19. Get Visitor Stats (Admin)
**Endpoint:** `GET /api/admin/visitors/stats`

**Description:** Retrieves detailed visitor statistics including IP addresses and user agents.

**Response:**
```json
{
  "today": 42,
  "this_week": 215,
  "this_month": 980,
  "this_year": 5432,
  "total": 12845,
  "unique_today": 38,
  "unique_this_week": 180,
  "unique_this_month": 850,
  "unique_this_year": 4200,
  "unique_total": 10000,
  "daily_chart": [...],
  "monthly_chart": [...],
  "top_pages": [...],
  "recent_visitors": [
    {
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "path": "/portfolio",
      "referer": "https://google.com",
      "country": "United States",
      "city": "San Francisco",
      "visited_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Database error

---

## Error Handling

### Standard Error Response Format
```json
{
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes
- **200 OK** - Request successful
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid request format or validation failed
- **401 Unauthorized** - Missing or invalid authentication
- **403 Forbidden** - Authenticated but not authorized (insufficient permissions)
- **404 Not Found** - Resource does not exist
- **500 Internal Server Error** - Server error (database, file system, etc.)

### Error Messages
- `Invalid credentials` - Login failed (wrong username or password)
- `Invalid ID` - Path parameter is not a valid unsigned integer
- `Invalid or expired refresh token` - Token refresh failed
- `Score must be between 0 and 100` - Skill score validation failed
- `Level must be one of: Beginner, Intermediate, Advanced, Expert` - Invalid skill level
- `Not allowed` - CSRF validation failed or missing CSRF token
- `Database connection failed` - Health check detected connection issue

---

## Data Models

### User
```json
{
  "id": 1,
  "username": "admin",
  "password": "hashed_password",
  "role": "admin",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### About
```json
{
  "id": 1,
  "bio": "string",
  "summary": "string",
  "email": "string",
  "phone": "string",
  "location": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Portfolio
```json
{
  "id": 1,
  "title": "string",
  "description": "string",
  "category": "string",
  "tags": "string (comma-separated)",
  "project_url": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "media_files": [
    {
      "id": 1,
      "portfolio_id": 1,
      "type": "image|video",
      "url": "/uploads/portfolio/filename",
      "order_index": 0
    }
  ]
}
```

### Skill
```json
{
  "id": 1,
  "name": "string",
  "level": "Beginner|Intermediate|Advanced|Expert",
  "score": 0-100,
  "category": "string",
  "icon": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Qualification
```json
{
  "id": 1,
  "type": "string",
  "title": "string",
  "issuer": "string",
  "issue_date": "2024-01-15",
  "expiry_date": "2025-01-15|null",
  "credential_url": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Visitor
```json
{
  "id": 1,
  "ip_address": "string",
  "user_agent": "string",
  "path": "string",
  "referer": "string",
  "country": "string",
  "city": "string",
  "visited_at": "2024-01-15T10:30:00Z"
}
```

---

## Example Usage

### 1. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "adminpassword"
  }' \
  -c cookies.txt
```

### 2. Get CSRF Token and Use Protected Route
```bash
# Get CSRF token from cookies and use it
curl -X POST http://localhost:8080/api/admin/about \
  -H "Content-Type: application/json" \
  -H "X-XSRF-TOKEN: $(grep csrf_token cookies.txt | awk '{print $NF}')" \
  -b cookies.txt \
  -d '{
    "bio": "I am a developer",
    "summary": "Full-stack developer",
    "email": "user@example.com",
    "phone": "+1234567890",
    "location": "San Francisco"
  }'
```

### 3. Create Portfolio with Files
```bash
curl -X POST http://localhost:8080/api/admin/portfolio \
  -H "X-XSRF-TOKEN: $(grep csrf_token cookies.txt | awk '{print $NF}')" \
  -b cookies.txt \
  -F "title=My Project" \
  -F "description=A cool project" \
  -F "category=Web Development" \
  -F "tags=Go,React" \
  -F "project_url=https://github.com/user/project" \
  -F "media_files=@screenshot.png" \
  -F "media_files=@demo.mp4"
```

### 4. Create Skill
```bash
curl -X POST http://localhost:8080/api/admin/skills \
  -H "Content-Type: application/json" \
  -H "X-XSRF-TOKEN: $(grep csrf_token cookies.txt | awk '{print $NF}')" \
  -b cookies.txt \
  -d '{
    "name": "Go",
    "level": "Expert",
    "score": 95,
    "category": "Backend",
    "icon": "go-logo"
  }'
```

### 5. Track Visitor
```bash
curl -X POST http://localhost:8080/api/visitors/track \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/portfolio",
    "referer": "https://google.com"
  }'
```

### 6. Get Public Visitor Stats
```bash
curl http://localhost:8080/api/visitors/stats
```

### 7. Logout
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "X-XSRF-TOKEN: $(grep csrf_token cookies.txt | awk '{print $NF}')" \
  -b cookies.txt
```

---

## Environment Variables

Key environment variables used by the API:

```bash
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=portfolio

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_REFRESH_SECRET=your_refresh_secret_key
ACCESS_TOKEN_EXPIRY=900          # 15 minutes in seconds
REFRESH_TOKEN_EXPIRY=604800      # 7 days in seconds

# CSRF & Cookies
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false              # Set to true in production
COOKIE_SAMESITE=Lax

# CORS
ALLOWED_ORIGINS=http://localhost:3000,https://example.com

# Security
RESET_SECRET=your_reset_secret   # For admin password reset
```

---

## Rate Limiting & Security

- **CORS Protection**: Restricted to configured origins
- **CSRF Protection**: X-XSRF-TOKEN header required for state-changing operations
- **Password Security**: Passwords hashed with bcrypt (cost: DefaultCost)
- **Token Security**: HTTP-only cookies prevent XSS attacks
- **SSL/TLS**: Recommended for production (COOKIE_SECURE=true)

---

## Changelog

### Version 2.0.0
- JWT-based authentication with token rotation
- CSRF protection on protected routes
- Cookie-based session management
- GeoIP tracking for visitors
- Visitor statistics and analytics
- Comprehensive error handling
- CORS middleware with configurable origins

---

## Support & Issues

For issues or questions about the API, please refer to the project README or create an issue in the repository.
