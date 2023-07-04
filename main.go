package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/goocd/goocd/targets"
)

func main() {

	targetListF := flag.Bool("target-list", false, "List all compiled-in targets")
	targetF := flag.String("target", "", "Select a target")
	loadF := flag.String("load", "", "Load file (.elfparser, .hex, .bin)")
	readmemu32 := flag.String("readmemu32", "", "Uint32 Hex Memory Address you wish to read")
	reset := flag.Bool("reset", false, "Issue Reset Command to target")
	//Halt := flag.Bool("halt", false, "Halts run time operation")
	//resume := flag.Bool("resume", false, "Resumes runtime operation")

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

	tgt := targets.TargetMap[*targetF]
	if tgt == nil {
		log.Fatalf("Unable to find target %q, try 'goocd -target-list' to see available targets.", *targetF)
	}

	args := targets.Args{}

	if tgt.SupportsReadMemU32 && *readmemu32 != "" {
		splitReadMem := strings.Split(*readmemu32, ",")
		base := 0
		if !strings.HasPrefix(splitReadMem[0], "0x") {
			base = 16
		}
		addr, err := strconv.ParseUint(*readmemu32, base, 64)
		if err != nil {
			log.Fatalf("Unable to parse %s into an address + count", *readmemu32)
		}
		count := int64(1)
		if len(splitReadMem) > 1 {
			count, err = strconv.ParseInt(splitReadMem[1], 0, 64)
			if err != nil {
				log.Fatalf("Unable to parse %s into an address + count", *readmemu32)
			}
		}
		args.ReadMemU32Addr = addr
		args.ReadMemU32Count = int(count)
	}

	if tgt.SupportsReset && *reset {
		args.Reset = *reset
	}

	if tgt.SupportsLoad && *loadF != "" {
		args.Load = *loadF
	}
	//args.Load = *loadF
	//args.ReadMem = *readmemu32
	//args.Halt = *Halt
	//args.Resume = *resume

	err := tgt.Run(&args)
	if err != nil {
		log.Fatal(err)
	}

}
