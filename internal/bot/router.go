package bot

import (
	"android-server-brain/internal/system"
	tele "gopkg.in/telebot.v3"
)

// RegisterHandlers maps commands to functions
func RegisterHandlers(b *tele.Bot) {
	// Standard command handler
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome to Android Server Brain. Use /status to check system health.")
	})

	// System monitoring handler
	b.Handle("/status", func(c tele.Context) error {
		status := system.GetSystemStatus()
		return c.Send(status, tele.ModeMarkdown)
	})

	// Placeholder for future keyboard implementation
	// b.Handle("/menu", showMenu)
}