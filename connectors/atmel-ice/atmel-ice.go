package atmel_ice

import (
	"encoding/binary"
	"github.com/goocd/goocd/protocols/cmsis"
)

const (
	vendorID  = 0x3EB
	productID = 0x2141
)

//func init() {
//	// Provide the vendorID + ProductID Pattern that would be expected for the atmel ICE
//	DebuggerMap[uint32(vendorID<<16|productID)] = NewAtmelICE
//}

//// NewAtmelICE returns the interface of the AtmelICE struct to provide simple exposure of basic methods.
//// You can always use the struct directly if more refined control is needed
//func NewAtmelICE() (Debugger, error) {
//	log.Printf("GOT HERE: Set Up Atmel Ice")
//	ice := new(AtmelICE)
//	// Initialize Protocols it implements
//	// Initialize HID Device
//	// New CMSIS struct on ICE
//	// Provide device implementation to CMSIS as ReadWriter interface
//	return ice, nil
//}

// AtmelICE Will Contain nil pointers to the implementations it supports. Initialized when the NewAtmelICE function is called
type AtmelICE struct {
	*cmsis.CMSISDAP
}

func (a *AtmelICE) Configure() error {
	err := a.DAPConnect(cmsis.SWD)
	if err != nil {
		return err
	}

	//fmt.SPrintf("%x\n", a.Buffer[:32])

	err = a.DAPSWDConfigure(0x0)
	if err != nil {
		return err
	}
	//fmt.Printf("%x\n", a.Buffer[:32])

	err = a.DAPSWJSequence(0x88, 0xff) // Might be a byte slice???
	if err != nil {
		return err
	}
	//fmt.Printf("%x\n", a.Buffer[:32])

	err = a.DAPSWJClock(0x80841e00)
	if err != nil {
		return err
	}
	//fmt.Printf("%x\n", a.Buffer[:32])
	err = a.DAPTransferConfigure(0x0, 0x4000, 0x0)
	if err != nil {
		return err
	}
	//fmt.Printf("%x\n", a.Buffer[:32])
	var data []byte
	_ = data
	//Todo: Use Data values correctly to identify errors
	data, err = a.DAPTransfer(0, 0x1, 0x2, nil)
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x00000000, 0x4200000, 0x50060400, 0x5006})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x0, 0x20040, 0x40650, 0x650})
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x1, 0x6, nil)
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x3, 0x6, []uint32{0x04000000, 0x50060000})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x3, 0x6, []uint32{0x04, 0x650})
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x3, 0x8, []uint32{0xF0000000, 0x0F0E0000})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x3, 0x8, []uint32{0xF0, 0xE0F})
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x8, 0x8, []uint32{0x00000000, 0x01200000, 0xA2050000, 0x0000030E, 0x08F00000, 0x00070E00})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x8, 0x8, []uint32{0x0, 0x2001, 0x05A2, 0xE030000, 0xF008, 0xE0700})
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x00000000, 0x01220000, 0xA20500ED, 0x00E00F0E})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x0, 0x2201, 0xED0005A2, 0xE0F00E00})
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x3, 0x5, []uint32{0x40EF00E0, 0x0F0E0000})
	//if err != nil {
	//	return err
	//}

	data, err = a.DAPTransfer(0, 0x3, 0x5, []uint32{0xE000EF40, 0xE0F})
	if err != nil {
		return err
	}

	return nil
}

// Program will actually accept a file stream to be written into board memory. Probably come from a form of parser
func (a *AtmelICE) ReadAddr32(addr uint32) (uint32, error) {
	data, err := a.DAPTransfer(0, 0x3, 0x5, []uint32{addr, 0xE0F})
	if err != nil {
		return 0, err
	}

	value := binary.BigEndian.Uint32(data[3:])
	return value, nil
}

func (a *AtmelICE) WriteAddr32(addr, value uint32) error {
	_, err := a.DAPTransfer(0, 0x2, 0x5, []uint32{addr, (value<<8)&0xFFFFFF00 | 0xD, (value >> 24) & 0xFF})
	if err != nil {
		return err
	}
	return nil
}
