package moderation

import (
	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)


var BAN_Command = &discordgo.ApplicationCommand{
	Name: "ban",
	Description: "Ban a member from the current guild",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: 9,
			Name: "user",
			Description: "The member to ban",
			Required: true,
		},
		{
			Type: 3,
			Name: "reason",
			Description: "Why this member is being banned",
			Required: false,
		},
		{
			Type: 5,
			Name: "delete-messages",
			Description: "Delete messages from x days ago, leave empty if none (max: 7)",
			Required: false,
			MaxValue: 7,
		},
	},
}

func BAN_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, BAN_Command.Name)

	if !ok {
		return
	}

	
}
