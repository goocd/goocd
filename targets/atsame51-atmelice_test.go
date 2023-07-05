package targets

import (
	"log"
	"testing"
)

func TestLoadSequence(t *testing.T) {
	// atsame51-atmelice -load=bin/readSRAMLoop-waldo_hardwired_r9.elf
	tgt := TargetMap["atsame51-atmelice"]
	if tgt == nil {
		t.Fatalf("Unable to find target %q, try 'goocd -target-list' to see available targets.", "Test")
	}

	args := Args{}
	args.Load = "C:\\Users\\AlexLeon\\onestepgps\\trackersw\\bin\\readSRAMLoop-waldo_hardwired_r9.elf"
	err := tgt.Run(&args)
	if err != nil {
		log.Fatal(err)
	}

}
