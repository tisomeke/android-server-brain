package config

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
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
	json.NewDecoder(file).Decode(cfg)

	setupDirectories(cfg.StorageDir)
	return cfg
}

func setupDirectories(storagePath string) {
	// 1. Creating directory
	home, _ := os.UserHomeDir()
	fullStoragePath := filepath.Join(home, storagePath)
	
	err := os.MkdirAll(fullStoragePath, 0755)
	if err != nil {
		log.Printf("Ошибка создания директории: %v", err)
	}

	// 2. symlink ~/server -> downloads/server
	linkPath := filepath.Join(home, "server")
	
	// checking
	if _, err := os.Lstat(linkPath); err == nil {
		os.Remove(linkPath) 
	}

	err = os.Symlink(fullStoragePath, linkPath)
	if err != nil {
		log.Printf("Couldn't create symlink: %v", err)
	}
}