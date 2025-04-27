package misc

import (
	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var TEST_Command = &discordgo.ApplicationCommand{
	Name:        "test",
	Description: "A test command to demonstrate CLI validation",
}

func TEST_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, TEST_Command.Name)
	if !ok {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Test command works!",
		},
	})
}
