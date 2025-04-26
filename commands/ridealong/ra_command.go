package ridealong

import (
	"github.com/bwmarrin/discordgo"
)

var (
	Command = &discordgo.ApplicationCommand{
		Name:        "ra",
		Description: "Ride along panel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "request",
				Description: "Request a ride along",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "log",
				Description: "Log a ride along session",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "officer",
						Description: "Officer's callsign",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		},
	}
)

func ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "ra" {
		return
	}

	switch data.Options[0].Name {
	case "request":
		handleRequest(s, i)
	case "log":
		handleLog(s, i)
	}
}


