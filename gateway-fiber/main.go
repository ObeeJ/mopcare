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

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "9090"
	}

	fmt.Printf("🚀 Fiber Gateway starting on port %s\n", port)
	app.Listen(":" + port)
}

func proxyHandler(c *fiber.Ctx) error {
	metrics.IncrementRequests()
	path := c.Path()
	method := c.Method()

	var targetURL string
	// Use Docker service names for container communication
	courseServiceURL := os.Getenv("COURSE_SERVICE_URL")
	if courseServiceURL == "" {
		courseServiceURL = "http://course-service:8081"
	}
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8082"
	}
	enrollmentServiceURL := os.Getenv("ENROLLMENT_SERVICE_URL")
	if enrollmentServiceURL == "" {
		enrollmentServiceURL = "http://enrollment-service:8083"
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