// Package fileformats exists as a general owner for
// raw file formatting utility funcs and interfaces.
package fileformats

import (
	"encoding/binary"
	"fmt"
)

// Todo: Maybe provide auto alignment rather than erroring?

// ConvertByteSliceUint32Slice is just a utility function to put elf parser in a usable format. All Roms that use this feature must be 32bit aligned or it will fail.
func ConvertByteSliceUint32Slice(rom []byte) ([]uint32, error) {
	if len(rom)%4 > 0 {
		return nil, fmt.Errorf("error: ConvertByteSliceUint32Slice() Misaligned Rom")
	}

	rom32 := make([]uint32, len(rom)/4)

	for i := range rom32 {
		rom32[i] = binary.LittleEndian.Uint32(rom[i*4:])
	}

	return rom32, nil
}
