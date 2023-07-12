package cortexm4

import (
	"encoding/binary"
	"goocd/protocols/cmsisdap"
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

// Configure
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

	// Debug Power Enable
	resp, err = d.DAPTransfer(0, 5, d.encodeDAPRequest([]request{
		{
			// Clear out the Selections Registers to known state
			requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8),
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

	// Enable Debug Power. There's a 100% chance this can be simplified with a loop and verified rather than just doing this repeatedly in sequence
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

	// Read Access Port IDCODE
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

	// Initialize Debugging in the AHB-AP with DataSize 32 bit
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

	return nil
}

// ReadAddr32 does direct memory access and reads a 32 bit value from a provided address
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
	// FixMe: Do this parsing properly
	return binary.LittleEndian.Uint32(resp[3:7]), nil
}

// WriteAddr32 is a simple way to write a value to a given address
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

// WriteSeqAddr32 does sequential write transactions to the AHB-AccessPort address provided based off how many values are in the buffer.
func (d *DAPTransferCoreAccess) WriteSeqAddr32(addr uint32, value []uint32) error {
	//fmt.Printf("Seq Write Request with Len: %d at Address: %x\n", len(value), addr)
	requestBuffer := make([]request, 0, len(value)+2)
	requestBuffer = append(requestBuffer, request{requestByte: byte(cmsisdap.DebugPort | cmsisdap.Write | cmsisdap.PortRegister8)})
	requestBuffer = append(requestBuffer, request{requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegister4), payload: addr})
	for _, val := range value {
		requestBuffer = append(requestBuffer, request{requestByte: byte(cmsisdap.AccessPort | cmsisdap.Write | cmsisdap.PortRegisterC), payload: val})
		//fmt.Printf("Appending Val: %x\n", val)
		if len(requestBuffer) >= 66 {
			//fmt.Printf("Sending Transfer with Len: %d\n", len(requestBuffer))
			_, err := d.DAPTransfer(0, uint8(len(requestBuffer)), d.encodeDAPRequest(requestBuffer))
			if err != nil {
				return err
			}
			requestBuffer = requestBuffer[:0]
		}
	}

	if len(requestBuffer) > 0 {
		//fmt.Printf("Sending Transfer with Len: %d\n", len(requestBuffer))
		_, err := d.DAPTransfer(0, uint8(len(requestBuffer)), d.encodeDAPRequest(requestBuffer))
		if err != nil {
			return err
		}
		requestBuffer = requestBuffer[:0]
	}

	return nil
}

// WriteTransfer32 A simple way to abstract doing a single write transaction rather than a complete write which does multiple commands at once
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

// ReadTransfer32 A simple Way to abstract doing a single transaction rather than a complete Read which alters a few other registers.
// Always returns a uint32, up to the Caller to ensure they read the correct value out of it.
func (d *DAPTransferCoreAccess) ReadTransfer32(port, portRegister byte) (uint32, error) {
	resp, err := d.DAPTransfer(0, 1, d.encodeDAPRequest([]request{
		{
			requestByte: port | cmsisdap.Read | portRegister,
		},
	}))
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(resp[3:7]), nil
}

// encodeDAPRequest Simple way to just simplify how we use DAP Transfer requests.
// TODO: Decide if it makes sense for this and the type to live here rather than in CMSISDAP
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
		//fmt.Printf("Request: %x, PayLoad: %x\n", req.requestByte, req.payload)
		//randoBuf := make([]byte, 4)
		//binary.LittleEndian.PutUint32(randoBuf, req.payload)
		//fmt.Printf("Little Endian conversion: %x\n", randoBuf)
		idx += 4
	}
	return d.encodingBuffer[:idx]
}

// Halt access the cortex DHCSR register and writes the Halt bits according to CorextM4 specifications
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
