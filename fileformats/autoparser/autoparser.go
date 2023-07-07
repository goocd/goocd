// Package autoparser examines a file and selects one of the
// other parsers.
package autoparser

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/goocd/goocd/fileformats"
	"github.com/goocd/goocd/fileformats/elfparser"
	"github.com/goocd/goocd/fileformats/hexparser"
	"github.com/goocd/goocd/fileformats/mempgm"
)

func ParseFromPath(filePath string, defaultStartAddr uint64) (fileformats.ProgramReader, error) {

	// TODO: we could do some "magic" file introspection later
	// if it seems useful (elf and hex files can be detected
	// by inspecting the beginning of the file)

	ext := strings.ToLower(path.Ext(filePath))
	switch ext {
	case ".elf":
		return elfparser.ParseFromPath(filePath)
	case ".hex":
		return hexparser.ParseFromPath(filePath)
	case ".bin":
		// didn't bother making an abstraction for bin files
		b, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return &mempgm.MemProgramReader{
			MemProgramList: []mempgm.MemProgram{
				{
					StartAddress: defaultStartAddr,
					ByteSlice:    b,
				},
			},
		}, nil
	case "":
		return nil, errors.New("no file extension found, cannot determine file type")
	}

	return nil, fmt.Errorf("unrecogninzed file extension %q", ext)
}
