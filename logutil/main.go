package logutil

import (
	"bufio"
	"os"

	"github.com/tsudoko/pullcord/logformat"
)

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

func GuildCache(fpath string, cache *map[string]map[string][]string) (err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := logformat.Read(scanner)
		switch entry[0] {
		case "add":
			if (*cache)[entry[1]] == nil {
				(*cache)[entry[1]] = make(map[string][]string)
			}
			(*cache)[entry[1]][entry[2]] = entry
		case "del":
			delete((*cache)[entry[1]], entry[2])
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
		if entry[0] == "add" && entry[1] == "message" {
			id = entry[2]
		}
	}

	err = scanner.Err()

	return
}
