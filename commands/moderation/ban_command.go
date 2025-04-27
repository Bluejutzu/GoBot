package moderation

import (
	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var BAN_Command = &discordgo.ApplicationCommand{
	Name:        "ban",
	Description: "Ban a member from the current guild",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        9,
			Name:        "user",
			Description: "The member to ban",
			Required:    true,
		},
		{
			Type:        3,
			Name:        "reason",
			Description: "Why this member is being banned",
			Required:    false,
		},
		{
			Type:        10,
			Name:        "delete-messages",
			Description: "Delete messages from x days ago, leave empty if none (max: 7)",
			Required:    false,
			MaxValue:    7,
		},
	},
}

func BAN_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, BAN_Command.Name)

	if !ok {
		return
	}

	// Command was invoked in a DM
	if i.Interaction.Member == nil || i.Interaction.User != nil {
		return
	}

	// Check permissions of Member
	if i.Interaction.Member.Permissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Insufficient Permissions",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	var (
		data        = i.ApplicationCommandData()
		memberToBan = data.Options[0].UserValue(s)
		reason      string
		days        int
	)
	if len(data.Options) > 1 {
		reason = data.Options[1].StringValue()
	}

	if len(data.Options) > 2 {
		days = int(data.Options[2].FloatValue())
	}

	guildID := i.GuildID

	err := s.GuildBanCreateWithReason(guildID, memberToBan.ID, reason, days)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to ban user: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully banned " + memberToBan.Username,
		},
	})
}
