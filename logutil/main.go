package logutil

import (
	"bufio"
	"bytes"
	"github.com/tsudoko/pullcord/logcache"
	"github.com/tsudoko/pullcord/logentry"
	"github.com/tsudoko/pullcord/tsv"
	"io"
	"os"
	"sync"
)

func LastMessageID(fpath string) (id string, err error) {
	f, err := os.Open(fpath)
	osstat, _ := os.Stat(fpath)
	if err != nil {
		return
	}

	scanner := NewScanner(f, osstat.Size())

	var line string
	for err == nil {
		line, _, err = scanner.Line()
		entry := tsv.ReadString(line)
		if len(entry) < logentry.HOp || len(entry) < logentry.HType {
			continue
		}
		if entry[logentry.HOp] == "add" && entry[logentry.HType] == "message" {
			id = entry[logentry.HID]
			break
		}

	}
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
		var val sync.Map
		valIface, ok := ids.Load(e[logentry.HType])
		if ok {
			val = valIface.(sync.Map) //(map[string]bool)
		} else {
			val = sync.Map{}
		}
		val.Store(e[logentry.HID], true)
		ids.Store(e[logentry.HType], val)
	}

	err = scanner.Err()
	return
}

type Scanner struct {
	r   io.ReaderAt
	pos int64
	err error
	buf []byte
}

func NewScanner(r io.ReaderAt, pos int64) *Scanner {
	return &Scanner{r: r, pos: pos}
}

func (s *Scanner) readMore() {
	if s.pos == 0 {
		s.err = io.EOF
		return
	}
	var size int64 = 1024
	if size > s.pos {
		size = s.pos
	}
	s.pos -= size
	buf2 := make([]byte, size, size+int64(len(s.buf)))

	// ReadAt attempts to read full buff!
	_, s.err = s.r.ReadAt(buf2, s.pos)
	if s.err == nil {
		s.buf = append(buf2, s.buf...)
	}
}

func (s *Scanner) Line() (line string, start int64, err error) {
	if s.err != nil {
		return "", 0, s.err
	}
	for {
		lineStart := bytes.LastIndexByte(s.buf, '\n')
		if lineStart >= 0 {
			// We have a complete line:
			var line string
			line, s.buf = string(dropCR(s.buf[lineStart+1:])), s.buf[:lineStart]
			return line, s.pos + int64(lineStart+1), nil
		}
		// Need more data:
		s.readMore()
		if s.err != nil {
			if s.err == io.EOF {
				if len(s.buf) > 0 {
					return string(dropCR(s.buf)), 0, nil
				}
			}
			return "", 0, s.err
		}
	}
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
