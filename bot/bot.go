package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bluejutzu/GoBot/bot/handlers"
	"github.com/bluejutzu/GoBot/commands/misc"
	"github.com/bluejutzu/GoBot/commands/moderation"
	"github.com/bluejutzu/GoBot/commands/ridealong"
	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string // Set in main.go

	commands = []*discordgo.ApplicationCommand{
		ridealong.Command,

		misc.ID_Commmand,
		misc.PING_Command,
		misc.EIGHTBALL_Command,
		misc.SAY_COMMAND,
		misc.SUGGEST_Command,

		moderation.BAN_Command,
		moderation.SOFTBAN_Command,
		moderation.MUTE_Command,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ra":            ridealong.ParseCommand,
		"ping":          misc.PING_ParseCommand,
		"what-is-my-id": misc.ID_ParseCommand,
		"8ball":         misc.EIGHTBALL_ParseCommand,
		"say":           misc.SAY_ParseCommand,
		"suggest":       misc.SUGGEST_ParseCommand,
		"ban":           moderation.BAN_ParseCommand,
		"softban":       moderation.SOFTBAN_ParseCommand,
		"mute":          moderation.MUTE_ParseCommand,
	}

	Green     = "\033[32m"
	Reset     = "\033[0m"
	BrightRed = "\033[31;1m"
)

func Run() {
	if BotToken == "" {
		log.Fatal("Bot token is not set")
	}

	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildPresences |
		discordgo.IntentsGuildMembers

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("Registering commands...")
	startTime := time.Now()
	registeredCommands, err := discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, "", commands)
	if err != nil {
		log.Fatalf("%vError registering commands: %v%v", BrightRed, err, Reset)
	}
	duration := time.Since(startTime).Abs()
	fmt.Printf("%vAll %d commands registered successfully! Took: %s%v\n", Green, len(registeredCommands), duration, Reset)

	discord.AddHandler(handlers.Router)

	defer func() {
		fmt.Println("\nCleaning up registered commands by bulk overwriting with an empty list...")
		startTime := time.Now()
		_, err := discord.ApplicationCommandBulkOverwrite(discord.State.User.ID, "", []*discordgo.ApplicationCommand{})
		if err != nil {
			log.Printf("%vError cleaning up (deleting) commands: %v%v", BrightRed, err, Reset)
		} else {
			duration := time.Since(startTime).Abs()
			fmt.Printf("%vCommands cleaned up successfully! Took: %s%v\n", Green, duration, Reset)
		}
		discord.Close()
	}()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
