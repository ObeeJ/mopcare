package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func dbValExists(c *gin.Context, val string) bool {
	dbVal, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available in context"})
		return false
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database connection"})
		return false
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM courses WHERE title = $1)", val).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	return exists
}

func dbVal(c *gin.Context) (*sql.DB, bool) {
	dbVal, ok := c.Get("db")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection not available in context"})
		return nil, false
	}
	db, ok := dbVal.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database connection type in context"})
		return nil, false
	}
	return db, true
}

func dbValInRange(val int64, min, max int64) bool {
	return val >= min && val <= max
}

func generateUniqueId() (string, error) {
	return uuid.New().String(), nil
}

func getNextEpisodeNumber(db *sql.DB, courseID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM series WHERE course_id = $1", courseID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}

func uploadToSupabase(c *gin.Context, file *multipart.FileHeader, bucket string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_PROJECT_URL")
	if supabaseURL == "" {
		return "", fmt.Errorf("SUPABASE_PROJECT_URL environment variable not set")
	}

	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
	if supabaseKey == "" {
		return "", fmt.Errorf("SUPABASE_SERVICE_KEY environment variable not set")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err = io.Copy(part, src); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), "POST",
		supabaseURL+"/storage/v1/object/"+bucket+"/"+file.Filename, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", fmt.Errorf("upload canceled by client")
		}
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return supabaseURL + "/storage/v1/object/public/" + bucket + "/" + file.Filename, nil
}

func nullIfEmpty(val string) interface{} {
	if val == "" {
		return nil
	}
	return val
}
