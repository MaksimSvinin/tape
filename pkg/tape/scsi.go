package tape

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/benmcclelland/sgio"
)

type SCSIInquiryReturn struct {
	PeripheralDeviceType uint8 `bitfield:"5"` // Byte 0
	PeripheralQualifier  uint8 `bitfield:"3"`
	Reserved0            uint8
	Version              uint8
	ReponseDataFormat    uint8 `bitfield:"4"`
	HiSup                uint8 `bitfield:"1"`
	NACA                 uint8 `bitfield:"1"`
	Obsolete0            uint8 `bitfield:"1"`
	Obsolete1            uint8 `bitfield:"1"`
	AdditionalLen        uint8
	Protect              uint8 `bitfield:"1"`
	Reserved1            uint8 `bitfield:"2"`
	ThreePC              uint8 `bitfield:"1"`
	TPGS                 uint8 `bitfield:"2"`
	ACC                  uint8 `bitfield:"1"`
	SCCS                 uint8 `bitfield:"1"`
	Osef0                uint8
	Osef1                uint8
	VendorID             [8]byte
	ProductID            [16]byte
	ProductRevision      [4]byte // YMDV(F63D), Y=15 M=6 D=3 V=D
	Reserved2            uint8
	Obsolete2            uint8
	MaxSpeed             uint8 `bitfield:"4"`
	ProtocolID           uint8 `bitfield:"4"`
	FIPS                 uint8 `bitfield:"2"`
	Reserved3            uint8 `bitfield:"5"`
	Restricted           uint8 `bitfield:"1"`
	Reserved4            uint8
	OEMSpecific          uint8
	OEMSpecificSubfield  uint8
	Reserved5            uint8
	Reserved6            uint32
	PartNumber           [8]byte
	Reserved7            uint8
	Reserved8            uint8
	Truc1                uint16
	Truc2                uint16
	Truc3                uint16
	Truc4                uint16
	Truc5                uint16
	Truc6                uint16
}

// copy of sg.OpenScsiDevice() but with RDONLY instead of O_RDWR.
func OpenScsiDeviceRO(fname string) (*os.File, error) {
	f, err := os.OpenFile(fname, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	var version uint32
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(sgio.SG_GET_VERSION_NUM),
		uintptr(unsafe.Pointer(&version)),
	)
	if errno != 0 {
		return nil, fmt.Errorf("failed to get version info from sg device (errno=%w)", errno)
	}
	if version < 30000 {
		return nil, fmt.Errorf("device does not appear to be an sg device")
	}
	return f, nil
}
