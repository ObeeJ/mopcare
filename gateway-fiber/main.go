package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type CacheEntry struct {
	Data      string
	ExpiresAt time.Time
}

type Cache struct {
	store sync.Map
	ttl   time.Duration
}

func NewCache() *Cache {
	return &Cache{ttl: 5 * time.Minute}
}

func (c *Cache) Get(key string) (string, bool) {
	if val, ok := c.store.Load(key); ok {
		entry := val.(CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, true
		}
		c.store.Delete(key)
	}
	return "", false
}

func (c *Cache) Set(key, data string) {
	c.store.Store(key, CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	})
}

type Metrics struct {
	TotalRequests int64
	CacheHits     int64
	CacheMisses   int64
	mu            sync.RWMutex
}

func (m *Metrics) IncrementRequests() {
	m.mu.Lock()
	m.TotalRequests++
	m.mu.Unlock()
}

func (m *Metrics) IncrementCacheHits() {
	m.mu.Lock()
	m.CacheHits++
	m.mu.Unlock()
}

func (m *Metrics) IncrementCacheMisses() {
	m.mu.Lock()
	m.CacheMisses++
	m.mu.Unlock()
}

var (
	cache   = NewCache()
	metrics = &Metrics{}
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       false, // Disabled for Docker compatibility
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Mopcare-Gateway",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "mopcare-api-gateway",
			"version": "1.0.0",
		})
	})

	app.Get("/metrics", func(c *fiber.Ctx) error {
		metrics.mu.RLock()
		defer metrics.mu.RUnlock()
		return c.JSON(fiber.Map{
			"gateway": fiber.Map{
				"total_requests": metrics.TotalRequests,
				"cache_hits":     metrics.CacheHits,
				"cache_misses":   metrics.CacheMisses,
			},
		})
	})

	app.Use(proxyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("GATEWAY_PORT")
		if port == "" {
			port = "10000"
		}
	}

	fmt.Printf("🚀 Fiber Gateway starting on port %s\n", port)
	app.Listen(":" + port)
}

func proxyHandler(c *fiber.Ctx) error {
	metrics.IncrementRequests()
	path := c.Path()
	method := c.Method()

	var targetURL string
	// For single-service deployment, all services run in the same container
	courseServiceURL := os.Getenv("COURSE_SERVICE_URL")
	if courseServiceURL == "" {
		courseServiceURL = "http://localhost:8081"
	}
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8082"
	}
	enrollmentServiceURL := os.Getenv("ENROLLMENT_SERVICE_URL")
	if enrollmentServiceURL == "" {
		enrollmentServiceURL = "http://localhost:8083"
	}

	// For Render deployment, return mock responses since services aren't running
	if os.Getenv("RENDER") != "" {
		return handleMockResponse(c, path, method)
	}

	if strings.HasPrefix(path, "/courses") || strings.HasPrefix(path, "/series") {
		targetURL = courseServiceURL
	} else if strings.HasPrefix(path, "/users") && !strings.Contains(path, "/enrollments") {
		targetURL = userServiceURL
	} else if strings.Contains(path, "/enrollments") {
		targetURL = enrollmentServiceURL
	} else {
		return c.Status(404).JSON(fiber.Map{"error": "Service not found"})
	}

	cacheKey := fmt.Sprintf("%s:%s", method, path)
	if method == "GET" {
		if cachedData, found := cache.Get(cacheKey); found {
			metrics.IncrementCacheHits()
			c.Set("X-Cache", "HIT")
			return c.SendString(cachedData)
		}
		metrics.IncrementCacheMisses()
		c.Set("X-Cache", "MISS")
	}

	return proxy.Do(c, targetURL+path)
}

func handleMockResponse(c *fiber.Ctx, path, method string) error {
	if method == "GET" && path == "/courses" {
		return c.JSON([]fiber.Map{
			{"id": 1, "title": "Managing Diabetes in Your Golden Years", "content": "Comprehensive guide to diabetes management for seniors", "unique_id": "diabetes-seniors-101"},
			{"id": 2, "title": "Heart Health After 65", "content": "Essential cardiovascular care for seniors", "unique_id": "heart-health-seniors"},
		})
	}
	if method == "POST" && path == "/courses" {
		return c.Status(201).JSON(fiber.Map{"id": 3, "title": "New Course", "message": "Course created successfully"})
	}
	if method == "GET" && path == "/users" {
		return c.JSON([]fiber.Map{{"id": 1, "first_name": "Margaret", "last_name": "Johnson", "email": "margaret.johnson@email.com"}})
	}
	if method == "POST" && path == "/users" {
		return c.Status(201).JSON(fiber.Map{"id": 2, "first_name": "New User", "message": "User created successfully"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Mock response for " + method + " " + path})
}