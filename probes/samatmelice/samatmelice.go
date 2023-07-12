package samatmelice

import "github.com/goocd/goocd/protocols/cmsisdap"

const (
	VendorID  = uint16(0x03eb)
	ProductID = uint16(0x2141)
)

var (
	IceParamaters = &cmsisdap.Parameters{
		SWDConfigClockCycles: cmsisdap.SWDConfigClockCycles1,
		SWDConfigDataPhase:   cmsisdap.SWDConfigNoDataPhase,
		SWJSeqCount:          0x88,
		SWJSeqData:           []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x9E, 0xE7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		DAPTransferCycles:    0x0,
		DAPWaitTime:          0xFFFF,
		DAPMatchTime:         0x0,
		DAPPort:              cmsisdap.DebugPort,
	}
)
