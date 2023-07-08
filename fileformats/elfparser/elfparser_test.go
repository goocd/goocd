package elfparser

import (
	"errors"
	"io"
	"testing"

	"github.com/goocd/goocd/fileformats"
)

func TestELFParser(t *testing.T) {

	// b, err := os.ReadFile("../testdata/example1.elf")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	addr, b, err := ExtractROM("../testdata/example1.elf")
	if err != nil {
		t.Fatal(err)
	}

	_, _ = addr, b
	t.Logf("addr=%X, len(b)=%d", addr, len(b))

	pr, err := ParseFromPath("../testdata/example1.elf")
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
