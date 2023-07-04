package atsame51

import (
	"encoding/binary"
	"fmt"
	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/protocols/cmsisdap"
	"time"
)

const (
	DSUBaseAddress = 0x41002000
	DSUCTRL        = 0x0
	DSUStatusA     = 0x1

	DSUReset = 0x1
	DSUDone  = 0x1
)

const (
	NVMCTRLBaseAddress = 0x41004000
	NVMCTRLCTRLA       = 0x0
	NVMCTRLCTRLB       = 0x4
	NVMCTRLNVMPARAM    = 0x8
	NVMCTRLINTENCLR    = 0xC
	NVMCTRLINTENSET    = 0xE
	NVMCTRLINTFLAG     = 0x10
	NVMCTRLSTATUS      = 0x12
	NVMCTRLADDR        = 0x14
	NVMCTRLRUNLOCK     = 0x18

	NVMPageSize = 0x200 // 512
)

const (
	SRAMBaseAddress = 0x20000000
)

type Atsame51 struct {
	*cmsisdap.CMSISDAP
	*cortexm4.DAPTransferCoreAccess
}

func (a *Atsame51) Configure(clockSpeed uint32) error {
	if a.CMSISDAP == nil || a.DAPTransferCoreAccess == nil {
		return fmt.Errorf("error: Atsame51.Configure() missing required dependencies")
	}

	// Reverse Endianness of the Clock Speed
	clockBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(clockBuf, clockSpeed)

	// CMSIS Configure for the Atsame51

	// Get to known state
	err := a.DAPDisconnect()
	if err != nil {
		return err
	}
	// Config Cycles
	err = a.DAPSWDConfigure(cmsisdap.SWDConfigClockCycles1 | cmsisdap.SWDConfigNoDataPhase)
	if err != nil {
		return err
	}
	// Todo: figure out what this actually means https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__SWJ__Sequence.html
	err = a.DAPSWJSequence(0x88, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x9E, 0xE7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	if err != nil {
		return err
	}
	// Set Clock Speed to 2Mhz. Note: This is debugger Clock Speed not to exceed 10x the chip Max ClockSpeed
	err = a.DAPSWJClock(binary.BigEndian.Uint32(clockBuf))
	if err != nil {
		return err
	}
	err = a.DAPTransferConfigure(0x0, 0x4000, 0x0)
	if err != nil {
		return err
	}
	// Connect to Chip
	err = a.DAPConnect(cmsisdap.SWDPort)
	if err != nil {
		return err
	}

	err = a.DAPTransferCoreAccess.Configure()
	if err != nil {
		return err
	}
	return nil
}

func (a *Atsame51) ClearRegionLock() error {
	// Check Locks
	resp, err := a.ReadAddr32(NVMCTRLBaseAddress|NVMCTRLRUNLOCK, 1)
	if err != nil {
		return err
	}

	if resp == 0xFFFFFFFF {
		return nil
	}

	//Todo: Clear Lock

	return nil
}

func (a *Atsame51) LoadProgram(startAddress uint32, rom []byte) error {
	rom32, err := a.ConvertByteSliceUint32Slice(rom)
	if err != nil {
		return err
	}

	err = a.Halt()
	if err != nil {
		return err
	}

	err = a.ClearRegionLock()
	if err != nil {
		return err
	}

	err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLADDR, 0x10000)
	if err != nil {
		return err
	}

	err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLCTRLB, 0xA501)
	if err != nil {
		return err
	}

	err = a.WaitForCMDClear()
	if err != nil {
		return err
	}

	err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 0x1)
	if err != nil {
		return err
	}

	for i := range rom32 {
		fmt.Printf("Rom32: Len %d\n", len(rom32))
		err = a.WriteAddr32(startAddress+uint32(i*4), rom32[i]) // AKA 0x80 AKA 512 bytes AKA 1 Page Size
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() WriteSeqAdd32 Error: %+v", err)
		}

		err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLCTRLB, 0xA503)
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() WriteAddr32 Error: %+v", err)
		}

		err = a.WaitForCMDClear()
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() WaitForCMDClear Error: %+v", err)
		}

		err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 0x1)
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() WriteAddr32 Error: %+v", err)
		}
	}
	return nil
}

func (a *Atsame51) ConvertByteSliceUint32Slice(rom []byte) ([]uint32, error) {
	if len(rom)%4 > 0 {
		return nil, fmt.Errorf("error: Atsame51.ConvertByteSliceUint32Slice() Misaligned Rom")
	}

	rom32 := make([]uint32, len(rom)/4)

	for i := range rom32 {
		rom32[i] = binary.LittleEndian.Uint32(rom[i*4:])
	}

	//romSegements := make([][]uint32, 0, (len(rom32)/128)+1)
	//
	//for i := 0; i <= len(rom32)/128; i++ {
	//	buffer := make([]uint32, 128)
	//	copy(buffer, rom32[i*128:])
	//	romSegements = append(romSegements, buffer)
	//}

	return rom32, nil
}

func (a *Atsame51) WaitForCMDClear() error {
	ti := time.Now()
	for {
		val, err := a.ReadAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 1)
		if err != nil {
			return err
		}

		if val&0x1 > 0 {
			break
		}

		val, err = a.ReadTransfer32(cmsisdap.DebugPort, cmsisdap.PortRegisterC)
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() transfer Error: %+v", err)
		}

		if val&0x1 > 0 {
			break
		}

		if time.Since(ti) > time.Second*3 {
			return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
		}
	}

	return nil
}

// Reset Via hardware reset
func (a *Atsame51) Reset() error {
	err := a.CMSISDAP.Reset()
	if err != nil {
		return err
	}
	return nil
}
