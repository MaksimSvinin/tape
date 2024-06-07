package tape

import (
	"errors"
	"os"
	"strings"

	"github.com/benmcclelland/sgio"
)

func GetAttribute(attr *CmAttr, infoDrive *os.File) error {
	senseBuf := make([]byte, sgio.SENSE_BUF_LEN)
	replyBuf := make([]byte, readAttReplyLen)

	/* READ ATTRIBUTE (8Ch)
		bits: 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0
	   byte0: --- OPERATION CODE (8Ch) ----
	   byte1: reserved  | SERVICE ACTION
	   byte2: obsolete
	   byte3: obsolete
	   byte4: obsolete
	   byte5: LOGICAL VOLUME NUMBER
	   byte6: reserved
	   byte7: PARTITION NUMBER
	   byte8: (MSB) <-- FIRST ATTRIBUTE
	   byte9:     IDENTIFIER      --> (LSB)
	  byte10: (MSB) <-- ALLOCATION
	  byte11:
	  byte12:
	  byte13:     LENGTH          --> (LSB)
	  byte14: reserved                | CACHE
	  byte15: CONTROL BYTE (00h)
	*/
	inqCmdBlk := []uint8{0x8C, 0, 0, 0, 0, 0, 0, 0, 0x04, 0x00, 0, 0, 159, 0, 0, 0}
	inqCmdBlk[8] = uint8(0xff & (attr.Command >> 8))
	inqCmdBlk[9] = uint8(0xff & attr.Command)
	inqCmdBlk[12] = uint8(0xff & attr.Len)

	ioHdr := &sgio.SgIoHdr{
		InterfaceID:    int32('S'),
		CmdLen:         uint8(len(inqCmdBlk)),
		MxSbLen:        sgio.SENSE_BUF_LEN,
		DxferDirection: sgio.SG_DXFER_FROM_DEV,
		DxferLen:       readAttReplyLen,
		Dxferp:         &replyBuf[0],
		Cmdp:           &inqCmdBlk[0],
		Sbp:            &senseBuf[0],
		Timeout:        sgio.TIMEOUT_20_SECS,
	}

	attr.IsValid = false

	err := sgio.SgioSyscall(infoDrive, ioHdr)
	if err != nil {
		return err
	}

	err = sgio.CheckSense(ioHdr, &senseBuf)
	if err != nil {
		return err
	}

	if attr.DataType == typeBinary {
		attr.DataInt = 0
		for i := 0; i < attr.Len; i++ {
			attr.DataInt *= 256
			attr.DataInt += uint64(replyBuf[9+i])
		}
		attr.IsValid = true
		return nil
	}

	if attr.DataType == typeASCII {
		attr.DataStr = string(replyBuf[9:(9 + attr.Len)])
		if !attr.NoTrim {
			attr.DataStr = strings.TrimRight(attr.DataStr, " ")
		}
		attr.IsValid = true
		return nil
	}

	return errors.New("invalid type")
}
