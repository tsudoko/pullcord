// Package pull contains functions related to downloading historical data.
package pull

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/entry"
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
		for i := len(msgs)-1; i >= 0; i-- {
			fmt.Println(entry.Message("add", msgs[i]))

			for _, e := range msgs[i].Embeds {
				fmt.Println(entry.Embed("add", msgs[i].ID, e))
			}

			for _, a := range msgs[i].Attachments {
				fmt.Println(entry.Attachment("add", msgs[i].ID, a))
			}

			// a bit more complicated, needs to iterate over d.MessageReactions(id, msgs[i].ID, emojiID, $limit)
			// might be expensive
			/*
			for _, r := range msgs[i].Reactions {

			}
			*/
		}

		log.Printf("[%s] downloaded %d messages, last id: %s with content %s", id, len(msgs), msgs[0].ID, msgs[0].Content)
	}
}

func Guild(d *discordgo.Session, id string) {

}