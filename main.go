package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

	// Course routes
	router.POST("/courses", createCourse)
	router.GET("/courses", getCourses)
	router.GET("/courses/:id", getCourse)
	router.PUT("/courses/:id", updateCourse)
	router.DELETE("/courses/:id", deleteCourse)

	// Series routes
	router.GET("/courses/:id/series", getSeriesForCourse)
	router.GET("/series/:id", getSeriesByID)
	router.POST("/courses/:id/series", createSeriesForCourse)
	router.PUT("/series/:id", updateSeries)
	router.DELETE("/series/:id", deleteSeries)

	// User routes
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.DELETE("/users/:id", deleteUser)

	// Enrollment routes
	router.GET("/users/:id/enrollments", getUserEnrollments)
	router.POST("/users/:id/enrollments", createUserEnrollment)
	router.DELETE("/enrollments/:id", deleteUserEnrollment)

	// Profile route
	router.GET("/users/:id/profile", getUserProfile)

	// Payment route
	router.PUT("/users/:id/payment", updateUserPayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
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

// createCourse handles the POST /courses endpoint to create a new course
func createCourse(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		handleCourseUpload(c)
		return
	}

	// Fallback to JSON creation
	var newCourse struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		OverviewVideoURL string `json:"overview_video_url"` // optional but allowed
		CoverImageURL    string `json:"cover_image_url"`    // optional but allowed
	}
	if err := c.BindJSON(&newCourse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newCourse.Title == "" || newCourse.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Check if title exists
	if dbValExists(c, newCourse.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course with this title already exists"})
		return
	}

	var id int
	var createdAt time.Time

	uniqueID, err := generateUniqueId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique ID"})
		return
	}

	// Store course including optional media URLs
	err = db.QueryRow(
		`INSERT INTO courses (title, content, overview_video_url, cover_image_url, unique_id) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, created_at`,
		newCourse.Title, newCourse.Content, newCourse.OverviewVideoURL, newCourse.CoverImageURL, uniqueID,
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	course := Course{
		ID:               id,
		Title:            newCourse.Title,
		Content:          newCourse.Content,
		OverviewVideoURL: newCourse.OverviewVideoURL,
		CoverImageURL:    newCourse.CoverImageURL,
		CreatedAt:        createdAt,
	}
	c.JSON(http.StatusCreated, course)
}

// getCourses handles the GET /courses endpoint to retrieve all courses
func getCourses(c *gin.Context) {
	dbVal, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available in context"})
		return
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database connection type in context"})
		return
	}
	rows, err := db.Query("SELECT id, title, content, overview_video_url, cover_image_url, created_at FROM courses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Title, &course.Content, &course.OverviewVideoURL, &course.CoverImageURL, &course.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		courses = append(courses, course)
	}
	c.JSON(http.StatusOK, courses)
}

// getCourse handles the GET /courses/:id endpoint to retrieve a specific course
func getCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if !dbValInRange(int64(id), 1, 1000000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID must be between 1 and 1,000,000"})
		return
	}

	// Retrieve the database connection from the context
	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Query the database for the course with the given ID
	var course Course
	err = db.QueryRow(
		"SELECT id, title, content, overview_video_url, cover_image_url, created_at FROM courses WHERE id = $1",
		id,
	).Scan(&course.ID, &course.Title, &course.Content, &course.OverviewVideoURL, &course.CoverImageURL, &course.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

// updateCourse handles the PUT /courses/:id endpoint to update an existing course
func updateCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse form"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
		return
	}

	// Optional file uploads
	coverFile, coverHeader, _ := c.Request.FormFile("cover_image")
	videoFile, videoHeader, _ := c.Request.FormFile("overview_video")

	var coverURL, videoURL string

	if coverFile != nil {
		coverURL, err = uploadToSupabase(c, coverHeader, "course-assets")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload cover image"})
			return
		}
	}

	if videoFile != nil {
		videoURL, err = uploadToSupabase(c, videoHeader, "course-assets")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload overview video"})
			return
		}
	}

	// Perform update
	query := `
		UPDATE courses 
		SET title = $1, content = $2,
		    cover_image_url = COALESCE($3, cover_image_url),
		    overview_video_url = COALESCE($4, overview_video_url)
		WHERE id = $5
	`
	_, err = db.Exec(query, title, content, nullIfEmpty(coverURL), nullIfEmpty(videoURL), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// deleteCourse handles the DELETE /courses/:id endpoint to delete a course
func deleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID must be a positive integer"})
		return
	}

	if id >= 1000000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is too large"})
		return
	}

	dbVal, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available in context"})
		return
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection type"})
		return
	}

	_, err = db.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// getSeriesForCourse handles GET /courses/:id/series to fetch all series for a course
type Series struct {
	ID            int       `json:"id"`
	CourseID      int       `json:"course_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	VideoURL      string    `json:"video_url"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	Duration      int       `json:"duration"`        // seconds
	IsFreePreview bool      `json:"is_free_preview"` // true or false
	CreatedAt     time.Time `json:"created_at"`
}

func getSeriesForCourse(c *gin.Context) {
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	rows, err := db.Query("SELECT id, course_id, title, description, created_at FROM series WHERE course_id = $1", courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		seriesList = append(seriesList, s)
	}
	c.JSON(http.StatusOK, seriesList)
}

func getSeriesByID(c *gin.Context) {
	idStr := c.Param("id")
	seriesID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	var s Series
	err = db.QueryRow("SELECT id, course_id, title, description, created_at FROM series WHERE id = $1", seriesID).Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func createSeriesForCourse(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		handleSeriesUpload(c)
		return
	}

	// Fallback to JSON creation
	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var newSeries struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		VideoURL      string `json:"video_url"`
		Duration      int    `json:"duration"`
		IsFreePreview bool   `json:"is_free_preview"`
	}
	if err := c.BindJSON(&newSeries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newSeries.Title == "" || newSeries.Description == "" || newSeries.VideoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description and video URL are required"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Generate episode number
	epNumber, err := getNextEpisodeNumber(db, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate episode label"})
		return
	}
	epTitle := fmt.Sprintf("Ep%d: %s", epNumber, newSeries.Title)

	var id int
	var createdAt time.Time
	err = db.QueryRow(
		"INSERT INTO series (course_id, title, description, video_url, duration, is_free_preview) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at",
		courseID, epTitle, newSeries.Description, newSeries.VideoURL, newSeries.Duration, newSeries.IsFreePreview,
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s := Series{
		ID:            id,
		CourseID:      courseID,
		Title:         epTitle,
		Description:   newSeries.Description,
		VideoURL:      newSeries.VideoURL,
		Duration:      newSeries.Duration,
		IsFreePreview: newSeries.IsFreePreview,
		CreatedAt:     createdAt,
	}

	c.JSON(http.StatusCreated, s)
}

func handleSeriesUpload(c *gin.Context) {
	// Parse multipart form with 32MB max memory
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
		return
	}

	courseIDStr := c.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Get form values
	title := c.PostForm("title")
	description := c.PostForm("description")
	durationStr := c.PostForm("duration")
	isFreePreviewStr := c.PostForm("is_free_preview")

	if title == "" || description == "" || durationStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description and duration are required"})
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration value"})
		return
	}

	isFreePreview := isFreePreviewStr == "true"

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Generate episode number
	epNumber, err := getNextEpisodeNumber(db, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate episode label"})
		return
	}
	epTitle := fmt.Sprintf("Ep%d: %s", epNumber, title)

	// Process uploaded files
	var videoURL, thumbnailURL string

	// Handle video upload
	videoFile, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file is required"})
		return
	}
	file, err := videoFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open video file"})
		return
	}
	defer file.Close()
	videoURL, err = uploadToSupabase(c, videoFile, "series-videos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Video upload failed: " + err.Error()})
		return
	}

	// Handle thumbnail upload (optional)
	thumbnailFile, err := c.FormFile("thumbnail")
	if err == nil {
		file, err := thumbnailFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open thumbnail file"})
			return
		}
		defer file.Close()
		thumbnailURL, err = uploadToSupabase(c, thumbnailFile, "series-thumbnails")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Thumbnail upload failed: " + err.Error()})
			return
		}
	}

	var id int
	var createdAt time.Time
	err = db.QueryRow(
		"INSERT INTO series (course_id, title, description, video_url, thumbnail_url, duration, is_free_preview) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at",
		courseID, epTitle, description, videoURL, thumbnailURL, duration, isFreePreview,
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s := Series{
		ID:            id,
		CourseID:      courseID,
		Title:         epTitle,
		Description:   description,
		VideoURL:      videoURL,
		ThumbnailURL:  thumbnailURL,
		Duration:      duration,
		IsFreePreview: isFreePreview,
		CreatedAt:     createdAt,
	}

	c.JSON(http.StatusCreated, s)
}

func updateSeries(c *gin.Context) {
	idStr := c.Param("id")
	seriesID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" && description == "" && c.Request.MultipartForm.File["video"] == nil && c.Request.MultipartForm.File["thumbnail"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update fields provided"})
		return
	}

	var videoURL, thumbURL string

	videoFile, err := c.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get video file"})
		return
	}
	if videoFile != nil {
		videoURL, err = uploadToSupabase(c, videoFile, "series-videos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
			return
		}
	}

	thumbFile, err := c.FormFile("thumbnail")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get thumbnail file"})
		return
	}
	if thumbFile != nil {
		thumbURL, err = uploadToSupabase(c, thumbFile, "series-thumbnails")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload thumbnail"})
			return
		}
	}

	query := `
		UPDATE series 
		SET 
			title = COALESCE(NULLIF($1, ''), title),
			description = COALESCE(NULLIF($2, ''), description),
			video_url = COALESCE(NULLIF($3, ''), video_url),
			thumbnail_url = COALESCE(NULLIF($4, ''), thumbnail_url)
		WHERE id = $5
	`

	_, err = db.Exec(query, title, description, nullIfEmpty(videoURL), nullIfEmpty(thumbURL), seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update series"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Series updated successfully"})
}

func deleteSeries(c *gin.Context) {
	idStr := c.Param("id")
	seriesID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	result, err := db.Exec("DELETE FROM series WHERE id = $1", seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete series"})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Series deleted successfully"})
}

// Course represents the structure of a course in the database
type Course struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	OverviewVideoURL string    `json:"overview_video_url"`
	CoverImageURL    string    `json:"cover_image_url"`
	CreatedAt        time.Time `json:"created_at"`
}

func handleCourseUpload(c *gin.Context) {
	// Parse multipart form with 32MB max memory
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
		return
	}

	// Get form values
	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and content are required"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Check for duplicate title
	if dbValExists(c, title) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course with this title already exists"})
		return
	}

	// Process uploaded files
	var overviewVideoURL, coverImageURL string

	// Handle overview video upload
	overviewVideoFile, overviewVideoHeader, err := c.Request.FormFile("overview_video")
	if err == nil {
		defer overviewVideoFile.Close()
		overviewVideoURL, err = uploadToSupabase(c, overviewVideoHeader, "course-assets")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Video upload failed: " + err.Error()})
			return
		}
	}

	// Handle cover image upload
	coverImageFile, coverImageHeader, err := c.Request.FormFile("cover_image")
	if err == nil {
		defer coverImageFile.Close()
		coverImageURL, err = uploadToSupabase(c, coverImageHeader, "course-assets")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cover image upload failed: " + err.Error()})
			return
		}
	}

	var id int
	var createdAt time.Time
	uniqueID, err := generateUniqueId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique ID"})
		return
	}

	err = db.QueryRow(
		`INSERT INTO courses (title, content, overview_video_url, cover_image_url, unique_id) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, created_at`,
		title, content, overviewVideoURL, coverImageURL, uniqueID,
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	course := Course{
		ID:               id,
		Title:            title,
		Content:          content,
		OverviewVideoURL: overviewVideoURL,
		CoverImageURL:    coverImageURL,
		CreatedAt:        createdAt,
	}
	c.JSON(http.StatusCreated, course)
}

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

func getUsers(c *gin.Context) {
	db, ok := dbVal(c)
	if !ok {
		return
	}

	rows, err := db.Query("SELECT id, first_name, last_name, email, total_amount_paid, created_at, enrolled_courses, completed_courses_count, state, city FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.TotalAmountPaid, &user.CreatedAt, &user.EnrolledCourses, &user.CompletedCoursesCount, &user.State, &user.City); err != nil {
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

type UserCourseEnrollment struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	CourseID int    `json:"course_id"`
	Status   string `json:"status"`
}

func getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	var user User
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, total_amount_paid, created_at, enrolled_courses, completed_courses_count, state, city FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.TotalAmountPaid, &user.CreatedAt, &user.EnrolledCourses, &user.CompletedCoursesCount, &user.State, &user.City)
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

	dbVal, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection type"})
		return
	}

	// Check if a user with the same email already exists
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

	// Insert the new user
	var id int
	var createdAt time.Time
	err = db.QueryRow(
		"INSERT INTO users (first_name, last_name, email, total_amount_paid, enrolled_courses, completed_courses_count, state, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at",
		user.FirstName, user.LastName, user.Email, user.TotalAmountPaid, user.EnrolledCourses, user.CompletedCoursesCount, user.State, user.City,
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

	db, ok := dbVal(c)
	if !ok {
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

func getUserEnrollments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
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
	// Validate status
	if enrollment.Status != "enrolled" && enrollment.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'enrolled' or 'completed'"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Check if user exists
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

	// Check if course exists
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

	// Check if enrollment already exists
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

	// Insert the enrollment
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

	db, ok := dbVal(c)
	if !ok {
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

func getUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db, ok := dbVal(c)
	if !ok {
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

	// Fetch state and city for the profile
	err = db.QueryRow("SELECT state, city FROM users WHERE id = $1", id).Scan(&user.State, &user.City)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	db, ok := dbVal(c)
	if !ok {
		return
	}

	// Check if user exists
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

	// Update total_amount_paid
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

