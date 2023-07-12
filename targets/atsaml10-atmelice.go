package targets

import (
	"fmt"
	"github.com/goocd/goocd/actions/samflash"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/fileformats/autoparser"
	"github.com/goocd/goocd/mcus/sam/atsaml10d16a"
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

			// Configure CMSIS + Cortex
			checkErr(cms.Configure(cmsisdap.ClockSpeed2Mhz, samatmelice.IceParamaters))
			checkErr(core.Configure())

			if args.WriteMemU32Count > 0 && args.WriteMemU32Addr != 0x804000 {
				err = core.WriteAddr32(uint32(args.WriteMemU32Addr), uint32(args.WriteMemU32Value))
				checkErr(err)
				fmt.Printf("WriteAddr32[Address: 0x%x, Value: 0x%x\n]", args.WriteMemU32Addr, args.WriteMemU32Value)
			}

			if args.ReadMemU32Count > 0 {
				val, err := core.ReadAddr32(uint32(args.ReadMemU32Addr), args.ReadMemU32Count)
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
					NVMControllerAddress:     atsaml10d16a.NVMCTRL_Addr,
					NVMSetWriteAddressOffset: atsaml10d16a.NVMCTRL_ADDR_Offset,
					NVMPARAMOffset:           atsaml10d16a.NVMCTRL_PARAM_Offset,
					NVMPageSizeMask:          atsaml10d16a.NVMCTRL_PARAM_PSZ_Msk,
					NVMPageSizePos:           atsaml10d16a.NVMCTRL_PARAM_PSZ_Pos,
					NVMPageCountMask:         atsaml10d16a.NVMCTRL_PARAM_FLASHP_Msk,
					NVMPageCountPos:          atsaml10d16a.NVMCTRL_PARAM_FLASHP_Pos,

					NVMClearReady:  false,
					NVMReadyOffSet: atsaml10d16a.NVMCTRL_STATUS_Offset,
					NVMReadyMask:   atsaml10d16a.NVMCTRL_STATUS_READY_Msk,
					NVMReadyVal:    atsaml10d16a.NVMCTRL_STATUS_READY,

					NVMCMDOffSet: atsaml10d16a.NVMCTRL_CTRLA_Offset,
					NVMCMDKey:    atsaml10d16a.NVMCTRL_CTRLA_CMDEX_KEY,
					NVMCMDKeyPos: atsaml10d16a.NVMCTRL_CTRLA_CMDEX_Pos,
					NVMEraseCMD:  atsaml10d16a.NVMCTRL_CTRLA_CMD_ER,
					NVMWriteCMD:  atsaml10d16a.NVMCTRL_CTRLA_CMD_WP,
				}
				err = nvm.LoadProgram(program.Bytes())
				checkErr(err)
				fmt.Printf("Successfully Flashed Rom\n")
			}

			if args.Reset {
				err = cms.Reset()
				checkErr(err)
				fmt.Printf("Successfully Reset\n")
			}

			_ = cms.DAPDisconnect()
			return nil
		},
	})

}
