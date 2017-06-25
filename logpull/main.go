// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/logformat"
)

func Channel(d *discordgo.Session, gid, id, after string) {
	filename := fmt.Sprintf("channels/%s/%s.tsv", gid, id)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[%s] error opening the log file: %v", id, err)
		return
	}

	for {
		msgs, err := d.ChannelMessages(id, 100, "", after, "")
		if err != nil {
			log.Printf("[%s] error getting messages from %s: %v", id, after, err)
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
				logformat.Write(f, logentry.Attachment("add", msgs[i].ID, a))
			}

			for _, r := range msgs[i].Reactions {
				users, err := d.MessageReactions(id, msgs[i].ID, r.Emoji.APIName(), 100)
				if err != nil {
					log.Printf("[%s] error getting users for reaction %s to %s: %v", id, r.Emoji.APIName(), msgs[i].ID, err)
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

		log.Printf("[%s] downloaded %d messages, last id: %s with content %s", id, len(msgs), msgs[0].ID, msgs[0].Content)
	}
}

func Guild(d *discordgo.Session, id string) {

}
