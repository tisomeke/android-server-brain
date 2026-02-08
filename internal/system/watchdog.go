package system

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"android-server-brain/config"

	tele "gopkg.in/telebot.v3"
)

// BatteryStatus represents the battery information from termux-api
type BatteryStatus struct {
	Health      string  `json:"health"`
	Percentage  float64 `json:"percentage"`
	Plugged     string  `json:"plugged"`
	Status      string  `json:"status"`
	Temperature float64 `json:"temperature"`
}

// Watchdog manages periodic system monitoring
type Watchdog struct {
	bot          *tele.Bot
	config       *config.Config
	interval     time.Duration
	lastNotified bool
	startTime    time.Time
}

// NewWatchdog creates a new watchdog instance
func NewWatchdog(bot *tele.Bot, cfg *config.Config, interval time.Duration) *Watchdog {
	return &Watchdog{
		bot:          bot,
		config:       cfg,
		interval:     interval,
		lastNotified: false,
		startTime:    time.Now(),
	}
}

// Start begins the watchdog monitoring loop
func (w *Watchdog) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Printf("Watchdog started with interval: %v", w.interval)

		for range ticker.C {
			w.checkBattery()
		}
	}()
}

// checkBattery monitors battery status and sends notifications when needed
func (w *Watchdog) checkBattery() {
	battery, err := getBatteryStatus()
	if err != nil {
		log.Printf("Failed to get battery status: %v", err)
		return
	}

	// Check if battery is low and not charging
	isLowBattery := battery.Percentage < 20
	isCharging := battery.Plugged == "PLUGGED_AC" || battery.Plugged == "PLUGGED_USB" || battery.Status == "CHARGING"

	// Send notification only if:
	// 1. Battery is low (<20%)
	// 2. Device is not charging
	// 3. We haven't notified recently (anti-spam)
	if isLowBattery && !isCharging && !w.lastNotified {
		message := fmt.Sprintf(
			"⚠️ *Low Battery Alert*\n\n"+
				"🔋 Current charge: %.1f%%\n"+
				"🌡 Temperature: %.1f°C\n"+
				"🔌 Status: %s\n\n"+
				"Please connect charger!",
			battery.Percentage,
			battery.Temperature,
			battery.Status,
		)

		_, err := w.bot.Send(&tele.User{ID: w.config.AdminID}, message, tele.ModeMarkdown)
		if err != nil {
			log.Printf("Failed to send battery alert: %v", err)
		} else {
			log.Printf("Sent low battery alert: %.1f%%", battery.Percentage)
			w.lastNotified = true
		}
	} else if !isLowBattery || isCharging {
		// Reset notification flag when battery is healthy or charging
		w.lastNotified = false
	}
}

// getBatteryStatus retrieves battery information from termux-api
func getBatteryStatus() (*BatteryStatus, error) {
	cmd := exec.Command("termux-battery-status")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute termux-battery-status: %w", err)
	}

	var battery BatteryStatus
	err = json.Unmarshal(output, &battery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse battery status: %w", err)
	}

	return &battery, nil
}

// GetWatchdogStatus returns formatted watchdog status information
func (w *Watchdog) GetStatus() string {
	uptime := time.Since(w.startTime)

	return fmt.Sprintf(
		"🐕 *Watchdog Status*\n\n"+
			"⏱ Uptime: %v\n"+
			"📅 Interval: %v\n"+
			"📢 Last notified: %v",
		uptime.Round(time.Second),
		w.interval,
		w.lastNotified,
	)
}

// GetBatteryInfo returns formatted battery information for manual checks
func GetBatteryInfo() string {
	battery, err := getBatteryStatus()
	if err != nil {
		return fmt.Sprintf("❌ Battery info unavailable: %v", err)
	}

	return fmt.Sprintf(
		"🔋 *Battery Status*\n\n"+
			"Charge: %.1f%%\n"+
			"Status: %s\n"+
			"Plugged: %s\n"+
			"Health: %s\n"+
			"Temperature: %.1f°C",
		battery.Percentage,
		battery.Status,
		battery.Plugged,
		battery.Health,
		battery.Temperature,
	)
}
