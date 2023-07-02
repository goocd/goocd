package cmsisdap

import (
	"encoding/binary"
)

/*
  All Specifications have been grabbed from https://arm-software.github.io/CMSIS_5/DAP/html/index.html
  This is designed to be an extremely basic implementations and allow users to just pass requests based on the spec.

  Returned
*/

type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type CMSISDAP struct {
	ReadWriter ReadWriter
	Buffer     [512]byte // Pre-Allocated to avoid potential os allocation issues
}

func (c *CMSISDAP) DAPInfo(info byte) error {
	c.zeroBuffer()
	c.Buffer[1] = DAPInfoCMD
	c.Buffer[2] = info
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPHostStatus(host byte, status byte) error {
	c.zeroBuffer()
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
	c.zeroBuffer()
	c.Buffer[1] = DAPConnectCMD
	c.Buffer[2] = port
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPDisconnect() error {
	c.zeroBuffer()
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
	c.zeroBuffer()
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
	c.zeroBuffer()
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
	c.zeroBuffer()
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
	c.zeroBuffer()
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
	c.zeroBuffer()
	c.Buffer[1] = DAPSWJClockCMD
	binary.BigEndian.PutUint32(c.Buffer[2:], clock)
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWJSequence(seq byte, data []byte) error {
	c.zeroBuffer()
	c.Buffer[1] = DAPSWJSequenceCMD
	c.Buffer[2] = seq
	copy(c.Buffer[3:], data)
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPSWDConfigure(config byte) error {
	c.zeroBuffer()
	c.Buffer[1] = DAPSWDConfigCMD
	c.Buffer[2] = config
	err := c.sendAndRead()
	if err != nil {
		return err
	}
	return nil
}

func (c *CMSISDAP) DAPTransferConfigure(cycles byte, wait uint16, match uint16) error {
	c.zeroBuffer()
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
	c.zeroBuffer()
	c.Buffer[1] = DAPTransferCMD
	c.Buffer[2] = dapidx // Note: Ignored when using SWD
	c.Buffer[3] = count
	copy(c.Buffer[4:], data)
	err := c.sendAndRead()
	if err != nil {
		return nil, err
	}

	if c.Buffer[2] != 0x1 {
		return nil, ErrBadDAPResponseStatus{}
	}

	return c.Buffer[:32], nil
}

func (c *CMSISDAP) zeroBuffer() {
	for i := range c.Buffer {
		c.Buffer[i] = 0
	}
}

// sendAndRead wraps the actual read/writes to the underlying device. The Read always follows the Write to accept the response from the connected device.
func (c *CMSISDAP) sendAndRead() error {
	//fmt.Printf("Out: %x\n", c.Buffer[:32])
	_, err := c.ReadWriter.Write(c.Buffer[:])
	if err != nil {
		return err
	}
	c.zeroBuffer()
	_, err = c.ReadWriter.Read(c.Buffer[:])
	if err != nil {
		return err
	}
	//fmt.Printf("In:  %x\n", c.Buffer[:32])
	return nil
}

type ErrBadDAPResponseStatus struct {
}

func (e ErrBadDAPResponseStatus) Error() string {
	return "Error: CMSIS DAP_Repsonse Status: Not ok"
}
