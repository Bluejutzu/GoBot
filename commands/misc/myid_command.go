package misc

import (
	"fmt"

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var ID_Commmand = &discordgo.ApplicationCommand{
	Name:        "what-is-my-id",
	Description: "Get your user ID.",
}

func ID_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, ID_Commmand.Name)

	if !ok {
		return
	}

	// Check if Member is nil before accessing
	if i.Interaction.Member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error: Could not retrieve user information",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	userID := i.Interaction.Member.User.ID
	guildID := "Direct Message"
	if i.Interaction.GuildID != "" {
		guildID = i.Interaction.GuildID
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@!%s> \nUser ID: `%s` \nServer ID: %s", userID, userID, guildID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
