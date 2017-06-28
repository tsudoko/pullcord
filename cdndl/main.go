package cdndl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
)

// Returned when the request gets a non-200 HTTP response.
type ErrNotOk struct {
	error
	StatusCode int
}

// discordgo uses EndpointAPI, which includes an extra "/api" path element
var EndpointCDNEmojis = discordgo.EndpointCDN + "emojis/"

var avatarFormats = []string{"gif", "png", "jpg"}

func absDL(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	fPath := u.Path[1:]
	if _, err := os.Stat(fPath); err == nil {
		return nil
	}

	r, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return ErrNotOk{errors.New("non-200 status code"), r.StatusCode}
	}

	if err = saveFile(r.Body, fPath); err != nil {
		return err
	}

	return nil
}

func saveFile(r io.Reader, fPath string) error {
	if err := os.MkdirAll(path.Dir(fPath), os.ModeDir|0755); err != nil {
		return err
	}

	f, err := os.Create(fPath + ".part")
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	if err = os.Rename(fPath+".part", fPath); err != nil {
		return err
	}

	return nil
}

func Avatar(uid, hash string) error {
	for _, ext := range avatarFormats {
		url := fmt.Sprintf("%s%s/%s.%s", discordgo.EndpointCDNAvatars, uid, hash, ext)
		err := absDL(url)
		if notOk, ok := err.(ErrNotOk); ok && notOk.StatusCode == 415 {
			continue
		} else if err == nil {
			break
		} else {
			return err
		}
	}

	return nil
}

func Emoji(id string) error {
	return absDL(EndpointCDNEmojis + id + ".png")
}

func Icon(gid, hash string) error {
	return absDL(discordgo.EndpointGuildIcon(gid, hash))
}

func Splash(gid, hash string) error {
	return absDL(discordgo.EndpointGuildSplash(gid, hash) + "?size=2048")
}

func Attachment(url string) error {
	return absDL(url)
}
