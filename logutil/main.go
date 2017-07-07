package logutil

import (
	"bufio"
	"os"

	"github.com/tsudoko/pullcord/logcache"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/tsv"
)

func LastMessageID(fpath string) (id string, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := tsv.Read(scanner)
		if entry[logentry.HOp] == "add" && entry[logentry.HType] == "message" {
			id = entry[logentry.HID]
		}
	}

	err = scanner.Err()
	return
}

func AllIDs(fpath string, ids *logcache.IDs) (err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		e := tsv.Read(scanner)
		if (*ids)[e[logentry.HType]] == nil {
			(*ids)[e[logentry.HType]] = make(map[string]bool)
		}
		(*ids)[e[logentry.HType]][e[logentry.HID]] = true
	}

	err = scanner.Err()
	return
}
