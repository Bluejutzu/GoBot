package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bluejutzu/GoBot/cli/version"
	"github.com/spf13/cobra"
)

type PkgVersionInfo struct {
	Version string `json:"version"`
}

// SelfUpdateCmd represents the self-update command
var SelfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update the CLI to the latest version",
	Long:  `Checks for and installs the latest version of the GoBot CLI from pkg.dev using go install.`,
	RunE:  runSelfUpdate,
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")

	// Get latest version from pkg.dev
	latestVersion, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %v", err)
	}

	currentVersion := version.Version
	if !isNewerVersion(currentVersion, latestVersion.Version) {
		fmt.Printf("You're already on the latest version (%s)\n", currentVersion)
		return nil
	}

	fmt.Printf("New version available: %s (current: %s)\n", latestVersion.Version, currentVersion)
	fmt.Println("Starting update process...")

	// Use go install to update the CLI
	updateCmd := exec.Command("go", "install", fmt.Sprintf("github.com/bluejutzu/GoBot/cli@%s", latestVersion.Version))
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update CLI: %v", err)
	}

	fmt.Printf("Successfully updated to version %s!\n", latestVersion.Version)
	return nil
}

func getLatestVersion() (*PkgVersionInfo, error) {
	resp, err := http.Get("https://pkg.go.dev/github.com/bluejutzu/Gobot@latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info PkgVersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func isNewerVersion(current, latest string) bool {
	// Strip v prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Split versions into parts
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// Compare each part
	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	// If we get here and latest has more parts, it's newer
	return len(latestParts) > len(currentParts)
}
