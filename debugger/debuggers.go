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

func findDebugger(info *hid.DeviceInfo) error {
	debuggerID := uint32(info.VendorID)<<16 | uint32(info.ProductID)
	_, ok := DebuggerMap[debuggerID]
	if ok {
		// Sanity Checks since Vendor ID + Product ID aren't guaranteed
		availableID = debuggerID
	}
	return nil
}
