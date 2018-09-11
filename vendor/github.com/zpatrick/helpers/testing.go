package helpers

import (
	"io/ioutil"
	"os"
	"testing"
)

func TempFileWithContent(t *testing.T, content []byte) (*os.File, func()) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}

	remove := func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return f, remove
}
