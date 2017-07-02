package logcache

import (
	"bufio"
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
