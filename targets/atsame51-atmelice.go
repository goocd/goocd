package targets

import (
	"fmt"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/protocols/cmsis"
	"github.com/sstallion/go-hid"
	"log"
)

func init() {
	addTarget(&Target{
		Name:               "atsame51-atmelice",
		Description:        "Atsame51 using AtemlIce over cmsis-dap",
		SupportsReadMemU32: true,
		Run: func(args *Args) error {

			if err := hid.Init(); err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := hid.Exit(); err != nil {
					log.Fatal(err)
				}
			}()
			// Open the device using the VID and PID.
			d, err := hid.OpenFirst(0x03eb, 0x2141) // Atmel-ICE VIP & PID
			checkErr(err)
			cms := &cmsis.CMSISDAP{ReadWriter: d}
			checkErr(cms.Configure())
			core := cortexm4.DAPTransferCoreAccess{DAPTransferer: cms}
			checkErr(core.Configure())
			if args.ReadMemU32Count > 0 {
				var val uint32
				val, err = core.ReadAddr32(uint32(args.ReadMemU32Addr), args.ReadMemU32Count)
				checkErr(err)
				fmt.Printf("ReadAddr32[Address: 0x%x, Value: 0x%x]", args.ReadMemU32Addr, val)
			}
			return nil
		},
	})
	//TargetMap["atsame51-atmelice"] = TargetFunc(func(args *Args) error {
	//	// Initialize the hid package.
	//	if err := hid.Init(); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	// Open the device using the VID and PID.
	//	d, err := hid.OpenFirst(0x03eb, 0x2141) // Atmel-ICE VIP & PID
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	cms := &cmsis.CMSISDAP{ReadWriter: d}
	//	ice := &atmel_ice.AtmelICE{CMSISDAP: cms}
	//	err = ice.Configure()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//log.Printf("GOT HERE: atsame51-cmsisdap")
	//
	//	// Open File, Buffer here etc.
	//	if args.Load != "" {
	//		//pgmSrc, err := parserany.Parse(*loadF)
	//		//chkerr(err)
	//		//nvm := nvmload.NVMLoader {
	//		//	ProgramSource: pgmSrc
	//		//	NVMAccess: reg,
	//		//}
	//		//chkerr(nvm.NVMLoad())
	//	}
	//
	//	if args.ReadMem != "" {
	//		data, err := ice.ReadAddr32(addr)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		log.Printf("ReadAddr32[Address: %x, Value: %x]", addr, data)
	//	}
	//
	//	if args.Halt {
	//		//err = ice.WriteAddr32(debugAddress, debugWriteKey|debugEnable|debugHalt)
	//		//if err != nil {
	//		//	log.Fatal(err)
	//		//}
	//	}
	//
	//	if args.Resume {
	//		//err = ice.WriteAddr32(debugAddress, debugWriteKey|debugEnable)
	//		//if err != nil {
	//		//	log.Fatal(err)
	//		//}
	//	}
	//
	//	return nil
	//})

}
