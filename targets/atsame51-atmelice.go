package targets

import (
	"encoding/binary"
	"fmt"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/protocols/cmsisdap"
	"github.com/sstallion/go-hid"
	"log"
)

func init() {
	addTarget(&Target{
		Name:               "atsame51-atmelice",
		Description:        "Atsame51 using AtemlIce over cmsisdap-dap",
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
			cms := &cmsisdap.CMSISDAP{ReadWriter: d}

			// Reverse Endianness of the Clock Speed
			clockBuf := make([]byte, 4)
			binary.LittleEndian.PutUint32(clockBuf, cmsisdap.ClockSpeed2Mhz)

			// CMSIS Configure
			// Get to known state
			checkErr(cms.DAPDisconnect())
			// Config Cycles
			checkErr(cms.DAPSWDConfigure(cmsisdap.SWDConfigClockCycles1 | cmsisdap.SWDConfigNoDataPhase))

			// Todo: figure out what this actually means https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__SWJ__Sequence.html
			checkErr(cms.DAPSWJSequence(0x88, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x9E, 0xE7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}))

			// Set Clock Speed to 2Mhz. Note: This is debugger Clock Speed not to exceed 10x the chip Max ClockSpeed
			checkErr(cms.DAPSWJClock(binary.BigEndian.Uint32(clockBuf)))
			checkErr(cms.DAPTransferConfigure(0x0, 0x4000, 0x0))

			// Connect to Chip
			checkErr(cms.DAPConnect(cmsisdap.SWDPort))

			// Pass Configured CMS to the Cortex Driver
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

}
