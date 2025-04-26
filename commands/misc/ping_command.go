package misc

import (
	"fmt"

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var PING_Command = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Get the Bots ping",
}

func PING_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, PING_Command.Name)

	if !ok {
		return
	}
	
	latency := s.HeartbeatLatency()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🏓 Pong! WebSocket latency: %dms", latency.Milliseconds()),
		},
	})
}
