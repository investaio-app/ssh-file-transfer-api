package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ssh-file-transfer-api/internal/models"
)

// Logger middleware logs all requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Determine final response status based on errors
		if len(c.Errors) > 0 {
			// Log errors
			gin.DefaultErrorWriter.Write([]byte(c.Errors.String()))
		}

		gin.DefaultWriter.Write([]byte(
			"[API] " + time.Now().Format("2006/01/02 - 15:04:05") + 
			" | " + method + " | " + path + 
			" | " + http.StatusText(statusCode) + 
			" | " + latency.String() + "\n",
		))
	}
}

// RateLimiter middleware limits request rates
func RateLimiter(requests int, duration time.Duration) gin.HandlerFunc {
	// In a real implementation, you'd use a proper rate limiter
	// This is a simplified example
	type client struct {
		count    int
		lastSeen time.Time
	}
	
	clients := make(map[string]*client)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		// Clean up old entries
		now := time.Now()
		for clientIP, clientData := range clients {
			if now.Sub(clientData.lastSeen) > duration {
				delete(clients, clientIP)
			}
		}
		
		// Check if client exists
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{count: 0, lastSeen: now}
		}
		
		// Update last seen
		clients[ip].lastSeen = now
		
		// Check if rate limit exceeded
		if clients[ip].count >= requests {
			c.JSON(http.StatusTooManyRequests, models.APIError{
				Code:    http.StatusTooManyRequests,
				Message: "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		
		// Increment request count
		clients[ip].count++
		
		c.Next()
	}
}

// ErrorHandler middleware handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()
			
			// Return appropriate error response
			c.JSON(http.StatusInternalServerError, models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
				Details: err.Error(),
			})
		}
	}
}
