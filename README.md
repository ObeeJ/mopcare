# Mopcare Learning & Course Management System

A comprehensive learning management system built with modern technologies, featuring course management, user enrollment, and integrated payment processing.

## 🚀 Features

- **Course Management**: Create, update, and manage courses with multimedia content
- **Series & Episodes**: Organize course content into structured series with video lessons
- **User Management**: Complete user registration, authentication, and profile management
- **Enrollment System**: Track user course enrollments and progress
- **Payment Integration**: Secure payment processing and transaction tracking
- **File Upload**: Support for video and image uploads via Supabase Storage
- **RESTful API**: Clean, well-documented API endpoints

## 🛠️ Tech Stack

### Backend
- **Go 1.24** - High-performance backend language
- **Gin Framework** - Fast HTTP web framework
- **PostgreSQL** - Robust relational database
- **Supabase** - Backend-as-a-Service for database and storage

### Frontend Integration
- **Supabase Auth** - User authentication and authorization
- **Payment Processing** - Integrated payment collection system

### DevOps & Deployment
- **Render** - Cloud deployment platform
- **GitHub Actions** - CI/CD pipeline automation
- **Git** - Version control

### Development Tools
- **godotenv** - Environment variable management
- **lib/pq** - PostgreSQL driver for Go
- **google/uuid** - UUID generation

## 📋 Prerequisites

- Go 1.24 or higher
- PostgreSQL database
- Supabase account
- Render account (for deployment)

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ObeeJ/go-gin-backend.git
   cd go-gin-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Configure the following variables:
   ```env
   SUPABASE_DB_URL=your_supabase_database_url
   SUPABASE_PROJECT_URL=your_supabase_project_url
   SUPABASE_SERVICE_KEY=your_supabase_service_key
   PORT=9090
   ```

4. **Run database migrations**
   Execute the SQL schema in your Supabase dashboard or PostgreSQL client.

5. **Start the application**
   ```bash
   go run main.go
   ```

## 🗄️ Database Schema

The system uses four main tables:
- `courses` - Course information and metadata
- `series` - Course episodes and lessons
- `users` - User profiles and data
- `user_course_enrollments` - Enrollment tracking

## 🔗 API Endpoints

### Courses
- `GET /courses` - List all courses
- `POST /courses` - Create new course
- `GET /courses/:id` - Get specific course
- `PUT /courses/:id` - Update course
- `DELETE /courses/:id` - Delete course

### Series
- `GET /courses/:id/series` - Get course series
- `POST /courses/:id/series` - Create series for course
- `GET /series/:id` - Get specific series
- `PUT /series/:id` - Update series
- `DELETE /series/:id` - Delete series

### Users
- `GET /users` - List users
- `POST /users` - Create user
- `GET /users/:id` - Get user profile
- `DELETE /users/:id` - Delete user

### Enrollments
- `GET /users/:id/enrollments` - Get user enrollments
- `POST /users/:id/enrollments` - Enroll user in course
- `DELETE /enrollments/:id` - Remove enrollment

## 🚀 Deployment

The application is configured for automatic deployment on Render with GitHub Actions:

1. **Push to main branch** triggers automatic deployment
2. **Render** builds and deploys the application
3. **Environment variables** are managed through Render dashboard

## 🤝 Architecture

- **Frontend**: Handles user interface, authentication (Supabase Auth), and payment collection
- **Backend**: Manages business logic, data processing, and API endpoints
- **Database**: PostgreSQL via Supabase for data persistence
- **Storage**: Supabase Storage for file uploads

## 📝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**ObeeJ** - [GitHub Profile](https://github.com/ObeeJ)

## 🙏 Acknowledgments

- Gin Framework community
- Supabase team
- Go community
- Render platform

---

**Mopcare Learning & Course Management System** - Empowering education through technology.