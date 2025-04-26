package misc

import (
	"fmt"
	"math/rand"

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var (
	EIGHTBALL_Command = &discordgo.ApplicationCommand{
		Name:        "8ball",
		Description: "What does your fate hold?",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "What do you want to ask the magical 8ball?",
				Required:    true,
			},
		},
	}

	eightball_responses = []string{
		"Certainly",
		"No doubt!",
		"Not sure on this one...",
		"Should rephrase your question",
		"100%! Never been more sure in my life",
		"Uhhh what",
		"No idea what that means",
		"Yes",
		"No",
		"50/50",
	}
)

func EIGHTBALL_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, EIGHTBALL_Command.Name)

	if !ok {
		return
	}

	if len(eightball_responses) <= 0 {
		return
	}
	data := i.ApplicationCommandData()

	curr_response := eightball_responses[rand.Intn(len(eightball_responses))]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("**%v** \n-# %v", curr_response, data.Options[0].Value),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
