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

	// Create inline keyboard markup
	markup := &tele.ReplyMarkup{}

	// Define inline buttons for reboot confirmation using markup helpers
	rebootConfirmBtn := markup.Data("✅ YES - Reboot System", "reboot_confirm", "confirm")
	rebootCancelBtn := markup.Data("❌ NO - Cancel", "reboot_cancel", "cancel")

	// Register button callback handlers
	b.Handle(&rebootConfirmBtn, func(c tele.Context) error {
		result, err := system.RebootSystem()
		if err != nil {
			return c.Send(result, tele.ModeMarkdown)
		}
		return c.Send(result, tele.ModeMarkdown)
	})

	b.Handle(&rebootCancelBtn, func(c tele.Context) error {
		return c.Send("❌ Reboot cancelled.", tele.ModeMarkdown)
	})

	// Reboot system handler with inline buttons
	b.Handle("/reboot", func(c tele.Context) error {
		markup.Inline(
			markup.Row(rebootConfirmBtn),
			markup.Row(rebootCancelBtn),
		)

		return c.Send("⚠️ *System Reboot Confirmation*\n\nAre you sure you want to reboot the system? This will disconnect all active sessions.", tele.ModeMarkdown, markup)
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

	// Update system handler
	b.Handle("/update", func(c tele.Context) error {
		args := c.Args()

		// If no arguments, check for updates
		if len(args) == 0 {
			c.Send("🔍 Checking for updates...", tele.ModeMarkdown)

			result, err := system.CheckForUpdates()
			if err != nil {
				return c.Send(fmt.Sprintf("❌ Error checking for updates: %v", err), tele.ModeMarkdown)
			}

			if !result.Success {
				return c.Send(result.Message, tele.ModeMarkdown)
			}

			// If updates are available, show update options
			if result.NewVersion != "" {
				message := fmt.Sprintf(
					"%s\n\n"+
						"*Current version:* `%s`\n"+
						"*Available version:* `%s`\n\n"+
						"Use `/update now` to install updates",
					result.Message,
					strings.TrimSpace(result.OldVersion),
					strings.TrimSpace(result.NewVersion),
				)
				return c.Send(message, tele.ModeMarkdown)
			}

			return c.Send(result.Message, tele.ModeMarkdown)
		}

		// If argument is "now", perform update
		if args[0] == "now" {
			c.Send("🔄 Starting update process...", tele.ModeMarkdown)

			// Perform update
			result, err := system.PerformUpdate()
			if err != nil {
				return c.Send(fmt.Sprintf("❌ Update failed: %v", err), tele.ModeMarkdown)
			}

			if !result.Success {
				return c.Send(result.Message, tele.ModeMarkdown)
			}

			// Show success message
			message := fmt.Sprintf(
				"%s\n\n"+
					"*Updated to version:* `%s`\n"+
					"Backup created at: `%s`\n\n"+
					"Restarting ASB service now...",
				result.Message,
				strings.TrimSpace(result.NewVersion),
				result.BackupPath,
			)

			c.Send(message, tele.ModeMarkdown)

			// Restart ASB service
			restartMsg, restartErr := system.RestartASB()
			if restartErr != nil {
				return c.Send(fmt.Sprintf("⚠️ %s\n\n%s", restartMsg, "Manual restart may be required."), tele.ModeMarkdown)
			}

			return c.Send(fmt.Sprintf("✅ %s", restartMsg), tele.ModeMarkdown)
		}

		// Invalid argument
		return c.Send("Usage:\n• `/update` - Check for updates\n• `/update now` - Install available updates", tele.ModeMarkdown)
	})
}
