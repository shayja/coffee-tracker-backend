# ☕ Coffee Tracker Backend

A clean-architecture Go backend for tracking coffee consumption, integrated with Supabase.

---

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer**: Business entities and repository interfaces
- **Use Cases Layer**: Application business logic
- **Infrastructure Layer**: External dependencies (database, HTTP, storage, etc.)

### Architecture Diagram

```mermaid
flowchart TD
    subgraph Domain
        E1[Entities]
        R1[Repository Interfaces]
    end

    subgraph UseCases
        UC1[Business Logic / Use Cases]
    end

    subgraph Infrastructure
        DB[Supabase Database]
        HTTP[HTTP Handlers & Middleware]
        Auth[JWT Service]
        Storage[Supabase Storage]
        RepoImpl[Repository Implementations]
    end

    subgraph Server
        Main[cmd/server/main.go]
        Routes[routes.go]
        ServerGo[server.go & dependencies.go]
    end

    E1 --> UC1
    R1 --> UC1
    UC1 --> RepoImpl
    RepoImpl --> DB
    HTTP --> UC1
    Auth --> HTTP
    Storage --> HTTP
    Main --> Routes
    Routes --> HTTP
    Routes --> ServerGo

🚀 Quick Start
Prerequisites
Go 1.24+
Supabase account
Fly.io account (optional, for deployment)
Setup Locally
Clone the repository:
git clone <your-repo>
cd backend
cp .env.example .env
Set up Supabase:
Create a new Supabase project
Run schema.sql in Supabase SQL Editor
Copy your Database URL from Settings → Database
Configure environment variables:
PORT=8080
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT_ID].supabase.co:5432/postgres
JWT_SECRET=your-super-secret-key
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h # 7 days
Install dependencies:
go mod tidy
Run the server:
go run cmd/server/main.go
📡 API Endpoints
Health Check
GET /health
Coffee Entries
POST /api/v1/entries
GET /api/v1/entries?date=2025-08-12&limit=20&offset=0
GET /api/v1/stats
Sample Request:
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Perfect morning coffee",
    "timestamp": "2025-10-07T08:00:00Z"
  }'
🚀 Deployment to Fly.io
Install Fly CLI and login:
curl -L https://fly.io/install.sh | sh
fly auth login
Deploy the application:
fly launch --no-deploy
fly secrets set DATABASE_URL="your_supabase_url"
fly secrets set JWT_SECRET="your_jwt_secret"
fly deploy
Test deployed API:
curl https://your-app-name.fly.dev/health
📁 Project Structure
backend/
├── cmd/server/                  # Application entry point and routes
│   ├── main.go
│   ├── server.go
│   ├── dependencies.go
│   └── routes.go
├── internal/
│   ├── contextkeys/             # Context helpers for middleware
│   ├── entities/                # Business entities
│   ├── repositories/            # Repository interfaces
│   ├── usecases/                # Application use cases
│   └── infrastructure/          # External dependencies
│       ├── auth/                # JWT service
│       ├── config/              # App configuration
│       ├── database/            # Database connections
│       ├── http/handlers/       # HTTP handlers
│       ├── http/middleware/     # Middleware
│       └── repositories/        # Repository implementations
├── schema.sql                    # Database schema
├── Dockerfile                    # Container configuration
├── fly.toml                      # Fly.io deployment config
└── README.md                     # Project documentation
🧪 Testing
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
🔧 Development Guide
Adding a New Feature
Define entities in internal/entities/
Add repository interface in internal/repositories/
Implement use case in internal/usecases/
Implement repository in internal/infrastructure/repositories/
Add HTTP handler in internal/infrastructure/http/handlers/
Register route in cmd/server/routes.go
Database Migrations
Add or update tables in schema.sql
Run changes in Supabase SQL Editor
```

📝 License
MIT License — see LICENSE for details.
