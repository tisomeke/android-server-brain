package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

// SaveTelegramFile downloads a file from Telegram and saves it to the target directory
func SaveTelegramFile(b *tele.Bot, doc *tele.Document, targetDir string) (string, error) {
	// Get the file path from Telegram servers
	file, err := b.FileByID(doc.FileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file by ID: %v", err)
	}

	// Prepare the full destination path
	home, _ := os.UserHomeDir()
	fullPath := filepath.Join(home, targetDir, doc.FileName)

	// Create the file on disk
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Construct the download URL
	// Note: telebot handles the direct download via b.Download
	err = b.Download(file, fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}

	// Set executable permissions (chmod +x)
	err = os.Chmod(fullPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to set executable permissions: %v", err)
	}

	return fullPath, nil
}