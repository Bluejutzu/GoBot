package moderation

import (
	"fmt"

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

	dev := false

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
		data             = i.ApplicationCommandData()
		memberToBan, err = s.GuildMember(i.GuildID, data.Options[0].UserValue(s).ID)
		reason           string
	)

	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to ban user: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if len(data.Options) > 1 {
		reason = data.Options[1].StringValue()
	}

	err = s.GuildBanCreateWithReason(i.GuildID, memberToBan.User.ID, reason, 7)
	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to ban user: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err = s.GuildBanDelete(i.GuildID, memberToBan.User.ID)
	if err != nil && !dev {
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
			Content: "Successfully softbanned " + memberToBan.User.Username + "\n DM'ing the user an invite link...",
		},
	})

	channel, err := s.UserChannelCreate(memberToBan.User.ID)

	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to Create user DM: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	guild, err := s.Guild(i.GuildID)

	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch guild: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	invite, err := s.ChannelInviteCreate(guild.SystemChannelID, discordgo.Invite{
		MaxAge:    86400,
		MaxUses:   1,
		Temporary: false,
	})

	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create invite: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("You've been banned from %v for %v. This is a new invite created to rejoin: %v", guild.Name, reason, invite))

	if err != nil && !dev {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to DM invite to user: " + err.Error() + fmt.Sprintf("\n-# %v", invite),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
}
