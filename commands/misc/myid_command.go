package misc

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ID_Commmand = &discordgo.ApplicationCommand{
	Name:        "what-is-my-id",
	Description: "Get your user ID.",
}

func ID_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "what-is-my-id" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@!%s> ID: `%v` \nCommand used in server: %v ", i.Interaction.User.ID, i.Interaction.User.ID, i.Interaction.GuildID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
