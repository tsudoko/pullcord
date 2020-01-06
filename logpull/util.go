package logpull

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isDiscordError(e error, code int) bool {
	r, ok := e.(*discordgo.RESTError)
	return ok && r.Message != nil && r.Message.Code == code
}

func isBotSession(d *discordgo.Session) bool {
	return strings.HasPrefix(strings.ToLower(d.Token), "bot ")
}
