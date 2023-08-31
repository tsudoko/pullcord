// Package tsv reads and writes tab-separated values (TSV) files.
package tsv

import (
	"bufio"
	"io"
	"strings"
)

var subs = [][]string{
	[]string{"\\", "\\\\"},
	[]string{"\n", "\\n"},
	[]string{"\t", "\\t"},
}

func Read(s *bufio.Scanner) []string {
	record := strings.Split(s.Text(), "\t")

	for i := range record {
		for j := len(subs) - 1; j >= 0; j-- {
			record[i] = strings.Replace(record[i], subs[j][1], subs[j][0], -1)
		}
	}

	return record
}

func ReadString(s string) []string {
	record := strings.Split(s, "\t")

	for i := range record {
		for j := len(subs) - 1; j >= 0; j-- {
			record[i] = strings.Replace(record[i], subs[j][1], subs[j][0], -1)
		}
	}

	return record
}

func Write(w io.Writer, record []string) error {
	for i := range record {
		for j := 0; j < len(subs); j++ {
			record[i] = strings.Replace(record[i], subs[j][0], subs[j][1], -1)
		}
	}

	_, err := w.Write([]byte(strings.Join(record, "\t") + "\n"))
	return err
}
