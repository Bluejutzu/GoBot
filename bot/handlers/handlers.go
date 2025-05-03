package handlers

import (
	"fmt"
	"strings"

	"github.com/bluejutzu/GoBot/commands/misc"
	"github.com/bluejutzu/GoBot/commands/moderation"
	"github.com/bluejutzu/GoBot/commands/ridealong"
	"github.com/bwmarrin/discordgo"
)

var (
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
)

func Router(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleCommand(s, i)
	case discordgo.InteractionModalSubmit:
		handleModal(s, i)
	case discordgo.InteractionMessageComponent:
		handleComponent(s, i)
	}
}

func handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	if h, ok := commandHandlers[name]; ok {
		h(s, i)
	} else {
		fmt.Println("Unhandled command:", name)
	}
}

func handleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.ModalSubmitData().CustomID
	switch {
	case strings.HasPrefix(customID, "ridealong_"):
		ridealong.HandleModal(s, i)
	case strings.HasPrefix(customID, "suggest_"):
		HandleSuggestionModal(s, i)
	default:
		fmt.Println("Unhandled modal ID:", customID)
	}
}

func handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	switch {
	case strings.HasPrefix(customID, "ridealong_"):
		ridealong.HandleButtons(s, i)
	case strings.HasPrefix(customID, "suggest_"):
		HandleSuggestionComponent(s, i)
	default:
		fmt.Println("Unhandled component ID:", customID)
	}
}
