# Hospital CMS - Golang

A hospital content management system built with Go, featuring role-based access control and modern authentication.

## Technology Stack

### Backend
- **Language**: Go 1.19+
- **Web Framework**: Fiber - Fast, flexible HTTP web framework
- **Database**: MySQL 8.0+ - Primary data store
- **Cache**: Redis - For session management and caching
- **Authentication**: JWT (JSON Web Tokens)
- **Migration**: golang-migrate - Database versioning

### Development Tools
- **Version Control**: Git
- **API Testing**: Postman
- **Mocking**: gomock - For unit testing
- **SQL Mocking**: sqlmock - For database testing
- **Testing**: Go testing package with testify

## Project Structure

```
hospital-cms-golang/
├── cmd/                    # Application entry points
│   └── main.go            # Main application file
├── config/                 # Configuration files
│   ├── config.go          # Configuration structures
│   └── env-local.yaml     # Environment-specific config
├── database/              # Database related files
│   └── migrations/        # SQL migration files
├── internal/              # Private application code
│   ├── dependency/        # Dependency injection
│   ├── middleware/        # HTTP middlewares
│   ├── module/           # Feature modules
│   │   ├── auth/         # Authentication module
│   │   │   ├── domain/   # Domain models and interfaces
│   │   │   ├── handler/  # HTTP handlers
│   │   │   ├── repository/ # Data access layer
│   │   │   └── usecase/  # Business logic
│   │   └── [other modules]
│   ├── response/         # Common HTTP responses
│   └── utils/           # Utility functions
├── postman/             # Postman collections
├── scripts/            # Utility scripts
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Architecture

The project follows Clean Architecture principles with the following layers:

1. **Domain Layer** (`domain/`)
   - Contains business logic interfaces
   - Defines domain models and entities
   - Houses repository and usecase interfaces

2. **Repository Layer** (`repository/`)
   - Implements data access logic
   - Handles database operations
   - Manages data persistence

3. **Usecase Layer** (`usecase/`)
   - Implements business logic
   - Orchestrates data flow between layers
   - Handles business rules and validations

4. **Handler Layer** (`handler/`)
   - Manages HTTP requests and responses
   - Handles input validation
   - Routes requests to appropriate usecases

## Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control
  - Session management with Redis
  - Password hashing and verification

- **User Management**
  - User registration and login
  - Profile management
  - Role assignment
  - Password reset functionality

- **Security**
  - Password hashing using bcrypt
  - Token-based authentication
  - Role-based access control
  - Request validation and sanitization

## Prerequisites

- Go 1.19 or higher
- MySQL 8.0 or higher
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Installation

1. Install golang-migrate:
```bash
# For macOS using Homebrew
brew install golang-migrate

# For other platforms, visit:
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```

2. Clone the repository:
```bash
git clone https://github.com/yourusername/hospital-cms-golang.git
cd hospital-cms-golang
```

3. Install dependencies:
```bash
go mod download
```

## Database Setup

1. Create a MySQL database:
```sql
CREATE DATABASE apexa;
```

2. Configure your database connection in `config/env-local.yaml`:
```yaml
database:
  driverName: "mysql"
  user: "root"
  password: "your_password"
  dbname: "apexa"
  network: "tcp"
  address: "localhost:3308"  # Adjust port if needed
```

3. Run database migrations:
```bash
# Apply all migrations
migrate -path database/migrations -database "mysql://root:password@tcp(localhost:3308)/apexa" up

# Rollback all migrations
migrate -path database/migrations -database "mysql://root:password@tcp(localhost:3308)/apexa" down

# Rollback specific number of migrations
migrate -path database/migrations -database "mysql://root:password@tcp(localhost:3308)/apexa" down N

# Force a specific version
migrate -path database/migrations -database "mysql://root:password@tcp(localhost:3308)/apexa" force VERSION

# Create a new migration
migrate create -ext sql -dir database/migrations -seq migration_name
```

## Migration Files

The project uses SQL migrations located in `database/migrations/`:

1. `000001_create_auth_tables.up.sql`: Creates core authentication tables
   - `users`: Stores user information
   - `roles`: Defines available roles
   - `user_roles`: Maps users to roles

2. `000002_seed_roles.up.sql`: Seeds initial role data
   - Adds default roles like admin, doctor, etc.

3. `000003_create_user_tokens_table.up.sql`: Creates authentication token table
   - Manages user sessions and authentication tokens

Each migration has a corresponding `.down.sql` file for rollback operations.

## Running the Application

1. Start the server:
```bash
go run main.go serve-rest-api
```

2. The API will be available at `http://localhost:8080/api`

## API Documentation

Import the Postman collection from `postman/Auth_API_Tests.postman_collection.json` for API documentation and testing.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
