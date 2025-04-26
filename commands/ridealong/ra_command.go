package ridealong

import (
	"github.com/bwmarrin/discordgo"
)

type RaLogData struct {
	RecruitID    string
	PassFail     string
	DrivingScore int
	GrammarScore int
	TotalScore   int
	Callsign     string
	Performance  string
	Notes        string
}

var raLogCache = make(map[string]RaLogData)

var Command = &discordgo.ApplicationCommand{
	Name:        "ra",
	Description: "Ride along panel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "request",
			Description: "Request a ride along",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "log",
			Description: "Log a ride along",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "recruit",
					Description: "The recruit you did an R/A with",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "pass-fail",
					Description: "Did the recruit fail or pass?",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Pass",
							Value: "pass",
						},
						{
							Name:  "Fail",
							Value: "fail",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "driving-score",
					Description: "Was the recruit abiding to road laws? Out of 10.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "grammar-score",
					Description: "Was the recruit using sufficient grammar? Out of 10.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "total-score",
					Description: "How would you rate the recruit overall? Out of 10.",
					Required:    true,
				},
			},
		},
	},
}

func ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "ra" {
		return
	}
	print(data.Options[0].Name)
	switch data.Options[0].Name {
	case "request":
		handleRequest(s, i)
	case "log":
		handleLog(s, i)
	}
}
