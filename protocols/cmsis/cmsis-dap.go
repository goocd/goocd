package cmsis

import (
	"encoding/binary"
)

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

// DAP Response Status
const (
	DAP_OK    = 0x0
	DAP_Error = 0xFF
)

type CMSISDAP struct {
	ReadWriter ReadWriter
	Buffer     [512]byte // Pre-Allocated to avoid potential os allocation issues
}

func (c *CMSISDAP) Configure() {
	// Do stuff
}

// exposes CMSIS-DAP stuff that follows https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__Transfer.html

func (c *CMSISDAP) DAPInfo(info DAPInfo) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPInfoCMD, byte(info)})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPHostStatus(h HostType, s HostStatus) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPHostStatusCMD, byte(h), byte(s)})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPConnect(p ConnectPort) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPConnectCMD, byte(p)})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPDisconnect() error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPDisconnectCMD})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPWriteAbort(w uint32) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPWriteAbortCMD})
	binary.BigEndian.PutUint32(c.Buffer[3:], w)
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPDelay(d uint16) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPDelay
	binary.BigEndian.PutUint16(c.Buffer[2:], d)
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPResetTarget() error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPResetTarget})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWJPins(out byte, sel PinSelect, dur uint32) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPSWJPinsCMD, out, byte(sel)})
	binary.BigEndian.PutUint32(c.Buffer[4:], dur)
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWJClock(clock uint32) error {
	c.ZeroBuffer()
	c.Buffer[1] = DAPSWJClockCMD
	binary.BigEndian.PutUint32(c.Buffer[2:], clock)
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWJSequence(seq byte, data byte) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPSWJSequenceCMD, seq, data})
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWDConfigure(config byte) error {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPSWDConfigCMD, config})
	err := c.SendAndRead()
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
	err := c.SendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPTransfer(idx uint8, count uint8, request byte, data []uint32) ([]byte, error) {
	c.ZeroBuffer()
	copy(c.Buffer[1:], []byte{DAPTransferCMD, idx, count, request})
	if len(data) > 0 {
		for i, tempWord := range data {
			binary.BigEndian.PutUint32(c.Buffer[5+(4*i):], tempWord)
		}
	}
	//fmt.Printf("O: %x\n", c.Buffer[:32])
	err := c.SendAndRead()
	if err != nil {
		return nil, err
	}
	if !c.AckOK() {
		return nil, ErrBadAck{}
	}

	return c.Buffer[:], nil
}

// but it ALSO exposes methods that are friendly to the next layer - addressing CoreSight debug registers in this case

func (c *CMSISDAP) CoreSightDebugSend() error {
	return nil
}

func (c *CMSISDAP) ZeroBuffer() {
	for i := range c.Buffer {
		c.Buffer[i] = 0
	}
}

func (c *CMSISDAP) SendAndRead() error {
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

func (c *CMSISDAP) AckOK() bool {
	return c.Buffer[2] == 0x1
}

type ErrBadAck struct {
}

func (e ErrBadAck) Error() string {
	return "Err Bad Ack"
}
