package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bluejutzu/GoBot/cli/version"
	"github.com/spf13/cobra"
)

type PkgVersionInfo struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// SelfUpdateCmd represents the self-update command
var SelfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update the CLI to the latest version",
	Long:  `Checks for and installs the latest version of the GoBot CLI from pkg.dev.`,
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

	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "gobot-update")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Download new version
	fmt.Println("Downloading new version...")
	binaryPath := filepath.Join(tempDir, "gobot")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	if err := downloadFile(binaryPath, latestVersion.URL); err != nil {
		return fmt.Errorf("failed to download update: %v", err)
	}

	// Make binary executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %v", err)
		}
	}

	// Get Go binary directory
	goBin := os.Getenv("GOBIN")
	if goBin == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
		}
		goBin = filepath.Join(gopath, "bin")
	}

	// Replace current binary
	currentBinary := filepath.Join(goBin, "gobot")
	if runtime.GOOS == "windows" {
		currentBinary += ".exe"
	}

	fmt.Println("Installing new version...")
	if err := os.Rename(binaryPath, currentBinary); err != nil {
		// Try copy if rename fails (might be on different device)
		if copyErr := copyFile(binaryPath, currentBinary); copyErr != nil {
			return fmt.Errorf("failed to install update: %v", copyErr)
		}
	}

	fmt.Printf("Successfully updated to version %s!\n", latestVersion.Version)
	return nil
}

func getLatestVersion() (*PkgVersionInfo, error) {
	// TODO: Replace this URL with your actual pkg.dev registry URL
	resp, err := http.Get("https://pkg.dev/bluejutzu/gobot/latest")
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

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
