package ridealong

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	switch i.ModalSubmitData().CustomID {
	case "ra_request_modal":
		handleRequestModalSubmit(s, i)
	}
}

func handleRequestModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	callsign := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	embed := &discordgo.MessageEmbed{
		Title:       "Ride Along Request Submitted",
		Description: fmt.Sprintf("Callsign: %sStatus: Pending", callsign),
		Color:       0xffaa00,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
