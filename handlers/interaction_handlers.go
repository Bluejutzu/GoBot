package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(c *discordgo.Session, i *discordgo.InteractionCreate) {
	print("Interaction detected: ", i.Interaction.Type, i.Interaction.ID, "\n")
}
