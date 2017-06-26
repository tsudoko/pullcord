package cdndl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	cdnUrl     = "https://cdn.discordapp.com"
	avatarPath = "/avatars"
	emojiPath  = "/emojis"
	iconPath   = "/icons"
)

var avatarFormats = []string{
	"%s/%s.gif",
	"%s/%s.png",
	"%s/%s.jpg",
}

func saveFile(r io.Reader, fPath string) error {
	if _, err := os.Stat(fPath); err == nil {
		// TODO: avoid making the HTTP request in the first place
		return nil
	}

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

func Avatar(uid, hash string) {
	url := cdnUrl + avatarPath

	for _, format := range avatarFormats {
		path := fmt.Sprintf(format, uid, hash)

		r, err := http.Get(url + "/" + path)
		if err != nil {
			log.Println("failed to get", url+"/"+path+":", err)
		}
		defer r.Body.Close()

		if r.StatusCode != 200 {
			continue
		}

		if err := saveFile(r.Body, "avatars/"+path); err != nil {
			log.Println("failed to save", path+":", err)
		} else {
			break
		}
	}
}

func Emoji(id string) {
	url := cdnUrl + emojiPath
	path := fmt.Sprintf("%s.png", id)

	r, err := http.Get(url + "/" + path)
	if err != nil {
		log.Println("failed to get", url+"/"+path+":", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Println("non-200 status code for", url+":", err)
	}

	if err = saveFile(r.Body, "emojis/"+path); err != nil {
		log.Println("failed to save", path+":", err)
	}
}

func Icon(gid, hash string) {
	url := cdnUrl + iconPath
	path := fmt.Sprintf("%s/%s.png", gid, hash)

	r, err := http.Get(url + "/" + path)
	if err != nil {
		log.Println("failed to get", url+"/"+path+":", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Println("non-200 status code for", url+":", err)
	}

	if err = saveFile(r.Body, "icons/"+path); err != nil {
		log.Println("failed to save", path+":", err)
	}
}

func Attachment(URL string) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Println("url parsing error:", err)
	}

	r, err := http.Get(URL)
	if err != nil {
		log.Println("failed to get", URL+":", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Println("non-200 status code for", URL+":", err)
	}

	if err = saveFile(r.Body, u.Path[1:]); err != nil {
		log.Println("failed to save", u.Path+":", err)
	}
}
