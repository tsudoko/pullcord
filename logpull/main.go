// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/cdndl"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/logformat"
)

func Guild(d *discordgo.Session, id string) {
	filename := fmt.Sprintf("channels/%s/guild.tsv", id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[%s] error opening the log file: %v", id, err)
		return
	}
	defer f.Close()

	guild, err := d.Guild(id)
	if err != nil {
		log.Printf("[%s] error getting guild info: %v", id, err)
		//goto members
	}

	if guild.Icon != "" {
		err := cdndl.Icon(id, guild.Icon)
		if err != nil {
			log.Printf("[%s] error downloading the guild icon: %v", id, err)
		}
	}

	logformat.Write(f, logentry.Guild("add", guild))

	for _, c := range guild.Channels {
		logformat.Write(f, logentry.Channel("add", c))

		for _, o := range c.PermissionOverwrites {
			logformat.Write(f, logentry.PermOverwrite("add", o))
		}
	}

	for _, r := range guild.Roles {
		logformat.Write(f, logentry.Role("add", r))
	}

	for _, e := range guild.Emojis {
		err := cdndl.Emoji(e.ID)
		if err != nil {
			log.Printf("[%s] error downloading emoji %s: %v", id, e.ID, err)
		}
		logformat.Write(f, logentry.Emoji("add", e))
	}

	after := "0"
	for {
		members, err := d.GuildMembers(id, after, 1000)
		if err != nil {
			log.Printf("[%s] error getting members from %s: %v", id, after, err)
			continue
		}

		if len(members) == 0 {
			break
		}

		for _, m := range members {
			after = m.User.ID

			if m.User.Avatar != "" {
				err := cdndl.Avatar(m.User.ID, m.User.Avatar)
				if err != nil {
					log.Printf("[%s] error downloading avatar for user %s: %v", id, m.User.ID, err)
				}
			}

			logformat.Write(f, logentry.User("add", m))
		}
	}
}

func Channel(d *discordgo.Session, gid, id, after string) {
	filename := fmt.Sprintf("channels/%s/%s.tsv", gid, id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[%s/%s] error opening the log file: %v", gid, id, err)
		return
	}
	defer f.Close()

	for {
		msgs, err := d.ChannelMessages(id, 100, "", after, "")
		if err != nil {
			log.Printf("[%s/%s] error getting messages from %s: %v", gid, id, after, err)
		}

		if len(msgs) == 0 {
			break
		}

		after = msgs[0].ID

		// messages are retrieved in descending order
		for i := len(msgs) - 1; i >= 0; i-- {
			logformat.Write(f, logentry.Message("add", msgs[i]))

			for _, e := range msgs[i].Embeds {
				logformat.Write(f, logentry.Embed("add", msgs[i].ID, e))
			}

			for _, a := range msgs[i].Attachments {
				log.Printf("[%s/%s] downloading attachment %s", gid, id, a.ID)
				err := cdndl.Attachment(a.URL)
				if err != nil {
					log.Printf("[%s/%s] error downloading attachment %s: %v", gid, id, a.ID, err)
				}
				logformat.Write(f, logentry.Attachment("add", msgs[i].ID, a))
			}

			for _, r := range msgs[i].Reactions {
				users, err := d.MessageReactions(id, msgs[i].ID, r.Emoji.APIName(), 100)
				if err != nil {
					log.Printf("[%s/%s] error getting users for reaction %s to %s: %v", gid, id, r.Emoji.APIName(), msgs[i].ID, err)
				}

				for _, u := range users {
					reaction := &discordgo.MessageReaction{
						UserID:    u.ID,
						MessageID: msgs[i].ID,
						Emoji:     *r.Emoji,
						ChannelID: id,
					}
					logformat.Write(f, logentry.Reaction("add", reaction))
				}

				if r.Count > 100 {
					reaction := &discordgo.MessageReaction{
						UserID:    "",
						MessageID: msgs[i].ID,
						Emoji:     *r.Emoji,
						ChannelID: id,
					}
					for i := 0; i < r.Count-100; i++ {
						logformat.Write(f, logentry.Reaction("add", reaction))
					}
				}
			}
		}

		log.Printf("[%s/%s] downloaded %d messages, last id: %s with content %s", gid, id, len(msgs), msgs[0].ID, msgs[0].Content)
	}
}
