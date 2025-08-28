# Go Currency API

A production-ready RESTful API for currency information of countries built with Go, featuring getting currency info of countries, caching with Redis, and PostgreSQL for data persistence.

## 🚀 Features

- **RESTful API** - Clean and well-documented endpoints
- **Currency information of Countries** - Up-to-date currency of countries
- **Redis Caching** - High-performance rate caching
- **PostgreSQL Database** - Reliable data persistence  
- **Docker Support** - Easy development setup
- **Comprehensive Testing** - Unit and integration tests
- **Health Monitoring** - Health check endpoints
- **CORS Support** - Frontend integration ready
- **Graceful Shutdown** - Proper resource cleanup

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin HTTP Framework
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **ORM**: GORM
- **Migration**: golang-migrate
- **Testing**: testify
- **Containerization**: Docker & Docker Compose

## 📁 Project Structure

```
go-currency-api/
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access layer
│   └── model/              # Data models
├── pkg/                    # Shared packages
├── migrations/             # Database migrations
├── docker/                 # Docker configurations
├── scripts/                # Utility scripts
└── docs/                   # Documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Tarifsiz/go-currency-api.git
   cd go-currency-api
   ```

2. **Initialize Go module and install dependencies**
   ```bash
   go mod init github.com/Tarifsiz/go-currency-api
   go mod tidy
   ```

3. **Start the infrastructure services**
   ```bash
   make docker-up
   # Or manually: docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the API server**
   ```bash
   make run
   # Or manually: go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Development with Tools

Start with database and Redis admin tools:
```bash
make docker-up-tools
```

This provides:
- **pgAdmin**: http://localhost:5050 (admin@admin.com / admin)
- **Redis Commander**: http://localhost:8081

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/currencies` | List supported currencies |

### Example Requests

**Health Check**
```bash
curl http://localhost:8080/health
```

**Get Currencies**
```bash
curl http://localhost:8080/api/v1/currencies
```

### Response Format

All API responses follow a consistent structure:
```json
{
  "success": true,
  "data": {...},
  "error": null,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 🧪 Testing

Run the test suite:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

## 🐳 Docker Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start PostgreSQL and Redis |
| `make docker-up-tools` | Start with admin tools |
| `make docker-down` | Stop all services |
| `make docker-logs` | View service logs |

## 🗄️ Database

### Connection Details
- **Host**: localhost:5432
- **Database**: currency_db
- **User**: currency_user
- **Password**: currency_pass

### Migrations

Create a new migration:
```bash
make migrate-create
```

Run migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

## ⚡ Redis Cache

- **Host**: localhost:6379
- **Database**: 0
- **Use**: Exchange rate caching, session storage

## 🔧 Configuration

Configure via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | API server port |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | currency_user | Database user |
| `DB_PASSWORD` | currency_pass | Database password |
| `DB_NAME` | currency_db | Database name |
| `REDIS_ADDR` | localhost:6379 | Redis address |

## 📈 Development Roadmap

This project follows a structured development approach with daily iterations:

- **Week 1**: Core infrastructure and basic API
- **Week 2**: Advanced features, testing, and optimization

See the [Development Guide](docs/DEVELOPMENT.md) for detailed iteration plans.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Commit Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `docs:` - Documentation updates
- `chore:` - Maintenance tasks

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Web framework by [Gin](https://gin-gonic.com/)
- Database by [PostgreSQL](https://www.postgresql.org/)
- Caching by [Redis](https://redis.io/)

---

**Made with ❤️ by [Tarifsiz](https://github.com/Tarifsiz)**