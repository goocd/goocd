package fileformats

// ProgramReader is anything that can read one or more Programs,
// typically from a file.
type ProgramReader interface {
	// will return the next program or error, io.EOF if no more programs
	NextProgram() (Program, error)
}

// Program is a starting memory address and data bytes.
// A program must be row/page/block/whatever aligned according
// to the target chip in order for it to load properly (typically
// the size of a "block erase" is the largest alignment requirement,
// but this value can vary greatly from one MCU to the next).
// For this reason, the regions of an ELF file are generally
// concatenated into one program in order to fulfill this
// alignment requirement, and any other formats that support
// multiple regions should be careful to not return small chunks
// of code at random offsets just because the underlying file
// stores them like that.
type Program interface {
	StartAddr() uint64 // start address, must be aligned per target chip
	Bytes() []byte     // raw program bytes, if not aligned to write size maybe zero-padded during loading
}
