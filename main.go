package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"android-server-brain/config"
	"android-server-brain/internal/bot"

	tele "gopkg.in/telebot.v3"
)

func main() {
	// Initialize configuration and directories
	cfg := config.LoadConfig()

	// Verification: check if termux-api is installed
	if _, err := exec.LookPath("termux-battery-status"); err != nil {
		log.Println("Warning: termux-api not found. Some monitoring features will be disabled.")
	}

	// Bot settings
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// Middleware: restrict access to AdminID only
	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID != cfg.AdminID {
				return nil // Ignore unauthorized users
			}
			return next(c)
		}
	})

	// Setup routes
	bot.RegisterHandlers(b)

	log.Printf("ASB Started: Admin ID %d", cfg.AdminID)
	b.Start()
}