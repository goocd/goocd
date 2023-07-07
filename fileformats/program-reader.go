package fileformats

// ProgramReader is anything that can read one or more Programs,
// typically from a file.
type ProgramReader interface {
	// will return the next program or error, io.EOF if no more programs
	NextProgram() (Program, error)
}

// Program is a starting memory address and data bytes.
type Program interface {
	StartAddr() uint64
	Bytes() []byte
}
