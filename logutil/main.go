package logutil

import (
	"bufio"
	"io"
	"os"

	"github.com/tsudoko/pullcord/logcache"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/tsv"
)

func WriteNew(w io.Writer, e []string, cache *logcache.Entries) {
	cacheEntry := (*cache)[e[logentry.HType]][e[logentry.HID]]

	if len(cacheEntry) < logentry.HTime+1 || len(e) < logentry.HTime+1 ||
		!Equals(cacheEntry[logentry.HTime+1:], e[logentry.HTime+1:]) {
		tsv.Write(w, e)
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
