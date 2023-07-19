package samflash

import (
	"encoding/binary"
	"fmt"
	"goocd/core/cortexm4"
	"goocd/protocols/cmsisdap"
	"time"
)

const (
	// 8 bytes
	NVMCTRL_PARAM_PSZ_8 = 0x0
	// 16 bytes
	NVMCTRL_PARAM_PSZ_16 = 0x1
	// 32 bytes
	NVMCTRL_PARAM_PSZ_32 = 0x2
	// 64 bytes
	NVMCTRL_PARAM_PSZ_64 = 0x3
	// 128 bytes
	NVMCTRL_PARAM_PSZ_128 = 0x4
	// 256 bytes
	NVMCTRL_PARAM_PSZ_256 = 0x5
	// 512 bytes
	NVMCTRL_PARAM_PSZ_512 = 0x6
	// 1024 bytes
	NVMCTRL_PARAM_PSZ_1024 = 0x7
)

//type Cortex interface {
//	Halt() error
//	ReadAddr32(addr uint32, count int) (value uint32, err error)
//	WriteAddr32(addr, value uint32) error
//	WriteSeqAddr32(addr uint32, value []uint32) error
//	WriteTransfer32(port, portRegister byte, value uint32) error
//}

type NVMFlash struct {
	*cmsisdap.CMSISDAP
	*cortexm4.DAPTransferCoreAccess

	Stats bool

	WriteSize       uint32
	PageCount       uint32
	EraseSize       uint32
	EraseMultiplyer uint32
	FlashSize       uint32

	WriteAddress             uint32
	NVMControllerAddress     uint32
	NVMSetWriteAddressOffset uint32

	NVMPARAMOffset   uint32
	NVMPageSizeMask  uint32
	NVMPageSizePos   uint32
	NVMPageCountMask uint32
	NVMPageCountPos  uint32
	NVMPageSizeBits  uint32

	NVMClearReady  bool
	NVMReadyOffSet uint32
	NVMReadyMask   uint32
	NVMReadyVal    uint32
	nvmReadyShift  uint32

	NVMCMDOffSet uint32
	NVMCMDKey    uint32
	NVMCMDKeyPos uint32
	NVMEraseCMD  uint32
	NVMWriteCMD  uint32
}

func (nvm *NVMFlash) LoadProgram(rom []byte) error {
	if nvm == nil {
		return nil
	}

	err := nvm.Configure()
	if err != nil {
		return err
	}

	err = nvm.MemoryChecks(len(rom))
	if err != nil {
		return err
	}

	err = nvm.Halt()
	if err != nil {
		return err
	}

	err = nvm.ClearRegionLock()
	if err != nil {
		return err
	}

	err = nvm.WriteTransfer32(cmsisdap.AccessPort, cmsisdap.PortRegister0, 0x52)
	if err != nil {
		return err
	}

	buffer := make([]uint32, 0, nvm.WriteSize/4)
	offset := uint32(0)
	nvm.nvmReadyShift = (nvm.NVMReadyOffSet % 4) * 8
	//fmt.Printf("NVMSHIFT: %d\n", nvm.nvmReadyShift)

	if nvm.Stats {
		fmt.Printf("Total ROM LEN: %d\n", len(rom))
		fmt.Printf("Total Flash Size: %d\n", nvm.FlashSize)
		fmt.Printf("Erase Row Size: %d\n", nvm.EraseSize)
		fmt.Printf("Page Size: %d\n", nvm.WriteSize)
	}

	for i := 0; i < len(rom); i += 4 {
		if i%int(nvm.EraseSize) == 0 {
			//fmt.Printf("Setting New Base: %x\n", nvm.WriteAddress+offset)
			err = nvm.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMSetWriteAddressOffset, nvm.WriteAddress+offset)
			if err != nil {
				return err
			}
			// Todo: See if this needs to update the Addr field in the NVM CTRL
			err = nvm.Erase()
			if err != nil {
				return err
			}
		}

		val := binary.LittleEndian.Uint32(rom[i:])
		buffer = append(buffer, val)

		if len(buffer) < int(nvm.WriteSize/4) {
			continue
		}
		//fmt.Printf("Write %x To Address: %x\n", buffer, nvm.WriteAddress+offset)
		err = nvm.WriteSeqAddr32(nvm.WriteAddress+offset, buffer)
		if err != nil {
			return err
		}
		buffer = buffer[:0]
		offset += nvm.WriteSize
		err = nvm.Commit()
		if err != nil {
			return err
		}
	}

	if len(buffer) > 0 { // Means not page aligned
		for len(buffer) < int(nvm.WriteSize/4) {
			buffer = append(buffer, 0) // 0 Pad it, to page align it
		}
		err = nvm.WriteSeqAddr32(nvm.WriteAddress+offset, buffer)
		if err != nil {
			return err
		}
		err = nvm.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (nvm *NVMFlash) ClearRegionLock() error {
	// Check Locks
	resp, err := nvm.ReadAddr32(nvm.NVMControllerAddress|0x18, 1)
	if err != nil {
		return err
	}

	if resp == 0xFFFFFFFF {
		return nil
	}

	//TODO: Clear Lock

	return nil
}

func (nvm *NVMFlash) Commit() error {
	//fmt.Printf("Issuing Write Command: %x\n", (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMWriteCMD)
	err := nvm.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMCMDOffSet, (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMWriteCMD)
	if err != nil {
		return err
	}
	err = nvm.WaitForReady()
	if err != nil {
		return err
	}
	return nil
}

func (nvm *NVMFlash) Erase() error {
	//	fmt.Printf("Issuing Erase Command: %x\n", (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMEraseCMD)
	err := nvm.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMCMDOffSet, (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMEraseCMD)
	if err != nil {
		return err
	}
	err = nvm.WaitForReady()
	if err != nil {
		return err
	}
	return nil
}

func (nvm *NVMFlash) WaitForReady() error {
	ti := time.Now()
	for {
		// Read Flag
		val, err := nvm.ReadAddr32(nvm.NVMControllerAddress+nvm.NVMReadyOffSet, 1)
		if err != nil {
			return err
		}

		val = val >> nvm.nvmReadyShift
		//fmt.Printf("ValAP: %x\n", val)
		// Bitwise Flag check
		if val&nvm.NVMReadyMask > 0 {
			break
		}

		// Timeout
		if time.Since(ti) > time.Second*1 {
			return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
		}
	}

	// Todo: See if this is even needed since this isn't the interrupt register
	if nvm.NVMClearReady {
		// Clear Interrupt
		err := nvm.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMReadyOffSet, nvm.NVMReadyVal)
		if err != nil {
			return err
		}

		ti = time.Now()
		for {
			// Wait For interrupt to be cleared after you cleared it
			val, err := nvm.ReadAddr32(nvm.NVMControllerAddress+nvm.NVMReadyOffSet, 1)
			if err != nil {
				return err
			}
			val = val >> nvm.nvmReadyShift
			//fmt.Printf("ValAP: %x\n", val)
			if val&nvm.NVMReadyMask == 0 {
				break
			}

			if time.Since(ti) > time.Second*1 {
				return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
			}
		}
	}
	return nil
}

func (nvm *NVMFlash) Configure() error {
	readVal, err := nvm.ReadAddr32(nvm.NVMControllerAddress+nvm.NVMPARAMOffset, 1)
	if err != nil {
		return err
	}

	nvm.PageCount = (readVal & nvm.NVMPageCountMask) >> nvm.NVMPageCountPos
	nvm.NVMPageSizeBits = (readVal & nvm.NVMPageSizeMask) >> nvm.NVMPageSizePos
	switch nvm.NVMPageSizeBits {
	case NVMCTRL_PARAM_PSZ_8:
		nvm.WriteSize = 8
	case NVMCTRL_PARAM_PSZ_16:
		nvm.WriteSize = 16
	case NVMCTRL_PARAM_PSZ_32:
		nvm.WriteSize = 32
	case NVMCTRL_PARAM_PSZ_64:
		nvm.WriteSize = 64
	case NVMCTRL_PARAM_PSZ_128:
		nvm.WriteSize = 128
	case NVMCTRL_PARAM_PSZ_256:
		nvm.WriteSize = 256
	case NVMCTRL_PARAM_PSZ_512:
		nvm.WriteSize = 512
	case NVMCTRL_PARAM_PSZ_1024:
		nvm.WriteSize = 1024
	default:
		nvm.WriteSize = 0
	}
	nvm.EraseSize = nvm.WriteSize * nvm.EraseMultiplyer
	nvm.FlashSize = nvm.WriteSize * nvm.PageCount
	return nil
}

func (nvm *NVMFlash) MemoryChecks(romLen int) error {

	if nvm.WriteSize == 0 {
		return fmt.Errorf("error: could not determine write size of chip with value: %d", nvm.NVMPageSizeBits)
	}

	// Rom Larger than Flash Size
	if romLen > int(nvm.FlashSize) {
		return fmt.Errorf("error: Rom of Len %d exceeds Flash of Size: %d", romLen, nvm.FlashSize)
	}

	// Rom from Address would overflow Flash Size
	if int(nvm.FlashSize) < int(nvm.WriteAddress)+romLen {
		return fmt.Errorf("error: Rom of Len %d exceeds Flash of Size: %d When starting at Address: %x", romLen, nvm.FlashSize, nvm.WriteAddress)
	}

	// If Base Address isn't Row/Block and Page aligned Error since it needs to be on a boundary
	if nvm.WriteAddress%nvm.EraseSize != 0 || nvm.WriteAddress%nvm.WriteSize != 0 {
		return fmt.Errorf("error: Rom Start Address: %x is not aligned with page Size: %d or Row/Block Size: %d", nvm.WriteAddress, nvm.WriteSize, nvm.EraseSize)
	}

	// Todo: Any other memory checks we might need

	return nil
}
