package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ssh-file-transfer-api/internal/config"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	
	server := &Server{
		router: router,
		config: cfg,
	}
	
	// Apply middlewares
	router.Use(Logger())
	router.Use(ErrorHandler())
	router.Use(RateLimiter(cfg.RateLimitRequests, cfg.RateLimitDuration))
	
	// Setup routes
	server.setupRoutes()
	
	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.HealthCheck)
	
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// File transfer endpoints
		transfers := v1.Group("/transfers")
		{
			transfers.POST("", s.TransferFile)
			transfers.GET("/:id", s.GetFileTransferStatus)
		}
	}
}

// Run starts the API server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
