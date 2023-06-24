package debugger

import "log"

const (
	vendorID  = 0x3EB
	productID = 0x2141
)

func init() {
	DebuggerMap[uint32(vendorID<<16|productID)] = NewAtmelICE
}

func NewAtmelICE() (Debugger, error) {
	log.Printf("GOT HERE: Set Up Atmel Ice")
	ice := new(AtmelICE)
	// Initialize Protocols it implements
	// Initialize HID Device
	// New CMSIS struct on ICE
	// Provide device implementation to CMSIS as ReadWriter interface
	return ice, nil
}

type AtmelICE struct {
}

func (a *AtmelICE) Program() error {
	log.Printf("GOT HERE: Atmel Ice Programming Start")
	return nil
}
