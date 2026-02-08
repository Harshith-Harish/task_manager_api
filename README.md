# Task Manager REST API

A production-ready REST API built with Go, Gin framework, and PostgreSQL for managing tasks with full CRUD operations.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/license-MIT-green)

## 🚀 Features

- ✅ RESTful API design with clean architecture
- ✅ Full CRUD operations for tasks
- ✅ PostgreSQL database with connection pooling
- ✅ Input validation and comprehensive error handling
- ✅ Request logging middleware
- ✅ CORS support for frontend integration
- ✅ Query filtering (status, priority, limit)
- ✅ Task statistics endpoint
- ✅ Docker support for easy deployment
- ✅ Unit tests for core functionality

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL 16
- **ORM/Query:** database/sql with lib/pq driver
- **Environment:** godotenv for configuration
- **Containerization:** Docker & Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Git

## ⚡ Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/Harshith-Harish/task-manager-api.git
cd task-manager-api
```

### 2. Start PostgreSQL Database
```bash
docker-compose up -d postgres
```

This starts PostgreSQL on `localhost:5432` with:
- Database: `taskdb`
- User: `postgres`
- Password: `postgres`

### 3. Set Up Environment
```bash
cp .env.example .env
```

### 4. Install Dependencies
```bash
go mod download
```

### 5. Run the Application
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Health Check
```http
GET /health
```

### Tasks

| Method | Endpoint | Description | Query Params |
|--------|----------|-------------|--------------|
| GET | `/api/v1/tasks` | Get all tasks | `status`, `priority`, `limit` |
| GET | `/api/v1/tasks/:id` | Get task by ID | - |
| POST | `/api/v1/tasks` | Create new task | - |
| PUT | `/api/v1/tasks/:id` | Update task | - |
| DELETE | `/api/v1/tasks/:id` | Delete task | - |

### Statistics
```http
GET /api/v1/stats
```

## 💡 Usage Examples

### Create a Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive README and API docs",
    "status": "pending",
    "priority": "high"
  }'
```

**Response:**
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "status": "pending",
  "priority": "high",
  "created_at": "2025-02-08T10:30:00Z",
  "updated_at": "2025-02-08T10:30:00Z"
}
```

### Get All Tasks
```bash
curl http://localhost:8080/api/v1/tasks
```

### Filter Tasks
```bash
# Get pending tasks
curl "http://localhost:8080/api/v1/tasks?status=pending"

# Get high priority tasks
curl "http://localhost:8080/api/v1/tasks?priority=high"

# Get in-progress tasks, limit to 10
curl "http://localhost:8080/api/v1/tasks?status=in_progress&limit=10"
```

### Update a Task
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

### Delete a Task
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

### Get Statistics
```bash
curl http://localhost:8080/api/v1/stats
```

**Response:**
```json
{
  "total": 25,
  "by_status": {
    "pending": 10,
    "in_progress": 8,
    "completed": 7
  },
  "by_priority": {
    "low": 5,
    "medium": 12,
    "high": 8
  },
  "completed_today": 3
}
```

## 🗂 Project Structure
```
task-manager-api/
├── main.go                 # Application entry point
├── models/                 # Data models and validation
│   ├── task.go
│   └── errors.go
├── handlers/               # HTTP request handlers
│   └── tasks.go
├── database/               # Database connection and setup
│   └── db.go
├── middleware/             # Custom middleware
│   ├── logger.go
│   └── cors.go
├── config/                 # Configuration management
│   └── config.go
├── tests/                  # Unit tests
│   └── tasks_test.go
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── .env.example            # Environment template
├── .gitignore             # Git ignore rules
├── docker-compose.yml     # Docker configuration
└── README.md              # This file
```

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | taskdb | Database name |
| `PORT` | 8080 | API server port |
| `GIN_MODE` | debug | Gin mode (debug/release) |

### Database Schema
```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'in_progress', 'completed')),
    priority VARCHAR(20) NOT NULL DEFAULT 'medium' 
        CHECK (priority IN ('low', 'medium', 'high')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
```

## 🧪 Testing

Run unit tests:
```bash
go test ./tests/... -v
```

Run tests with coverage:
```bash
go test ./tests/... -cover
```

## 🐳 Docker Commands

Start all services (PostgreSQL + pgAdmin):
```bash
docker-compose up -d
```

Stop all services:
```bash
docker-compose down
```

View logs:
```bash
docker-compose logs -f postgres
```

Access pgAdmin (optional):
- URL: http://localhost:5050
- Email: admin@admin.com
- Password: admin

## 🔒 Data Model

### Task Object
```json
{
  "id": 1,
  "title": "string (required, max 255 chars)",
  "description": "string (optional)",
  "status": "pending|in_progress|completed (required)",
  "priority": "low|medium|high (required)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Validation Rules

- `title`: Required, 1-255 characters
- `status`: Must be one of: `pending`, `in_progress`, `completed`
- `priority`: Must be one of: `low`, `medium`, `high`
- `description`: Optional, no length limit

## 🚧 Roadmap

- [ ] Add authentication (JWT)
- [ ] Add pagination for task lists
- [ ] Add task categories/tags
- [ ] Add due dates and reminders
- [ ] Add user assignment
- [ ] Add task comments
- [ ] Add file attachments
- [ ] Add search functionality
- [ ] Add GraphQL API
- [ ] Add rate limiting
