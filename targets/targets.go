// Package targets contains the "wiring" for each supported target.
package targets

// TargetMap is where each target registers itself.
var TargetMap = make(map[string]Target)

// Args is the options that come in from the command line
// and tell a target what to do.
type Args struct {
	Load string // file path to load (elf, hex, bin)
}

// Target is anything that can be "Run" as a target.
type Target interface {
	Run(args *Args) error
}

// TargetFunc let's us describe a target as a simple function.
type TargetFunc func(args *Args) error

// Run implements Target.
func (f TargetFunc) Run(args *Args) error {
	return f(args)
}
