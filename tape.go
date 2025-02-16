package poker

import (
	"io"
	"os"
)

type Tape struct {
	File *os.File
}

func (t *Tape) Write(p []byte) (n int, err error) {
	err = t.File.Truncate(0)
	if err != nil {
		return 0, err
	}
	_, err = t.File.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return t.File.Write(p)
}
