package logpull

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/tsudoko/pullcord/cdndl"
)

const (
	_ = iota
	cdnIcon
	cdnSplash
	cdnAvatar
	cdnChannelIcon
)

// TODO: move somewhere else
func handleDLError(err error) error {
	if cerr, ok := err.(cdndl.ErrNotOk); ok && cerr.StatusCode == 404 {
		log.Printf("warning: skipping %s (404)", cerr.URL)
		return nil
	} else {
		// TODO: retry?
		return err
	}
}

func (p *Puller) cdnDL(v interface{}, subtype int) error {
	if p.lightMode {
		return nil
	}

	var err error

	switch v.(type) {
	case *discordgo.MessageAttachment:
		return cdndl.Attachment(v.(*discordgo.MessageAttachment).URL)
	case *discordgo.Guild:
		g := v.(*discordgo.Guild)
		switch subtype {
		case cdnIcon:
			err = cdndl.Icon(g.ID, g.Icon)
		case cdnSplash:
			err = cdndl.Splash(g.ID, g.Splash)
		default:
			panic("unsupported subtype")
		}
	case *discordgo.User:
		if subtype == cdnAvatar {
			err = cdndl.Avatar(v.(*discordgo.User))
		} else {
			panic("unsupported subtype")
		}
	case *discordgo.Channel:
		c := v.(*discordgo.Channel)
		if subtype == cdnChannelIcon {
			err = cdndl.ChannelIcon(c.ID, c.Icon)
		} else {
			panic("unsupported subtype")
		}
	case *discordgo.Emoji:
		e := v.(*discordgo.Emoji)
		err = cdndl.Emoji(e.ID, e.Animated)
	default:
		panic("unsupported type")
	}

	return handleDLError(err)
}
