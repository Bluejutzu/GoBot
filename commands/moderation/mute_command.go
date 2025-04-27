package moderation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bluejutzu/GoBot/helpers"
	"github.com/bwmarrin/discordgo"
)

var MUTE_Command = &discordgo.ApplicationCommand{
	Name:        "mute",
	Description: "Mute a member from the current guild",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        9,
			Name:        "user",
			Description: "The member to mute",
			Required:    true,
		},
		{
			Type:        3,
			Name:        "reason",
			Description: "Why this member is being muted",
			Required:    false,
		},
		{
			Type:        3,
			Name:        "duration",
			Description: "Duration of the mute (ex: 1h, 30min, 0.5min, 10d)",
			Required:    false,
		},
	},
}

func MUTE_ParseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ok := helpers.SafeCommandParse(i, MUTE_Command.Name)

	if !ok {
		return
	}

	// Command was invoked in a DM
	if i.Interaction.Member == nil || i.Interaction.User != nil {
		return
	}

	if i.Interaction.Member.Permissions&discordgo.PermissionModerateMembers != discordgo.PermissionModerateMembers {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Insufficient Permissions",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var (
		data         discordgo.ApplicationCommandInteractionData = i.ApplicationCommandData()
		memberToMute *discordgo.User                             = data.Options[0].UserValue(s)
		reason       string
		duration     string
		guildID      = i.Interaction.Member.GuildID
	)

	if len(data.Options) > 1 {
		reason = data.Options[1].StringValue()
	}

	if len(data.Options) > 2 {
		duration = data.Options[2].StringValue()
	}

	timeoutUntil, err := parseDuration(duration)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid duration format: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	err = s.GuildMemberTimeout(guildID, memberToMute.ID, timeoutUntil)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to timeout user: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	content := fmt.Sprintf("Successfully muted %s", memberToMute.Username)
	if duration != "" {
		content += fmt.Sprintf(" for %s", duration)
	}
	if reason != "" {
		content += fmt.Sprintf(" (Reason: %s)", reason)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

// parseDuration converts duration strings like "30min", "1h", "2d" into a *time.Time
func parseDuration(durStr string) (*time.Time, error) {
	if durStr == "" {
		return nil, fmt.Errorf("duration cannot be empty")
	}

	re := regexp.MustCompile(`^(\d*\.?\d+)([a-zA-Z]+)$`)
	matches := re.FindStringSubmatch(durStr)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid duration format. Expected format: number + unit (e.g., 30min, 1h, 2d)")
	}

	value := matches[1]
	unit := strings.ToLower(matches[2])

	// Parse the numeric value
	var duration time.Duration
	var err error

	switch unit {
	case "s", "sec", "seconds":
		duration, err = time.ParseDuration(value + "s")
	case "min", "mins":
		duration, err = time.ParseDuration(value + "m")
	case "h", "hr", "hrs", "hour", "hours":
		duration, err = time.ParseDuration(value + "h")
	case "d", "day", "days":
		hours, err := time.ParseDuration(value + "h")
		if err == nil {
			duration = hours * 24
		}
	default:
		return nil, fmt.Errorf("unsupported duration unit: %s", unit)
	}

	if err != nil {
		return nil, err
	}

	timeoutTime := time.Now().Add(duration)
	return &timeoutTime, nil
}
