package logentry

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Message(op string, m *discordgo.Message) []string {
	var (
		ts  string
		tts string
	)

	switch op {
	case "edit":
		ts = string(m.EditedTimestamp)
	case "del": // exact deletions can't be retrieved from audit logs or the history
		// TODO: move to the outer handler?
		ts = time.Now().Format("2006-01-02T15:04:05.000000-07:00")
	}

	if m.Tts {
		tts = "tts"
	}

	return []string{op, "message", m.ID, m.Author.ID, ts, tts, m.Content}
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

func User(op string, m *discordgo.Member) []string {
	var name string

	if m.Nick != "" {
		name = m.Nick
	} else {
		name = m.User.Username
	}

	return []string{op, "user", m.User.ID, name, m.User.Avatar}
}
