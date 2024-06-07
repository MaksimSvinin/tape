package tape

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/HewlettPackard/structex"
	"github.com/MaximSvinin/tape/pkg/model"
	"github.com/benmcclelland/mt"
	"github.com/benmcclelland/sgio"
)

const (
	devSt  = "/dev/st0"
	mtPath = "/bin/mt"
)

type StorageInfo struct {
	Vendor   string
	Model    string
	Firmware string

	Attributes *model.TapeInfo
}

type Tape interface {
	Info() (*StorageInfo, error)

	Write(file io.Reader) (int64, error)
	Read() (io.ReadCloser, error)

	Erase() error
	Rewind() error
	Eject() error
}

type tape struct {
	mt *mt.Drive
	cm *Cm
	m  sync.Mutex
}

// локальный стример ленты.
func NewTapeStorage() (Tape, error) {
	return &tape{
		cm: NewCm(),
		mt: mt.NewDriveCmd(devSt, mtPath),
	}, nil
}

func (t *tape) Info() (*StorageInfo, error) {
	t.m.Lock()
	defer t.m.Unlock()

	infoDrive, err := OpenScsiDeviceRO(devSt)
	if err != nil {
		return nil, err
	}
	defer infoDrive.Close()

	err = sgio.TestUnitReady(infoDrive)
	if err != nil {
		return nil, err
	}

	replyBuf := []byte{1, 128, 3, 2, 91, 0, 1, 48, 72, 80, 32, 32, 32, 32, 32, 32, 85, 108,
		116, 114, 105, 117, 109, 32, 50, 45, 83, 67, 83, 73, 32, 32, 70, 54, 51, 68, 0,
		0, 0, 0, 0, 12, 0, 36, 68, 82, 45, 49, 48, 0, 0, 0, 0, 0, 0, 0, 12, 0, 0, 84, 11,
		28, 2, 119, 2, 28, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	var parsed = new(SCSIInquiryReturn)
	err = structex.Decode(bytes.NewReader(replyBuf), parsed)
	if err != nil {
		return nil, err
	}

	for i := range t.cm.attributes {
		err = GetAttribute(t.cm.attributes[i], infoDrive)
		if err != nil {
			continue
		}
	}

	return &StorageInfo{
		Vendor:   strings.Trim(string(parsed.VendorID[:]), " \u0000"),
		Model:    strings.Trim(string(parsed.ProductID[:]), " \u0000"),
		Firmware: strings.Trim(string(parsed.ProductRevision[:]), " \u0000"),

		Attributes: t.cm.GetOptions(),
	}, nil
}

func (t *tape) Write(file io.Reader) (int64, error) {
	t.m.Lock()
	defer t.m.Unlock()

	f, _, err := OpenTapeWriteOnly(devSt)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.ReadFrom(file)
}

func (t *tape) Read() (io.ReadCloser, error) {
	t.m.Lock()
	defer t.m.Unlock()

	f, _, err := OpenTapeReadOnly(devSt)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (t *tape) Erase() error {
	return t.mt.Erase()
}

func (t *tape) Rewind() error {
	return t.mt.Rewind()
}

func (t *tape) Eject() error {
	return t.mt.Eject()
}
