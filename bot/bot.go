package bot

import (
    "fmt"
    "log"
    "os"
    "os/signal"

    "github.com/bluejutzu/GoBot/commands/ridealong"
    "github.com/bwmarrin/discordgo"
)

var (
    BotToken string
    commands = []*discordgo.ApplicationCommand{
        ridealong.Command,
    }

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "ra": ridealong.ParseCommand,
    }
)

func Run() {
    // Check if token is empty
    if BotToken == "" {
        log.Fatal("Bot token is not set")
    }

    // Create Discord session
    discord, err := discordgo.New("Bot " + BotToken)
    if err != nil {
        log.Fatal("Error creating Discord session:", err)
    }

    // Open connection to Discord
    err = discord.Open()
    if err != nil {
        log.Fatal("Error opening connection:", err)
    }

    // Register commands after establishing connection
    registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
    for i, command := range commands {
        cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
        if err != nil {
            log.Fatalf("Error creating command %v: %v", command.Name, err)
        }
        registeredCommands[i] = cmd
    }

    // Add handlers
    discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
            h(s, i)
        }
    })

    discord.AddHandler(ridealong.HandleButtons)
    discord.AddHandler(ridealong.HandleModal)

    defer func() {
        // Cleanup commands on shutdown
        for _, cmd := range registeredCommands {
            err := discord.ApplicationCommandDelete(discord.State.User.ID, "", cmd.ID)
            if err != nil {
                log.Printf("Cannot delete command %v: %v", cmd.Name, err)
            }
        }
        discord.Close()
    }()

    fmt.Println("Bot is now running. Press CTRL-C to exit.")

    // Wait for interrupt signal
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    <-stop
}
