package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"github.com/yourusername/ssh-file-transfer-api/internal/models"
)

// Client manages SSH connections and file transfers
type Client struct {
	config *ssh.ClientConfig
}

// NewClient creates a new SSH client
func NewClient(username, password, keyPath string) (*Client, error) {
	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use ssh.FixedHostKey or ssh.KnownHosts
		Timeout:         30 * time.Second,
	}

	// Use password if provided
	if password != "" {
		config.Auth = []ssh.AuthMethod{
			ssh.Password(password),
		}
	} else if keyPath != "" {
		// Otherwise use private key
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}

		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else {
		return nil, fmt.Errorf("either password or keyPath must be provided")
	}

	return &Client{config: config}, nil
}

// TransferFile transfers a file to a remote server
func (c *Client) TransferFile(req models.FileTransferRequest) (*models.FileTransferResponse, error) {
	startTime := time.Now()
	
	// Override client config with request-specific authentication if provided
	config := c.config
	if req.Username != "" || req.Password != "" || req.PrivateKeyPath != "" {
		config = &ssh.ClientConfig{
			User:            req.Username,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		
		if req.Password != "" {
			config.Auth = []ssh.AuthMethod{
				ssh.Password(req.Password),
			}
		} else if req.PrivateKeyPath != "" {
			key, err := os.ReadFile(req.PrivateKeyPath)
			if err != nil {
				return nil, fmt.Errorf("unable to read private key: %v", err)
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return nil, fmt.Errorf("unable to parse private key: %v", err)
			}

			config.Auth = []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			}
		}
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", req.TargetHost, req.TargetPort)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to connect to SSH server: %v", err),
		}, err
	}
	defer conn.Close()

	// Create SFTP client
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to create SFTP client: %v", err),
		}, err
	}
	defer sftpClient.Close()

	// Open local file
	srcFile, err := os.Open(req.SourceFilePath)
	if err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to open source file: %v", err),
		}, err
	}
	defer srcFile.Close()

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(req.TargetFilePath)
	if err := sftpClient.MkdirAll(targetDir); err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to create target directory: %v", err),
		}, err
	}

	// Create remote file
	dstFile, err := sftpClient.Create(req.TargetFilePath)
	if err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to create target file: %v", err),
		}, err
	}
	defer dstFile.Close()

	// Copy file contents
	bytesWritten, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return &models.FileTransferResponse{
			Status:    "failed",
			SourceFile: req.SourceFilePath,
			TargetFile: req.TargetFilePath,
			TargetHost: req.TargetHost,
			BytesWritten: bytesWritten,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  time.Since(startTime).String(),
			Error:     fmt.Sprintf("failed to copy file contents: %v", err),
		}, err
	}

	// Success response
	endTime := time.Now()
	return &models.FileTransferResponse{
		ID:         fmt.Sprintf("transfer-%d", time.Now().UnixNano()),
		Status:     "completed",
		SourceFile: req.SourceFilePath,
		TargetFile: req.TargetFilePath,
		TargetHost: req.TargetHost,
		BytesWritten: bytesWritten,
		StartTime:  startTime,
		EndTime:    endTime,
		Duration:   endTime.Sub(startTime).String(),
	}, nil
}
