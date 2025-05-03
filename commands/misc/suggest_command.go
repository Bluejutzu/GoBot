package misc

import (
	"fmt"

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var SUGGEST_Command = &discordgo.ApplicationCommand{
	Name:        "suggest",
	Description: "Suggest something like yeah",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        3,
			Name:        "for",
			Description: "What are you suggesting for",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Discord",
					Value: "discord",
				},
			},
		},
		{
			Type:        3,
			Name:        "suggestion",
			Description: "What are you suggesting",
			Required:    true,
		},
	},
}

func SUGGEST_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, SUGGEST_Command.Name)
	if !ok {
		return
	}

	var (
		suggestion string
		choice     string = "general" // default
	)

	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "suggestion":
			suggestion = opt.StringValue()
		case "for":
			choice = opt.StringValue()
		}
	}

	user := i.Member.User
	embed := &discordgo.MessageEmbed{
		Title:       "New Suggestion",
		Description: suggestion,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Category", Value: choice, Inline: true},
			{Name: "From", Value: fmt.Sprintf("<@%s>", user.ID), Inline: true},
		},
		Color: 0x00bfff,
	}

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Approve",
					Style:    discordgo.SuccessButton,
					CustomID: fmt.Sprintf("suggest_approve_%s", user.ID),
				},
				discordgo.Button{
					Label:    "Deny",
					Style:    discordgo.DangerButton,
					CustomID: fmt.Sprintf("suggest_deny_%s", user.ID),
				},
				discordgo.Button{
					Label:    "Ask for More Info",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("suggest_ask_%s", user.ID),
				},
			},
		},
	}

	// Send the embed to a dedicated suggestion channel (replace channel ID)
	s.ChannelMessageSendComplex("SUGGESTION_CHANNEL_ID", &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: buttons,
	})

	// Acknowledge user
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thanks for your suggestion!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
