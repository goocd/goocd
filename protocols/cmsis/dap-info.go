package cmsis

type DAPInfo byte

const (
	VendorName              = DAPInfo(0x1)
	ProductName             = DAPInfo(0x2)
	SerialNumber            = DAPInfo(0x3)
	CMSISDAPProtocolVersion = DAPInfo(0x4)
	TargetDeviceVendor      = DAPInfo(0x5)
	TargetDeviceName        = DAPInfo(0x6)
	TargetBoardvendor       = DAPInfo(0x7)
	TargetBoardName         = DAPInfo(0x8)
	ProductFirmwareVersion  = DAPInfo(0x9)
	Capabilities            = DAPInfo(0xF0)
	TestDomainTimer         = DAPInfo(0xF1)
	UARTReceiveBufferSize   = DAPInfo(0xFB)
	UARTTransmiteBufferSize = DAPInfo(0xFC)
	SWOTraceBufferSize      = DAPInfo(0xFD)
	PacketCount             = DAPInfo(0xFE)
	PacketSize              = DAPInfo(0xFF)
)
