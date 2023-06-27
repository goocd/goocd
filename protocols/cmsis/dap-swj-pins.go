package cmsis

type PinSelect byte

const (
	PinSelectSWCLKTCK = 0x1
	PinSelectSWDIOTMS = 0x2
	PinSelectTDI      = 0x4
	PinSelectTDO      = 0x8
	PinSelectNTRST    = 0x20
	PinSelectNReset   = 0x80
)
