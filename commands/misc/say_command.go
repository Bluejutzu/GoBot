package misc

import (

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var SAY_COMMAND = &discordgo.ApplicationCommand{
	Name:        "say",
	Description: "Send a message via the bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        3,
			Name:        "content",
			Description: "What do you want to send",
			Required:    true,
		},
		{
			Type:        3,
			Name:        "reference-id",
			Description: "ID of the message to reply to (Optional)",
			Required:    false,
		},
		{
			Type:        3,
			Name:        "channel-id",
			Description: "ID of the channel of the message to reply to (Required if reference-id)",
			Required:    false,
		},
	},
}

func SAY_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, SAY_COMMAND.Name)

	if !ok {
		return
	}

	// Command was invoked in a DM
	isDM := i.Interaction.Member == nil
	if !isDM && i.Interaction.Member.Permissions&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
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
		content              string
		message_reference_id string
		channel_reference_id string
	)

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	if opt, ok := optionMap["content"]; ok {
		content = opt.StringValue()
	}
	if opt, ok := optionMap["reference-id"]; ok {
		message_reference_id = opt.StringValue()
	}
	if opt, ok := optionMap["channel-id"]; ok {
		channel_reference_id = opt.StringValue()
	}

	if message_reference_id != "" && channel_reference_id == "" {
		return
	} else if channel_reference_id != "" && message_reference_id == "" {
		return
	}

	if message_reference_id != "" {
		_, err := s.ChannelMessageSendReply(channel_reference_id, content, &discordgo.MessageReference{
			MessageID: message_reference_id,
		})

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to reply to message " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully sent message",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		_, err := s.ChannelMessageSend(i.ChannelID, content)

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to send message " + err.Error(),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully sent message",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
