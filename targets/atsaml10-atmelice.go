package targets

import (
	"fmt"
	"github.com/goocd/goocd/actions/samflash"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/fileformats/autoparser"
	"github.com/goocd/goocd/mcus/atsame51"
	"github.com/goocd/goocd/mcus/sam/atsaml10e16a"
	"github.com/goocd/goocd/probes/samatmelice"
	"github.com/goocd/goocd/protocols/cmsisdap"
	"github.com/goocd/goocd/protocols/usbhid"
)

func init() {

	addTarget(&Target{
		Name:                "atsaml10-atmelice",
		Description:         "Atsaml10 using AtemlIce over cmsisdap-dap",
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

			if args.WriteMemU32Addr == 0x804000 {
				//, 0x8DCB8CED
				urowVals := []uint32{0xB08F437F, 0xFFFFF8BB, 0x7F0800FF, 0x00000001, 0x00000000, 0x00000000, 0x00003FFF, 0xC2BB6920}
				baseAddr := 0x804000

				for i := range urowVals {
					addr := baseAddr + (i * 4)
					val, err := atsam.ReadAddr32(uint32(addr), 1)
					checkErr(err)
					fmt.Printf("ReadAddr32[Address: 0x%x, Value: 0x%x]\n", uint32(addr), val)
				}

				err = atsam.WriteAddr32(uint32(atsaml10e16a.NVMCTRL_Addr+0x101C), uint32(baseAddr))
				checkErr(err)
				fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x\n]", uint32(atsaml10e16a.NVMCTRL_Addr+0x101C), uint32(baseAddr))

				for i, val := range urowVals {
					addr := baseAddr + (i * 4)
					err = atsam.WriteAddr32(uint32(addr), uint32(val))
					checkErr(err)
					fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x\n]", addr, val)
				}

				err = atsam.WriteAddr32(uint32(atsaml10e16a.NVMCTRL_Addr+0x1000), 0xA504)
				checkErr(err)
				fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x\n]", uint32(atsaml10e16a.NVMCTRL_Addr+0x1000), 0xA504)

				for i := range urowVals {
					addr := baseAddr + (i * 4)
					val, err := atsam.ReadAddr32(uint32(addr), 1)
					checkErr(err)
					fmt.Printf("ReadAddr32[Address: 0x%x, Value: 0x%x]\n", uint32(addr), val)
				}

			}

			if args.WriteMemU32Count > 0 && args.WriteMemU32Addr != 0x804000 {
				err = atsam.WriteAddr32(uint32(args.WriteMemU32Addr), uint32(args.WriteMemU32Value))
				checkErr(err)
				fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x\n]", args.WriteMemU32Addr, args.WriteMemU32Value)
			}

			if args.ReadMemU32Count > 0 {
				val, err := atsam.ReadAddr32(uint32(args.ReadMemU32Addr), args.ReadMemU32Count)
				checkErr(err)
				fmt.Printf("ReadAddr32[Address: 0x%x, Value: 0x%x]\n", args.ReadMemU32Addr, val)
			}

			if args.Load != "" {
				programReader, err := autoparser.ParseFromPath(args.Load, 0x0)
				checkErr(err)
				program, err := programReader.NextProgram()
				checkErr(err)
				nvm := &samflash.NVMFlash{
					CMSISDAP:                 cms,
					DAPTransferCoreAccess:    core,
					WriteAddress:             uint32(program.StartAddr()),
					EraseMultiplyer:          4, // Not easily Parsable, but in the Data sheet for the chip in the memory organization  section of NVMController
					NVMControllerAddress:     atsaml10e16a.NVMCTRL_Addr + 0x1000,
					NVMSetWriteAddressOffset: atsaml10e16a.NVMCTRL_ADDR_Offset,
					NVMPARAMOffset:           atsaml10e16a.NVMCTRL_PARAM_Offset,
					NVMPageSizeMask:          atsaml10e16a.NVMCTRL_PARAM_PSZ_Msk,
					NVMPageSizePos:           atsaml10e16a.NVMCTRL_PARAM_PSZ_Pos,
					NVMPageCountMask:         atsaml10e16a.NVMCTRL_PARAM_FLASHP_Msk,
					NVMPageCountPos:          atsaml10e16a.NVMCTRL_PARAM_FLASHP_Pos,

					NVMReadyOffSet: atsaml10e16a.NVMCTRL_STATUS_Offset,
					NVMReadyMask:   atsaml10e16a.NVMCTRL_STATUS_READY_Msk,
					NVMReadyVal:    atsaml10e16a.NVMCTRL_STATUS_READY,

					NVMCMDOffSet: atsaml10e16a.NVMCTRL_CTRLA_Offset,
					NVMCMDKey:    atsaml10e16a.NVMCTRL_CTRLA_CMDEX_KEY,
					NVMCMDKeyPos: atsaml10e16a.NVMCTRL_CTRLA_CMDEX_Pos,
					NVMEraseCMD:  atsaml10e16a.NVMCTRL_CTRLA_CMD_ER,
					NVMWriteCMD:  atsaml10e16a.NVMCTRL_CTRLA_CMD_WP,
				}
				err = nvm.LoadProgram(program.Bytes())
				checkErr(err)
				fmt.Printf("Successfully Flashed Rom\n")
			}

			if args.Reset {
				err = atsam.Reset()
				checkErr(err)
				fmt.Printf("Successfully Reset\n")
			}

			_ = cms.DAPDisconnect()
			return nil
		},
	})

}
