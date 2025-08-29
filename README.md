# Coffee Tracker Backend

A clean architecture Go backend for tracking coffee consumption with Supabase database integration.

## 🏗️ Architecture

This project follows Clean Architecture principles:

- **Domain Layer**: Business entities and repository interfaces
- **Use Cases Layer**: Application business logic
- **Infrastructure Layer**: External dependencies (database, HTTP, etc.)

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Supabase account
- Fly.io account (for deployment)

### Setup

1. **Clone and setup**:

```bash
git clone <your-repo>
cd backend
cp .env.example .env
```

2. **Set up Supabase**:

   - Create a new Supabase project
   - Run the SQL schema from `schema.sql` in your Supabase SQL Editor
   - Get your database URL from Settings → Database

3. **Configure environment**:

```bash
# Edit .env with your values
PORT=8080
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT_ID].supabase.co:5432/postgres
JWT_SECRET=your-super-secret-key
```

4. **Install dependencies**:

```bash
go mod tidy
```

5. **Run locally**:

```bash
go run cmd/server/main.go
```

### API Endpoints

- `GET /health` - Health check
- `POST /api/v1/entries` - Create coffee entry
- `GET /api/v1/entries?date=2025-08-12limit=20&offset=0` - Get coffee entries
- `GET /api/v1/stats` - Get coffee statistics

### Sample Request

```bash
# Create a coffee entry
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Perfect morning coffee",
    "timestamp": {current_date_time}
  }'
```

## 🚀 Deployment

### Deploy to Fly.io

1. **Install Fly CLI and login**:

```bash
curl -L https://fly.io/install.sh | sh
fly auth login
```

2. **Deploy**:

```bash
fly launch --no-deploy
fly secrets set DATABASE_URL="your_supabase_url"
fly secrets set JWT_SECRET="your_jwt_secret"
fly deploy
```

3. **Test deployed API**:

```bash
curl https://your-app-name.fly.dev/health
```

## 📁 Project Structure

```
backend/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── domain/                     # Business logic layer
│   │   ├── entities/               # Business entities
│   │   └── repositories/           # Repository interfaces
│   ├── infrastructure/             # External dependencies
│   │   ├── config/                 # Configuration
│   │   ├── database/               # Database connections
│   │   ├── http/handlers/          # HTTP handlers
│   │   └── repositories/           # Repository implementations
│   └── usecases/                   # Application business rules
├── schema.sql                      # Database schema
├── Dockerfile                      # Container configuration
├── fly.toml                        # Fly.io deployment config
└── README.md                       # This file
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 🔧 Development

### Adding New Features

1. Define entities in `internal/domain/entities/`
2. Create repository interface in `internal/domain/repositories/`
3. Implement use case in `internal/usecases/`
4. Create repository implementation in `internal/infrastructure/repositories/`
5. Add HTTP handler in `internal/infrastructure/http/handlers/`
6. Register routes in `cmd/server/routes.go`

### Database Migrations

Add new tables/columns by updating `schema.sql` and running in Supabase SQL Editor.

## 📝 License

MIT License - see LICENSE file for details.
