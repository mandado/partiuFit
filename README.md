# PartiuFit 💪

A Go-based fitness tracking API that allows users to manage their workouts, track exercises, and maintain their fitness journey. Built with a clean architecture using Go, PostgreSQL, and Chi router.

## 🚀 Features

- **User Management**: User registration, authentication, and profile management
- **Workout Tracking**: Create, read, update, and delete workouts
- **Exercise Management**: Track individual exercises within workouts
- **Authentication**: Token-based authentication system
- **Database Migrations**: Automated database schema management
- **Health Monitoring**: Built-in health check endpoints
- **Hot Reload**: Development environment with auto-reload

## 🏗️ Project Architecture

```
partiuFit/
├── internal/
│   ├── app/                    # Application initialization and configuration
│   ├── database/               # Database connection and utilities
│   ├── handlers/               # HTTP request handlers
│   ├── middlewares/            # HTTP middlewares (auth, error handling)
│   ├── requests/               # Request validation structures
│   ├── routes/                 # API route definitions
│   ├── store/                  # Data access layer
│   ├── tokens/                 # Token management
│   ├── utils/                  # Utility functions
│   └── valueObjects/           # Domain value objects
├── migrations/                 # Database migration files
├── config/                     # Configuration files
├── bin/                        # Compiled binaries
└── tmp/                        # Temporary files (development)
```

## 🔧 Prerequisites

Before running this application, ensure you have the following installed:

### Required Dependencies
- **Go 1.24+** - [Install Go](https://golang.org/doc/install)
- **PostgreSQL 14+** - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)

### Development Tools (Recommended)
- **Air** - Hot reloading for Go apps
  ```bash
  go install github.com/air-verse/air@latest
  ```
- **golangci-lint** - Go linter
  ```bash
  # Linux/macOS
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
  
  # Or via Homebrew (macOS)
  brew install golangci-lint
  ```
- **goimports** - Import formatting
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  ```

## ⚙️ Environment Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd partiuFit
   ```

2. **Environment Configuration**
   
   Create a `.env` file in the root directory:
   ```bash
   cp .env.testing .env
   ```
   
   Update the `.env` file with your configuration:
   ```env
   DATABASE_URL="postgres://postgres:postgres@localhost:5432/partiufit?sslmode=disable"
   PORT=8080
   APP_ENV=development
   ```

3. **Database Setup**
   
   **Option A: Using Docker (Recommended)**
   ```bash
   # Start PostgreSQL with Docker Compose
   docker-compose up db -d
   ```
   
   **Option B: Local PostgreSQL**
   ```bash
   # Create database manually
   createdb partiufit
   ```

## 🚀 Quick Start

### Using Docker (Recommended)
```bash
# Start all services (database + application)
docker-compose up

# Or run in background
docker-compose up -d
```

### Local Development
```bash
# Install dependencies
go mod download

# Run database migrations (automatically handled by app)
# Start the application with hot reload
make run

# Or build and run manually
make build
./bin/partiuFit
```

The API will be available at `http://localhost:8080`

## 📋 Available Make Commands

```bash
make help      # Show available commands
make format    # Format Go code using gofmt and goimports
make lint      # Run golangci-lint
make run       # Run with hot reload (uses Air)
make build     # Build the application binary
make test      # Run all tests
make clean     # Clean build artifacts
```

## 🔌 API Endpoints

### Health Check
- `GET /health` - Check application health status

### Authentication
- `POST /tokens` - Generate authentication token (login)

### User Management
- `POST /users` - Register a new user
- `PUT /users` - Update user profile (requires authentication)

### Workout Management (Authentication Required)
- `GET /workouts` - Get all user workouts
- `POST /workouts` - Create a new workout
- `GET /workouts/{id}` - Get specific workout by ID
- `PUT /workouts/{id}` - Update specific workout
- `DELETE /workouts/{id}` - Delete specific workout

## 🗄️ Database Schema

The application uses PostgreSQL with the following main entities:

- **Users**: User accounts and profiles
- **Workouts**: Workout sessions
- **Workout_Entries**: Individual exercises within workouts
- **Tokens**: Authentication tokens

Migrations are automatically applied on application startup.

## 🧪 Testing

Run the test suite:
```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/store
```

## 🔧 Development Workflow

1. **Start the development environment**
   ```bash
   docker-compose up db -d  # Start database
   make run                 # Start app with hot reload
   ```

2. **Code formatting and linting**
   ```bash
   make format  # Format code
   make lint    # Run linter
   ```

3. **Running tests**
   ```bash
   make test
   ```

## 📊 Key Dependencies

- **Web Framework**: [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- **Database**: [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- **Migrations**: [Goose](https://github.com/pressly/goose) - Database migration tool
- **Validation**: [validator](https://github.com/go-playground/validator) - Struct validation
- **Logging**: [Zap](https://github.com/uber-go/zap) - Structured logging
- **Environment**: [godotenv](https://github.com/joho/godotenv) - Environment variable loading
- **Passwords**: [bcrypt](https://golang.org/x/crypto/bcrypt) - Password hashing
- **Testing**: [Testify](https://github.com/stretchr/testify) - Testing toolkit

## 🌐 Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | PostgreSQL connection string | - | ✅ |
| `PORT` | Server port | `8080` | ✅ |
| `APP_ENV` | Application environment | `development` | ❌ |

## 🚀 Production Deployment

1. **Build the application**
   ```bash
   make build
   ```

2. **Set environment variables**
   ```bash
   export APP_ENV=production
   export DATABASE_URL="your-production-db-url"
   export PORT=8080
   ```

3. **Run the binary**
   ```bash
   ./bin/partiuFit
   ```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make format`)
6. Lint code (`make lint`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## 📝 Project Status

This is an active fitness tracking API project built with modern Go practices. The application follows clean architecture principles with clear separation of concerns.

## 🔒 Security

- Passwords are hashed using bcrypt
- Token-based authentication
- Input validation on all endpoints
- SQL injection prevention through parameterized queries

## 📞 Support

If you encounter any issues or have questions, please open an issue on the repository.

---

**Made with ❤️ using Go**