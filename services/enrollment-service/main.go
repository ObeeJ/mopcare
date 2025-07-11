package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type UserCourseEnrollment struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	CourseID int    `json:"course_id"`
	Status   string `json:"status"`
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "enrollment-service", "status": "running"})
	})

	router.GET("/users/:id/enrollments", getUserEnrollments)
	router.POST("/users/:id/enrollments", createUserEnrollment)
	router.DELETE("/enrollments/:id", deleteUserEnrollment)

	port := os.Getenv("ENROLLMENT_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}
	log.Printf("Enrollment service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start enrollment service: %v", err)
	}
}

func connectDB() (*sql.DB, error) {
	connStr := os.Getenv("SUPABASE_DB_URL")
	if connStr == "" {
		if err := godotenv.Load(); err != nil {
			return nil, errors.New("supabase_db_url environment variable is not set and .env file not found")
		}
		connStr = os.Getenv("SUPABASE_DB_URL")
		if connStr == "" {
			return nil, errors.New("supabase_db_url not found in environment or .env")
		}
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}
	fmt.Println("Database connection established successfully.")
	return db, nil
}

func getUserEnrollments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := getDB(c)
	if db == nil {
		return
	}

	rows, err := db.Query("SELECT id, user_id, course_id, status FROM user_course_enrollments WHERE user_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var enrollments []UserCourseEnrollment
	for rows.Next() {
		var enrollment UserCourseEnrollment
		if err := rows.Scan(&enrollment.ID, &enrollment.UserID, &enrollment.CourseID, &enrollment.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		enrollments = append(enrollments, enrollment)
	}
	if len(enrollments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No enrollments found for this user"})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

func createUserEnrollment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var enrollment UserCourseEnrollment
	if err := c.BindJSON(&enrollment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if enrollment.UserID == 0 || enrollment.CourseID == 0 || enrollment.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, Course ID, and Status are required"})
		return
	}
	if enrollment.UserID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID in body must match URL parameter"})
		return
	}
	if enrollment.Status != "enrolled" && enrollment.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'enrolled' or 'completed'"})
		return
	}

	db := getDB(c)
	if db == nil {
		return
	}

	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", enrollment.UserID).Scan(&userExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	var courseExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", enrollment.CourseID).Scan(&courseExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !courseExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course does not exist"})
		return
	}

	var enrollmentExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_course_enrollments WHERE user_id = $1 AND course_id = $2)", enrollment.UserID, enrollment.CourseID).Scan(&enrollmentExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if enrollmentExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already enrolled in this course"})
		return
	}

	var enrollmentID int
	err = db.QueryRow(
		"INSERT INTO user_course_enrollments (user_id, course_id, status) VALUES ($1, $2, $3) RETURNING id",
		enrollment.UserID, enrollment.CourseID, enrollment.Status,
	).Scan(&enrollmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	enrollment.ID = enrollmentID
	c.JSON(http.StatusCreated, enrollment)
}

func deleteUserEnrollment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	db := getDB(c)
	if db == nil {
		return
	}

	result, err := db.Exec("DELETE FROM user_course_enrollments WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Enrollment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Enrollment deleted successfully"})
}

func getDB(c *gin.Context) *sql.DB {
	dbVal, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available"})
		return nil
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database connection type"})
		return nil
	}
	return db
}