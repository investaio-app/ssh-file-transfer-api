package models

import "time"

// FileTransferRequest defines the request payload for file transfer
type FileTransferRequest struct {
	TargetHost     string `json:"target_host" binding:"required"`
	TargetPort     int    `json:"target_port" binding:"required"`
	SourceFilePath string `json:"source_file_path" binding:"required"`
	TargetFilePath string `json:"target_file_path" binding:"required"`
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	PrivateKeyPath string `json:"private_key_path,omitempty"`
}

// FileTransferResponse defines the response payload for file transfer
type FileTransferResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	SourceFile   string    `json:"source_file"`
	TargetFile   string    `json:"target_file"`
	TargetHost   string    `json:"target_host"`
	BytesWritten int64     `json:"bytes_written"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     string    `json:"duration"`
	Error        string    `json:"error,omitempty"`
}

// FileTransferStatus represents a status update for a file transfer operation
type FileTransferStatus struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	PercentComplete float64   `json:"percent_complete"`
	BytesTransferred int64     `json:"bytes_transferred"`
	StartTime     time.Time `json:"start_time"`
	LastUpdated   time.Time `json:"last_updated"`
	Error         string    `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
