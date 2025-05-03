package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleSuggestionComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	parts := strings.Split(data.CustomID, "_")
	if len(parts) < 3 || parts[0] != "suggest" {
		return
	}

	action := parts[1]
	userID := parts[2]

	switch action {
	case "approve", "deny":
		status := map[string]string{
			"approve": "✅ Approved",
			"deny":    "❌ Denied",
		}[action]

		embed := i.Message.Embeds[0]
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Status",
			Value: status,
		})
		embed.Color = map[string]int{
			"approve": 0x00ff00,
			"deny":    0xff0000,
		}[action]

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{}},
				},
			},
		})

	case "ask":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				Title:    "Ask for more info",
				CustomID: fmt.Sprintf("suggest_ask_modal_%s", userID),
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "reason",
								Label:       "Message to user",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Could you elaborate on your suggestion?",
								Required:    true,
							},
						},
					},
				},
			},
		})
	}
}

func HandleSuggestionModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	if !strings.HasPrefix(data.CustomID, "suggest_ask_modal_") {
		return
	}

	userID := strings.TrimPrefix(data.CustomID, "suggest_ask_modal_")
	message := ""

	for _, row := range data.Components {
		for _, comp := range row.(*discordgo.ActionsRow).Components {
			if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == "reason" {
				message = input.Value
			}
		}
	}

	channel, err := s.UserChannelCreate(userID)
	if err == nil {
		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Content: fmt.Sprintf("📩 Admin has responded to your suggestion:\n> %s", message),
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Message sent to the user.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
