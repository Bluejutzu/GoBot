package ridealong

import (
	"github.com/bwmarrin/discordgo"
)

func HandleButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	switch i.MessageComponentData().CustomID {
	case "ra_request":
		showRequestModal(s, i)
	case "ra_log_complete":
		completeLog(s, i)
	}
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

func completeLog(s *discordgo.Session, i* discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ride along completed",
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
