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
	loadF := flag.String("load", "", "Load program file (.elf, .hex, .bin) to flash, base address implied from file or defaults based on target")
	readmemu32 := flag.String("readmemu32", "", "uint32 memory address you wish to read followed by optional 32-bit word count, e.g. '0x20004000,5'")
	writememu32 := flag.String("writememu32", "", "uint32 memory address and value you wish to write and optional 32-bit word count, comma separated, e.g. '0x20004000,0xF0E0D0C0,1'")
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
		addr, err := strconv.ParseUint(splitReadMem[0], 0, 64) // supports hex, dec, oct, bin
		if err != nil {
			log.Fatalf("Unable to parse %q into an address + count: %v", *readmemu32, err)
		}
		count := int64(1)
		if len(splitReadMem) > 1 {
			count, err = strconv.ParseInt(splitReadMem[1], 0, 64)
			if err != nil {
				log.Fatalf("Unable to parse %q into an address + count: %v", *readmemu32, err)
			}
			if count < 0 {
				log.Fatalf("Invalid count %d", count)
			}
		}
		args.ReadMemU32Addr = addr
		args.ReadMemU32Count = int(count)
	}

	if tgt.SupportsWriteMemU32 && *writememu32 != "" {
		splitReadMem := strings.Split(*writememu32, ",")
		if len(splitReadMem) < 2 {
			log.Fatalf("Unable to parse %q into an address + value + count properly", *writememu32)
		}
		addr, err := strconv.ParseUint(splitReadMem[0], 0, 64)
		if err != nil {
			log.Fatalf("Unable to parse %q into an address + value + count properly: %v", *writememu32, err)
		}

		value, err := strconv.ParseUint(splitReadMem[1], 0, 64)
		if err != nil {
			log.Fatalf("Unable to parse %q into an address + value + count properly: %v", *writememu32, err)
		}

		count := int64(1)
		if len(splitReadMem) > 2 {
			count, err = strconv.ParseInt(splitReadMem[2], 0, 64)
			if err != nil {
				log.Fatalf("Unable to parse %s into an address + count", *readmemu32)
			}

			if count > 0 {
				log.Fatalf("Unable to parse %s into an address + count", *readmemu32)
			}
		}
		args.WriteMemU32Addr = addr
		args.WriteMemU32Value = value
		args.WriteMemU32Count = int(count)
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
