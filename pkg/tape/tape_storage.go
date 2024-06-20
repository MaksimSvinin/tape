package tape

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/HewlettPackard/structex"
	"github.com/benmcclelland/mtio"
	"github.com/benmcclelland/sgio"
	"github.com/rs/zerolog/log"

	"github.com/MaximSvinin/tape/pkg/model"
)

const (
	devNst = "/dev/nst0"

	bufCount = 32768
)

type Tape interface {
	Info() (*model.TapeInfo, error)

	Write(file io.Reader) (*model.FileWriteInfo, error)
	Read(fileNumbers []int, patch string) error

	Erase() error
	Eject() error
}

type tape struct {
	cm *Cm

	m         sync.RWMutex
	operation model.Operation
}

// локальный стример ленты.
func NewTapeStorage() Tape {
	return &tape{
		cm:        NewCm(),
		m:         sync.RWMutex{},
		operation: model.Unknown,
	}
}

func (t *tape) GetOperation() string {
	t.m.RLock()
	defer t.m.RUnlock()

	return t.operation.String()
}

func (t *tape) Info() (*model.TapeInfo, error) {
	t.m.Lock()
	defer func() {
		t.operation = model.Unknown
		t.m.Unlock()
	}()
	t.operation = model.Info

	infoDrive, err := OpenScsiDeviceRO(devNst)
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

	err = t.cm.ReadAttributes(infoDrive)
	if err != nil {
		log.Warn().Err(err).Msg("warn read attributes")
	}

	return &model.TapeInfo{
		Vendor:   strings.Trim(string(parsed.VendorID[:]), " \u0000"),
		Model:    strings.Trim(string(parsed.ProductID[:]), " \u0000"),
		Firmware: strings.Trim(string(parsed.ProductRevision[:]), " \u0000"),

		Attributes: t.cm.GetAttributes(),
	}, nil
}

func (t *tape) Write(file io.Reader) (*model.FileWriteInfo, error) {
	t.m.Lock()
	defer t.unlock()
	t.operation = model.Write

	f, _, err := OpenTapeWriteOnly(devNst)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = t.eom(f)
	if err != nil {
		return nil, err
	}

	n, err := f.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	err = t.eof(f)
	if err != nil {
		return nil, err
	}

	status, err := mtio.GetStatus(f)
	if err != nil {
		return nil, err
	}
	return &model.FileWriteInfo{
		FileNo:     status.FileNo,
		BytesWrite: n,
	}, nil
}

func (t *tape) Read(fileNumbers []int, patch string) error {
	t.m.Lock()
	defer t.m.Unlock()
	t.operation = model.Read

	f, _, err := OpenTapeReadOnly(devNst)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.rewind(f)
	if err != nil {
		return err
	}

	if len(fileNumbers) == 0 {
		return t.extractAllFiles(f, patch)
	}

	fileNumbersMap := make(map[int]struct{})
	maxFileNumber := 0
	for i := range fileNumbers {
		fileNumbersMap[fileNumbers[i]] = struct{}{}
		if maxFileNumber < fileNumbers[i] {
			maxFileNumber = fileNumbers[i]
		}
	}
	return t.extractFiles(f, patch, maxFileNumber, fileNumbersMap)
}

func (t *tape) Erase() error {
	t.m.Lock()
	defer t.unlock()
	t.operation = model.Erase

	f, _, err := OpenTapeWriteOnly(devNst)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.rewind(f)
	if err != nil {
		return err
	}
	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTERASE), mtio.WithCount(1)))
}

func (t *tape) Eject() error {
	t.m.Lock()
	defer t.unlock()

	f, _, err := OpenTapeReadOnly(devNst)
	if err != nil {
		return err
	}
	defer f.Close()

	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTOFFL)))
}

func (t *tape) rewind(f *os.File) error {
	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTREW)))
}

func (t *tape) eom(f *os.File) error {
	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTEOM)))
}

func (t *tape) eof(f *os.File) error {
	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTWEOF)))
}

func (t *tape) fsf(f *os.File) error {
	return mtio.DoOp(f, mtio.NewMtOp(mtio.WithOperation(mtio.MTFSF)))
}

func (t *tape) extractAllFiles(f *os.File, patch string) error {
	i := 1
	for {
		outPath := path.Join(patch, fmt.Sprintf("file%d", i))
		log.Info().Str("outPath", outPath).Msg("create out file")

		stop, err := t.readData(f, outPath)
		if err != nil {
			return err
		}
		if stop {
			err = os.Remove(outPath)
			if err != nil {
				return err
			}
			break
		}
		i++
	}
	return nil
}

func (t *tape) extractFiles(
	f *os.File,
	patch string,
	maxFileNumber int,
	fileNumbersMap map[int]struct{},
) error {
	for i := range maxFileNumber {
		i++

		if _, ok := fileNumbersMap[i]; ok {
			outPath := path.Join(patch, fmt.Sprintf("file%d", i))
			log.Info().Str("outPath", outPath).Msg("create out file")

			_, err := t.readData(f, outPath)
			if err != nil {
				return err
			}
		} else {
			err := t.fsf(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *tape) readData(f *os.File, outPath string) (bool, error) {
	outFile, err := os.Create(outPath)
	if err != nil {
		return false, err
	}
	defer outFile.Close()

	buf := make([]byte, bufCount)
	start := true

	for {
		n, err := f.Read(buf) //nolint:govet //TODO
		if start && n == 0 {
			return true, nil
		}
		start = false

		if errors.Is(err, io.EOF) {
			if n != 0 {
				_, err = outFile.Write(buf[:n])
				if err != nil {
					return false, err
				}
				break
			}
			return false, err
		} else if err != nil {
			return false, err
		}

		_, err = outFile.Write(buf[:n])
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (t *tape) unlock() {
	t.operation = model.Unknown
	t.m.Unlock()
}
