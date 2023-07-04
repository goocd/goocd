package targets

import (
	"log"
	"testing"
)

func TestAtsame51AtmelICERun(t *testing.T) {
	tgt := TargetMap["atsame51-samatmelice"]
	if tgt == nil {
		log.Fatalf("Unable to find target %q, try 'goocd -target-list' to see available targets.", "test")
	}
	args := Args{}
	args.ReadMemU32Addr = 0x806018
	err := tgt.Run(&args)
	if err != nil {
		t.Fatal(err)
	}
}
