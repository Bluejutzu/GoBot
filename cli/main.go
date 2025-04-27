package cli

import (
	"log"

	"github.com/bluejutzu/GoBot/cli/cmd"
	"github.com/bluejutzu/GoBot/cli/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
// It serves as the entry point for the GoBot CLI application.
var rootCmd = &cobra.Command{
	Use:   "gobot",
	Short: "GoBot CLI - A tool for managing GoBot development",
	Long:  `GoBot CLI provides development tools and utilities for the GoBot Discord bot project.`,
}

// Execute adds all child commands to the root command and sets up the CLI application.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// If an error occurs, it will be returned and the application will exit with a non-zero status.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(cmd.DevCmd)
	rootCmd.AddCommand(cmd.PublishCmd)
	rootCmd.AddCommand(cmd.SelfUpdateCmd)
}

func initConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Could not load config file: %v", err)
		return
	}

	// Make config available to commands
	cmd.SetConfig(cfg)
}
