package util

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func SaveUploadedFile(file *multipart.FileHeader, path string) error {
	// Open the uploaded file.
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file for the uploaded content.
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded content to the destination file.
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func ValidateImageFile(file multipart.File) (bool, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, errors.New("failed to read file")
	}

	// Deteksi tipe MIME
	fileType := http.DetectContentType(buffer)

	// Daftar tipe file yang diizinkan
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range allowedTypes {
		if fileType == t {
			// Jika valid, reset file pointer
			file.Seek(0, 0)
			return true, nil
		}
	}

	// Jika tipe file tidak valid
	return false, errors.New("only image files are allowed")
}
