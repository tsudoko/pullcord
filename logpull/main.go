// Package logpull contains functions related to downloading historical data.
package logpull

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"

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

	gLog *os.File

	// per-guild caches
	gCache   logcache.Entries // for tracking changes between different pulls
	gEver    logcache.IDs     // for determining if there's a need to add an entry for an outside entity, i.e. a user who left
	gDeleted logcache.IDs     // for tracking deletions between different pulls, gCache could be used for that as well
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
		p.gEver = make(logcache.IDs)
		if _, err := os.Stat(guildFilename); err == nil {
			if err := logcache.NewEntries(guildFilename, &p.gCache); err != nil {
				return fmt.Errorf("error reconstructing guild state: %v", err)
			}

			if err := logutil.AllIDs(guildFilename, &p.gEver); err != nil {
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

func (p *Puller) Close() error {
	return p.gLog.Close()
}

func (p *Puller) openLog(id string) error {
	if p.gLog != nil {
		return nil
	}

	filename := fmt.Sprintf("channels/%s/guild.tsv", id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	p.gLog = f
	return err
}

func (p *Puller) pullGuild(id string) error {
	p.gDeleted = p.gCache.IDs()
	err := p.openLog(id)
	if err != nil {
		return fmt.Errorf("error opening the log file: %v", err)
	}

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

		p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", guild))
		delete(p.gDeleted[logentry.Type(guild)], guild.ID)
	}

	for _, c := range guild.Channels {
		p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", c))
		delete(p.gDeleted[logentry.Type(c)], c.ID)

		for _, o := range c.PermissionOverwrites {
			p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", o))
			delete(p.gDeleted[logentry.Type(o)], o.ID)
		}
	}

	for _, r := range guild.Roles {
		p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", r))
		delete(p.gDeleted[logentry.Type(r)], r.ID)
	}

	for _, e := range guild.Emojis {
		err := cdndl.Emoji(e.ID)
		if err != nil {
			return fmt.Errorf("error downloading emoji %s: %v", e.ID, err)
		}
		p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", e))
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

			if p.gEver["member"] == nil {
				p.gEver["member"] = make(map[string]bool)
			}
			p.gEver["member"][m.User.ID] = true

			p.gCache.WriteNew(p.gLog, logentry.Make("history", "add", m))
			delete(p.gDeleted[logentry.Type(m)], m.User.ID)
		}

		log.Printf("[%s] downloaded %d members, last id %s with name %s", id, len(members), after, members[len(members)-1].User.Username)
	}

	p.gLog.Sync()

	for etype, ids := range p.gDeleted {
		for id := range ids {
			entry := p.gCache[etype][id]
			entry[logentry.HTime] = logentry.Timestamp()
			entry[logentry.HOp] = "del"
			tsv.Write(p.gLog, entry)
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
			log.Println(p.gEver)

			if p.gEver["member"] == nil {
				p.gEver["member"] = make(map[string]bool)
			}

			if !p.gEver["member"][msgs[i].Author.ID] {
				member := &discordgo.Member{User: msgs[i].Author}
				p.gCache.WriteNew(p.gLog, logentry.Make("history", "del", member))
				p.gEver["member"][msgs[i].Author.ID] = true
			}

			for _, u := range msgs[i].Mentions {
				if !p.gEver["member"][u.ID] {
					member := &discordgo.Member{User: u}
					p.gCache.WriteNew(p.gLog, logentry.Make("history", "del", member))
					p.gEver["member"][u.ID] = true
				}
			}

			for _, match := range regexp.MustCompile("<:[^:]+:([0-9]+)>").FindAllStringSubmatch(msgs[i].Content, -1) {
				err := cdndl.Emoji(match[1])
				if err != nil {
					return fmt.Errorf("error downloading outside emoji %s: %v", match[1], err)
				}
			}

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
				if r.Emoji.ID != "" {
					err := cdndl.Emoji(r.Emoji.ID)
					if err != nil {
						return fmt.Errorf("error downloading outside emoji %s: %v", r.Emoji.ID, err)
					}
				}

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

	p.gLog.Sync()
	return nil
}
