package logutil

import (
	"bufio"
	"io"
	"os"

	"github.com/tsudoko/pullcord/logformat"
)

const (
	HTime = iota
	HFetchType
	HOp
	HType
	HID
)

type EntryCache map[string]map[string][]string
type IDCache map[string]map[string]bool

func (ec *EntryCache) IDCache() IDCache {
	ic := make(IDCache)

	for etype, ids := range *ec {
		ic[etype] = make(map[string]bool)
		for id := range ids {
			ic[etype][id] = true
		}
	}

	return ic
}

func WriteNew(w io.Writer, e []string, cache *EntryCache) {
	if !Equals((*cache)[e[HType]][e[HID]], e[1:]) {
		logformat.Write(w, e)
	}
}

func Equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func GuildCache(fpath string, cache *EntryCache) (err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := logformat.Read(scanner)
		switch entry[HOp] {
		case "add":
			if (*cache)[entry[HType]] == nil {
				(*cache)[entry[HType]] = make(map[string][]string)
			}
			(*cache)[entry[HType]][entry[HID]] = entry[1:]
		case "del":
			delete((*cache)[entry[HType]], entry[HID])
		}
	}

	err = scanner.Err()
	return
}

func LastMessageID(fpath string) (id string, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := logformat.Read(scanner)
		if entry[HOp] == "add" && entry[HType] == "message" {
			id = entry[HID]
		}
	}

	err = scanner.Err()
	return
}
