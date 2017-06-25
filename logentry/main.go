package logentry

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func formatBool(name string, variable bool) string {
	if variable {
		return name
	} else {
		return ""
	}
}

func Message(op string, m *discordgo.Message) []string {
	var ts string

	// exact deletions can't be retrieved from audit logs or the history
	// TODO: move to the outer handler?
	if op == "del" {
		ts = time.Now().Format("2006-01-02T15:04:05.000000-07:00")
	} else if m.EditedTimestamp != "" {
		ts = string(m.EditedTimestamp)
	}

	return []string{op, "message", m.ID, m.Author.ID, ts, formatBool("tts", m.Tts), m.Content}
}

func Attachment(op string, messageID string, a *discordgo.MessageAttachment) []string {
	return []string{op, "attachment", messageID, a.ID}
}

func Reaction(op string, r *discordgo.MessageReaction) []string {
	return []string{op, "reaction", r.MessageID, r.UserID, r.Emoji.APIName()}
}

func Embed(op string, messageID string, e *discordgo.MessageEmbed) []string {
	j, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return []string{op, "embed", messageID, string(j)}
}

func Guild(op string, g *discordgo.Guild) []string {
	return []string{
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
		//g.MFALevel,
	}
}

func User(op string, m *discordgo.Member) []string {
	var name string

	if m.Nick != "" {
		name = m.Nick
	} else {
		name = m.User.Username
	}

	return []string{op, "user", m.User.ID, name, m.User.Avatar}
}

func Role(op string, r *discordgo.Role) []string {
	return []string{
		op,
		"role",
		r.ID,
		r.Name,
		strconv.Itoa(r.Color),
		strconv.Itoa(r.Position),
		strconv.Itoa(r.Permissions),
		formatBool("hoist", r.Hoist),
	}
}

func RoleAssign(op string, uid, rid string) []string {
	return []string{op, "roleassign", uid, rid}
}

func Channel(op string, c *discordgo.Channel) []string {
	return []string{
		op,
		"channel",
		c.ID,
		c.Type,
		strconv.Itoa(c.Position),
		c.Name,
		c.Topic,
	}
}

func PermOverwrite(op string, o *discordgo.PermissionOverwrite) []string {
	return []string{
		op,
		"permoverwrite",
		o.ID,
		o.Type,
		strconv.Itoa(o.Allow),
		strconv.Itoa(o.Deny),
	}
}

func Emoji(op string, e *discordgo.Emoji) []string {
	return []string{
		op,
		"emoji",
		e.ID,
		e.Name,
		formatBool("nocolons", !e.RequireColons),
	}
}
