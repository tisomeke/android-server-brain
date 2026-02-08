package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetSystemStatus collects data from Termux and OS
func GetSystemStatus() string {
	// Battery info from termux-api
	battery, err := exec.Command("termux-battery-status").Output()
	if err != nil {
		battery = []byte("Unavailable (API missing)")
	}

	// Storage info for the data partition
	df, err := exec.Command("sh", "-c", "df -h /data | tail -1 | awk '{print $4}'").Output()
	if err != nil {
		df = []byte("N/A")
	}

	// Uptime info
	uptime, err := exec.Command("uptime", "-p").Output()
	if err != nil {
		uptime = []byte("N/A")
	}

	return fmt.Sprintf(
		"📊 *System Status*\n\n"+
			"🔋 *Battery:* %s\n"+
			"💾 *Free Space:* %s\n"+
			"⏱ *Uptime:* %s",
		strings.TrimSpace(string(battery)),
		strings.TrimSpace(string(df)),
		strings.TrimSpace(string(uptime)),
	)
}