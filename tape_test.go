package main

import (
	"io"
	"testing"
)

func TestTape(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	_, err := tape.Write([]byte("abc"))
	if err != nil {
		t.Fatalf("tape error during write: %v", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("tape error during seek: %v", err)
	}
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("\ngot: \t%q\nwant:\t%q", got, want)
	}
}
