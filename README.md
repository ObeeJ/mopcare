# 🚀 Mopcare Learning & Course Management System (LMS & CMS)

**Mopcare** is a modern, microservices-based Learning Management System (LMS) and Course Management System (CMS) built for scalable, enterprise-grade education platforms. Leveraging high-performance Go frameworks and cloud-native principles, Mopcare delivers secure, resilient, and performant learning experiences.

---

## 🌟 Key Features

- Microservices architecture for scalability and maintainability
- High-performance backend using **Go Fiber** and **Gin**
- Secure authentication and authorization with **Supabase**
- Intelligent API Gateway with caching, load balancing, health checks, and metrics
- Comprehensive course & series management with PostgreSQL integration
- User management including profiles, payment tracking, and analytics
- Enrollment tracking with progress and completion status
- Real-time monitoring and performance analytics
- Future Volunteer Management System to enhance community engagement

---

## 🏗️ Architecture Overview


+----------------------+
           |   Client Applications |  ← (Web, Mobile, Admin UI)
           +----------+-----------+
                      |
                      v
          +-------------------------+
          |      API Gateway         |  (Go Fiber)
          |  • Request Routing       |
          |  • Caching               |
          |  • Load Balancing        |
          |  • Health Checks         |
          |  • Metrics & Logging     |
          +-----------+-------------+
                      |
    +-----------------+--------------------+
    |                  |                   |
    v                  v                   v
+-----------+    +-------------+     +--------------+
| Course    |    | User        |     | Enrollment   |
| Service   |    | Service     |     | Service      |
| (Go Fiber)|    | (Go Gin)    |     | (Go Fiber)   |
| - Courses |    | - Auth via  |     | - Enrollment |
| - Series  |    |   Supabase  |     | - Progress   |
| - CRUD    |    | - Profiles  |     | - Completion |
+-----------+    +-------------+     +--------------+

---

## 🛠️ Services

| Service            | Description                                  | Tech Stack     | Port  | Database   |
|--------------------|----------------------------------------------|----------------|-------|------------|
| API Gateway        | Central routing, caching, load balancing, health & metrics | Go Fiber / Gin | 9090  | -          |
| Course Service     | Course & series management (CRUD operations) | Go Fiber       | 8081  | PostgreSQL |
| User Service       | User profiles, authentication (via Supabase), payments, analytics | Go Gin         | 8082  | PostgreSQL |
| Enrollment Service | Enrollment tracking, progress & completion   | Go Fiber       | 8083  | PostgreSQL |

---

## 🔐 Authentication & Security

- Supabase-managed authentication with JWT and RBAC (Role-Based Access Control)
- Environment variables to secure sensitive credentials
- Validated inputs and secure inter-service communication
- Health checks and monitoring to ensure high availability

---

## 📈 Observability & Performance

- Real-time metrics via `/metrics` endpoint on API Gateway
- Cache hit/miss statistics to optimize responsiveness
- Load balancing to ensure fault tolerance and scale
- Health check endpoints for automated monitoring and recovery

---

## 🚀 Getting Started

### Development

```bash
# Start all services (Windows)
./run-services.bat

# Or individually:
cd services/course-service && go run main.go
cd services/user-service && go run main.go
cd services/enrollment-service && go run main.go
cd gateway-fiber && go run main.go

Production (Docker)

docker-compose up --build


---

🔧 Configuration

Create a .env file with the following variables:

SUPABASE_DB_URL=your_supabase_db_url
SUPABASE_PROJECT_URL=your_supabase_project_url
SUPABASE_SERVICE_KEY=your_supabase_service_key
GATEWAY_PORT=9090


---

📡 API Endpoints (Base URL: http://localhost:9090)

Category	Endpoint	Description

System	GET /health	Health check
System	GET /metrics	Performance metrics
Courses	GET /courses	List all courses
Courses	POST /courses	Create a new course
Courses	GET /courses/:id	Get course details
Courses	PUT /courses/:id	Update course
Courses	DELETE /courses/:id	Delete course
Series	GET /courses/:id/series	List series in a course
Series	POST /courses/:id/series	Add series to a course
Series	GET /series/:id	Get series details
Series	PUT /series/:id	Update series
Series	DELETE /series/:id	Delete series
Users	GET /users	List all users
Users	POST /users	Register a new user
Users	GET /users/:id	Get user profile
Users	PUT /users/:id/payment	Update payment info
Users	DELETE /users/:id	Delete user
Enrollments	GET /users/:id/enrollments	Get user enrollments
Enrollments	POST /users/:id/enrollments	Enroll user in course
Enrollments	DELETE /enrollments/:id	Remove enrollment



---

🔗 Live Demo & Resources

Live API: https://mopcare-x0vw.onrender.com/

Health Check: GET /health

Metrics: GET /metrics



---

🚀 Roadmap

Volunteer Management System integration for community engagement

Enhanced analytics dashboard for admins and educators

Mobile app support and offline mode

Rich multimedia content support in CMS



---

🤝 Contributing

Contributions are welcome! Please fork the repo and submit pull requests.
