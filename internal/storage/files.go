package storage

import (
	"fmt"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v3"
)

// SaveTelegramFile downloads a file from Telegram and saves it to the asb_files directory
func SaveTelegramFile(b *tele.Bot, doc *tele.Document, targetDir string) (string, error) {
	// Get the file path from Telegram servers
	file, err := b.FileByID(doc.FileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file by ID: %v", err)
	}

	// Prepare the full destination path in ~/asb_files (symlinked to /storage/emulated/0/Download/asb_files)
	home, _ := os.UserHomeDir()
	asbFilesDir := filepath.Join(home, "asb_files")
	
	// Ensure the directory exists
	err = os.MkdirAll(asbFilesDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create asb_files directory: %v", err)
	}
	
	fullPath := filepath.Join(asbFilesDir, doc.FileName)

	// Create the file on disk
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Download the file directly using telebot
	err = b.Download(&file, fullPath)
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
