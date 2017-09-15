package logcache

import (
	"bufio"
	"io"
	"os"

	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/tsv"
)

type Entries map[string]map[string][]string
type IDs map[string]map[string]bool

func (ec *Entries) IDs() IDs {
	ic := make(IDs)

	for etype, ids := range *ec {
		ic[etype] = make(map[string]bool)
		for id := range ids {
			ic[etype][id] = true
		}
	}

	return ic
}

func NewEntries(fpath string, cache *Entries) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		e := tsv.Read(scanner)

		switch e[logentry.HOp] {
		case "add":
			if (*cache)[e[logentry.HType]] == nil {
				(*cache)[e[logentry.HType]] = make(map[string][]string)
			}
			(*cache)[e[logentry.HType]][e[logentry.HID]] = e
		case "del":
			delete((*cache)[e[logentry.HType]], e[logentry.HID])
		}
	}

	return scanner.Err()
}

func (cache *Entries) WriteNew(w io.Writer, e []string) {
	cacheEntry := (*cache)[e[logentry.HType]][e[logentry.HID]]

	if len(cacheEntry) < logentry.HTime+1 || len(e) < logentry.HTime+1 ||
		!entryEquals(cacheEntry[logentry.HTime+1:], e[logentry.HTime+1:]) {
		tsv.Write(w, e)
	}
}

func entryEquals(a, b []string) bool {
	larger := b
	smaller := a

	if len(a) > len(b) {
		larger = a
		smaller = b
	}

	for i := range smaller {
		if a[i] != b[i] {
			return false
		}
	}

	// For backwards compatibility with older entry formats
	for _, field := range larger[len(smaller):] {
		if field != "" {
			return false
		}
	}

	return true
}
