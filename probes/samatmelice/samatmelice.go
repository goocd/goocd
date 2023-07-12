package samatmelice

import "goocd/protocols/cmsisdap"

const (
	VendorID  = uint16(0x03eb)
	ProductID = uint16(0x2141)
)

var (
	IceParamaters = &cmsisdap.Parameters{
		SWDConfigClockCycles: cmsisdap.SWDConfigClockCycles1,
		SWDConfigDataPhase:   cmsisdap.SWDConfigNoDataPhase,
		// Todo: figure out what this actually means https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__SWJ__Sequence.html
		SWJSeqCount:       0x88,
		SWJSeqData:        []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x9E, 0xE7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		DAPTransferCycles: 0x0,
		// TODO: Tune this in. This extreme example was simple to let large transfers finish before returning
		DAPWaitTime:  0xFFFF,
		DAPMatchTime: 0x0,
		DAPPort:      cmsisdap.DebugPort,
	}
)
