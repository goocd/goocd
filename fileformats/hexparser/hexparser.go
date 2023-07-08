package hexparser

import (
	"errors"
	"os"

	"github.com/goocd/goocd/fileformats"
	"github.com/goocd/goocd/fileformats/mempgm"
	"github.com/marcinbor85/gohex"
)

func ParseFromPath(filePath string) (fileformats.ProgramReader, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(f)
	if err != nil {
		return nil, err
	}

	ret := &mempgm.MemProgramReader{}

	// NOTE: We don't do any alignment checks here -
	// if this turns out to be a problem and we need
	// to merge different chunks (like is done in
	// elfparser) the we can burn that bridge when
	// we cross it.  But so far, this hasn't been
	// an issue.

	for _, s := range mem.GetDataSegments() {
		ret.MemProgramList = append(ret.MemProgramList, mempgm.MemProgram{
			StartAddress: uint64(s.Address),
			ByteSlice:    s.Data,
		})
	}

	// zero data we do not consider valid
	if len(ret.MemProgramList) == 0 {
		return nil, errors.New("parse succeeded but no programs found in hex file")
	}

	return ret, nil
}
