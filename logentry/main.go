// Package logentry describes the format of log entries.
package logentry

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const timeFormat = "2006-01-02T15:04:05.000000-07:00"

func formatBool(name string, variable bool) string {
	if variable {
		return name
	} else {
		return ""
	}
}

func wrap(in []string) []string {
	return append([]string{time.Now().Format(timeFormat)}, in...)
}

func Message(ftype, op string, m *discordgo.Message) []string {
	row := []string{
		ftype,
		op,
		"message",
		m.ID,
		m.Author.ID,
		string(m.EditedTimestamp),
		formatBool("tts", m.Tts),
		m.Content,
	}
	return wrap(row)
}

func Attachment(ftype, op string, messageID string, a *discordgo.MessageAttachment) []string {
	row := []string{ftype, op, "attachment", a.ID, messageID}
	return wrap(row)
}

func Reaction(ftype, op string, r *discordgo.MessageReaction) []string {
	row := []string{ftype, op, "reaction", r.UserID, r.MessageID, r.Emoji.APIName()}
	return wrap(row)
}

func Embed(ftype, op string, messageID string, e *discordgo.MessageEmbed) []string {
	j, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	row := []string{ftype, op, "embed", messageID, string(j)}
	return wrap(row)
}

func Guild(ftype, op string, g *discordgo.Guild) []string {
	row := []string{
		ftype,
		op,
		"guild",
		g.ID,
		g.Name,
		g.Icon,
		g.Splash,
		g.OwnerID,
		g.AfkChannelID,
		strconv.Itoa(g.AfkTimeout),
		formatBool("embeddable", g.EmbedEnabled),
		g.EmbedChannelID,
	}
	return wrap(row)
}

func User(ftype, op string, m *discordgo.Member) []string {
	row := []string{
		ftype,
		op,
		"user",
		m.User.ID,
		m.User.Username,
		m.Nick,
		m.User.Discriminator,
		m.User.Avatar,
		strings.Join(m.Roles, ","),
	}
	return wrap(row)
}

func Role(ftype, op string, r *discordgo.Role) []string {
	row := []string{
		ftype,
		op,
		"role",
		r.ID,
		r.Name,
		strconv.Itoa(r.Color),
		strconv.Itoa(r.Position),
		strconv.Itoa(r.Permissions),
		formatBool("hoist", r.Hoist),
	}
	return wrap(row)
}

func Channel(ftype, op string, c *discordgo.Channel) []string {
	row := []string{
		ftype,
		op,
		"channel",
		c.ID,
		c.Type,
		strconv.Itoa(c.Position),
		c.Name,
		c.Topic,
	}
	return wrap(row)
}

func PermOverwrite(ftype, op string, o *discordgo.PermissionOverwrite) []string {
	row := []string{
		ftype,
		op,
		"permoverwrite",
		o.ID,
		o.Type,
		strconv.Itoa(o.Allow),
		strconv.Itoa(o.Deny),
	}
	return wrap(row)
}

func Emoji(ftype, op string, e *discordgo.Emoji) []string {
	row := []string{
		ftype,
		op,
		"emoji",
		e.ID,
		e.Name,
		formatBool("nocolons", !e.RequireColons),
	}
	return wrap(row)
}
