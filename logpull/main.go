// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/cdndl"
	"github.com/tsudoko/pullcord/logcache"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/logutil"
	"github.com/tsudoko/pullcord/tsv"
)

type Puller struct {
	PulledGuilds map[string]bool

	d *discordgo.Session

	// per-guild caches
	gCache   logcache.Entries
	gDeleted logcache.IDs
}

func NewPuller(d *discordgo.Session) *Puller {
	return &Puller{
		PulledGuilds: make(map[string]bool),
		d:            d,
	}
}

func (p *Puller) Pull(c discordgo.Channel) error {
	last := "0"
	filename := fmt.Sprintf("channels/%s/%s.tsv", c.GuildID, c.ID)
	guildFilename := fmt.Sprintf("channels/%s/guild.tsv", c.GuildID)

	if err := os.MkdirAll(path.Dir(filename), os.ModeDir|0755); err != nil {
		return errors.New("creating the guild dir failed")
	}

	if !p.PulledGuilds[c.GuildID] {
		p.PulledGuilds[c.GuildID] = true
		p.gCache = make(logcache.Entries)
		if _, err := os.Stat(guildFilename); err == nil {
			if err := logcache.NewEntries(guildFilename, &p.gCache); err != nil {
				return fmt.Errorf("error reconstructing guild state: %v", err)
			}
		}
		err := p.pullGuild(c.GuildID)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(filename); err == nil {
		last, err = logutil.LastMessageID(filename)
		if err != nil {
			return fmt.Errorf("error getting last message id: %v", err)
		}
	}

	err := p.pullChannel(&c, last)
	if err != nil {
		return err
	}

	return nil
}

func (p *Puller) pullGuild(id string) error {
	p.gDeleted = p.gCache.IDs()
	filename := fmt.Sprintf("channels/%s/guild.tsv", id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening the log file: %v", err)
	}
	defer f.Close()

	guild, err := p.d.Guild(id)
	if err != nil {
		return fmt.Errorf("error getting guild info: %v", err)
	} else {
		if guild.Icon != "" {
			err := cdndl.Icon(id, guild.Icon)
			if err != nil {
				return fmt.Errorf("error downloading the guild icon: %v", err)
			}
		}

		if guild.Splash != "" {
			err := cdndl.Splash(id, guild.Splash)
			if err != nil {
				return fmt.Errorf("error downloading the guild splash: %v", err)
			}
		}

		gEntry := logentry.Make("history", "add", guild)
		p.gCache.WriteNew(f, gEntry)
		delete(p.gDeleted[logentry.Type(guild)], guild.ID)
	}

	for _, c := range guild.Channels {
		cEntry := logentry.Make("history", "add", c)
		p.gCache.WriteNew(f, cEntry)
		delete(p.gDeleted[logentry.Type(c)], c.ID)

		for _, o := range c.PermissionOverwrites {
			oEntry := logentry.Make("history", "add", o)
			p.gCache.WriteNew(f, oEntry)
			delete(p.gDeleted[logentry.Type(o)], o.ID)
		}
	}

	for _, r := range guild.Roles {
		rEntry := logentry.Make("history", "add", r)
		p.gCache.WriteNew(f, rEntry)
		delete(p.gDeleted[logentry.Type(r)], r.ID)
	}

	for _, e := range guild.Emojis {
		err := cdndl.Emoji(e.ID)
		if err != nil {
			return fmt.Errorf("error downloading emoji %s: %v", e.ID, err)
		}
		eEntry := logentry.Make("history", "add", e)
		p.gCache.WriteNew(f, eEntry)
		delete(p.gDeleted[logentry.Type(e)], e.ID)
	}

	after := "0"
	for {
		members, err := p.d.GuildMembers(id, after, 1000)
		if err != nil {
			return fmt.Errorf("error getting members from %s: %v", after, err)
		}

		if len(members) == 0 {
			break
		}

		for _, m := range members {
			after = m.User.ID

			if m.User.Avatar != "" {
				err := cdndl.Avatar(m.User)
				log.Printf("[%s] downloading avatar for user %s (%s)", id, m.User.ID, m.User.Username)
				if err != nil {
					return fmt.Errorf("error downloading avatar for user %s: %v", m.User.ID, err)
				}
			}

			uEntry := logentry.Make("history", "add", m)
			p.gCache.WriteNew(f, uEntry)
			delete(p.gDeleted[logentry.Type(m)], m.User.ID)
		}

		log.Printf("[%s] downloaded %d members, last id %s with name %s", id, len(members), after, members[len(members)-1].User.Username)
	}

	for etype, ids := range p.gDeleted {
		for id := range ids {
			entry := p.gCache[etype][id]
			entry[logentry.HTime] = logentry.Timestamp()
			entry[logentry.HOp] = "del"
			tsv.Write(f, entry)
		}
	}

	return nil
}

func (p *Puller) pullChannel(c *discordgo.Channel, after string) error {
	filename := fmt.Sprintf("channels/%s/%s.tsv", c.GuildID, c.ID)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening the log file: %v", err)
	}
	defer f.Close()

	for {
		msgs, err := p.d.ChannelMessages(c.ID, 100, "", after, "")
		if err != nil {
			return fmt.Errorf("error getting messages from %s: %v", after, err)
		}

		if len(msgs) == 0 {
			break
		}

		after = msgs[0].ID

		// messages are retrieved in descending order
		for i := len(msgs) - 1; i >= 0; i-- {
			tsv.Write(f, logentry.Make("history", "add", msgs[i]))

			for _, e := range msgs[i].Embeds {
				tsv.Write(f, logentry.Make("history", "add", &logentry.Embed{*e, msgs[i].ID}))
			}

			for _, a := range msgs[i].Attachments {
				log.Printf("[%s/%s] downloading attachment %s", c.GuildID, c.ID, a.ID)
				err := cdndl.Attachment(a.URL)
				if err != nil {
					return fmt.Errorf("error downloading attachment %s: %v", a.ID, err)
				}
				tsv.Write(f, logentry.Make("history", "add", &logentry.Attachment{*a, msgs[i].ID}))
			}

			for _, r := range msgs[i].Reactions {
				users, err := p.d.MessageReactions(c.ID, msgs[i].ID, r.Emoji.APIName(), 100)
				if err != nil {
					return fmt.Errorf("error getting users for reaction %s to %s: %v", r.Emoji.APIName(), msgs[i].ID, err)
				}

				for _, u := range users {
					reaction := &logentry.Reaction{
						discordgo.MessageReaction{u.ID, msgs[i].ID, *r.Emoji, c.ID},
						1,
					}

					tsv.Write(f, logentry.Make("history", "add", reaction))
				}

				if r.Count > 100 {
					reaction := &logentry.Reaction{
						discordgo.MessageReaction{"", msgs[i].ID, *r.Emoji, c.ID},
						r.Count - 100,
					}
					tsv.Write(f, logentry.Make("history", "add", reaction))
				}
			}
		}

		log.Printf("[%s/%s] downloaded %d messages, last id %s with content %s", c.GuildID, c.ID, len(msgs), msgs[0].ID, msgs[0].Content)
	}

	return nil
}
