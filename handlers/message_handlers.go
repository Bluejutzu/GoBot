package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func MessageCreate(c *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == c.State.User.ID {
		return
	}
	print("Message from ", m.Author.Username, ": ", m.Content, "\n")

	if strings.HasPrefix(m.Content, "!") {
		print("Command detected: ", m.Content, "\n")
		c.ChannelMessageSendReply(m.ChannelID, "Hello!", m.Reference())
	}
}
	