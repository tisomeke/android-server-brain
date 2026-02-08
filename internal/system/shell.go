package system

import (
	"context"
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
