package atmel_ice

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/goocd/goocd/protocols/cmsis"
	"github.com/sstallion/go-hid"
	"log"
	"testing"
	"time"
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

	halt(t, ice)
	resume(t, ice)
	halt(t, ice)
	resume(t, ice)
	//data, err := ice.DAPTransfer(0x0, 0x3, 0x8, []uint32{0x1000F0, 0x0F})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])
	//
	//data, err = ice.DAPTransfer(0x0, 0x3, 0x8, []uint32{0x0000F0, 0x0F})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])
	//
	//data, err = ice.DAPTransfer(0x0, 0x3, 0x8, []uint32{0x0, 0xF0ED0005, 0xE0})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])

	//data, err = ice.DAPTransfer(0x0, 0x1, 0x5, []uint32{0xE000EDF0})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])
	//
	//data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0xA05F0001})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])
	//
	//data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0xA05F0003})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Logf("%x", data[:32])

	//var val uint32
	//val, err = ice.ReadAddr32(0x008061FC)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Printf("I: %x\n", val)
	//
	//val, err = ice.ReadAddr32(0x00806010)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Printf("I: %x\n", val)
	//
	//val, err = ice.ReadAddr32(0x00806014)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Printf("I: %x\n", val)
	//
	//val, err = ice.ReadAddr32(0x00806018)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Printf("I: %x\n", val)
	//
	//err = ice.WriteAddr32(0x20030000, 0xA501)
	//if err != nil {
	//	t.Fatal(err)
	//}

	if err := d.Close(); err != nil {
		log.Fatal(err)
	}
	// Finalize the hid package.
	if err := hid.Exit(); err != nil {
		log.Fatal(err)
	}
}

func resume(t *testing.T, ice *AtmelICE) {
	data, err := ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0x01005FA0})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", data[:32])
	ti := time.Now()
	for time.Since(ti) < time.Second*20 {
		data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0x01005FA0})
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%x", data[:32])
		time.Sleep(time.Second)
	}
}

func halt(t *testing.T, ice *AtmelICE) {
	data, err := ice.DAPTransfer(0x0, 0x1, 0x8, []uint32{0x0})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%x", data[:32])

	data, err = ice.DAPTransfer(0x0, 0x1, 0x5, []uint32{0xF0ED00E0})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%x", data[:32])

	data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0x01005FA0})
	if err != nil { // A05F0001
		t.Fatal(err)
	}

	t.Logf("%x", data[:32])

	data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0x03005FA0})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%x", data[:32])

	ti := time.Now()
	for time.Since(ti) < time.Second*20 {
		data, err = ice.DAPTransfer(0x0, 0x1, 0xD, []uint32{0x03005FA0})
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%x", data[:32])
		time.Sleep(time.Second)
	}
}

func TestBinFlipping(t *testing.T) {
	cms := &cmsis.CMSISDAP{}
	ice := &AtmelICE{CMSISDAP: cms}
	err := ice.WriteAddr32(0x41004004, 0xA501)
	if err != nil {
		t.Fatal(err)
	}

	if hex.EncodeToString(ice.Buffer[:12]) != "0005000205044000410d01a5" {
		t.Fatalf("\nGot:      %x\nExpected: 0005000205044000410d01a5", ice.Buffer[:12])
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

/*

Things to keep track of
Bank Selection:





*/
