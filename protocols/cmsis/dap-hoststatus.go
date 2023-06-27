package cmsis

type HostType byte
type HostStatus byte

const (
	HostConnect = HostType(0x0)
	HostRunning = HostType(0x1)

	StatusOff = HostStatus(0x0)
	StatusOn  = HostStatus(0x1)
)
