// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/logentry"
)

func Channel(d *discordgo.Session, id string) {
	msgid := "0"
	for {
		msgs, err := d.ChannelMessages(id, 100, "", msgid, "")
		if err != nil {
			log.Printf("[%s] error getting messages from %s: %v", id, msgid, err)
		}

		if len(msgs) == 0 {
			break
		}

		msgid = msgs[0].ID

		// messages are retrieved in descending order
		for i := len(msgs) - 1; i >= 0; i-- {
			fmt.Println(logentry.Message("add", msgs[i]))

			for _, e := range msgs[i].Embeds {
				fmt.Println(logentry.Embed("add", msgs[i].ID, e))
			}

			for _, a := range msgs[i].Attachments {
				fmt.Println(logentry.Attachment("add", msgs[i].ID, a))
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
					fmt.Println(logentry.Reaction("add", reaction))
				}

				if r.Count > 100 {
					reaction := &discordgo.MessageReaction{
						UserID:    "",
						MessageID: msgs[i].ID,
						Emoji:     *r.Emoji,
						ChannelID: id,
					}
					for i := 0; i < r.Count-100; i++ {
						fmt.Println(logentry.Reaction("add", reaction))
					}
				}
			}
		}

		log.Printf("[%s] downloaded %d messages, last id: %s with content %s", id, len(msgs), msgs[0].ID, msgs[0].Content)
	}
}

func Guild(d *discordgo.Session, id string) {

}
