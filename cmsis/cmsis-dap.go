package cmsis

type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}
type CMSISDAP struct {
	ReadWriter ReadWriter
	buffer     [512]byte // Pre-Allocated to avoid potential os allocation issues
}

// exposes CMSIS-DAP stuff that follows https://arm-software.github.io/CMSIS_5/DAP/html/group__DAP__Transfer.html

func (c *CMSISDAP) DAPInfo(info DAP_Info) error {
	return nil
}

func (c *CMSISDAP) DAPTransferConfigure() error {
	return nil
}

func (c *CMSISDAP) DAPTransfer(idx uint8, count uint8) error {
	return nil
}

// but it ALSO exposes methods that are friendly to the next layer - addressing CoreSight debug registers in this case

func (c *CMSISDAP) CoreSightDebugSend() error {
	return nil
}
