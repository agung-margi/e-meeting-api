package util

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveBase64Image decodes Base64 image and saves it. Returns the public URL.
func SaveBase64Image(base64Image string, uploadDir string) (string, error) {
	dataParts := strings.Split(base64Image, ",")
	if len(dataParts) != 2 {
		return "", fmt.Errorf("invalid Base64 format")
	}

	// Validate image format
	if !(strings.Contains(dataParts[0], "image/jpeg") || strings.Contains(dataParts[0], "image/png")) {
		return "", fmt.Errorf("only JPEG and PNG formats are supported")
	}

	decodedData, err := base64.StdEncoding.DecodeString(dataParts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64: %v", err)
	}

	// Determine file extension
	ext := ".png"
	if strings.Contains(dataParts[0], "image/jpeg") {
		ext = ".jpg"
	}

	// Generate filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(uploadDir, filename)

	// Create directory if it doesn't exist
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Save file to disk
	err = os.WriteFile(fullPath, decodedData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	//GetDomain
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return "", fmt.Errorf("environment variable DOMAIN is not set")
	}

	// Return the public URL
	return fmt.Sprintf("%s/uploads/rooms/%s", domain, filename), nil
}
