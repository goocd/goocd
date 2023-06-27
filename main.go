package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goocd/goocd/targets"
)

func main() {

	targetListF := flag.Bool("target-list", false, "List all compiled-in targets")
	targetF := flag.String("target", "", "Select a target")
	loadF := flag.String("load", "", "Load file (.elf, .hex, .bin)")
	readmemu32 := flag.String("readmemu32", "", "Uint32 Hex Memory Address you wish to read")

	flag.Usage = func() {
		// TODO: customize as needed
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *targetListF {
		// TODO: sort and see if we want to implement any help info
		fmt.Printf("Targets:\n%+v", targets.TargetMap)
		return
	}

	args := targets.Args{}
	args.Load = *loadF
	args.ReadMem = *readmemu32

	tgt := targets.TargetMap[*targetF]
	if tgt == nil {
		log.Fatalf("Unable to find target %q, try 'goocd -target-list' to see available targets.", *targetF)
	}

	err := tgt.Run(&args)
	if err != nil {
		log.Fatal(err)
	}

}
