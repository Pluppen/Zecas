package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware adds CORS headers to responses
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		log.Println("Cors working as hard as it can")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogger logs all requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		log.Printf("[%s] %s %s %d %s - %s", method, path, clientIP, statusCode, latency, c.Errors.String())
	}
}

// RateLimiter limits requests based on IP address
// This is a simple example; in production use a proper rate limiting library
func RateLimiter() gin.HandlerFunc {
	// Store IPs and their request times
	ipMap := make(map[string][]time.Time)
	const maxRequests = 60 // Max requests per minute
	const duration = time.Minute

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Clean up old requests
		now := time.Now()
		newRequests := []time.Time{}
		for _, t := range ipMap[ip] {
			if now.Sub(t) < duration {
				newRequests = append(newRequests, t)
			}
		}
		ipMap[ip] = newRequests

		// Check if too many requests
		if len(ipMap[ip]) >= maxRequests {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		// Add this request
		ipMap[ip] = append(ipMap[ip], now)

		c.Next()
	}
}

// ErrorHandler catches any panics and returns a 500 error
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(500, gin.H{"error": "Internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Authentication middleware (placeholder - would need to be properly implemented)
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		token := c.GetHeader("Authorization")

		// In a real implementation, validate the token here
		// For now, just check if it exists
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user ID in context (placeholder - real implementation would extract from token)
		c.Set("userID", "example-user-id")

		c.Next()
	}
}
