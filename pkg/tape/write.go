package tape

import "os"

func OpenTapeWriteOnly(drive string) (*os.File, bool, error) {
	stat, err := os.Stat(drive)

	var isRegular bool
	if err == nil {
		isRegular = stat.Mode().IsRegular()
	} else {
		if os.IsNotExist(err) {
			isRegular = true
		} else {
			return nil, false, err
		}
	}

	var f *os.File
	if isRegular {
		f, err = os.OpenFile(drive, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, false, err
		}

		// Clear the file's content
		if err = f.Truncate(0); err != nil {
			return nil, false, err
		}

		if err = f.Close(); err != nil {
			return nil, false, err
		}
	}

	if isRegular {
		f, err = os.OpenFile(drive, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, false, err
		}

		// No need to go to end manually due to `os.O_APPEND`
	} else {
		f, err = os.OpenFile(drive, os.O_APPEND|os.O_WRONLY, os.ModeCharDevice)
		if err != nil {
			return nil, false, err
		}
	}

	return f, isRegular, nil
}
