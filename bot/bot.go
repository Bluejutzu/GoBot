package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bluejutzu/GoBot/commands/misc"
	"github.com/bluejutzu/GoBot/commands/moderation"
	"github.com/bluejutzu/GoBot/commands/ridealong"
	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	commands = []*discordgo.ApplicationCommand{
		ridealong.Command,

		misc.ID_Commmand,
		misc.PING_Command,
		misc.EIGHTBALL_Command,

		moderation.BAN_Command,
	}

	/*
		CommandHandlers stores the mapping between Discord slash commands and their handler functions.

		map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
		  - string: The command name that triggered the interaction
		  - func: The handler function that processes the interaction
	*/
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ra":            ridealong.ParseCommand,

		"ping":          misc.PING_ParseCommand,
		"what-is-my-id": misc.ID_ParseCommand,
		"8ball":         misc.EIGHTBALL_ParseCommand,

		"ban":           moderation.BAN_ParseCommand,
	}
)

func Run() {
	if BotToken == "" {
		log.Fatal("Bot token is not set")
	}

	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuildPresences

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
		if err != nil {
			log.Fatalf("Error creating command %v: %v", command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			ridealong.HandleModal(s, i)
		case discordgo.InteractionMessageComponent:
			ridealong.HandleButtons(s, i)
		}
	})

	defer func() {
		for _, cmd := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("Cannot delete command %v: %v", cmd.Name, err)
			}
		}
		discord.Close()
	}()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
