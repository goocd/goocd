package targets

import (
	"fmt"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/fileformats/elfparser"
	"github.com/goocd/goocd/mcus/atsame51"
	"github.com/goocd/goocd/probes/samatmelice"
	"github.com/goocd/goocd/protocols/cmsisdap"
	"github.com/goocd/goocd/protocols/usbhid"
)

func init() {
	addTarget(&Target{
		Name:                "atsame51-atmelice",
		Description:         "Atsame51 using AtemlIce over cmsisdap-dap",
		SupportsReadMemU32:  true,
		SupportsWriteMemU32: true,
		SupportsReset:       true,
		SupportsLoad:        true,
		Run: func(args *Args) error {
			d, err := usbhid.OpenFirstHid(samatmelice.VendorID, samatmelice.ProductID)
			checkErr(err)
			defer d.CleanUp()
			// Pass CMSIS the USBHID Device
			cms := &cmsisdap.CMSISDAP{ReadWriter: d}

			// Pass CMS to the Cortex Driver
			core := &cortexm4.DAPTransferCoreAccess{DAPTransferer: cms}

			// Pass Both to Atsame51 struct
			atsam := &atsame51.Atsame51{CMSISDAP: cms, DAPTransferCoreAccess: core}
			checkErr(atsam.Configure(cmsisdap.ClockSpeed2Mhz))

			if args.WriteMemU32Count > 0 {
				err := atsam.WriteAddr32(uint32(args.WriteMemU32Addr), uint32(args.WriteMemU32Value))
				checkErr(err)
				fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x]", args.ReadMemU32Addr, args.WriteMemU32Value)
			}

			if args.ReadMemU32Count > 0 {
				val, err := atsam.ReadAddr32(uint32(args.ReadMemU32Addr), args.ReadMemU32Count)
				checkErr(err)
				fmt.Printf("ReadAddr32[Address: 0x%x, Value: 0x%x]", args.ReadMemU32Addr, val)
			}

			if args.Load != "" {
				addr, rom, err := elfparser.ExtractROM(args.Load)
				checkErr(err)
				err = atsam.LoadProgram(uint32(addr), rom)
				checkErr(err)
				fmt.Printf("Successfully Flashed Rom")
			}

			if args.Reset {
				err = atsam.Reset()
				checkErr(err)
				fmt.Printf("Successfully Reset")
			}

			_ = cms.DAPDisconnect()
			return nil
		},
	})
}
