package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/logpull"
)

var (
	username = flag.String("user", "", "email address")
	password = flag.String("pass", "", "password")
	token    = flag.String("t", "", "access token")

	cid  = flag.String("c", "", "comma-separated channel IDs to include")
	gid  = flag.String("s", "", "comma-separated server IDs to include")
	xcid = flag.String("C", "", "comma-separated channel IDs to exclude")
	xgid = flag.String("S", "", "comma-separated server IDs to exclude")

	cids, gids, xcids, xgids map[string]bool

	historyMode = flag.Bool("history", false, "download the whole history")

	dlDM = flag.Bool("dm", false, "download DMs")
	// not fully implemented yet, we currently don't check if all emoji/attachments/etc with log entries have been downloaded
	//lightMode = flag.Bool("light", false, "skip downloading non-textual data such as attachments or emoji")
)

func do(d *discordgo.Session, _ *discordgo.Ready) {
	pullers := make(map[string]*logpull.Puller)
	channels := wantedChannels(d)

	if *historyMode {
		for _, c := range channels {
			if pullers[c.GuildID] == nil {
				p, err := logpull.NewPuller(d, c.GuildID)
				if err != nil {
					log.Fatalf("[%s] %v", c.GuildID, err)
				}

				pullers[c.GuildID] = p
				err = p.PullGuild(c.GuildID)
				if err != nil {
					log.Fatalf("[%s] %v", c.GuildID, err)
				}
			}

			err := pullers[c.GuildID].PullChannel(&c)
			if err != nil {
				log.Fatalf("[%s/%s] %v", c.GuildID, c.ID, err)
			}
		}

		if *dlDM {
			p, err := logpull.NewPuller(d, "@me")
			if err != nil {
				log.Fatalf("[@me] %v", err)
			}

			pullers["@me"] = p
			err = p.PullDMGuild()
			if err != nil {
				log.Fatalf("[@me] %v", err)
			}

			uc, err := d.UserChannels()
			if err != nil {
				log.Fatalf("[@me] error getting user channels: %v", err)
			}

			for _, c := range uc {
				c.GuildID = "@me"
				err := pullers[c.GuildID].PullChannel(c)
				if err != nil {
					log.Fatalf("[%s/%s] %v", c.GuildID, c.ID, err)
				}
			}
		}

		for id, p := range pullers {
			if err := p.Close(); err != nil {
				log.Printf("[%s] error closing log file: %v", id, err)
			}
		}
	}

	os.Exit(0)
}

func main() {
	flag.Parse()

	cids = makeWanted(*cid)
	gids = makeWanted(*gid)
	xcids = makeWanted(*xcid)
	xgids = makeWanted(*xgid)

	if !*historyMode {
		log.Fatal("no modes specified, nothing to do")
	}

	d, err := discordgo.New(*username, *password, *token)
	if err != nil {
		log.Fatal("login failed:", err)
	}

	d.AddHandler(do)

	err = d.Open()
	defer d.Close()
	if err != nil {
		log.Fatal("opening the websocket connection failed:", err)
	}

	if *token == "" {
		log.Println("token:", d.Token)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
