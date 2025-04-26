package helpers

import "github.com/bwmarrin/discordgo"

func SafeCommandParse(i *discordgo.InteractionCreate, n string) (ok bool) {
	return i != nil && // Is there a interaction
	i.Interaction != nil && // Does the interaction have the proper values
	i.Type == discordgo.InteractionApplicationCommand && // Is the interaction a Command
	i.ApplicationCommandData().Name == n // Is the command name right
}
