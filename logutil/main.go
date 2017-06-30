package logutil

import (
	"bufio"
	"io"
	"os"

	"github.com/tsudoko/pullcord/logformat"
)

const (
	hTime = iota
	hFetchType
	hOp
	hType
	hID
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
	if !Equals((*cache)[e[hType]][e[hID]], e[1:]) {
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
		switch entry[hOp] {
		case "add":
			if (*cache)[entry[hType]] == nil {
				(*cache)[entry[hType]] = make(map[string][]string)
			}
			(*cache)[entry[hType]][entry[hID]] = entry[1:]
		case "del":
			delete((*cache)[entry[hType]], entry[hID])
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
		if entry[hOp] == "add" && entry[hType] == "message" {
			id = entry[hID]
		}
	}

	err = scanner.Err()
	return
}
