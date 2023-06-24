package cmsis

type DAP_Info byte

const (
	VendorName                 = DAP_Info(0x1)
	ProductName                = DAP_Info(0x2)
	SerialNumber               = DAP_Info(0x3)
	CMSISDAP_ProtocolVersion   = DAP_Info(0x4)
	Target_Device_Vendor       = DAP_Info(0x5)
	Target_Device_Name         = DAP_Info(0x6)
	Target_Board_vendor        = DAP_Info(0x7)
	Target_Board_Name          = DAP_Info(0x8)
	Product_Firmware_Version   = DAP_Info(0x9)
	Capabilities               = DAP_Info(0xF0)
	Test_Domain_Timer          = DAP_Info(0xF1)
	UART_Receive_Buffer_Size   = DAP_Info(0xFB)
	UART_Transmite_Buffer_Size = DAP_Info(0xFC)
	SWO_Trace_Buffer_Size      = DAP_Info(0xFD)
	Packet_Count               = DAP_Info(0xFE)
	Packet_Size                = DAP_Info(0xFF)
)
