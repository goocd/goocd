package atsame51

// TODO: split this apart and remove this file+dir, constants
// can come from sam/atsame51j20a instead

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/goocd/goocd/core/cortexm4"
	"github.com/goocd/goocd/fileformats"
	"github.com/goocd/goocd/protocols/cmsisdap"
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

	NVMKey      = 0xA500
	NVMPageSize = 0x200 // 512
)

const (
	SRAMBaseAddress = 0x20000000
)

type Atsame51 struct {
	*cmsisdap.CMSISDAP
	*cortexm4.DAPTransferCoreAccess
}

// Configure Debugger with specs that work on the Atsame51
func (a *Atsame51) Configure(clockSpeed uint32) error {
	if a.CMSISDAP == nil || a.DAPTransferCoreAccess == nil {
		return fmt.Errorf("error: Atsame51.Configure() missing required dependencies")
	}

	// Reverse Endianness of the Clock Speed
	clockBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(clockBuf, clockSpeed)

	// Debugger/CMSIS Configure for the Atsame51

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
	// TODO: Tune this in. This extreme example was simple to let large transfers finish before returning
	err = a.DAPTransferConfigure(0x0, 0xFFFF, 0x0)
	if err != nil {
		return err
	}
	// Connect to Chip
	err = a.DAPConnect(cmsisdap.SWDPort)
	if err != nil {
		return err
	}

	// Cortex Configuration
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

	//TODO: Clear Lock

	return nil
}

func (a *Atsame51) LoadProgram(startAddress uint32, rom []byte) error {
	// Convert Rom to useable state
	rom32, err := fileformats.ConvertByteSliceUint32Slice(rom)
	if err != nil {
		return err
	}

	// Halt MCU
	err = a.Halt()
	if err != nil {
		return err
	}

	// Clear Locks
	err = a.ClearRegionLock()
	if err != nil {
		return err
	}

	// Sets Start Address of NVM CMDEX Commands.
	// [In English] Basically just makes the NVM Controller Point to the address you provide as the target for future operations.
	err = a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLADDR, 0x10000)
	if err != nil {
		return err
	}

	// EraseBlock
	err = a.NVMCMD(0x1)
	if err != nil {
		return err
	}

	err = a.WriteTransfer32(cmsisdap.AccessPort, cmsisdap.PortRegister0, 0x52)
	if err != nil {
		return err
	}

	buffer := make([]uint32, 0, 64)
	i := 0
	initialize := true
	for _, val := range rom32 {
		buffer = append(buffer, val)
		if len(buffer) < 64 {
			continue
		}

		//fmt.Printf("Index: %d\n", i)
		err = a.WriteSeqAddr32(initialize, startAddress+(uint32(i)*4), buffer) // AKA 0x80 AKA 512 bytes AKA 1 Page Size
		if err != nil {
			return fmt.Errorf("error: Atsame51.LoadProgram() WriteSeqAdd32 Error: %+v", err)
		}
		initialize = false
		i += 64
		buffer = buffer[:0]

		if i%128 == 0 { // 128*4 = 512 == 1 page
			err = a.NVMCMD(0x3)
			if err != nil {
				return fmt.Errorf("error: Atsame51.LoadProgram() WriteAddr32 Error: %+v", err)
			}
			initialize = true
		}
	}

	err = a.WriteSeqAddr32(initialize, startAddress+(uint32(i)*4), buffer) // AKA 0x80 AKA 512 bytes AKA 1 Page Size
	if err != nil {
		return fmt.Errorf("error: Atsame51.LoadProgram() WriteSeqAdd32 Error: %+v", err)
	}

	err = a.NVMCMD(0x3)
	if err != nil {
		return fmt.Errorf("error: Atsame51.LoadProgram() WriteAddr32 Error: %+v", err)
	}

	return nil
}

// Todo: Maybe add custom time out options

// WaitForNVMCMDClear reads the NVM CMD register till it's complete or times out
// Then it writes the first bit to clear the flag for the next time it needs to wait
func (a *Atsame51) WaitForNVMCMDClear() error {
	ti := time.Now()
	for {
		// Read Flag
		val, err := a.ReadAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 1)
		if err != nil {
			return err
		}

		//fmt.Printf("ValAP: %x\n", val)
		// Bitwise Flag check
		if val&0x1 > 0 {
			break
		}

		// Timeout
		if time.Since(ti) > time.Second*1 {
			return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
		}
	}

	// Clear Interrupt
	err := a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 0x1)
	if err != nil {
		return err
	}

	ti = time.Now()
	for {
		// Wait For interrupt to be cleared after you wrote to it
		val, err := a.ReadAddr32(NVMCTRLBaseAddress|NVMCTRLINTFLAG, 1)
		if err != nil {
			return err
		}

		//fmt.Printf("ValAP: %x\n", val)
		if val&0x1 == 0 {
			break
		}

		if time.Since(ti) > time.Second*1 {
			return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
		}
	}

	return nil
}

func (a *Atsame51) NVMCMD(cmd uint32) error {
	err := a.WriteAddr32(NVMCTRLBaseAddress|NVMCTRLCTRLB, NVMKey|cmd)
	if err != nil {
		return err
	}
	err = a.WaitForNVMCMDClear()
	if err != nil {
		return err
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
