package ridealong

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handleRequest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ride along request system coming soon",
		},
	})
}

func handleLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData().Options[0]
	recruit := data.Options[0].UserValue(s)
	passFail := data.Options[1].StringValue()
	driving := int(data.Options[2].IntValue())
	grammar := int(data.Options[3].IntValue())
	total := int(data.Options[4].IntValue())

	cacheKey := fmt.Sprintf("%s_%s", i.Member.User.ID, i.ID)
	raLogCache[cacheKey] = RaLogData{
		RecruitID:    recruit.ID,
		PassFail:     passFail,
		DrivingScore: driving,
		GrammarScore: grammar,
		TotalScore:   total,
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "performance",
					Label:       "Performance",
					Style:       discordgo.TextInputParagraph,
					Required:    true,
					MaxLength:   500,
					Placeholder: "How did the cadet perform?",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "notes",
					Label:       "Notes (Optional)",
					Style:       discordgo.TextInputShort,
					Required:    false,
					Placeholder: "This will be DMed to the Cadet",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "callsign",
					Label:       "Callsign",
					Style:       discordgo.TextInputShort,
					Required:    true,
					Placeholder: "Only the numbers: 000 - if you failed the cadet then put in N/A",
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   cacheKey,
			Title:      "R/A Logging Results",
			Components: components,
		},
	})
}
