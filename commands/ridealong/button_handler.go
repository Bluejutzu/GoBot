package ridealong

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	cacheKey := parts[1]

	_, exists := raLogCache[cacheKey]
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This session has expired or data was not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var color int
	var status string

	switch action {
	case "approve_promotion":
		color = 0x00FF00
		status = "Status: Approved"
	case "reject_promotion":
		color = 0xFF0000
		status = "Status: Rejected"
	default:
		return
	}

	message := i.Message
	if len(message.Embeds) == 0 {
		return
	}

	embed := message.Embeds[0]
	newEmbed := &discordgo.MessageEmbed{
		Title:       embed.Title,
		Description: embed.Description,
		Color:       color,
		Author: &discordgo.MessageEmbedAuthor{
			Name: status,
		},
	}

	delete(raLogCache, cacheKey)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{newEmbed},
			Components: []discordgo.MessageComponent{},
		},
	})
}

func showRequestModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modal := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "ra_request_modal",
			Title:    "Ride Along Request",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "callsign",
							Label:     "Your Callsign",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MinLength: 2,
							MaxLength: 10,
						},
					},
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, modal)
}

func completeLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ride along completed",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
