package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func wantedChannels(d *discordgo.Session) []discordgo.Channel {
	channels := make([]discordgo.Channel, 0)
	gid := ""
	for {
		guilds, err := d.UserGuilds(100, "", gid)
		if err != nil {
			log.Fatal("error getting guilds:", err)
		}

		if len(guilds) == 0 {
			break
		}

		for _, g := range guilds {
			gid = g.ID
			if !wantedGuild(gid) {
				continue
			}

			gch, err := d.GuildChannels(gid)
			if err != nil {
				log.Fatalf("error getting channels for %s: %s", gid, err)
				continue
			}

			for _, c := range gch {
				if !wantedChannel(c.ID) || c.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				channels = append(channels, *c)
			}
		}
	}
	return channels
}

func wantedChannel(id string) bool {
	if len(cids) != 0 {
		return cids[id] && !xcids[id]
	} else {
		return !xcids[id]
	}
}

func wantedGuild(id string) bool {
	if len(gids) != 0 {
		return gids[id] && !xgids[id]
	} else {
		return !xgids[id]
	}
}

func makeWanted(idstr string) map[string]bool {
	ids := make(map[string]bool)
	for _, i := range strings.Split(idstr, ",") {
		if i != "" {
			ids[i] = true
		}
	}
	return ids
}
