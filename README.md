# 🚀 Mopcare Learning & Course Management System

A modern microservices-based LMS built with Go Fiber, featuring intelligent API gateway, caching, and comprehensive course management.

## ✨ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client App    │───▶│  API Gateway     │───▶│  Microservices  │
└─────────────────┘    │  (Go/Fiber)      │    └─────────────────┘
                       │                  │
                       │ • Caching        │
                       │ • Load Balancing │
                       │ • Metrics        │
                       │ • Health Checks  │
                       └──────────────────┘
```

## 🏗️ Services

### API Gateway (Port 9090)
- **Intelligent Routing** - Routes requests to appropriate microservices
- **Caching** - Sub-second response times for GET requests
- **Metrics** - Real-time performance analytics
- **Health Monitoring** - Service availability tracking

### Course Service (Port 8081)
- Course & Series Management
- PostgreSQL database integration
- CRUD operations for courses and series

### User Service (Port 8082)
- User Authentication & Profiles
- Payment tracking
- User analytics

### Enrollment Service (Port 8083)
- Course enrollments
- Progress tracking
- Completion status

## 🚀 Quick Start

### Development
```bash
# Start all services
run-services.bat

# Or manually:
cd services/course-service && go run main.go
cd services/user-service && go run main.go  
cd services/enrollment-service && go run main.go
cd gateway-fiber && go run main.go
```

### Production (Docker)
```bash
docker-compose up --build
```

## 📡 API Endpoints

**Base URL:** `http://localhost:9090`

### System
- `GET /health` - Gateway health check
- `GET /metrics` - Performance metrics

### Courses
- `GET /courses` - List all courses
- `POST /courses` - Create new course
- `GET /courses/:id` - View course
- `PUT /courses/:id` - Update course
- `DELETE /courses/:id` - Delete course

### Series
- `GET /courses/:id/series` - List series in course
- `POST /courses/:id/series` - Create series
- `GET /series/:id` - View series
- `PUT /series/:id` - Update series
- `DELETE /series/:id` - Delete series

### Users
- `GET /users` - List all users
- `POST /users` - Create new user
- `GET /users/:id` - View user
- `DELETE /users/:id` - Delete user
- `PUT /users/:id/payment` - Update payment info

### Enrollments
- `GET /users/:id/enrollments` - View user's enrollments
- `POST /users/:id/enrollments` - Enroll in course
- `DELETE /enrollments/:id` - Remove enrollment

## 🔧 Configuration

Create `.env` file with:
```env
SUPABASE_DB_URL=your_database_connection_string
SUPABASE_PROJECT_URL=your_supabase_project_url  
SUPABASE_SERVICE_KEY=your_supabase_service_key
GATEWAY_PORT=9090
```

## 📊 Performance Metrics

Access real-time metrics at: `GET /metrics`

```json
{
  "gateway": {
    "total_requests": 1000,
    "cache_hits": 300,
    "cache_misses": 700
  }
}
```

## 🛡️ Security

- Environment variables for sensitive data
- PostgreSQL database connections
- Input validation and error handling
- Secure service-to-service communication

## 🏭 Production Deployment

- **Docker** - Multi-stage builds for optimization
- **Render** - Cloud deployment ready
- **Health Checks** - Automatic service monitoring  
- **Caching** - Improved response times
- **Metrics** - Performance monitoring

## 📁 Project Structure

```
├── gateway-fiber/           # API Gateway (Fiber)
├── services/
│   ├── course-service/      # Course management
│   ├── user-service/        # User management
│   └── enrollment-service/  # Enrollment management
├── docker-compose.yml       # Container orchestration
├── run-services.bat         # Development startup script
├── render.yaml             # Render deployment config
└── README.md               # This file
```

## 🔗 Links

- **Live API:** [https://mopcare-x0vw.onrender.com/](https://mopcare-x0vw.onrender.com/)
- **Health Check:** `GET /health`
- **Metrics:** `GET /metrics`
