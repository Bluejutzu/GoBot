package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bluejutzu/GoBot/cli/version"
	"github.com/spf13/cobra"
)

// PublishCmd represents the command that handles versioning and publishing of the CLI tool.
// It manages version tagging, packaging the CLI source code, and publishing to pkg.dev.
var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a new version to pkg.dev",
	Long: `Publishes a new version of the GoBot CLI source code to pkg.dev.
This command will:
1. Create and push a new git tag
2. Package the CLI source code
3. Publish the package to pkg.dev`,
	RunE: runPublish,
}

func init() {
	PublishCmd.Flags().String("tag", "", "Version tag to publish (required)")
	PublishCmd.MarkFlagRequired("tag")
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

	// Package and publish CLI source code
	fmt.Println("Packaging CLI source code...")
	packageData, err := packageCLISource()
	if err != nil {
		return fmt.Errorf("failed to package CLI source: %v", err)
	}

	// Publish to pkg.dev
	fmt.Println("Publishing to pkg.dev...")
	if err := publishToPkgDev(tag, packageData); err != nil {
		return fmt.Errorf("failed to publish to pkg.dev: %v", err)
	}

	fmt.Printf("Successfully published version %s to pkg.dev\n", tag)
	fmt.Println("Users can now update using: go install github.com/bluejutzu/GoBot/cli@latest")

	return nil
}

// packageCLISource creates a tar.gz archive of the CLI source code
func packageCLISource() ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Walk through the cli directory
	err := filepath.Walk("./cli", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a regular file
		if !info.Mode().IsRegular() {
			return nil
		}

		// Create relative path for the tar
		relPath, err := filepath.Rel("./cli", path)
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Write file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	// Close writers
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// publishToPkgDev handles the package publishing to pkg.dev
func publishToPkgDev(version string, sourcePackage []byte) error {
	// Create form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add version field
	if err := writer.WriteField("version", version); err != nil {
		return fmt.Errorf("failed to write version field: %v", err)
	}

	// Add package file
	part, err := writer.CreateFormFile("package", "cli.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := part.Write(sourcePackage); err != nil {
		return fmt.Errorf("failed to write package data: %v", err)
	}

	writer.Close()

	// Send request
	resp, err := http.Post(
		"https://pkg.dev/bluejutzu/gobot/publish",
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return fmt.Errorf("failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("publish failed with status: %s", resp.Status)
	}

	return nil
}

// isValidVersion validates that the provided version string follows
// the expected format (starts with 'v').
func isValidVersion(v string) bool {
	return len(v) > 0 && v[0] == 'v'
}
