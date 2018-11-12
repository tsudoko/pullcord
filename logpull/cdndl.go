package logpull

import (
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

func (p *Puller) cdnDL(v interface{}, subtype int) error {
	if p.lightMode {
		return nil
	}

	switch v.(type) {
	case *discordgo.MessageAttachment:
		return cdndl.Attachment(v.(*discordgo.MessageAttachment).URL)
	case *discordgo.Guild:
		g := v.(*discordgo.Guild)
		switch subtype {
		case cdnIcon:
			return cdndl.Icon(g.ID, g.Icon)
		case cdnSplash:
			return cdndl.Splash(g.ID, g.Splash)
		default:
			panic("unsupported subtype")
		}
	case *discordgo.User:
		if subtype == cdnAvatar {
			return cdndl.Avatar(v.(*discordgo.User))
		} else {
			panic("unsupported subtype")
		}
	case *discordgo.Channel:
		c := v.(*discordgo.Channel)
		if subtype == cdnChannelIcon {
			return cdndl.ChannelIcon(c.ID, c.Icon)
		} else {
			panic("unsupported subtype")
		}
	case *discordgo.Emoji:
		e := v.(*discordgo.Emoji)
		return cdndl.Emoji(e.ID, e.Animated)
	default:
		panic("unsupported type")
	}
}
