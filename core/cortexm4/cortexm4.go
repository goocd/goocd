package cortexm4

import (
	"encoding/binary"
	"fmt"
	"github.com/goocd/goocd/protocols/cmsisdap"
)

// Cortex M4 specific Register data
const ()

// Debug Port IDCODE Masks
const (
	VersionMask    = 0xF0000000
	PartNumberMask = 0xFFFF000
	DesignerMask   = 0xFFE
)

const (
	FlashPatchCTRLRegister = 0xE0002000
	FlashPatchCTRLEnable   = 0x1
	FlashPatchCTRLKey      = 0x2
)

const (
	DebugHaltingControlStatusRegister = 0xE000EDF0
	DebugHaltingControlStatusKey      = 0xA05F0000
	DebugHaltingControlStatusEnable   = 0x1
	DebugHaltingControlStatusHalt     = 0x2
	DebugHaltingControlStatusStep     = 0x3
)

// Debug Port CTRL Register Mappings
const (
	CSYSPWRUPREQEnable  = 0x40000000
	CSYSPWRUPREQDisable = 0x0

	CDBGPWRUPREQEnable  = 0x10000000
	CDBGPWRUPREQDisable = 0x0

	CDBGRSTREQEnable  = 0x4000000
	CDBGRSTREQDisable = 0x0

	TRNCNTMask0x1FF000
	MASKLANEMask = 0xF00

	ORUNDETECTEnable
	ORUNDETECTDisable
)

// Useful Consts
const (
	APSELPOS     = 0x24
	APBANKSELPOS = 0x4
)

const (
	AHBAPDAPEnable   = 0x40
	AHBAPEnableDebug = 0x20000000
	DataSizeuint8    = 0x0
	DataSizeuint16   = 0x1
	DataSizeuint32   = 0x2
)

// Port Banks
const (
	Bank0 = uint32(iota)
	Bank1
	Bank2
	Bank3
	Bank4
	Bank5
	Bank6
	Bank7
	Bank8
	Bank9
	BankA
	BankB
	BankC
	BankD
	BankE
	BankF
	BankPos = 0x4
)

type DAPTransferer interface {
	DAPTransfer(dapidx uint8, count uint8, data []byte) ([]byte, error)
}

type DAPTransferCoreAccess struct {
	DAPTransferer
	encodingBuffer [512]byte
}

type request struct {
	requestByte byte
	payload     uint32
}

func (d *DAPTransferCoreAccess) Configure() (err error) {
	// Read Debug Port IDCODE
	resp, err := d.DAPTransfer(0, 1, d.encodeDAPRequest([]request{
		{
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister0),
		},
	}))
	if err != nil {
		return err
	}
	_ = resp
	// Todo: Validate
	//fmt.Printf("%x\n", resp)

	// Debug Power Enable
	resp, err = d.DAPTransfer(0, 5, d.encodeDAPRequest([]request{
		{
			// Clear out the Selections Registers to known state
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
		},
		{
			// Enable Debug Power
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     CSYSPWRUPREQEnable | CDBGPWRUPREQEnable | 0x20,
		},
		{
			// Read Back Register for validation
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister4),
		},
		{
			// Enable Debug Power
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     CSYSPWRUPREQEnable | CDBGPWRUPREQEnable,
		},
		{
			// Read Back Register for validation
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister4),
		},
	}))
	if err != nil {
		return err
	}

	// Todo: Validate
	//fmt.Printf("%x\n", resp)

	resp, err = d.DAPTransfer(0, 1, d.encodeDAPRequest([]request{
		{
			// Read Back Register for validation
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister4),
		},
	}))
	if err != nil {
		return err
	}

	// Todo: Validate
	//fmt.Printf("%x\n", resp)

	resp, err = d.DAPTransfer(0, 3, d.encodeDAPRequest([]request{
		{
			// Read Back Register for validation
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister4),
		},
		{
			// Enable Debug Power
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     CSYSPWRUPREQEnable | CDBGPWRUPREQEnable,
		},
		{
			// Read Back Register for validation
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Read | cmsisdap.PortRegister4),
		},
	}))
	if err != nil {
		return err
	}

	// Todo: Validate
	//fmt.Printf("%x\n", resp)
	resp, err = d.DAPTransfer(0, 2, d.encodeDAPRequest([]request{
		{
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
			payload:     BankF << BankPos,
		},
		{
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Read | cmsisdap.PortRegisterC),
		},
	}))
	if err != nil {
		return err
	}
	// Todo: Validate
	//fmt.Printf("%x\n", resp)
	resp, err = d.DAPTransfer(0, 2, d.encodeDAPRequest([]request{
		{
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
			payload:     Bank0,
		},
		{
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister0),
			payload:     AHBAPEnableDebug | AHBAPDAPEnable | DataSizeuint32,
		},
	}))
	if err != nil {
		return err
	}
	// Todo: Validate
	//fmt.Printf("%x\n", resp)

	return nil
}

func (d *DAPTransferCoreAccess) ReadAddr32(addr uint32, count int) (value uint32, err error) {
	// Todo: Take into account the count and continuous read
	resp, err := d.DAPTransfer(0, 3, d.encodeDAPRequest([]request{
		{
			// Clear out the Selections Registers to known state
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
		},
		{
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     addr,
		},
		{
			// Read Data Register
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Read | cmsisdap.PortRegisterC),
		},
	}))
	if err != nil {
		return 0, err
	}
	// Todo: Validate
	//fmt.Printf("%x", resp)
	// FixMe: Do this parsing properly
	return binary.BigEndian.Uint32(resp[3:7]), nil
}

func (d *DAPTransferCoreAccess) WriteAddr32(addr, value uint32) error {
	_, err := d.DAPTransfer(0, 3, d.encodeDAPRequest([]request{
		{
			// Clear out the Selections Registers to known state
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
		},
		{
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     addr,
		},
		{
			// Read Data Register
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegisterC),
			payload:     value,
		},
	}))
	if err != nil {
		return err
	}
	return nil
}

func (d *DAPTransferCoreAccess) WriteSeqAddr32(addr uint32, value []uint32) error {
	if len(value) > 101 {
		return fmt.Errorf("error: DAPTransferCoreAccess.WriteSeqAddr32() len of values is too large")
	}

	requestBuffer := make([]request, 2, len(value)+2)
	requestBuffer[0] = request{requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8)}
	requestBuffer[1] = request{requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister4), payload: addr}

	for _, val := range value {
		requestBuffer = append(requestBuffer, request{requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegisterC), payload: val})
	}

	_, err := d.DAPTransfer(0, uint8(len(requestBuffer)), d.encodeDAPRequest(requestBuffer))
	if err != nil {
		return err
	}

	return nil
}

func (d *DAPTransferCoreAccess) WriteTransfer32(port, portRegister byte, value uint32) error {
	_, err := d.DAPTransfer(0, 1, d.encodeDAPRequest([]request{
		{
			requestByte: port | cmsisdap.Write | portRegister,
			payload:     value,
		},
	}))
	if err != nil {
		return err
	}
	return nil
}

func (d *DAPTransferCoreAccess) ReadTransfer32(port, portRegister byte) (uint32, error) {
	resp, err := d.DAPTransfer(0, 1, d.encodeDAPRequest([]request{
		{
			requestByte: port | cmsisdap.Read | portRegister,
		},
	}))
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(resp[3:7]), nil
}

func (d *DAPTransferCoreAccess) encodeDAPRequest(requests []request) []byte {
	idx := 0
	for _, req := range requests {
		d.encodingBuffer[idx] = req.requestByte
		idx++

		if req.requestByte&cmsisdap.Read > 0 {
			// No more is needed for reads
			continue
		}
		binary.LittleEndian.PutUint32(d.encodingBuffer[idx:], req.payload)
		idx += 4
	}
	return d.encodingBuffer[:idx]
}

func (d *DAPTransferCoreAccess) Halt() error {
	_, err := d.DAPTransfer(0, 3, d.encodeDAPRequest([]request{
		{
			// Clear out the Selections Registers to known state
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
		},
		{
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister4),
			payload:     DebugHaltingControlStatusRegister,
		},
		{
			// Read Data Register
			requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister0),
			payload:     DebugHaltingControlStatusKey | DebugHaltingControlStatusHalt | DebugHaltingControlStatusEnable,
		},
	}))
	if err != nil {
		return err
	}
	return nil
}
