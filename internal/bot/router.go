package bot

import (
	"fmt"
	"android-server-brain/internal/storage"
	"android-server-brain/internal/system"
	"android-server-brain/config"
	
	tele "gopkg.in/telebot.v3"
)

func RegisterHandlers(b *tele.Bot, cfg *config.Config) {
	// ... previous handlers (/start, /status)

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
}