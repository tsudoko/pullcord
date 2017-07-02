package logutil

import (
	"bufio"
	"os"

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
