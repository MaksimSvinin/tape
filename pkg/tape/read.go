package tape

import (
	"os"
	"syscall"
)

func OpenTapeReadOnly(drive string) (*os.File, bool, error) {
	fileDescription, err := os.Stat(drive)
	if err != nil {
		return nil, false, err
	}

	var f *os.File
	isRegular := fileDescription.Mode().IsRegular()
	if isRegular {
		f, err = os.Open(drive)
		if err != nil {
			return f, isRegular, err
		}

		return f, isRegular, nil
	}

	f, err = os.OpenFile(drive, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeCharDevice)
	if err != nil {
		return f, isRegular, err
	}

	return f, isRegular, nil
}
