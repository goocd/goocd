package cmsisdap

// TODO: JTAG stuff

// DAP Commands
const (
	DAPInfoCMD           = 0x0
	DAPHostStatusCMD     = 0x1
	DAPConnectCMD        = 0x2
	DAPDisconnectCMD     = 0x3
	DAPTransferConfigCMD = 0x4
	DAPTransferCMD       = 0x5
	DAPWriteAbortCMD     = 0x8
	DAPDelay             = 0x9
	DAPResetTarget       = 0xA
	DAPSWJPinsCMD        = 0x10
	DAPSWJClockCMD       = 0x11
	DAPSWJSequenceCMD    = 0x12
	DAPSWDConfigCMD      = 0x13
	DAPSWDSequenceCMD    = 0x1D
)

// DAP Connect
const (
	DefaultPort = 0x0
	SWDPort     = 0x1
	JTAGPort    = 0x2
)

// DAP Host/Status
const (
	HostConnect = 0x0
	HostRunning = 0x1

	StatusOff = 0x0
	StatusOn  = 0x1
)

// PinOut Mask
const (
	PinMaskSWCLKTCK = 0x1
	PinMaskSWDIOTMS = 0x2
	PinMaskTDI      = 0x4
	PinMaskTDO      = 0x8
	PinMaskNTRST    = 0x20
	PinMaskNReset   = 0x80
)

// DAP Info
const (
	VendorName              = 0x1
	ProductName             = 0x2
	SerialNumber            = 0x3
	CMSISDAPProtocolVersion = 0x4
	TargetDeviceVendor      = 0x5
	TargetDeviceName        = 0x6
	TargetBoardvendor       = 0x7
	TargetBoardName         = 0x8
	ProductFirmwareVersion  = 0x9
	Capabilities            = 0xF0
	TestDomainTimer         = 0xF1
	UARTReceiveBufferSize   = 0xFB
	UARTTransmiteBufferSize = 0xFC
	SWOTraceBufferSize      = 0xFD
	PacketCount             = 0xFE
	PacketSize              = 0xFF
)

// DAP Transfer
const (
	DebugPort  = 0x0
	AccessPort = 0x1

	Read  = 0x2
	Write = 0x0

	PortRegister0 = 0x0
	PortRegister4 = 0x4
	PortRegister8 = 0x8
	PortRegisterC = 0xC

	ValueMatch = 0x10
	MatchMask  = 0x20
	TimeStamp  = 0x80
)

// DAP Response Status
const (
	DAP_OK    = 0x0
	DAP_Error = 0xFF
)

// Todo: Figure out if Little Endian is always expcted as this is the Big Endian Representation

// Clock Speeds
const (
	ClockSpeed2Mhz = uint32(0x1E8480)
	ClockSpeed4Mhz = uint32(0x3D0900)
)

// DAP SWD Configs
const (
	SWDConfigClockCycles1 = 0x0
	SWDConfigClockCycles2 = 0x1
	SWDConfigClockCycles3 = 0x2
	SWDConfigClockCycles4 = 0x3

	SWDConfigNoDataPhase     = 0x0
	SWDConfigAlwaysDataPhase = 0x4
)
