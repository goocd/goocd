package targets

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/goocd/goocd/connectors/atmel-ice"
	"github.com/goocd/goocd/protocols/cmsis"
	"github.com/sstallion/go-hid"
	"log"
	"strings"
)

func init() {

	TargetMap["atsame51-atmelice"] = TargetFunc(func(args *Args) error {
		// Initialize the hid package.
		if err := hid.Init(); err != nil {
			log.Fatal(err)
		}

		// Open the device using the VID and PID.
		d, err := hid.OpenFirst(0x03eb, 0x2141) // Atmel-ICE VIP & PID
		if err != nil {
			log.Fatal(err)
		}

		cms := &cmsis.CMSISDAP{ReadWriter: d}
		ice := &atmel_ice.AtmelICE{CMSISDAP: cms}
		err = ice.Configure()
		if err != nil {
			log.Fatal(err)
		}

		//log.Printf("GOT HERE: atsame51-cmsisdap")

		// Open File, Buffer here etc.
		if args.Load != "" {
			//pgmSrc, err := parserany.Parse(*loadF)
			//chkerr(err)
			//nvm := nvmload.NVMLoader {
			//	ProgramSource: pgmSrc
			//	NVMAccess: reg,
			//}
			//chkerr(nvm.NVMLoad())
		}

		if args.ReadMem != "" {

			if !strings.HasPrefix(args.ReadMem, "0x") || len(args.ReadMem) > 10 {
				log.Fatal(fmt.Errorf("invalid Hex-Value Provided"))
			}
			args.ReadMem = args.ReadMem[2:]
			dst := make([]byte, hex.DecodedLen(len(args.ReadMem)))
			n, err := hex.Decode(dst, []byte(args.ReadMem))
			if err != nil {
				log.Fatal(err)
			}
			if n > 4 {
				log.Fatal("this should've been caught earlier on")
			}
			res := make([]byte, 4) // Gaurentee 4 bytes
			copy(res[4-len(dst):], dst)
			addr := binary.BigEndian.Uint32(res)

			data, err := ice.ReadAddr32(addr)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("ReadAddr32[Address: %x, Value: %x]", addr, data)
		}

		return nil
	})

}
