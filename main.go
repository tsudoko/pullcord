package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	username = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	token    = flag.String("t", "", "access token")
)

func main() {
	flag.Parse()
	d, err := discordgo.New(*username, *password, *token)
	if err != nil {
		log.Fatal("login failed:", err)
	}

	if *token == "" {
		log.Println("token:", d.Token)
	}

	gid := ""
	for {
		guilds, err := d.UserGuilds(100, "", gid)
		if err != nil {
			log.Fatal("error getting guilds:", err)
		}

		if len(guilds) == 0 {
			break
		}

		for _, i := range guilds {
			fmt.Println(i.ID, i.Name)
			gid = i.ID
		}
	}
}
