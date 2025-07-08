import React from "react";

const ReadmeComponent = () => { return ( <div className="p-6 max-w-4xl mx-auto text-gray-800"> <h1 className="text-3xl font-bold mb-4">📚 Mopcare Learning & Course Management API</h1> <p className="mb-6"> A modern, RESTful backend built with <strong>Go (Gin)</strong> for managing courses, series (episodes), users, and enrollments — complete with Supabase integration and ready for production deployment. </p> <div className="mb-6"> <span className="font-semibold">🔗 Live API Base URL:</span> <a
href="https://go-gin-backend-t6d2.onrender.com"
className="text-blue-600 underline ml-2"
> https://go-gin-backend-t6d2.onrender.com </a> </div>

<section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">🚀 Features</h2>
    <ul className="list-disc list-inside">
      <li>Full CRUD for Courses, Series, and Users</li>
      <li>User enrollment system</li>
      <li>Supabase Storage integration (for videos/images)</li>
      <li>PostgreSQL database</li>
      <li>Production-ready structure (Docker & .env)</li>
      <li>RESTful API with clean endpoints</li>
    </ul>
  </section>

  <section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">🛠 Tech Stack</h2>
    <ul className="list-disc list-inside">
      <li>Go 1.24</li>
      <li>Gin Framework</li>
      <li>PostgreSQL via Supabase</li>
      <li>Render (Cloud Deployment)</li>
      <li>Supabase Auth & Storage</li>
      <li>godotenv, uuid, pq</li>
    </ul>
  </section>

  <section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">📁 Project Structure</h2>
    <pre className="bg-gray-100 p-4 rounded text-sm">

{.  ├── controllers/ ├── models/ ├── routes/ ├── utils/ ├── main.go ├── .env.example ├── go.mod └── README.md} </pre> </section>

<section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">📦 Installation</h2>
    <ol className="list-decimal list-inside">
      <li>Clone the repository</li>
      <pre className="bg-gray-100 p-2 rounded">
        git clone https://github.com/ObeeJ/go-gin-backend.git
        cd go-gin-backend
      </pre>
      <li>Install dependencies</li>
      <pre className="bg-gray-100 p-2 rounded">go mod tidy</pre>
      <li>Set up environment variables</li>
      <pre className="bg-gray-100 p-2 rounded">cp .env.example .env</pre>
      <li>Run the server</li>
      <pre className="bg-gray-100 p-2 rounded">go run main.go</pre>
    </ol>
  </section>

  <section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">🔗 Main Endpoints</h2>
    <div className="space-y-3">
      <div>
        <strong>GET /</strong> – Health check
      </div>
      <div>
        <strong>GET /courses</strong> – List all courses
      </div>
      <div>
        <strong>POST /courses</strong> – Create a course
      </div>
      <div>
        <strong>GET /courses/:id</strong> – Get specific course
      </div>
      <div>
        <strong>PUT /courses/:id</strong> – Update course
      </div>
      <div>
        <strong>DELETE /courses/:id</strong> – Delete course
      </div>
    </div>
  </section>

  <section className="mb-8">
    <h2 className="text-xl font-semibold mb-2">👤 Users & Enrollments</h2>
    <ul className="list-disc list-inside">
      <li>GET /users</li>
      <li>POST /users</li>
      <li>GET /users/:id</li>
      <li>DELETE /users/:id</li>
      <li>GET /users/:id/profile</li>
      <li>GET /users/:id/enrollments</li>
      <li>POST /users/:id/enrollments</li>
      <li>DELETE /enrollments/:id</li>
    </ul>
  </section>

  <footer className="mt-10 border-t pt-4 text-sm text-gray-600">
    Built with ❤️ by <a href="https://github.com/ObeeJ" className="text-blue-500">ObeeJ</a>
  </footer>
</div>

); };

export default ReadmeComponent;

