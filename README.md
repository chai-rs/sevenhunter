# SevenHunter Take-Home Assignment: User Management API

## Overview
This project implements a RESTful API for user authentication and management. The system provides secure user registration, authentication using JWT tokens, and comprehensive user profile management capabilities.

## Core Functionality
The API provides:
- User registration with email validation and password hashing
- Secure authentication with JWT-based access and refresh tokens
- User profile management (read, update, delete)
- User listing with cursor-based pagination
- User statistics and counting
- Complete API documentation with Swagger

## Tech Stack

### Backend
- **Language**: Go 1.25.3
- **Web Framework**: Fiber v2 (Express-inspired framework for Go)
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Logging**: Zerolog
- **Validation**: ozzo-validation

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database UI**: Mongo Express
- **Migrations**: migrate-mongo

## API Features

### Authentication Endpoints
- `POST /v1/api/auth/register` - Create a new user account
- `POST /v1/api/auth/login` - Authenticate and receive tokens
- `POST /v1/api/auth/refresh` - Refresh access token using refresh token

### User Management Endpoints (Protected)
- `GET /v1/api/users/profile` - Get current user profile
- `PUT /v1/api/users/profile` - Update current user profile
- `DELETE /v1/api/users/profile` - Delete current user account
- `GET /v1/api/users` - List all users (paginated)
- `GET /v1/api/users/count` - Get total user count

## Security Features

### Authentication
- Password hashing using bcrypt
- JWT-based authentication with separate access and refresh tokens
- Configurable token expiration times
- Bearer token authentication for protected routes

### Validation
- Email format validation
- Password length requirements (8-64 characters)
- Username length requirements (2-32 characters)
- Request body validation

## Project Structure
```
├── cmd/
│   └── api/              # API server entry point
├── internal/
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # Authentication middleware
│   ├── model/            # Database models
│   ├── repo/             # Repository layer (data access)
│   ├── router/           # Route definitions
│   └── service/          # Business logic layer
├── pkg/
│   ├── fiber/            # Fiber utilities
│   └── logger/           # Logging utilities
├── docker/               # Docker configurations
├── migrations/           # Database migrations
├── docs/                 # Swagger documentation
├── docker-compose.yml    # Multi-service orchestration
├── makefile             # Build commands
└── README.md            # This file
```

## Setup and Running

### Prerequisites
- Docker and Docker Compose
- (Optional) Go 1.25.3 for local development

### Using Docker Compose (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd sevenhunter

# Copy environment variables
cp .env.example .env

# Start all services
docker-compose up
```

This will start:
- **API Server**: `http://localhost:8080`
- **Swagger Documentation**: `http://localhost:8080/swagger`
- **MongoDB**: `localhost:27017`
- **Mongo Express UI**: `http://localhost:8081`

### Environment Variables
Configure the following variables in `.env`:

```bash
# Application Settings
APP_PORT="8080"
APP_CORS_ALLOWED_ORIGINS="*"
APP_CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE,OPTIONS"

# JWT Configuration
AUTH_SECRET="your_secret_key"
AUTH_ACCESS_TOKEN_TTL="15m"
AUTH_REFRESH_TOKEN_TTL="168h"

# Database Configuration
MONGO_URI="mongodb://admin:admin@localhost:27017/"
```

### Local Development
```bash
# Install dependencies
go mod download

# Run migrations
# (MongoDB will be auto-migrated via docker-compose)

# Run the application
go run cmd/api/main.go
```

## API Documentation

### Swagger UI
Access interactive API documentation at `http://localhost:8080/swagger` after starting the server.

### Authentication Flow
1. **Register**: Create a new account with email, name, and password
2. **Login**: Authenticate to receive access and refresh tokens
3. **Access Protected Routes**: Include `Authorization: Bearer <access_token>` in request headers
4. **Refresh Token**: Use refresh token to get new access token when expired

### Example API Calls

#### Register
```bash
curl -X POST http://localhost:8080/v1/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "securepassword123"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/v1/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

#### Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/v1/api/users/profile \
  -H "Authorization: Bearer <your_access_token>"
```

## Docker Services

### MongoDB
- Default credentials: `admin/admin`
- Port: `27017`
- Persistent storage via Docker volumes

### Mongo Express
- Web-based MongoDB admin interface
- Access: `http://localhost:8081`
- Credentials: `admin/admin`

### API Server
- Built using multi-stage Docker build
- Includes health checks
- Auto-restarts on failure

## Development Features

### Architecture
- Clean architecture with separation of concerns
- Repository pattern for data access
- Service layer for business logic
- DTO pattern for request/response handling

### Code Quality
- Comprehensive error handling
- Request validation
- Structured logging
- Type-safe MongoDB operations

### Testing
- Unit tests support with testify
- Mock generation with mockery
- Test configuration available

## Key Design Decisions

### Token Management
- Separate access (15m) and refresh tokens (7d) for security
- Refresh tokens allow seamless re-authentication
- Configurable token expiration via environment variables

### Password Security
- Bcrypt hashing for password storage
- Never store plain-text passwords
- Minimum password length enforcement

### Database
- MongoDB for flexible schema and scalability
- Indexed email field for fast lookups
- Atomic operations for data consistency

### Pagination
- Cursor-based pagination for user listing
- Configurable page size (max 100 items)
- Efficient for large datasets

## API Response Format
All API responses follow a consistent format:
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "result": {
    // Response data here
  }
}
```

## Error Handling
- 400: Bad Request - Validation errors or invalid input
- 401: Unauthorized - Invalid or missing authentication
- 404: Not Found - Resource not found
- 409: Conflict - Duplicate resource (e.g., email already exists)
- 500: Internal Server Error - Server-side errors

---

**Developer**: Thanatorn Kanthala
**Contact**: tk.thanatorn@gmail.com
