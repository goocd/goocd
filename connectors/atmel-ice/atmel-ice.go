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

	data, err = a.DAPTransfer(0, 0x2, 0x8, []uint32{0x00000000, 0x04200000, 0x50060400, 0x5006})
	//[0x0, 0x0, 0x0, 0x0], [0x4], [0x50, 0x00, 0x00, 0x20], [0x06], [0x06, 0x50, 0x00, 0x04]
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x1, 0x6, nil)
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x3, 0x6, []uint32{04000000, 0x50060000})
	//	[0x0, 0x0, 0x0, 0x4], [0x50], [0x0, 0x0, 0x0, 0x6]
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x3, 0x8, []uint32{0xF0000000, 0x0F000000})
	//[0x0, 0x0, 0x0, 0xF0] [0xF] // Ends with Read
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x8, 0x8, []uint32{0x00000000, 0x01200000, 0xA2050000, 0x0000030E, 0x08F00000, 0x00070E00})
	//[0x0, 0x0, 0x0, 0x0], [0x1] , [0xA2 0x00, 0x00, 0x20] , [0x5], [0x00, 0x00, 0x00, 0x00], [0x3], [0x0, 0xF0, 0x08, 0x0E], [0x0], [ 0x0, 0x0E, 0x07, 0x00]
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x00000000, 0x01220000, 0xA20500ED, 0x00E00F})
	//                                               [0x0, 0x0, 0x0, 0x0], [0x1], [0xA2, 0x00, 0x00, 0x22], [0x5], [0xE0, 0x00, 0xED, 0x00]  [0xF] //  Ends with Debug Read
	if err != nil {
		return err
	}

	data, err = a.DAPTransfer(0, 0x3, 0x5, []uint32{0x40EF00E0, 0x0F0E0000})
	//[0xE0, 0x00, 0xEF, 0x40] [0xF]
	if err != nil {
		return err
	}

	//data, err = a.DAPTransfer(0, 0x1, 0x2, nil)
	//if err != nil {
	//	return err
	//}
	//// READ DP reg 0 0
	//fmt.Printf("%x\n", data[:24])
	//// Write DP Reg 8 0
	//// Write DP Reg 4 0x50000020
	//// Read DP Reg 4 Match Value? [ 0x6500004 ]
	//data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x0, 0x20040, 0x40650, 0x650})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//// Read DP Reg 4 0
	//data, err = a.DAPTransfer(0, 0x1, 0x6, nil)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//// Read DP Reg 4 Match value?  [ 0x6500004 ]
	//data, err = a.DAPTransfer(0, 0x3, 0x6, []uint32{0x04, 0x650})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//// Write DP Reg 8 F0
	//// Write AP Reg C 0
	//data, err = a.DAPTransfer(0, 0x3, 0x8, []uint32{0xF0, 0xF})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//
	//data, err = a.DAPTransfer(0, 0x8, 0x8, []uint32{0x0, 0x2001, 0x05A2, 0xE030000, 0xF008, 0xE0700})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x0, 0x2201, 0xED0005A2, 0xE0F00E00})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	//data, err = a.DAPTransfer(0, 0x3, 0x5, []uint32{0xE000EF40, 0xF})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%x\n", data[:24])
	/*
			//Todo: Use Data values correctly to identify errors
			data, err = a.DAPTransfer(0, 0x1, 0x2, nil)
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x2, 0x8, []uint32{0x00000000, 0x04 20 00 00, 0x5 0 06 0400, 0x5006})
															[0x0, 0x0, 0x0, 0x0], [0x4], [0x50, 0x00, 0x00, 0x20], [0x06], [0x06, 0x50, 0x00, 0x04]
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x1, 0x6, nil)
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x3, 0x6, []uint32{04000000, 0x50060000})
													[0x0, 0x0, 0x0, 0x4], [0x50], [0x0, 0x0, 0x0, 0x6]
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x3, 0x8, []uint32{0xF0000000, 0x0F000000})
															[0x0, 0x0, 0x0, 0xF0] [0xF] // Ends with Read
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x8, 0x8, []uint32{0x00000000, 0x0120 00 00, 0xA205 00 00, 0x00 00 030E, 0x08 F0 00 00, 0x00070E 00})
															[0x0, 0x0, 0x0, 0x0], [0x1] , [0xA2 0x00, 0x00, 0x20] , [0x5], [0x00, 0x00, 0x00, 0x00], [0x3], [0x0, 0xF0, 0x08, 0x0E], [0x0], [ 0x0, 0x0E, 0x07, 0x00]
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x5, 0x8, []uint32{0x00000000, 0x01220000, 0xA20500ED, 0x00E0 0F})
			//                                               [0x0, 0x0, 0x0, 0x0], [0x1], [0xA2, 0x00, 0x00, 0x22], [0x5], [0xE0, 0x00, 0xED, 0x00]  [0xF] //  Ends with Debug Read
			if err != nil {
				return err
			}

			data, err = a.DAPTransfer(0, 0x3, 0x5, []uint32{0x40EF00E0, 0x0F0E0000})
													[0xE0, 0x00, 0xEF, 0x40] [0xF]
			if err != nil {
				return err
			}

		help

			0x3 [0x5]: [0x0...], [0xB]:[0x0..], [0xE]

	*/

	return nil
}

// Wrie AP Reg 4
//

// Program will actually accept a file stream to be written into board memory. Probably come from a form of parser
func (a *AtmelICE) ReadAddr32(addr uint32) (uint32, error) {
	data, err := a.DAPTransfer(0, 0x3, 0x5, []uint32{addr, 0xF})
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

/*



 */
