# 🚀 Mopcare Learning & Course Management System (LMS & CMS)

**Mopcare** is a modern, microservices-based Learning Management System (LMS) and Course Management System (CMS) designed for scalable, enterprise-grade education platforms. Leveraging high-performance Go frameworks and cloud-native principles, Mopcare delivers secure, resilient, and performant learning experiences.

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
- Upcoming Volunteer Management System to boost community engagement  

---

## 🏗️ Architecture Overview

```svg
<svg width="720" height="400" viewBox="0 0 720 400" xmlns="http://www.w3.org/2000/svg" font-family="Segoe UI, Tahoma, Geneva, Verdana, sans-serif">

  <!-- Title -->
  <text x="360" y="30" text-anchor="middle" font-size="18" font-weight="700" fill="#2c3e50">Mopcare LMS & CMS Architecture Overview</text>

  <!-- Client Apps -->
  <rect x="300" y="50" width="120" height="40" rx="6" ry="6" fill="#3498db" />
  <text x="360" y="75" fill="white" font-weight="600" font-size="14" text-anchor="middle">Client Applications</text>
  <text x="360" y="92" fill="white" font-size="12" text-anchor="middle">(Web, Mobile, Admin UI)</text>

  <!-- API Gateway -->
  <rect x="280" y="120" width="160" height="70" rx="8" ry="8" fill="#1abc9c" />
  <text x="360" y="145" fill="white" font-weight="700" font-size="14" text-anchor="middle">API Gateway</text>
  <text x="360" y="165" fill="white" font-size="12" text-anchor="middle">
    Smart Routing | Caching | Load Balancing
  </text>
  <text x="360" y="180" fill="white" font-size="12" text-anchor="middle">
    Health Checks | Metrics | Logging
  </text>

  <!-- Lines from Client Apps to API Gateway -->
  <line x1="360" y1="90" x2="360" y2="120" stroke="#2980b9" stroke-width="2" marker-end="url(#arrowhead)"/>

  <!-- Services Container -->
  <rect x="40" y="210" width="640" height="140" fill="#ecf0f1" rx="12" ry="12" stroke="#bdc3c7" stroke-width="1"/>

  <!-- Service: Course -->
  <rect x="60" y="230" width="180" height="100" rx="8" ry="8" fill="#e67e22" />
  <text x="150" y="255" fill="white" font-weight="700" font-size="14" text-anchor="middle">Course Service</text>
  <text x="150" y="275" fill="white" font-size="12" text-anchor="middle">Go Fiber | PostgreSQL</text>
  <text x="150" y="295" fill="white" font-size="12" text-anchor="middle">Courses, Series, Content</text>

  <!-- Service: User -->
  <rect x="270" y="230" width="180" height="100" rx="8" ry="8" fill="#9b59b6" />
  <text x="360" y="255" fill="white" font-weight="700" font-size="14" text-anchor="middle">User Service</text>
  <text x="360" y="275" fill="white" font-size="12" text-anchor="middle">Go Gin | PostgreSQL</text>
  <text x="360" y="295" fill="white" font-size="12" text-anchor="middle">Auth via Supabase, Profiles, Payments</text>

  <!-- Service: Enrollment -->
  <rect x="480" y="230" width="180" height="100" rx="8" ry="8" fill="#16a085" />
  <text x="570" y="255" fill="white" font-weight="700" font-size="14" text-anchor="middle">Enrollment Service</text>
  <text x="570" y="275" fill="white" font-size="12" text-anchor="middle">Go Fiber | PostgreSQL</text>
  <text x="570" y="295" fill="white" font-size="12" text-anchor="middle">Enrollments, Progress, Completion</text>

  <!-- Lines from API Gateway to Services -->
  <line x1="280" y1="190" x2="150" y2="230" stroke="#16a085" stroke-width="2" marker-end="url(#arrowhead)"/>
  <line x1="360" y1="190" x2="360" y2="230" stroke="#9b59b6" stroke-width="2" marker-end="url(#arrowhead)"/>
  <line x1="440" y1="190" x2="570" y2="230" stroke="#e67e22" stroke-width="2" marker-end="url(#arrowhead)"/>

  <!-- Arrowhead Definition -->
  <defs>
    <marker id="arrowhead" markerWidth="10" markerHeight="7" 
      refX="10" refY="3.5" orient="auto" fill="#34495e">
      <polygon points="0 0, 10 3.5, 0 7" />
    </marker>
  </defs>

</svg>

Save the above SVG code as architecture.svg and embed it in your README with:
![Mopcare Architecture](./architecture.svg)


---

🛠️ Services

Service	Description	Tech Stack	Port	Database

API Gateway	Central routing, caching, load balancing, health & metrics	Go Fiber / Gin	9090	-
Course Service	Course & series management (CRUD operations)	Go Fiber	8081	PostgreSQL
User Service	User profiles, authentication (via Supabase), payments, analytics	Go Gin	8082	PostgreSQL
Enrollment Service	Enrollment tracking, progress & completion	Go Fiber	8083	PostgreSQL



---

🔐 Authentication & Security

Supabase-managed authentication with JWT and RBAC (Role-Based Access Control)

Environment variables to secure sensitive credentials

Validated inputs and secure inter-service communication

Health checks and monitoring for high availability



---

📈 Observability & Performance

Real-time metrics via /metrics endpoint on API Gateway

Cache hit/miss statistics to optimize responsiveness

Load balancing for fault tolerance and horizontal scalability

Health check endpoints for automated monitoring and graceful recovery



---

🚀 Getting Started

Development

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

Create a .env file with:

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

Volunteer Management System integration for enhanced community engagement

Advanced analytics dashboard for admins and educators

Mobile app support with offline capabilities

Rich multimedia content support in CMS



---

🤝 Contributing

Contributions are highly welcome! Please fork the repository, implement your changes, run tests, and submit a pull request.