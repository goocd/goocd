package cmsis

import (
	"encoding/binary"
)

/*
  All Specifications have been grabbed from https://arm-software.github.io/CMSIS_5/DAP/html/index.html and are designed to be extremely basic implementations with as minimal understand of the protocol as needed
*/

type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

// DAP Commands
const (
	DAPInfoCMD           = 0x0
	DAPHostStatusCMD     = 0x1
	DAPConnectCMD        = 0x2
	DAPDisconnectCMD     = 0x3
	DAPTransferConfigCMD = 0x4
	DAPTransferCMD       = 0x5
	DAPWriteAbortCMD     = 0x8
	DAPDelay             = 0x9
	DAPResetTarget       = 0xA
	DAPSWJPinsCMD        = 0x10
	DAPSWJClockCMD       = 0x11
	DAPSWJSequenceCMD    = 0x12
	DAPSWDConfigCMD      = 0x13
	DAPSWDSequenceCMD    = 0x1D
)

// DAP Connect
const (
	DefaultPort = 0x0
	SWDPort     = 0x1
	JTAGPort    = 0x2
)

// DAP Host/Status
const (
	HostConnect = 0x0
	HostRunning = 0x1

	StatusOff = 0x0
	StatusOn  = 0x1
)

// PinOut Mask
const (
	PinMaskSWCLKTCK = 0x1
	PinMaskSWDIOTMS = 0x2
	PinMaskTDI      = 0x4
	PinMaskTDO      = 0x8
	PinMaskNTRST    = 0x20
	PinMaskNReset   = 0x80
)

// DAP Info
const (
	VendorName              = 0x1
	ProductName             = 0x2
	SerialNumber            = 0x3
	CMSISDAPProtocolVersion = 0x4
	TargetDeviceVendor      = 0x5
	TargetDeviceName        = 0x6
	TargetBoardvendor       = 0x7
	TargetBoardName         = 0x8
	ProductFirmwareVersion  = 0x9
	Capabilities            = 0xF0
	TestDomainTimer         = 0xF1
	UARTReceiveBufferSize   = 0xFB
	UARTTransmiteBufferSize = 0xFC
	SWOTraceBufferSize      = 0xFD
	PacketCount             = 0xFE
	PacketSize              = 0xFF
)

// DAP Transfer
const (
	DebugPort  = 0x0
	AccessPort = 0x1

	Read  = 0x2
	Write = 0x0

	PortRegister0 = 0x0
	PortRegister4 = 0x4
	PortRegister8 = 0x8
	PortRegisterC = 0xC

	ValueMatch = 0x10
	MatchMask  = 0x20
	TimeStamp  = 0x80
)

// DAP Response Status
const (
	DAP_OK    = 0x0
	DAP_Error = 0xFF
)

type CMSISDAP struct {
	ReadWriter ReadWriter
	Buffer     [512]byte // Pre-Allocated to avoid potential os allocation issues
}

func (c *CMSISDAP) Configure() error {
	// Do stuff
	return nil
}

// exposes CMSIS-DAP stuff that follows https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__Transfer.html

func (c *CMSISDAP) DAPInfo(info byte) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPInfoCMD
	c.Buffer[2] = info
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPHostStatus(host byte, status byte) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPHostStatusCMD
	c.Buffer[2] = host
	c.Buffer[3] = status
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPConnect(port byte) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPConnectCMD
	c.Buffer[2] = port
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPDisconnect() error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPDisconnectCMD
	err := c.sendAndRead()
	if err != nil {
		return err
	}

	if c.Buffer[1] != DAP_OK {
		return ErrBadDAPResponseStatus{}
	}

	return nil
}

func (c *CMSISDAP) DAPWriteAbort(dapindex byte, word uint32) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPWriteAbortCMD
	c.Buffer[2] = dapindex // Note: Ignored when using SWD
	binary.BigEndian.PutUint32(c.Buffer[3:], word)
	err := c.sendAndRead()
	if err != nil {
		return err
	}

	if c.Buffer[1] != DAP_OK {
		return ErrBadDAPResponseStatus{}
	}
	return nil
}

func (c *CMSISDAP) DAPDelay(d uint16) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPDelay
	binary.BigEndian.PutUint16(c.Buffer[2:], d)
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	if c.Buffer[1] != DAP_OK {
		return ErrBadDAPResponseStatus{}
	}

	return nil
}

func (c *CMSISDAP) DAPResetTarget() error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPResetTarget
	err := c.sendAndRead()
	if err != nil {
		return err
	}

	if c.Buffer[1] != DAP_OK {
		return ErrBadDAPResponseStatus{}
	}

	return nil
}

func (c *CMSISDAP) DAPSWJPins(out byte, sel byte, waitDur uint32) (byte, error) {
	c.ZeroBuffer()
	c.Buffer[1] = DAPSWJPinsCMD
	c.Buffer[2] = out
	c.Buffer[3] = sel
	binary.BigEndian.PutUint32(c.Buffer[4:], waitDur)
	err := c.sendAndRead()
	if err != nil {
		return 0, err
	}

	return c.Buffer[1], nil
}

func (c *CMSISDAP) DAPSWJClock(clock uint32) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPSWJClockCMD
	binary.BigEndian.PutUint32(c.Buffer[2:], clock)
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWJSequence(seq byte, data byte) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPSWJSequenceCMD
	c.Buffer[2] = seq
	c.Buffer[3] = data
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWDConfigure(config byte) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPSWDConfigCMD
	c.Buffer[2] = config
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPTransferConfigure(cycles byte, wait uint16, match uint16) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPTransferConfigCMD
	c.Buffer[2] = cycles
	binary.BigEndian.PutUint16(c.Buffer[3:], wait)
	binary.BigEndian.PutUint16(c.Buffer[5:], match)
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

// DAPTransfer implements the DAP transfer protocol to the spec. The exact sequence and endianess of the data is generally MCU specific so it's up to the caller to pass in data structed accordingly.
func (c *CMSISDAP) DAPTransfer(dapidx uint8, count uint8, data []byte) ([]byte, error) {
	c.ZeroBuffer()
	c.Buffer[1] = DAPTransferCMD
	c.Buffer[2] = dapidx // Note: Ignored when using SWD
	c.Buffer[3] = count
	copy(c.Buffer[4:], data)
	err := c.sendAndRead()
	if err != nil {
		return nil, err
	}

	return c.Buffer[:], nil
}

func (c *CMSISDAP) ZeroBuffer() {
	for i := range c.Buffer {
		c.Buffer[i] = 0
	}
}

// sendAndRead wraps the actual read/writes to the underlying device. The Read always follows the Write to accept the response from the connected device.
func (c *CMSISDAP) sendAndRead() error {
	_, err := c.ReadWriter.Write(c.Buffer[:])
	if err != nil {
		return err
	}
	c.ZeroBuffer()
	_, err = c.ReadWriter.Read(c.Buffer[:])
	if err != nil {
		return err
	}

	return nil
}

type ErrBadDAPResponseStatus struct {
}

func (e ErrBadDAPResponseStatus) Error() string {
	return "Error: CMSIS DAP_Repsonse Status: Not ok"
}
