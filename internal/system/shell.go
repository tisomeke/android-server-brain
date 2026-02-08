package system

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ExecuteCommand runs a bash command with a timeout and returns combined output
func ExecuteCommand(command string) (string, error) {
	// Create a context with a 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute via 'sh -c' to support pipes and redirects
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// CombinedOutput returns both stdout and stderr
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return string(output) + "\n❌ Error: Command timed out", ctx.Err()
	}

	return string(output), err
}

// RebootSystem initiates system reboot
func RebootSystem() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "reboot")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("❌ Reboot command failed: %v\nOutput: %s", err, string(output)), err
	}

	return "🔄 System reboot initiated...", nil
}

// RestartService restarts a specific service
func RestartService(serviceName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try systemctl first (if available)
	cmd := exec.CommandContext(ctx, "systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()

	// If systemctl fails, try service command
	if err != nil {
		cmd = exec.CommandContext(ctx, "service", serviceName, "restart")
		output, err = cmd.CombinedOutput()
	}

	// If both fail, try direct kill and restart approach
	if err != nil {
		result := fmt.Sprintf("⚠️ Service restart attempt 1 failed: %v\n", err)

		// Try to kill the process
		killCmd := exec.CommandContext(ctx, "pkill", "-f", serviceName)
		killOutput, killErr := killCmd.CombinedOutput()
		result += fmt.Sprintf("Kill attempt output: %s\n", string(killOutput))

		if killErr == nil {
			result += "✅ Process killed. Service should restart automatically."
			return result, nil
		}

		return result + fmt.Sprintf("❌ All restart methods failed: %v", err), err
	}

	return fmt.Sprintf("✅ Service '%s' restarted successfully\nOutput: %s", serviceName, string(output)), nil
}

// ListServices returns available services that can be managed
func ListServices() string {
	return "📋 *Available Services for Restart:*\n\n" +
		"Common Android/Termux services:\n" +
		"• ssh (SSH server)\n" +
		"• nginx (Web server)\n" +
		"• apache2 (Web server)\n" +
		"• mysql (Database)\n" +
		"• postgresql (Database)\n" +
		"• redis (Cache server)\n" +
		"• docker (Container engine)\n" +
		"• cron (Scheduler)\n" +
		"• asb (This ASB service)\n\n" +
		"Usage: `/restart <service_name>`\n" +
		"Example: `/restart ssh`"
}
