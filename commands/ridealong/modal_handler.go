package ridealong

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := i.ModalSubmitData()
	cached, exists := raLogCache[data.CustomID]
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

	performance := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	notes := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	if notes == "" {
		notes = "No notes provided."
	}
	callsign := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	cached.Performance = performance
	cached.Notes = notes
	cached.Callsign = callsign
	raLogCache[data.CustomID] = cached

	color := 0x0000FF
	status := "Status: Pending Promotion"
	if cached.PassFail == "fail" {
		color = 0xFF0000
		status = "Status: Failed"
	}

	embed := &discordgo.MessageEmbed{
		Title: "Ride Along Results",
		Description: fmt.Sprintf("**Pass/Fail:** %s\n**Recruit:** <@%s>\n**Overall Score:** %d\n**Driving:** %d\n**Grammar:** %d\n**Field‐Officer:** <@%s>\n**Performance Notes:** %s\n**Additional Notes:** %s\n\nTime logged: %s",
			cached.PassFail, cached.RecruitID, cached.TotalScore, cached.DrivingScore, cached.GrammarScore,
			i.Member.User.ID, performance, notes, time.Now().Format(time.RFC1123)),
		Color: color,
		Author: &discordgo.MessageEmbedAuthor{
			Name: status,
		},
	}

	var components []discordgo.MessageComponent
	if cached.PassFail == "pass" {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Approve Promotion",
						Style:    discordgo.SuccessButton,
						CustomID: fmt.Sprintf("approve_promotion:%s", data.CustomID),
					},
					discordgo.Button{
						Label:    "Reject Promotion",
						Style:    discordgo.DangerButton,
						CustomID: fmt.Sprintf("reject_promotion:%s", data.CustomID),
					},
				},
			},
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}
