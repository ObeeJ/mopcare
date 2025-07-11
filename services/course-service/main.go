package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Course struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	OverviewVideoURL string    `json:"overview_video_url"`
	CoverImageURL    string    `json:"cover_image_url"`
	UniqueID         string    `json:"unique_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type Series struct {
	ID          int       `json:"id"`
	CourseID    int       `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		Prefork:      false, // Disabled for Docker compatibility
		ServerHeader: "Course-Service",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"service": "course-service", "status": "running"})
	})

	app.Post("/courses", createCourse)
	app.Get("/courses", getCourses)
	app.Get("/courses/:id", getCourse)
	app.Put("/courses/:id", updateCourse)
	app.Delete("/courses/:id", deleteCourse)

	app.Get("/courses/:id/series", getSeriesForCourse)
	app.Get("/series/:id", getSeriesByID)
	app.Post("/courses/:id/series", createSeriesForCourse)
	app.Put("/series/:id", updateSeries)
	app.Delete("/series/:id", deleteSeries)

	port := os.Getenv("COURSE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Course service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func connectDB() (*sql.DB, error) {
	connStr := os.Getenv("SUPABASE_DB_URL")
	if connStr == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("supabase_db_url environment variable is not set and .env file not found")
		}
		connStr = os.Getenv("SUPABASE_DB_URL")
		if connStr == "" {
			return nil, fmt.Errorf("supabase_db_url not found in environment or .env")
		}
	}
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %v", err)
	}
	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}
	fmt.Println("Database connection established successfully.")
	return database, nil
}

func createCourse(c *fiber.Ctx) error {
	var newCourse struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		OverviewVideoURL string `json:"overview_video_url"`
		CoverImageURL    string `json:"cover_image_url"`
		UniqueID         string `json:"unique_id"`
	}
	if err := c.BodyParser(&newCourse); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if newCourse.Title == "" || newCourse.Content == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title and content are required"})
	}

	var id int
	var createdAt time.Time
	err := db.QueryRow(
		`INSERT INTO courses (title, content, overview_video_url, cover_image_url, unique_id) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		newCourse.Title, newCourse.Content, newCourse.OverviewVideoURL, newCourse.CoverImageURL, newCourse.UniqueID,
	).Scan(&id, &createdAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	course := Course{
		ID:               id,
		Title:            newCourse.Title,
		Content:          newCourse.Content,
		OverviewVideoURL: newCourse.OverviewVideoURL,
		CoverImageURL:    newCourse.CoverImageURL,
		UniqueID:         newCourse.UniqueID,
		CreatedAt:        createdAt,
	}
	return c.Status(201).JSON(course)
}

func getCourses(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, title, content, overview_video_url, cover_image_url, unique_id, created_at FROM courses")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Title, &course.Content, &course.OverviewVideoURL, &course.CoverImageURL, &course.UniqueID, &course.CreatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		courses = append(courses, course)
	}
	return c.JSON(courses)
}

func getCourse(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var course Course
	err = db.QueryRow(
		"SELECT id, title, content, overview_video_url, cover_image_url, unique_id, created_at FROM courses WHERE id = $1",
		id,
	).Scan(&course.ID, &course.Title, &course.Content, &course.OverviewVideoURL, &course.CoverImageURL, &course.UniqueID, &course.CreatedAt)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Course not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(course)
}

func updateCourse(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var updateData struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		OverviewVideoURL string `json:"overview_video_url"`
		CoverImageURL    string `json:"cover_image_url"`
		UniqueID         string `json:"unique_id"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = db.Exec(
		`UPDATE courses SET title = $1, content = $2, overview_video_url = $3, cover_image_url = $4, unique_id = $5 WHERE id = $6`,
		updateData.Title, updateData.Content, updateData.OverviewVideoURL, updateData.CoverImageURL, updateData.UniqueID, id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update course"})
	}

	return c.JSON(fiber.Map{"message": "Course updated successfully"})
}

func deleteCourse(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	_, err = db.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete course"})
	}

	return c.JSON(fiber.Map{"message": "Course deleted successfully"})
}

func getSeriesForCourse(c *fiber.Ctx) error {
	courseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	rows, err := db.Query("SELECT id, course_id, title, description, created_at FROM series WHERE course_id = $1", courseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.CreatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		seriesList = append(seriesList, s)
	}
	return c.JSON(seriesList)
}

func getSeriesByID(c *fiber.Ctx) error {
	seriesID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid series ID"})
	}

	var s Series
	err = db.QueryRow("SELECT id, course_id, title, description, created_at FROM series WHERE id = $1", seriesID).
		Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Series not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

func createSeriesForCourse(c *fiber.Ctx) error {
	courseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var newSeries struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&newSeries); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if newSeries.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	var id int
	var createdAt time.Time
	err = db.QueryRow(
		`INSERT INTO series (course_id, title, description) VALUES ($1, $2, $3) RETURNING id, created_at`,
		courseID, newSeries.Title, newSeries.Description,
	).Scan(&id, &createdAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	series := Series{
		ID:          id,
		CourseID:    courseID,
		Title:       newSeries.Title,
		Description: newSeries.Description,
		CreatedAt:   createdAt,
	}
	return c.Status(201).JSON(series)
}

func updateSeries(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid series ID"})
	}

	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = db.Exec(
		`UPDATE series SET title = $1, description = $2 WHERE id = $3`,
		updateData.Title, updateData.Description, id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update series"})
	}

	return c.JSON(fiber.Map{"message": "Series updated successfully"})
}

func deleteSeries(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid series ID"})
	}

	_, err = db.Exec("DELETE FROM series WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete series"})
	}

	return c.JSON(fiber.Map{"message": "Series deleted successfully"})
}