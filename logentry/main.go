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

type Attachment struct {
	discordgo.MessageAttachment
	MessageID string
}

type Embed struct {
	discordgo.MessageEmbed
	MessageID string
}

func formatBool(name string, variable bool) string {
	if variable {
		return name
	} else {
		return ""
	}
}

func Timestamp(row []string) []string {
	return append([]string{time.Now().Format(timeFormat)}, row...)
}

func Type(v interface{}) string {
	switch v.(type) {
	case *discordgo.Message:
		return "message"
	case *Attachment:
		return "attachment"
	case *discordgo.MessageReaction:
		return "reaction"
	case *Embed:
		return "embed"
	case *discordgo.Guild:
		return "guild"
	case *discordgo.Member:
		return "user"
	case *discordgo.Role:
		return "role"
	case *discordgo.Channel:
		return "channel"
	case *discordgo.PermissionOverwrite:
		return "permoverwrite"
	case *discordgo.Emoji:
		return "emoji"
	default:
		panic("unsupported type")
	}
}

func Make(ftype, op string, v interface{}) []string {
	var row []string

	switch v := v.(type) {
	case *discordgo.Message:
		row = []string{
			v.ID,
			v.Author.ID,
			string(v.EditedTimestamp),
			formatBool("tts", v.Tts),
			v.Content,
		}
	case *Attachment:
		row = []string{v.ID, v.MessageID}
	case *discordgo.MessageReaction:
		row = []string{v.UserID, v.MessageID, v.Emoji.APIName()}
	case *Embed:
		j, err := json.Marshal(v.MessageEmbed)
		if err != nil {
			panic(err)
		}

		row = []string{v.MessageID, string(j)}
	case *discordgo.Guild:
		row = []string{
			v.ID,
			v.Name,
			v.Icon,
			v.Splash,
			v.OwnerID,
			v.AfkChannelID,
			strconv.Itoa(v.AfkTimeout),
			formatBool("embeddable", v.EmbedEnabled),
			v.EmbedChannelID,
		}
	case *discordgo.Member:
		row = []string{
			v.User.ID,
			v.User.Username,
			v.Nick,
			v.User.Discriminator,
			v.User.Avatar,
			strings.Join(v.Roles, ","),
		}
	case *discordgo.Role:
		row = []string{
			v.ID,
			v.Name,
			strconv.Itoa(v.Color),
			strconv.Itoa(v.Position),
			strconv.Itoa(v.Permissions),
			formatBool("hoist", v.Hoist),
		}
	case *discordgo.Channel:
		row = []string{
			v.ID,
			v.Type,
			strconv.Itoa(v.Position),
			v.Name,
			v.Topic,
		}
	case *discordgo.PermissionOverwrite:
		row = []string{
			v.ID,
			v.Type,
			strconv.Itoa(v.Allow),
			strconv.Itoa(v.Deny),
		}
	case *discordgo.Emoji:
		row = []string{
			v.ID,
			v.Name,
			formatBool("nocolons", !v.RequireColons),
		}
	default:
		panic("unsupported type")
	}

	return Timestamp(append([]string{ftype, op, Type(v)}, row...))
}
