package hexparser

import (
	"errors"
	"io"
	"testing"

	"github.com/goocd/goocd/fileformats"
)

func TestParseFromPath(t *testing.T) {

	pr, err := ParseFromPath("../testdata/example2.hex")
	if err != nil {
		t.Fatal(err)
	}

	var p fileformats.Program
	for {
		p, err = pr.NextProgram()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		t.Logf("PGM, addr=0x%X, len=%d", p.StartAddr(), len(p.Bytes()))
	}

}
