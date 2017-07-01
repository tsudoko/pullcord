// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/cdndl"
	"github.com/tsudoko/pullcord/logcache"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/logformat"
	"github.com/tsudoko/pullcord/logutil"
)

func Pull(d *discordgo.Session, c discordgo.Channel, fetchedGuilds *map[string]bool) {
	last := "0"
	filename := fmt.Sprintf("channels/%s/%s.tsv", c.GuildID, c.ID)
	guildFilename := fmt.Sprintf("channels/%s/guild.tsv", c.GuildID)

	if err := os.MkdirAll(path.Dir(filename), os.ModeDir|0755); err != nil {
		log.Printf("[%s/%s] creating the guild dir failed", c.GuildID, c.ID)
		return
	}

	if !(*fetchedGuilds)[c.GuildID] {
		(*fetchedGuilds)[c.GuildID] = true
		cache := make(logcache.Entries)
		if _, err := os.Stat(guildFilename); err == nil {
			if err := logcache.NewEntries(guildFilename, &cache); err != nil {
				log.Printf("[%s] error reconstructing guild state, skipping (%v)", c.GuildID, err)
				return
			}
		}
		pullGuild(d, c.GuildID, cache)
	}

	if _, err := os.Stat(filename); err == nil {
		last, err = logutil.LastMessageID(filename)
		if err != nil {
			log.Printf("[%s/%s] error getting last message id, skipping (%v)", c.GuildID, c.ID, err)
			return
		}
	}

	pullChannel(d, c.GuildID, c.ID, last)
}

func pullGuild(d *discordgo.Session, id string, cache logcache.Entries) {
	deleted := cache.IDs()
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
	} else {
		if guild.Icon != "" {
			err := cdndl.Icon(id, guild.Icon)
			if err != nil {
				log.Printf("[%s] error downloading the guild icon: %v", id, err)
			}
		}

		if guild.Splash != "" {
			err := cdndl.Splash(id, guild.Splash)
			if err != nil {
				log.Printf("[%s] error downloading the guild splash: %v", id, err)
			}
		}

		gEntry := logentry.Make("history", "add", guild)
		logutil.WriteNew(f, gEntry, &cache)
		delete(deleted[logentry.Type(guild)], guild.ID)
	}

	for _, c := range guild.Channels {
		cEntry := logentry.Make("history", "add", c)
		logutil.WriteNew(f, cEntry, &cache)
		delete(deleted[logentry.Type(c)], c.ID)

		for _, o := range c.PermissionOverwrites {
			oEntry := logentry.Make("history", "add", o)
			logutil.WriteNew(f, oEntry, &cache)
			delete(deleted[logentry.Type(o)], o.ID)
		}
	}

	for _, r := range guild.Roles {
		rEntry := logentry.Make("history", "add", r)
		logutil.WriteNew(f, rEntry, &cache)
		delete(deleted[logentry.Type(r)], r.ID)
	}

	for _, e := range guild.Emojis {
		err := cdndl.Emoji(e.ID)
		if err != nil {
			log.Printf("[%s] error downloading emoji %s: %v", id, e.ID, err)
		}
		eEntry := logentry.Make("history", "add", e)
		logutil.WriteNew(f, eEntry, &cache)
		delete(deleted[logentry.Type(e)], e.ID)
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
				err := cdndl.Avatar(m.User)
				if err != nil {
					log.Printf("[%s] error downloading avatar for user %s: %v", id, m.User.ID, err)
				}
			}

			uEntry := logentry.Make("history", "add", m)
			logutil.WriteNew(f, uEntry, &cache)
			delete(deleted[logentry.Type(m)], m.User.ID)
		}

		log.Printf("[%s] downloaded %d members, last id %s with name %s", id, len(members), after, members[len(members)-1].User.Username)
	}

	for etype, ids := range deleted {
		for id := range ids {
			entry := cache[etype][id]
			entry[logentry.HTime] = logentry.Timestamp()
			entry[logentry.HOp] = "del"
			logformat.Write(f, entry)
		}
	}
}

func pullChannel(d *discordgo.Session, gid, id, after string) {
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
			logformat.Write(f, logentry.Make("history", "add", msgs[i]))

			for _, e := range msgs[i].Embeds {
				logformat.Write(f, logentry.Make("history", "add", &logentry.Embed{*e, msgs[i].ID}))
			}

			for _, a := range msgs[i].Attachments {
				log.Printf("[%s/%s] downloading attachment %s", gid, id, a.ID)
				err := cdndl.Attachment(a.URL)
				if err != nil {
					log.Printf("[%s/%s] error downloading attachment %s: %v", gid, id, a.ID, err)
				}
				logformat.Write(f, logentry.Make("history", "add", &logentry.Attachment{*a, msgs[i].ID}))
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
					logformat.Write(f, logentry.Make("history", "add", reaction))
				}

				if r.Count > 100 {
					reaction := &discordgo.MessageReaction{
						UserID:    "",
						MessageID: msgs[i].ID,
						Emoji:     *r.Emoji,
						ChannelID: id,
					}
					for i := 0; i < r.Count-100; i++ {
						logformat.Write(f, logentry.Make("history", "add", reaction))
					}
				}
			}
		}

		log.Printf("[%s/%s] downloaded %d messages, last id %s with content %s", gid, id, len(msgs), msgs[0].ID, msgs[0].Content)
	}
}
