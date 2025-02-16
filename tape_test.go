package poker_test

import (
	"io"
	"testing"

	poker "webapp/v2"
)

func TestTape(t *testing.T) {
	file, clean := poker.CreateTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{file}

	_, err := tape.Write([]byte("abc"))
	poker.AssertNoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	poker.AssertNoError(t, err)

	newFileContents, err := io.ReadAll(file)
	poker.AssertNoError(t, err)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("\ngot: \t%q\nwant:\t%q", got, want)
	}
}
