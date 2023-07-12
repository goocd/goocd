// Package mempgm implements the fileformats.ProgramReader interface
// using a simple in memory struct.  Other formats can parse into this
// if it's helpful.
package mempgm

import (
	"io"

	"goocd/fileformats"
)

// MemProgramReader implements fileformats.ProgramReader.
// It is designed to be easy to construct an implementation
// of ProgramReader directly using in-memory data.
type MemProgramReader struct {
	MemProgramList []MemProgram
	Index          int
}

func (m *MemProgramReader) NextProgram() (fileformats.Program, error) {
	if m.Index >= len(m.MemProgramList) {
		return nil, io.EOF
	}
	ret := &m.MemProgramList[m.Index]
	m.Index++
	return ret, nil
}

// MemProgram implements fileformats.Program.
type MemProgram struct {
	StartAddress uint64
	ByteSlice    []byte
}

func (mp *MemProgram) StartAddr() uint64 {
	return mp.StartAddress
}

func (mp *MemProgram) Bytes() []byte {
	return mp.ByteSlice
}
