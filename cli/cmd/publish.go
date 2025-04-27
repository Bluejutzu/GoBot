package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/bluejutzu/GoBot/cli/version"
	"github.com/spf13/cobra"
)

// PublishCmd represents the command that handles versioning and publishing of the CLI tool.
// It manages version tagging, building with version information, and publishing to GitHub/pkg.dev.
var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a new version to GitHub and pkg.dev",
	Long: `Publishes a new version of GoBot to GitHub and pkg.dev with proper versioning.
This command will:
1. Create and push a new git tag
2. Build the binary with version information
3. Publish the binary to pkg.dev`,
	RunE: runPublish,
}

func init() {
	PublishCmd.Flags().String("tag", "", "Version tag to publish (required)")
	PublishCmd.MarkFlagRequired("tag")
}

type PkgPublishRequest struct {
	Version string `json:"version"`
	Binary  []byte `json:"binary"`
}

// runPublish handles the publishing process including version validation,
// git operations, and publishing to pkg.dev
func runPublish(cmd *cobra.Command, args []string) error {
	tag, _ := cmd.Flags().GetString("tag")

	// Validate version format
	if !isValidVersion(tag) {
		return fmt.Errorf("invalid version format. Must be in format v0.0.0")
	}

	fmt.Printf("Starting publish process for version %s...\n", tag)

	// Get current git commit
	commitHash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("failed to get git commit: %v", err)
	}

	// Update version information
	version.Version = tag
	version.Commit = string(commitHash)
	version.Date = time.Now().Format(time.RFC3339)

	// Build with version info
	fmt.Println("Building binary with version information...")
	buildCmd := exec.Command("go", "build",
		"-ldflags", fmt.Sprintf("-X github.com/bluejutzu/GoBot/cli/version.Version=%s -X github.com/bluejutzu/GoBot/cli/version.Commit=%s -X github.com/bluejutzu/GoBot/cli/version.Date=%s",
			version.Version, version.Commit, version.Date),
		"-o", "gobot.exe", "./gobotcli")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build: %v", err)
	}

	// Stage all changes
	fmt.Println("Staging changes...")
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %v", err)
	}

	// Commit changes
	fmt.Println("Creating release commit...")
	commitMsg := fmt.Sprintf("Release version %s", tag)
	if err := exec.Command("git", "commit", "-m", commitMsg).Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}

	// Create git tag
	fmt.Println("Creating git tag...")
	if err := exec.Command("git", "tag", "-a", tag, "-m", fmt.Sprintf("Release %s", tag)).Run(); err != nil {
		return fmt.Errorf("failed to create git tag: %v", err)
	}

	// Push commits and tags
	fmt.Println("Pushing to GitHub...")
	if err := exec.Command("git", "push", "origin", "main").Run(); err != nil {
		return fmt.Errorf("failed to push commits: %v", err)
	}
	if err := exec.Command("git", "push", "origin", tag).Run(); err != nil {
		return fmt.Errorf("failed to push tag: %v", err)
	}

	fmt.Printf("Successfully published version %s to GitHub\n", tag)

	// Read the binary
	fmt.Println("Reading binary for pkg.dev publishing...")
	binary, err := os.ReadFile("gobot.exe")
	if err != nil {
		return fmt.Errorf("failed to read binary: %v", err)
	}

	// Publish to pkg.dev
	fmt.Println("Publishing to pkg.dev...")
	if err := publishToPkgDev(tag, binary); err != nil {
		return fmt.Errorf("failed to publish to pkg.dev: %v", err)
	}

	fmt.Printf("Successfully published version %s to pkg.dev\n", tag)
	fmt.Println("Users can now update using: gobot self-update")

	return nil
}

// publishToPkgDev handles the package publishing to pkg.dev
func publishToPkgDev(version string, binary []byte) error {
	request := PkgPublishRequest{
		Version: version,
		Binary:  binary,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// TODO: Replace with your actual pkg.dev publish endpoint
	resp, err := http.Post("https://pkg.dev/bluejutzu/gobot/publish", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pkg.dev publish failed with status: %s", resp.Status)
	}

	return nil
}

// isValidVersion validates that the provided version string follows
// the expected format (starts with 'v').
func isValidVersion(v string) bool {
	// Simple version validation - could be more robust
	return len(v) > 0 && v[0] == 'v'
}
