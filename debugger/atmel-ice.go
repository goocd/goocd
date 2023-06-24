package debugger

import "log"

const (
	vendorID  = 0x3EB
	productID = 0x2141
)

func init() {
	// Provide the vendorID + ProductID Pattern that would be expected for the atmel ICE
	DebuggerMap[uint32(vendorID<<16|productID)] = NewAtmelICE
}

// NewAtmelICE returns the interface of the AtmelICE struct to provide simple exposure of basic methods.
// You can always use the struct directly if more refined control is needed
func NewAtmelICE() (Debugger, error) {
	log.Printf("GOT HERE: Set Up Atmel Ice")
	ice := new(AtmelICE)
	// Initialize Protocols it implements
	// Initialize HID Device
	// New CMSIS struct on ICE
	// Provide device implementation to CMSIS as ReadWriter interface
	return ice, nil
}

// AtmelICE Will Contain nil pointers to the implementations it supports. Initialized when the NewAtmelICE function is called
type AtmelICE struct {
}

// Program will actually accept a file stream to be written into board memory. Probably come from a form of parser
func (a *AtmelICE) Program() error {
	log.Printf("GOT HERE: Atmel Ice Programming Start")
	return nil
}
