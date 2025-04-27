package moderation

import (
	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var SOFTBAN_Command = &discordgo.ApplicationCommand{
	Name:        "softban",
	Description: "Ban a user imemdiately to purge their messages and then unban them",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        9,
			Name:        "user",
			Description: "The member to softban",
			Required:    true,
		},
		{
			Type:        3,
			Name:        "reason",
			Description: "Why this member is being softbanned",
			Required:    false,
		},
	},
}

func SOFTBAN_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, SOFTBAN_Command.Name)

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
	)

	if len(data.Options) > 1 {
		reason = data.Options[1].StringValue()
	}

	err := s.GuildBanCreateWithReason(i.GuildID, memberToBan.ID, reason, 7)
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

	err = s.GuildBanDelete(i.GuildID, memberToBan.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to unban user: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully softbanned " + memberToBan.Username,
		},
	})
}
