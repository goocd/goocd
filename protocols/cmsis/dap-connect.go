package cmsis

type ConnectPort byte

const (
	Default = ConnectPort(0x0)
	SWD     = ConnectPort(0x1)
	JTAG    = ConnectPort(0x2)
)
