package debugger

import (
	"fmt"
	"github.com/sstallion/go-hid"
)

var DebuggerMap = make(map[uint32]func() (Debugger, error))

var (
	availableID = uint32(0x0)
)

type Debugger interface {
	Program() error
}

// TODO: Add direct way of getting the correct debugger if we provide command line tooling for it

// GetDebugger is a generalized way to find any implemented debugger connected to the system
func GetDebugger() (Debugger, error) {
	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		return nil, err
	}

	if err := hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, findDebugger); err != nil {
		return nil, err
	}

	fnc, ok := DebuggerMap[availableID]
	if !ok || fnc == nil {
		return nil, fmt.Errorf("No Supported Debugger Found")
	}

	dbg, err := fnc()
	if err != nil {
		return nil, err
	}

	return dbg, nil
}

// findDebugger is the callback functino provided to the enumerator in order to more closely identify a valid usb connection.
// Todo: Potentially interface these away if we need to support more than just USB
func findDebugger(info *hid.DeviceInfo) error {
	// Check VendorID + ProductID pairings for all usb devices
	debuggerID := uint32(info.VendorID)<<16 | uint32(info.ProductID)
	_, ok := DebuggerMap[debuggerID]
	if ok {
		availableID = debuggerID
	}
	return nil
}
