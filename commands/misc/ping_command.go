package ping

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Command = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Get the Bots ping",
}

func ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "ping" {
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
