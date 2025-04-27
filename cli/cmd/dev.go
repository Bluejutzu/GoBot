package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bluejutzu/GoBot/cli/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	// DevCmd represents the development command that starts the bot in development mode
	// with features like file watching and command validation.
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Start development mode",
		Long:  `Starts the bot in development mode with file watching and command validation.`,
		Run:   runDev,
	}

	// cfg holds the current configuration loaded from Gobot.config.json
	cfg *config.Config
)

// SetConfig makes the loaded configuration available to commands
func SetConfig(conf *config.Config) {
	cfg = conf
}

// runDev executes the development mode of the bot.
// It sets up file watching on configured directories and validates Discord commands.
func runDev(cmd *cobra.Command, args []string) {
	if cfg == nil {
		cfg = &config.DefaultConfig
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch configured directories
	dirsToWatch := []string{
		filepath.Dir(cfg.MainFile),
		filepath.Dir(cfg.BotFile),
		cfg.CommandsDir,
	}

	for _, dir := range dirsToWatch {
		err = watcher.Add(dir)
		if err != nil {
			log.Printf("Error watching directory %s: %v", dir, err)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if filepath.Base(event.Name) == filepath.Base(cfg.BotFile) {
						validateCommands()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Run main.go
	go runMainGo()

	// Initial command validation
	validateCommands()

	// Block forever
	select {}
}

// runMainGo executes the main.go file of the bot
func runMainGo() {
	cmd := exec.Command("go", "run", cfg.MainFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running %s: %v", cfg.MainFile, err)
	}
}

// validateCommands checks if all commands in the commands directory
// are properly registered in bot.go
func validateCommands() {
	// Read bot.go
	botContent, err := os.ReadFile(cfg.BotFile)
	if err != nil {
		log.Printf("Error reading %s: %v", cfg.BotFile, err)
		return
	}

	// Read all command files
	cmdPattern := filepath.Join(cfg.CommandsDir, "**", "*_command.go")
	commandFiles, err := filepath.Glob(cmdPattern)
	if err != nil {
		log.Printf("Error finding command files: %v", err)
		return
	}

	for _, file := range commandFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Extract command name from file content
		cmdName := extractCommandName(string(content))
		if cmdName == "" {
			continue
		}

		botContentStr := string(botContent)

		commandVariations := []string{
			fmt.Sprintf(`%s.Command`, strings.Split(filepath.Base(file), "_command.go")[0]),
			fmt.Sprintf(`%s_Command`,  strings.Split(filepath.Base(file), "_command.go")[0]),
			fmt.Sprintf(`%sCommand`,  strings.Split(filepath.Base(file), "_command.go")[0]),
			fmt.Sprintf(`%s_Command`, strings.ToUpper(cmdName)),
			fmt.Sprintf(`%sCommand`, strings.Title(cmdName)),
			cmdName,
			strings.Split(filepath.Base(file), "_command.go")[0],
		}

		inCommandArray := false
		for _, variation := range commandVariations {
			if strings.Contains(botContentStr, variation) {
				inCommandArray = true
				break
			}
		}

		inHandlerMap := strings.Contains(botContentStr, fmt.Sprintf(`"%s":`, cmdName))

		if !inCommandArray {
			fmt.Printf("\033[33m⚠ Command %s (%s) is not included in the commands array in %s\033[0m\n",
				cmdName, file, cfg.BotFile)
		}
		if !inHandlerMap {
			fmt.Printf("\033[33m⚠ Command %s (%s) is not included in commandHandlers in %s\033[0m\n",
				cmdName, file, cfg.BotFile)
		}
	}
}

func extractCommandName(content string) string {
	if strings.Contains(content, "Name:") {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Name:") {
				parts := strings.Split(line, `"`)
				if len(parts) >= 3 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}
	return ""
}
