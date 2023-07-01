package cortexm4

import (
	"fmt"
	"github.com/goocd/goocd/protocols/cmsis"
)

// Cortex M4 specific Register data
const ()

type DAPTransferer interface {
	DAPTransfer(dapidx uint8, count uint8, data []byte) ([]byte, error)
}

type DAPTransferCoreAccess struct {
	DAPTransferer
}

func (d *DAPTransferCoreAccess) Configure() (err error) {
	// Figure out encoder stuff
	request := byte(cmsis.DebugPort | cmsis.Read | cmsis.PortRegister0)
	resp, err := d.DAPTransfer(0, 1, []byte{request})
	if err != nil {
		return err
	}
	fmt.Printf("%x", resp)
	return nil
}

func (t *DAPTransferCoreAccess) ReadAddr32(addr uint32, count int) (value uint32, err error) {
	return 0, nil
}
