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
)

var avatarFormats = []string{
	"%s/%s.gif",
	"%s/%s.png",
	"%s/%s.jpg",
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

	if err = os.Rename(fPath + ".part", fPath); err != nil {
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

		if err := saveFile(r.Body, path); err != nil {
			log.Println("failed to save", path+":", err)
		} else {
			break
		}
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

	if err = saveFile(r.Body, u.Path); err != nil {
		log.Println("failed to save", u.Path+":", err)
	}
}
