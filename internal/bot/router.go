package bot

import (
	"android-server-brain/config"
	"android-server-brain/internal/storage"
	"android-server-brain/internal/system"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func RegisterHandlers(b *tele.Bot, cfg *config.Config, watchdog *system.Watchdog) {
	// Standard command handler
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome to Android Server Brain. Use /status to check system health.")
	})

	// System monitoring handler
	b.Handle("/status", func(c tele.Context) error {
		status := system.GetSystemStatus()
		return c.Send(status, tele.ModeMarkdown)
	})

	// Handle incoming documents (files)
	b.Handle(tele.OnDocument, func(c tele.Context) error {
		doc := c.Message().Document

		c.Send(fmt.Sprintf("📥 Receiving file: %s...", doc.FileName))
		filePath, err := storage.SaveTelegramFile(b, doc, cfg.StorageDir)
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Error saving file: %v", err))
		}

		return c.Send(fmt.Sprintf("✅ File saved and made executable:\n`%s` \n\nYou can run it from `~/server/%s`", filePath, doc.FileName), tele.ModeMarkdown)
	})

	// Battery status handler
	b.Handle("/battery", func(c tele.Context) error {
		status := system.GetBatteryInfo()
		return c.Send(status, tele.ModeMarkdown)
	})

	// Watchdog status handler
	b.Handle("/watchdog", func(c tele.Context) error {
		status := watchdog.GetStatus()
		return c.Send(status, tele.ModeMarkdown)
	})

	// Command execution handler
	b.Handle("/exec", func(c tele.Context) error {
		// Extract command from message (remove "/exec " prefix)
		args := c.Args()
		if len(args) == 0 {
			return c.Send("Usage: `/exec <command>`", tele.ModeMarkdown)
		}

		fullCommand := strings.Join(args, " ")
		c.Send(fmt.Sprintf("⏳ Executing: `%s`...", fullCommand), tele.ModeMarkdown)

		// Run the command
		output, err := system.ExecuteCommand(fullCommand)

		// If output is empty, provide a fallback message
		if strings.TrimSpace(output) == "" {
			if err != nil {
				output = "Error: " + err.Error()
			} else {
				output = "Command executed successfully (no output)."
			}
		}

		// Wrap output in code blocks for readability
		return c.Send(fmt.Sprintf("📝 *Output:*\n```\n%s\n```", output), tele.ModeMarkdown)
	})

	// Reboot system handler with simple confirmation
	b.Handle("/reboot", func(c tele.Context) error {
		return c.Send("⚠️ *System Reboot Confirmation*\n\nAre you sure you want to reboot the system? This will disconnect all active sessions.\n\nType `YES` to confirm or anything else to cancel.", tele.ModeMarkdown)
	})

	// Handle reboot confirmation text
	b.Handle(tele.OnText, func(c tele.Context) error {
		if strings.ToUpper(c.Text()) == "YES" {
			// Check if this is in response to reboot command
			// For simplicity, we'll assume any "YES" means reboot confirmation
			result, err := system.RebootSystem()
			if err != nil {
				return c.Send(result, tele.ModeMarkdown)
			}
			return c.Send(result, tele.ModeMarkdown)
		}
		return nil // Ignore other text messages
	})

	// Restart service handler
	b.Handle("/restart", func(c tele.Context) error {
		args := c.Args()
		if len(args) == 0 {
			return c.Send(system.ListServices(), tele.ModeMarkdown)
		}

		serviceName := args[0]
		c.Send(fmt.Sprintf("⏳ Restarting service: `%s`...", serviceName), tele.ModeMarkdown)

		result, err := system.RestartService(serviceName)
		if err != nil {
			return c.Send(result, tele.ModeMarkdown)
		}

		return c.Send(result, tele.ModeMarkdown)
	})
}
