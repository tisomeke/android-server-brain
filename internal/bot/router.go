package bot

import (
	"fmt"
	"strings"
	"android-server-brain/internal/storage"
	"android-server-brain/internal/system"
	"android-server-brain/config"
	
	tele "gopkg.in/telebot.v3"
)

func RegisterHandlers(b *tele.Bot, cfg *config.Config) {
	// ... (previous handlers: /start, /status, tele.OnDocument)

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
}