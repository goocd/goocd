package cortexm4

import (
	"encoding/binary"
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
