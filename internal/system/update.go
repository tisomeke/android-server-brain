package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	Success    bool
	Message    string
	OldVersion string
	NewVersion string
	Error      error
	BackupPath string
}

// CheckForUpdates checks if there are updates available from the git repository
func CheckForUpdates() (*UpdateResult, error) {
	result := &UpdateResult{}

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Failed to get working directory: %v", err)
		return result, nil
	}

	// Check if we're in a git repository
	gitDir := filepath.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		result.Success = false
		result.Message = "❌ Not a git repository. Cannot check for updates."
		return result, nil
	}

	// Fetch latest changes
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fetchCmd := exec.CommandContext(ctx, "git", "fetch")
	fetchOutput, err := fetchCmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Failed to fetch updates: %v\nOutput: %s", err, string(fetchOutput))
		return result, nil
	}

	// Check if there are updates available
	statusCmd := exec.CommandContext(ctx, "git", "status", "-uno")
	statusOutput, err := statusCmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Failed to check status: %v\nOutput: %s", err, string(statusOutput))
		return result, nil
	}

	statusStr := string(statusOutput)
	if contains(statusStr, "Your branch is behind") || contains(statusStr, "can be fast-forwarded") {
		result.Success = true
		result.Message = "✅ Updates are available!"

		// Get current version/commit
		currentCmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD")
		currentOutput, _ := currentCmd.Output()
		result.OldVersion = string(currentOutput)

		// Get remote version/commit
		remoteCmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "@{u}")
		remoteOutput, _ := remoteCmd.Output()
		result.NewVersion = string(remoteOutput)

		return result, nil
	}

	result.Success = true
	result.Message = "✅ Already up to date. No updates available."
	return result, nil
}

// PerformUpdate performs the actual update process with backup
func PerformUpdate() (*UpdateResult, error) {
	result := &UpdateResult{}

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Failed to get working directory: %v", err)
		return result, nil
	}

	// Create backup
	backupPath := filepath.Join(wd, fmt.Sprintf("backup_%d", time.Now().Unix()))
	err = exec.Command("cp", "-r", wd, backupPath).Run()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Failed to create backup: %v", err)
		return result, nil
	}
	result.BackupPath = backupPath

	// Perform git pull
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pullCmd := exec.CommandContext(ctx, "git", "pull", "--ff-only")
	pullOutput, err := pullCmd.CombinedOutput()
	if err != nil {
		result.Success = false
		result.Error = err
		result.Message = fmt.Sprintf("❌ Update failed: %v\nOutput: %s", err, string(pullOutput))

		// Attempt rollback
		rollbackCmd := exec.Command("cp", "-r", backupPath+"/*", wd+"/")
		rollbackCmd.Run() // Ignore rollback errors

		return result, nil
	}

	// Get version info
	currentCmd := exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD")
	currentOutput, _ := currentCmd.Output()
	result.NewVersion = string(currentOutput)

	result.Success = true
	result.Message = fmt.Sprintf("✅ Successfully updated!\nOutput: %s", string(pullOutput))
	return result, nil
}

// RestartASB restarts the ASB service after update
func RestartASB() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to find and restart the ASB process
	cmd := exec.CommandContext(ctx, "pkill", "-f", "asb")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("⚠️ Could not terminate old ASB process: %v\nOutput: %s", err, string(output)), err
	}

	// Start new ASB process
	startCmd := exec.CommandContext(ctx, "./asb")
	startCmd.Dir, _ = os.Getwd()

	err = startCmd.Start()
	if err != nil {
		return fmt.Sprintf("❌ Failed to start new ASB process: %v", err), err
	}

	return "✅ ASB restarted successfully with new version!", nil
}

// CleanupBackup removes backup files/directories
func CleanupBackup(backupPath string) error {
	return os.RemoveAll(backupPath)
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
