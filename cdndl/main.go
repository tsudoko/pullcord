package cdndl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const maxSize = "4096"
const longNamePlaceholder = "_long_filename"

// Returned when the request gets a non-200 HTTP response.
type ErrNotOk struct {
	error
	URL        string
	StatusCode int
}

// discordgo uses EndpointAPI, which includes an extra "/api" path element
var EndpointCDNEmojis = discordgo.EndpointCDN + "emojis/"

func NewErrNotOk(URL string, code int) error {
	return ErrNotOk{fmt.Errorf("non-200 status code: %d", code), URL, code}
}

func absDL(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	fPath := u.Path[1:]
	if _, err := os.Stat(fPath); err == nil {
		return nil
	}

	return absDLTo(URL, fPath)
}

func absDLTo(URL string, target string) error {
	log.Printf("downloading %s", URL)

	r, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return NewErrNotOk(URL, r.StatusCode)
	}

	err = saveFile(r.Body, target)
	if cerr, ok := err.(*os.PathError); ok {
		if errno, ok := cerr.Err.(syscall.Errno); ok && errno == syscall.ENAMETOOLONG {
			dir, name := filepath.Split(target)
			newname := longNamePlaceholder + filepath.Ext(name)
			log.Printf("warning: %s: file name too long, renaming to %s", name, newname)
			return absDLTo(URL, filepath.Join(dir, newname))

		}
	}

	return err
}

func saveFile(r io.Reader, fPath string) error {
	if err := os.MkdirAll(path.Dir(fPath), os.ModeDir|0755); err != nil {
		return err
	}

	tempPath := fPath + ".part"

	f, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	f.Close()

	if err = os.Rename(tempPath, fPath); err != nil {
		return err
	}

	return nil
}

func Avatar(u *discordgo.User) error {
	return absDL(u.AvatarURL(maxSize))
}

func Emoji(id string, animated bool) error {
	var ext string
	if animated {
		ext = "gif"
	} else {
		ext = "png"
	}
	err := absDL(fmt.Sprintf("%s%s.%s?size=%s", EndpointCDNEmojis, id, ext, maxSize))
	if cerr, ok := err.(ErrNotOk); ok && cerr.StatusCode == 415 && ext == "gif" {
		log.Printf("warning: animated version of emoji %s doesn't exist, trying png", id)
		ext = "png"
		err = absDL(fmt.Sprintf("%s%s.%s?size=%s", EndpointCDNEmojis, id, ext, maxSize))
	}
	return err
}

func Icon(gid, hash string) error {
	return absDL(discordgo.EndpointGuildIcon(gid, hash) + "?size=" + maxSize)
}

func ChannelIcon(cid, hash string) error {
	return absDL(discordgo.EndpointGroupIcon(cid, hash) + "?size=" + maxSize)
}

func Splash(gid, hash string) error {
	return absDL(discordgo.EndpointGuildSplash(gid, hash) + "?size=" + maxSize)
}

func Attachment(url string) error {
	return absDL(url)
}
