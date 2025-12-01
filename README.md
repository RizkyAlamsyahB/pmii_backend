# PMII Backend

Backend system untuk aplikasi PMII (Pergerakan Mahasiswa Islam Indonesia) menggunakan arsitektur microservices dengan Go.

## 📋 Table of Contents

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Services](#services)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)

## 🛠 Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: MySQL 8.0
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing package with httptest

## 📁 Project Structure

```
pmii_backend/
├── services/
│   └── user_service/          # User & Authentication Service
│       ├── cmd/api/           # Application entry point
│       ├── config/            # Configuration management
│       ├── internal/
│       │   ├── delivery/      # HTTP handlers, middleware, routes
│       │   ├── domain/        # Entities & repository interfaces
│       │   ├── infrastructure/# Database, messaging, storage
│       │   ├── repository/    # Repository implementations
│       │   └── usecase/       # Business logic
│       ├── migrations/        # Database migrations
│       └── tests/             # Integration & E2E tests
├── shared/                    # Shared modules across services
│   ├── jwt/                   # JWT utilities
│   ├── logger/                # Logging utilities
│   └── utils/                 # Common response helpers
├── docker-compose.yml         # Docker services configuration
└── go.work                    # Go workspace file
```

## 🚀 Services

### User Service (Port 8081)

Service untuk manajemen user dan autentikasi.

**Features:**
- User registration
- User login with JWT
- User profile management
- Password management
- User CRUD operations
- Role-based access control (Admin/Author)

## 🏃 Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.25+ (untuk development lokal)

### Quick Start

1. **Clone repository**
   ```bash
   git clone https://github.com/RizkyAlamsyahB/pmii_backend.git
   cd pmii_backend
   ```

2. **Setup environment variables**
   
   File `.env` sudah tersedia di root dan `services/user_service/.env`

3. **Start services with Docker**
   ```bash
   docker compose up -d
   ```

4. **Run database migrations**
   ```bash
   docker compose exec -T db mysql -u root -proot pmii < services/user_service/migrations/000001_create_users_table.up.sql
   docker compose exec -T db mysql -u root -proot pmii < services/user_service/migrations/seed_user_admin.sql
   ```

5. **Verify services are running**
   ```bash
   curl http://localhost:8081/health
   ```

### Default Accounts

Setelah seed data:
- **Admin**: `admin@pmii.or.id` / `admin123`
- **Author**: `author@pmii.or.id` / `admin123`

⚠️ **Ganti password default di production!**

## 📚 API Documentation

### Base URL
```
http://localhost:8081
```

### Endpoints

#### Authentication & User Management

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/v1/auth/register` | No | Register new user |
| POST | `/v1/auth/login` | No | User login |
| POST | `/v1/auth/logout` | No | User logout |
| GET | `/v1/auth/me` | Yes | Get current user profile |
| GET | `/v1/users` | Yes | List all users (paginated) |
| GET | `/v1/users/:id` | Yes | Get user by ID |

### Example Requests

**Register**
```bash
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "John Doe",
    "email": "john@example.com",
    "password": "secret123",
    "level": "2"
  }'
```

**Login**
```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Get Profile**
```bash
curl http://localhost:8081/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Response Format

**Success Response:**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Operation successful"
  },
  "data": {
    // response data
  }
}
```

**Error Response:**
```json
{
  "meta": {
    "code": 400,
    "status": "error",
    "message": "Error message"
  },
  "data": ""
}
```

## 💻 Development

### Run Locally (without Docker)

1. **Setup database**
   ```bash
   # Ensure MySQL is running
   mysql -u root -p < services/user_service/migrations/000001_create_users_table.up.sql
   ```

2. **Update .env** di `services/user_service/`
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   ```

3. **Run service**
   ```bash
   cd services/user_service
   go run ./cmd/api
   ```

### Rebuild Docker Images

```bash
docker compose up -d --build user-service
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f user-service
```

## 🧪 Testing

### Run Integration Tests

```bash
cd services/user_service
go test ./tests/integration
```

### Run with Verbose Output

```bash
go test -v ./tests/integration
```

### Test Coverage

```bash
go test -cover ./tests/integration
```

## 🔑 Environment Variables

### User Service

| Variable | Description | Default |
|----------|-------------|---------|
| APP_NAME | Application name | PMII User Service |
| APP_ENV | Environment (development/production) | development |
| APP_PORT | Service port | 8081 |
| DB_HOST | Database host | db |
| DB_PORT | Database port | 3306 |
| DB_USER | Database user | root |
| DB_PASS | Database password | root |
| DB_NAME | Database name | pmii |
| JWT_SECRET | JWT signing secret | your-secret-key-change-in-production |
| JWT_EXPIRE_HOURS | JWT expiration in hours | 24 |

## 🐳 Docker Services

| Service | Port | Description |
|---------|------|-------------|
| user-service | 8081 | User & Auth API |
| db | 3306 | MySQL Database |
| phpmyadmin | 8080 | Database Management UI |

**Access phpMyAdmin**: http://localhost:8080
- Server: `db`
- Username: `root`
- Password: `root`

## 📝 Development Notes

### API Contract Guidelines

1. **Empty Data**: Gunakan `""` (empty string), bukan `null`
2. **Image URLs**: Selalu absolute path (`https://...`), bukan relative
3. **Error Messages**: Jelas dan informatif
4. **HTTP Status Codes**: Gunakan sesuai RESTful convention

### Code Style

- Follow Go conventions
- Use `gofmt` untuk formatting
- Commit message menggunakan [Conventional Commits](https://www.conventionalcommits.org/)

### Branching Strategy

- `main`: Production-ready code
- `dev`: Development branch
- Feature branches: `feat/feature-name`
- Bugfix branches: `fix/bug-name`

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feat/amazing-feature`)
3. Commit changes (`git commit -m 'feat: add amazing feature'`)
4. Push to branch (`git push origin feat/amazing-feature`)
5. Open Pull Request ke branch `dev`

## 📄 License

This project is private and proprietary.

## 👥 Team

Developed by PMII Development Team

---

**Version**: 1.0.0  
**Last Updated**: November 28, 2025
