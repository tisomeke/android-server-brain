package system

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetSystemStatus() string {
	// getting access to battery level with termux-api
	battery, _ := exec.Command("termux-battery-status").Output()
	
	// accessing available free space
	df, _ := exec.Command("sh", "-c", "df -h /data | tail -1 | awk '{print $4}'").Output()
	
	// accessing uptime
	uptime, _ := exec.Command("uptime", "-p").Output()

	status := fmt.Sprintf(
		"📊 **ASB Status**\n\n"+
		"🔋 Battery: %s\n"+
		"💾 Free Space: %s\n"+
		"⏱ Uptime: %s",
		string(battery),
		strings.TrimSpace(string(df)),
		strings.TrimSpace(string(uptime)),
	)
	return status
}