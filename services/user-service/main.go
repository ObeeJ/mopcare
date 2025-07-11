package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID                    int       `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	Email                 string    `json:"email"`
	TotalAmountPaid       float64   `json:"total_amount_paid"`
	CreatedAt             time.Time `json:"created_at"`
	EnrolledCourses       string    `json:"enrolled_courses,omitempty"`
	CompletedCoursesCount int64     `json:"completed_courses_count,omitempty"`
	State                 string    `json:"state,omitempty"`
	City                  string    `json:"city,omitempty"`
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
		c.JSON(http.StatusOK, gin.H{"service": "user-service", "status": "running"})
	})

	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/:id/profile", getUserProfile)
	router.PUT("/users/:id/payment", updateUserPayment)

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start user service: %v", err)
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

func getUsers(c *gin.Context) {
	db := getDB(c)
	if db == nil {
		return
	}

	rows, err := db.Query("SELECT id, first_name, last_name, email, total_amount_paid, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.TotalAmountPaid, &user.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
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

	var user User
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, total_amount_paid, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.TotalAmountPaid, &user.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name, last name, and email are required"})
		return
	}

	db := getDB(c)
	if db == nil {
		return
	}

	var emailExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&emailExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if emailExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	var id int
	var createdAt time.Time
	err = db.QueryRow(
		"INSERT INTO users (first_name, last_name, email, total_amount_paid) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		user.FirstName, user.LastName, user.Email, user.TotalAmountPaid,
	).Scan(&id, &createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	user.CreatedAt = createdAt
	c.JSON(http.StatusCreated, user)
}

func deleteUser(c *gin.Context) {
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

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
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
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func getUserProfile(c *gin.Context) {
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

	var user User
	err = db.QueryRow("SELECT id, first_name, last_name, email, total_amount_paid, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.TotalAmountPaid, &user.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	enrolledCoursesCount, err := getEnrolledCoursesCount(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count enrollments"})
		return
	}

	completedCoursesCount, err := getCompletedCoursesCount(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count completed courses"})
		return
	}

	// Skip state/city for now as columns may not exist
	// err = db.QueryRow("SELECT state, city FROM users WHERE id = $1", id).Scan(&user.State, &user.City)

	type userProfile struct {
		FirstName             string  `json:"first_name"`
		LastName              string  `json:"last_name"`
		Email                 string  `json:"email"`
		TotalAmountPaid       float64 `json:"total_amount_paid"`
		EnrolledCoursesCount  int64   `json:"enrolled_courses_count"`
		CompletedCoursesCount int64   `json:"completed_courses_count"`
		State                 string  `json:"state,omitempty"`
		City                  string  `json:"city,omitempty"`
	}

	c.JSON(http.StatusOK, userProfile{
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Email:                 user.Email,
		TotalAmountPaid:       user.TotalAmountPaid,
		EnrolledCoursesCount:  enrolledCoursesCount,
		CompletedCoursesCount: completedCoursesCount,
		State:                 user.State,
		City:                  user.City,
	})
}

func updateUserPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var payment struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if payment.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	db := getDB(c)
	if db == nil {
		return
	}

	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&userExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	result, err := db.Exec("UPDATE users SET total_amount_paid = total_amount_paid + $1 WHERE id = $2", payment.Amount, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully"})
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

func getEnrolledCoursesCount(db *sql.DB, userID int) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) FROM user_course_enrollments WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getCompletedCoursesCount(db *sql.DB, userID int) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) FROM user_course_enrollments WHERE user_id = $1 AND status = 'completed'", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}