package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ssh-file-transfer-api/internal/models"
	"github.com/yourusername/ssh-file-transfer-api/internal/ssh"
)

// TransferFile handles file transfer requests
func (s *Server) TransferFile(c *gin.Context) {
	var req models.FileTransferRequest
	
	// Validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	// Use default SSH port if not specified
	if req.TargetPort == 0 {
		req.TargetPort = 22
	}

	// Create SSH client
	username := req.Username
	if username == "" {
		username = s.config.SSHUsername
	}

	password := req.Password
	if password == "" {
		password = s.config.SSHPassword
	}

	keyPath := req.PrivateKeyPath
	if keyPath == "" {
		keyPath = s.config.SSHKeyPath
	}

	client, err := ssh.NewClient(username, password, keyPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create SSH client",
			Details: err.Error(),
		})
		return
	}

	// Execute file transfer
	response, err := client.TransferFile(req)
	if err != nil {
		// If response already contains error details, use it
		if response != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to transfer file",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFileTransferStatus retrieves the status of a file transfer
func (s *Server) GetFileTransferStatus(c *gin.Context) {
	id := c.Param("id")
	
	// In a real implementation, you would look up the transfer status from a database
	// This is a simplified example
	c.JSON(http.StatusOK, models.FileTransferStatus{
		ID:               id,
		Status:           "completed",
		PercentComplete:  100.0,
		BytesTransferred: 1024,
		StartTime:        time.Now().Add(-1 * time.Minute),
		LastUpdated:      time.Now(),
	})
}

// HealthCheck provides a basic health check endpoint
func (s *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
