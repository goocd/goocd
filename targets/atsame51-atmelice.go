package targets

import (
	"fmt"
	"github.com/goocd/goocd/actions/samflash"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/fileformats/autoparser"
	"github.com/goocd/goocd/mcus/atsame51"
	"github.com/goocd/goocd/mcus/sam/atsame51j20a"
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
					EraseMultiplyer:          16, // Not easily Parsable, but in the Data sheet for the chip in the memory organization  section of NVMController
					NVMControllerAddress:     atsame51j20a.NVMCTRL_Addr,
					NVMSetWriteAddressOffset: atsame51j20a.NVMCTRL_ADDR_Offset,
					NVMPARAMOffset:           atsame51j20a.NVMCTRL_PARAM_Offset,
					NVMPageSizeMask:          atsame51j20a.NVMCTRL_PARAM_PSZ_Msk,
					NVMPageSizePos:           atsame51j20a.NVMCTRL_PARAM_PSZ_Pos,
					NVMPageCountMask:         atsame51j20a.NVMCTRL_PARAM_NVMP_Msk,
					NVMPageCountPos:          atsame51j20a.NVMCTRL_PARAM_NVMP_Pos,

					NVMClearReady:  true,
					NVMReadyOffSet: atsame51j20a.NVMCTRL_STATUS_Offset,
					NVMReadyMask:   atsame51j20a.NVMCTRL_STATUS_READY_Msk,
					NVMReadyVal:    atsame51j20a.NVMCTRL_STATUS_READY,

					NVMCMDOffSet: atsame51j20a.NVMCTRL_CTRLB_Offset,
					NVMCMDKey:    atsame51j20a.NVMCTRL_CTRLB_CMDEX_KEY,
					NVMCMDKeyPos: atsame51j20a.NVMCTRL_CTRLB_CMDEX_Pos,
					NVMEraseCMD:  atsame51j20a.NVMCTRL_CTRLB_CMD_EB,
					NVMWriteCMD:  atsame51j20a.NVMCTRL_CTRLB_CMD_WP,
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
