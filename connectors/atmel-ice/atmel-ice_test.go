package atmel_ice

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/goocd/goocd/protocols/cmsis"
	"github.com/sstallion/go-hid"
	"log"
	"testing"
)

func TestAtmelICE(t *testing.T) {
	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	//if err := hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, logHID); err != nil {
	//	t.Fatal(err)
	//}

	// Open the device using the VID and PID.
	d, err := hid.OpenFirst(0x03eb, 0x2141) // Atmel-ICE VIP & PID
	if err != nil {
		t.Fatal(err)
	}

	cms := &cmsis.CMSISDAP{ReadWriter: d}
	ice := &AtmelICE{CMSISDAP: cms}
	err = ice.Configure()
	if err != nil {
		t.Fatal(err)
	}
	var val uint32
	val, err = ice.ReadAddr32(0x008061FC)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("I: %x\n", val)

	val, err = ice.ReadAddr32(0x00806010)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("I: %x\n", val)

	val, err = ice.ReadAddr32(0x00806014)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("I: %x\n", val)

	val, err = ice.ReadAddr32(0x00806018)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("I: %x\n", val)

	err = ice.WriteAddr32(0x20030000, 0xA501)
	if err != nil {
		t.Fatal(err)
	}

	if err := d.Close(); err != nil {
		log.Fatal(err)
	}
	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}

func TestBinFlipping(t *testing.T) {
	cms := &cmsis.CMSISDAP{}
	ice := &AtmelICE{CMSISDAP: cms}
	err := ice.WriteAddr32(0x41004004, 0xA501)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ice.ReadAddr32(0x41004004)
	if err != nil {
		t.Fatal(err)
	}

}

func TestConvertHexString(t *testing.T) {
	value := "0x806018"
	value = value[2:]
	dst := make([]byte, hex.DecodedLen(len(value)))
	n, err := hex.Decode(dst, []byte(value))
	if err != nil {
		t.Fatal(err)
	}
	if n > 4 {
		t.Fatal("IDK WHAT HAPPENED")
	}
	res := make([]byte, 4) // Gaurentee 4 bytes
	copy(res[4-len(dst):], dst)
	addr := binary.BigEndian.Uint32(res)
	t.Logf("%x", addr)
}
