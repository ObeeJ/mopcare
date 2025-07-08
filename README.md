<!DOCTYPE html><html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Mopcare LMS</title>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
  <script src="https://kit.fontawesome.com/a076d05399.js" crossorigin="anonymous"></script>
  <style>
    body {
      font-family: 'Inter', sans-serif;
      line-height: 1.6;
      margin: 2rem auto;
      max-width: 900px;
      padding: 0 1rem;
      color: #1e293b;
    }
    h1, h2, h3 {
      color: #1e3a8a;
    }
    code, pre {
      background: #f3f4f6;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      font-family: monospace;
    }
    .icon {
      margin-right: 0.5rem;
      color: #10b981;
    }
    .tech-stack i {
      font-size: 1.25rem;
      margin-right: 0.25rem;
      color: #6366f1;
    }
    ul {
      padding-left: 1rem;
    }
  </style>
</head>
<body>
  <h1><i class="fas fa-graduation-cap icon"></i> Mopcare Learning & Course Management System</h1>
  <p>A comprehensive LMS built with modern technologies for course management, enrollments, and payments.</p>  <h2><i class="fas fa-bolt icon"></i> Features</h2>
  <ul>
    <li>Course & Series Management with video support</li>
    <li>User Authentication & Profile System</li>
    <li>Enrollment & Progress Tracking</li>
    <li>Payment Integration</li>
    <li>Supabase File Uploads</li>
    <li>RESTful API with clear documentation</li>
  </ul>  <h2><i class="fas fa-cogs icon"></i> Tech Stack</h2>
  <h3>Backend</h3>
  <ul class="tech-stack">
    <li><i class="fab fa-golang"></i> Go 1.24</li>
    <li><i class="fas fa-fire"></i> Gin Framework</li>
    <li><i class="fas fa-database"></i> PostgreSQL</li>
    <li><i class="fas fa-cloud"></i> Supabase</li>
  </ul>  <h3>Frontend Integration</h3>
  <ul class="tech-stack">
    <li><i class="fas fa-lock"></i> Supabase Auth</li>
    <li><i class="fas fa-credit-card"></i> Payment Gateway</li>
  </ul>  <h3>Deployment</h3>
  <ul class="tech-stack">
    <li><i class="fab fa-github"></i> GitHub Actions</li>
    <li><i class="fas fa-server"></i> Render</li>
  </ul>  <h2><i class="fas fa-tools icon"></i> Setup Instructions</h2>
  <ol>
    <li>Clone the repository:<br><code>git clone https://github.com/ObeeJ/go-gin-backend.git</code></li>
    <li>Navigate into the project:<br><code>cd go-gin-backend</code></li>
    <li>Install dependencies:<br><code>go mod tidy</code></li>
    <li>Configure <code>.env</code> using <code>.env.example</code>:</li>
    <pre><code>SUPABASE_DB_URL=...
SUPABASE_PROJECT_URL=...
SUPABASE_SERVICE_KEY=...
PORT=9090</code></pre>
    <li>Run migrations & start the app:<br><code>go run main.go</code></li>
  </ol>  <h2><i class="fas fa-route icon"></i> API Endpoints</h2>
  <h3>System</h3>
  <ul>
    <li><code>GET /</code> - Health check</li>
  </ul>  <h3>Courses</h3>
  <ul>
    <li><code>GET /courses</code> - List all courses</li>
    <li><code>POST /courses</code> - Create new course</li>
    <li><code>GET /courses/:id</code> - View course</li>
    <li><code>PUT /courses/:id</code> - Update course</li>
    <li><code>DELETE /courses/:id</code> - Delete course</li>
  </ul>  <h3>Series</h3>
  <ul>
    <li><code>GET /courses/:id/series</code> - List series in course</li>
    <li><code>POST /courses/:id/series</code> - Create series</li>
    <li><code>GET /series/:id</code> - View series</li>
    <li><code>PUT /series/:id</code> - Update series</li>
    <li><code>DELETE /series/:id</code> - Delete series</li>
  </ul>  <h3>Users</h3>
  <ul>
    <li><code>GET /users</code> - List all users</li>
    <li><code>POST /users</code> - Create new user</li>
    <li><code>GET /users/:id</code> - View user</li>
    <li><code>DELETE /users/:id</code> - Delete user</li>
  </ul>  <h3>Enrollments</h3>
  <ul>
    <li><code>GET /users/:id/enrollments</code> - View user's enrollments</li>
    <li><code>POST /users/:id/enrollments</code> - Enroll in course</li>
    <li><code>DELETE /enrollments/:id</code> - Remove enrollment</li>
  </ul>  <h3>Payments</h3>
  <ul>
    <li><code>PUT /users/:id/payment</code> - Update payment info</li>
  </ul>  <p style="margin-top:2rem;"><strong>Base URL:</strong> <a href="https://go-gin-backend-t6d2.onrender.com" target="_blank">https://go-gin-backend-t6d2.onrender.com</a></p>
</body>
</html>