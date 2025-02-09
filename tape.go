package main

import (
	"io"
)

type tape struct {
	file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	_, err = t.file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return t.file.Write(p)
}
