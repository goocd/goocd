package usbhid

import (
	"github.com/sstallion/go-hid"
	"log"
)

// Wrapper for the hid library

type HidDevice struct {
	*hid.Device
}

func OpenFirstHid(vendorid, productid uint16) (*HidDevice, error) {

	if err := hid.Init(); err != nil {
		return nil, err
	}

	// Open the device using the VID and PID.
	d, err := hid.OpenFirst(vendorid, productid) // Atmel-ICE VIP & PID
	if err != nil {
		if hidErr := hid.Exit(); hidErr != nil {
			log.Printf("Failed To close HID: %+v", hidErr)
		}
		return nil, err
	}

	dev := &HidDevice{
		d,
	}

	return dev, nil
}

func (d *HidDevice) CleanUp() error {
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
	return nil
}
