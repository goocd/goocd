package cmsis

import (
	"strconv"
	"strings"
	"testing"
)

type fauxDevice struct {
	rbuf [64]byte
	wbuf [64]byte
}

func (f *fauxDevice) Read(i []byte) (int, error) {
	f.rbuf[0] = f.wbuf[0]
	// TODO procotol specific responses
	copy(i, f.rbuf[:])
	return len(f.rbuf), nil
}

func (f *fauxDevice) Write(i []byte) (int, error) {
	copy(f.wbuf[:], i)
	return len(i), nil
}

func TestCMSISDAP_DAPSWJPins(t *testing.T) {
	cmsis := new(CMSISDAP)
	d := new(fauxDevice)
	cmsis.ReadWriter = d
	_, err := cmsis.DAPSWJPins(0, 0, 0)
	if err != nil {
		t.Fatalf("Err: %+v", err)
	}
}

func TestRandom(t *testing.T) {
	exampleOne := "0xFFFF"
	exampleTwo := "0xFFFF,3"
	exampleOne = strings.Split(exampleOne, ",")[0]
	u, err := strconv.ParseUint(exampleOne, 0, 64)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Logf("ExampleOne: %s transleted to %x", exampleOne, u)
	exampleTwo = strings.Split(exampleTwo, ",")[0]
	u, err = strconv.ParseUint(exampleTwo, 0, 64)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Logf("ExampleOne: %s transleted to %x", exampleTwo, u)

}
