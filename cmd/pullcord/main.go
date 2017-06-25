package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/logpull"
	"github.com/tsudoko/pullcord/logutil"
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
)

func do(d *discordgo.Session, event *discordgo.Ready) {
	channels := wantedChannels(d)

	for _, c := range channels {
		last := "0"

		if err := os.MkdirAll("channels/"+c.GuildID, os.ModeDir|0755); err != nil {
			log.Printf("creating the guild dir for %s failed", c.GuildID)
			continue
		}

		filename := fmt.Sprintf("channels/%s/%s.tsv", c.GuildID, c.ID)

		if _, err := os.Stat(filename); err == nil {
			last, err = logutil.LastMessageID(filename)
			if err != nil {
				log.Printf("[%s] error getting last message id, cancelling: %v", c.ID, err)
				continue
			}
		}

		log.Printf("[%s] last downloaded message id: %s", c.ID, last)
		//go logpull.Channel(d, c.GuildID, c.ID, "0")
		logpull.Channel(d, c.GuildID, c.ID, last)
	}

	os.Exit(0)
}

func main() {
	flag.Parse()

	cids = makeWanted(*cid)
	gids = makeWanted(*gid)
	xcids = makeWanted(*xcid)
	xgids = makeWanted(*xgid)

	d, err := discordgo.New(*username, *password, *token)
	if err != nil {
		log.Fatal("login failed:", err)
	}

	err = d.Open()
	defer d.Close()
	if err != nil {
		log.Fatal("opening the websocket connection failed:", err)
	}

	d.AddHandler(do)

	if *token == "" {
		log.Println("token:", d.Token)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
