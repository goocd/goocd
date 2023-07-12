// Package targets contains the "wiring" for each supported target.
package targets

import "log"

// TargetMap is where each target registers itself.
var TargetMap = make(map[string]*Target)

// Args is the options that come in from the command line
// and tell a target what to do.
type Args struct {
	Load  string // file path to load (elfparser, hex, bin)
	Reset bool
	Stats bool
	// -readmemu32=0xF0000000,5
	// -readmemu32=0xF0000000 (count=1 implied)
	ReadMemU32Addr   uint64
	ReadMemU32Count  int
	WriteMemU32Addr  uint64
	WriteMemU32Value uint64
	WriteMemU32Count int
}

// Target is anything that can be "Run" as a target.
type Target struct {
	Name                string
	Description         string
	SupportsReadMemU32  bool
	SupportsWriteMemU32 bool
	SupportsReset       bool
	SupportsLoad        bool
	Run                 func(args *Args) error
}

// TargetFunc let's us describe a target as a simple function.
type TargetFunc func(args *Args) error

// Run implements Target.
func (f TargetFunc) Run(args *Args) error {
	return f(args)
}

func addTarget(tar *Target) {
	TargetMap[tar.Name] = tar
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
}
