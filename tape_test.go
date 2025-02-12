package poker

import (
	"io"
	"testing"
)

func TestTape(t *testing.T) {
	file, clean := CreateTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	_, err := tape.Write([]byte("abc"))
	AssertNoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	AssertNoError(t, err)

	newFileContents, err := io.ReadAll(file)
	AssertNoError(t, err)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("\ngot: \t%q\nwant:\t%q", got, want)
	}
}
