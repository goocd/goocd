package samflash

import (
	"encoding/binary"
	"fmt"
	"github.com/goocd/goocd/protocols/cmsisdap"
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

type Cortex interface {
	Halt() error
	ReadAddr32(addr uint32, count int) (value uint32, err error)
	WriteAddr32(addr, value uint32) error
	WriteSeqAddr32(addr uint32, value []uint32) error
	WriteTransfer32(port, portRegister byte, value uint32) error
}

type NVMFlash struct {
	Cortex

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

	NVMReadyOffSet uint32
	NVMReadyMask   uint32
	NVMReadyVal    uint32

	NVMCMDOffSet uint32
	NVMCMDKey    uint32
	NVMCMDKeyPos uint32
	NVMEraseCMD  uint32
	NVMWriteCMD  uint32
}

func (nvm *NVMFlash) LoadProgram(rom []byte) error {
	if nvm.Cortex == nil {
		return nil
	}

	err := nvm.Halt()
	if err != nil {
		return err
	}

	err = nvm.Configure()
	if err != nil {
		return err
	}

	err = nvm.MemoryChecks(len(rom))
	if err != nil {
		return err
	}
	err = nvm.WriteTransfer32(cmsisdap.AccessPort, cmsisdap.PortRegister0, 0xA2000022)
	if err != nil {
		return err
	}
	fmt.Printf("Debug NVM Data: %+v\n", *nvm)
	// Set Base Address
	err = nvm.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMSetWriteAddressOffset, nvm.WriteAddress)
	if err != nil {
		return err
	}

	err = nvm.Erase()
	if err != nil {
		return err
	}
	//rom[0] = 0x10
	//rom[1] = 0x10
	//rom[2] = 0x10
	//rom[3] = 0x10
	buffer := make([]uint32, 0, nvm.WriteSize/4)
	offset := uint32(0)
	for i := 0; i < len(rom); i += 4 {
		//if i%int(nvm.EraseSize) == 0 {
		//	// Todo: See if this needs to update the Addr field in the NVM CTRL
		//	err = nvm.Erase()
		//	if err != nil {
		//		return err
		//	}
		//}

		val := binary.LittleEndian.Uint32(rom[i:])
		fmt.Printf("Writing Val: %x\n", val)
		fmt.Printf("B[0] = %x\n", rom[i])
		fmt.Printf("B[1] = %x\n", rom[i+1])
		fmt.Printf("B[2] = %x\n", rom[i+2])
		fmt.Printf("B[3] = %x\n", rom[i+3])
		buffer = append(buffer, val)

		if len(buffer) < int(nvm.WriteSize/4) {
			continue
		}
		//fmt.Printf("Write To Address: %x\n", nvm.WriteAddress+offset)
		err = nvm.Cortex.WriteSeqAddr32(nvm.WriteAddress+offset, buffer)
		if err != nil {
			return err
		}
		err = nvm.Commit()
		if err != nil {
			return err
		}
		buffer = buffer[:0]
		offset += nvm.WriteSize
	}

	if len(buffer) > 0 { // Means not page aligned
		for len(buffer) < int(nvm.WriteSize/4) {
			buffer = append(buffer, 0) // 0 Pad it, to page align it
		}
		err = nvm.Cortex.WriteSeqAddr32(nvm.WriteAddress, buffer)
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

func (nvm *NVMFlash) Commit() error {
	fmt.Printf("Issuing Write Command: %x\n", (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMWriteCMD)
	err := nvm.Cortex.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMCMDOffSet, (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMWriteCMD)
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
	fmt.Printf("Issuing Erase Command: %x\n", (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMEraseCMD)
	err := nvm.Cortex.WriteAddr32(nvm.NVMControllerAddress+nvm.NVMCMDOffSet, (nvm.NVMCMDKey<<nvm.NVMCMDKeyPos)|nvm.NVMEraseCMD)
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
		val, err := nvm.ReadAddr32(0x41004010, 1)
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
		time.Sleep(time.Millisecond * 100)
	}

	// Clear Interrupt
	err := nvm.WriteAddr32(0x41004010, 0x1)
	if err != nil {
		return err
	}

	ti = time.Now()
	for {
		// Wait For interrupt to be cleared after you wrote to it
		val, err := nvm.ReadAddr32(0x41004010, 1)
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
		time.Sleep(time.Millisecond * 100)
	}

	//#######################################################

	//ti := time.Now()
	//
	//for {
	//	// Read Flag
	//	val, err := nvm.Cortex.ReadAddr32(nvm.NVMControllerAddress+0x10, 1)
	//	if err != nil {
	//		return err
	//	}
	//	// Need to shift it down according to its word alignment
	//	fmt.Printf("NVM Status Ready Address: %x\n", nvm.NVMControllerAddress+nvm.NVMReadyOffSet)
	//	fmt.Printf("NVM Status Ready: %x\n", val)
	//	fmt.Printf("NVM Status [Ready Mask: %x Ready Val: %x]\n", nvm.NVMReadyMask, nvm.NVMReadyVal)
	//	fmt.Printf("Debug attempts: %d\n", nvm.NVMReadyOffSet%4)
	//	// Need to shift it down according to its word alignment
	//	//val = val >> 16
	//	fmt.Printf("Second NVM Status Ready: %x\n", val)
	//
	//	// Bitwise Flag check
	//	if val&nvm.NVMReadyMask == nvm.NVMReadyVal {
	//		break
	//	}
	//
	//	// Timeout
	//	if time.Since(ti) > time.Second*1 {
	//		return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
	//	}
	//}
	//
	//// Clear Interrupt
	//err := nvm.WriteAddr32(nvm.NVMControllerAddress+0x10, 0x1)
	//if err != nil {
	//	return err
	//}
	//
	//ti = time.Now()
	//for {
	//	// Read Flag
	//	val, err := nvm.Cortex.ReadAddr32(nvm.NVMControllerAddress+0x10, 1)
	//	if err != nil {
	//		return err
	//	}
	//	// Need to shift it down according to its word alignment
	//	fmt.Printf("NVM Status Ready Address: %x\n", nvm.NVMControllerAddress+nvm.NVMReadyOffSet)
	//	fmt.Printf("NVM Status Ready: %x\n", val)
	//	fmt.Printf("NVM Status [Ready Mask: %x Ready Val: %x]\n", nvm.NVMReadyMask, nvm.NVMReadyVal)
	//	fmt.Printf("Debug attempts: %d\n", nvm.NVMReadyOffSet%4)
	//	// Need to shift it down according to its word alignment
	//	//val = val >> 16
	//	fmt.Printf("Second NVM Status Ready: %x\n", val)
	//
	//	// Bitwise Flag check
	//	if val&nvm.NVMReadyMask != nvm.NVMReadyVal {
	//		break
	//	}
	//
	//	// Timeout
	//	if time.Since(ti) > time.Second*1 {
	//		return fmt.Errorf("error: Atsame51.LoadProgram() timedout waiting for CMD clear")
	//	}
	//}

	return nil
}

func (nvm *NVMFlash) Configure() error {
	readVal, err := nvm.Cortex.ReadAddr32(nvm.NVMControllerAddress+nvm.NVMPARAMOffset, 1)
	if err != nil {
		return err
	}

	nvm.PageCount = (readVal & nvm.NVMPageCountMask) >> nvm.NVMPageCountPos
	pgSizeVal := (readVal & nvm.NVMPageSizeMask) >> nvm.NVMPageSizePos
	switch pgSizeVal {
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
