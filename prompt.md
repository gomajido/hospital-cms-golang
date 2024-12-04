# Hospital CMS Project Structure Guide

## Project Overview
This is a Hospital Content Management System built in Go, following clean architecture principles. The project uses a modular structure with clear separation of concerns.

## Tech Stack
- Backend: Go 1.22.4 with Fiber framework
- Database: MySQL 8.0+ with Redis caching
- Cloud: AWS (S3, SES, Secrets Manager)
- Authentication: JWT-based

## Directory Structure Explanation

### Root Level
```
hospital-cms-golang/
├── cmd/          # Entry points for the application
├── config/       # Configuration management
├── database/     # SQL migrations and schemas
├── internal/     # Core application code (private)
├── pkg/          # Shared libraries (public)
└── postman/      # API documentation
```

### Internal Module Structure
```
internal/
├── common/           # Common repositories
│   ├── whatsapp/    # Example common module
│   │   ├── domain/      # Interfaces
│   │   └── repository/  # Implementation
├── constant/         # Global constants
│   ├── auth.go         # Authentication constants
│   ├── constant.go     # General constants
│   ├── error.go        # Error codes and messages
│   ├── http.go         # HTTP-related constants
│   └── storage.go      # Storage constants
├── middleware/      # HTTP middleware components
│   └── auth_middleware.go  # Authentication & authorization middleware
├── module/           # Business modules
│   ├── appointment/  # Example module
│   │   ├── constant/     # Module constants
│   │   ├── domain/      # Entities & interfaces
│   │   ├── handler/     # HTTP handlers
│   │   ├── repository/  # Data access
│   │   ├── router/      # Route definitions
│   │   └── usecase/     # Business logic
├── dependency/      # Dependency injection
└── router/         # Main router setup
```

## Key Concepts

1. Module Organization:
   - Each business feature (appointment, article, etc.) is a separate module
   - Modules are self-contained with their own domain, handlers, and repositories
   - All modules follow the same structure pattern

2. Global Constants:
   - Located in `/internal/constant`
   - Centralized place for application-wide constants
   - Organized by domain (auth, http, error, etc.)
   - Includes:
     - Authentication constants (token types, durations, roles)
     - HTTP-related constants (methods, headers, content types)
     - Error codes and messages
     - Storage configurations
     - Validation messages
   - Module-specific constants stay in their respective module's constant folder

3. Middleware Components:
   - Located in `/internal/middleware`
   - Handles cross-cutting concerns across the application
   - Authentication middleware features:
     - Token validation and protection
     - Role-based access control
     - Ability-based authorization
     - Support for:
       - Single ability check (HasAbility)
       - Any ability check (HasAnyAbility)
       - All abilities check (HasAllAbilities)
   - Can be applied globally or to specific route groups
   - Integrates with the auth module's usecase layer

4. Common Repositories:
   - Located in `/internal/common`
   - Implements shared functionality like WhatsApp messaging, email services, etc.
   - Follows simplified module structure (domain + repository)
   - No handlers or routes as they're used internally by other modules

5. Dependency Flow:
   ```
   Router → Handler → Usecase → Repository → Database
                            ↘ Common Repository → External Services
   ```

6. Repository Pattern for data access
   - Middleware for cross-cutting concerns
   - Dependency Injection for services
   - Domain-Driven Design principles

## Code Examples

1. Route Registration:
```go
func RegisterAppointmentRoutes(router fiber.Router, handler *handler.AppointmentHandler, auth *middleware.AuthMiddleware) {
    appointmentRouter := router.Group("/appointments")
    appointmentRouter.Use(auth.Protected())
    {
        appointmentRouter.Post("/", handler.Create)
        appointmentRouter.Get("/me", handler.GetByUserID)
        // ... more routes
    }
}
```

2. Repository Setup:
```go
type AppRepositories struct {
    AuthRepo        domain.AuthRepository
    ArticleRepo     articleDomain.ArticleRepository
    DoctorRepo      doctorDomain.DoctorRepository
    AppointmentRepo appointmentDomain.AppointmentRepository
}
```

## Key Files to Understand
1. `/internal/router/router.go`: Main routing setup
2. `/internal/dependency/repositories.go`: Dependency injection
3. `/internal/module/*/router/*_router.go`: Module-specific routes
4. `/internal/module/*/domain/*.go`: Business entities and interfaces

## Common Operations
1. Adding a new feature:
   - Create new module in `/internal/module/`
   - Implement domain, handler, repository, router
   - Register routes in main router
   - Add repository to dependency container

2. Authentication:
   - Protected routes use auth middleware
   - JWT-based authentication
   - Role-based access control

## Best Practices
1. Keep modules independent
2. Use interfaces for dependency injection
3. Follow clean architecture principles
4. Maintain consistent error handling
5. Use middleware for cross-cutting concerns