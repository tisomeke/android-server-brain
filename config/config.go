package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	TelegramToken string `json:"telegram_token"`
	AdminID       int64  `json:"admin_id"`
	StorageDir    string `json:"storage_dir"` // downloads/server
}

func LoadConfig() *Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Please create config.json as in the template")
	}
	defer file.Close()

	cfg := &Config{}
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		log.Fatalf("Failed to parse config.json: %v", err)
	}

	// Validate required fields
	if cfg.TelegramToken == "" {
		log.Fatal("telegram_token is required in config.json")
	}
	if cfg.AdminID == 0 {
		log.Fatal("admin_id is required in config.json")
	}
	if cfg.StorageDir == "" {
		cfg.StorageDir = "downloads/server" // default value
	}

	setupDirectories(cfg.StorageDir)
	return cfg
}

func setupDirectories(storagePath string) {
	home, _ := os.UserHomeDir()

	// New storage structure:
	// 1. Use Android's system Downloads folder via termux-storage-create
	// 2. Create asb_files directory in Downloads
	// 3. Create symlink ~/asb_files -> /storage/emulated/0/Download/asb_files

	// Path to Android's system Downloads folder
	systemDownloadsPath := "/storage/emulated/0/Download"
	asbFilesDir := filepath.Join(systemDownloadsPath, "asb_files")

	// Create asb_files directory in system Downloads
	err := os.MkdirAll(asbFilesDir, 0755)
	if err != nil {
		log.Printf("Error creating asb_files in Downloads: %v", err)
		log.Printf("Falling back to home directory storage")
		// Fallback to home directory if system Downloads not accessible
		asbFilesDir = filepath.Join(home, "asb_files")
		os.MkdirAll(asbFilesDir, 0755)
	}

	// Create symlink ~/asb_files -> asbFilesDir
	linkPath := filepath.Join(home, "asb_files")

	// Remove existing symlink if present
	if _, err := os.Lstat(linkPath); err == nil {
		os.Remove(linkPath)
	}

	err = os.Symlink(asbFilesDir, linkPath)
	if err != nil {
		log.Printf("Couldn't create symlink: %v", err)
	}

	log.Printf("Storage setup: %s -> %s", linkPath, asbFilesDir)
}
