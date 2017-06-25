package logutil

import (
	"bufio"
	"os"

	"github.com/tsudoko/pullcord/logformat"
)

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
