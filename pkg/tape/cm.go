package tape

import (
	"fmt"
	"strings"
	"time"

	"github.com/MaximSvinin/tape/pkg/model"
)

// vim: ts=4:sts=4:

const (
	typeBinary = 0x00
	typeASCII  = 0x01

	readAttReplyLen = 512
	writeAttCmdLen  = 16
)

type CmAttr struct {
	IsValid  bool
	Name     string
	Command  int
	Len      int
	DataType int
	DataInt  uint64
	DataStr  string
	NoTrim   bool
	MockInt  uint64
	MockStr  string
}

type Cm struct {
	attributes []*CmAttr

	PartCapRemain *CmAttr //
	PartCapMax    *CmAttr //

	TapeAlertFlags        *CmAttr
	LoadCount             *CmAttr //
	MAMSpaceRemaining     *CmAttr
	AssigningOrganization *CmAttr //
	FormattedDensityCode  *CmAttr //
	InitializationCount   *CmAttr // err
	Identifier            *CmAttr // err
	VolumeChangeReference *CmAttr // err

	DeviceAtLoadN0             *CmAttr //
	DeviceAtLoadN1             *CmAttr //
	DeviceAtLoadN2             *CmAttr //
	DeviceAtLoadN3             *CmAttr //
	TotalWritten               *CmAttr //
	TotalRead                  *CmAttr //
	TotalWrittenSession        *CmAttr //
	TotalReadSession           *CmAttr //
	LogicalPosFirstEncrypted   *CmAttr // err
	LogicalPosFirstUnencrypted *CmAttr // err

	UsageHistory     *CmAttr
	PartUsageHistory *CmAttr

	Manufacturer             *CmAttr //
	SerialNo                 *CmAttr //
	Length                   *CmAttr //
	Width                    *CmAttr //
	AssigningOrg             *CmAttr //
	MediumDensity            *CmAttr //
	ManufactureDate          *CmAttr //
	MAMCapacity              *CmAttr //
	Type                     *CmAttr //
	TypeInformation          *CmAttr //
	UserText                 *CmAttr
	DateTimeLastWritten      *CmAttr // err
	TextLocalizationID       *CmAttr // err
	Barcode                  *CmAttr // err
	OwningHostTextualName    *CmAttr // err
	MediaPool                *CmAttr // err
	ApplicationFormatVersion *CmAttr // err
	MediumGloballyUniqID     *CmAttr // err
	MediaPoolGloballyUniqID  *CmAttr // err
}

type SpecsType struct {
	IsValid         bool
	NativeCap       int
	CompressedCap   int
	NativeSpeed     int
	CompressedSpeed int
	FullTapeMinutes int
	CompressFactor  string
	CanWORM         bool
	CanEncrypt      bool
	PartitionNumber int
	BandsPerTape    int
	WrapsPerBand    int
	TracksPerWrap   int
}

// https://github.com/hreinecke/sg3_utils/issues/18
func cmDensityFriendly(d int) (string, SpecsType) {
	friendlyName := "Unknown"
	var specs SpecsType
	switch d {
	case 0x40:
		friendlyName = "LTO-1"
		specs = SpecsType{true, 100, 200, 20, 40, 60 + 23, "2:1", false, false, 1, 4, 12, 8}
	case 0x42:
		friendlyName = "LTO-2"
		specs = SpecsType{true, 200, 400, 40, 80, 60 + 23, "2:1", false, false, 1, 4, 16, 8}
	case 0x44:
		friendlyName = "LTO-3"
		specs = SpecsType{true, 400, 800, 80, 160, 60 + 23, "2:1", true, false, 1, 4, 11, 16}
	case 0x46:
		friendlyName = "LTO-4"
		specs = SpecsType{true, 800, 1600, 120, 240, 60 + 51, "2:1", true, true, 1, 4, 14, 16}
	case 0x58:
		friendlyName = "LTO-5"
		specs = SpecsType{true, 1500, 3000, 140, 280, 60*3 + 10, "2:1", true, true, 2, 4, 20, 16}
	case 0x5A:
		friendlyName = "LTO-6"
		specs = SpecsType{true, 2500, 6250, 160, 400, 60*4 + 20, "2.5:1", true, true, 4, 4, 34, 16}
	case 0x5C:
		friendlyName = "LTO-7"
		specs = SpecsType{true, 6000, 15000, 300, 750, 60*5 + 33, "2.5:1", true, true, 4, 4, 28, 32}
	case 0x5D:
		friendlyName = "LTO-M8"
		specs = SpecsType{true, 9000, 22500, 300, 750, 60*8 + 20, "2.5:1", false, true, 4, 4, 42, 32}
	case 0x5E:
		friendlyName = "LTO-8"
		specs = SpecsType{true, 12000, 30000, 360, 900, 60*9 + 16, "2.5:1", true, true, 4, 4, 52, 32}
	case 0x60: /* guessed, to check FIXME */
		friendlyName = "LTO-9"
		specs = SpecsType{true, 18000, 45000, 400, 1000, 60*12 + 30, "2.5:1", true, true, 4, 0, 0, 0} /* FIXME */
	}
	return friendlyName, specs
}

func cmAttrNew(name string, command int, length int, datatype int, mock interface{}) *CmAttr {
	cmAttr := &CmAttr{
		Name:     name,
		Command:  command,
		Len:      length,
		DataType: datatype,
	}
	cmAttr.MockStr, _ = mock.(string)
	if mockInt, ok := mock.(int); ok {
		cmAttr.MockInt = uint64(mockInt)
	}
	return cmAttr
}

//nolint:funlen // TODO
func NewCm() *Cm {
	attributes := make([]*CmAttr, 0, 10)

	partCapRemain := cmAttrNew("Remaining capacity in partition (MiB)", 0x0000, 8, typeBinary, 198423)
	attributes = append(attributes, partCapRemain)

	partCapMax := cmAttrNew("Maximum capacity in partition (MiB)", 0x0001, 8, typeBinary, 200448)
	attributes = append(attributes, partCapMax)

	tapeAlertFlags := cmAttrNew("Tape alert flags", 0x0002, 8, typeBinary, 0)
	attributes = append(attributes, tapeAlertFlags)

	loadCount := cmAttrNew("Load count", 0x0003, 8, typeBinary, 42)
	attributes = append(attributes, loadCount)

	mAMSpaceRemaining := cmAttrNew("MAM space remaining (bytes)", 0x0004, 8, typeBinary, 850)
	attributes = append(attributes, mAMSpaceRemaining)

	assigningOrganization := cmAttrNew("Assigning organization", 0x0005, 8, typeASCII, "LTO-FAKE")
	attributes = append(attributes, assigningOrganization)

	formattedDensityCode := cmAttrNew("Formatted density code", 0x0006, 1, typeBinary, 66)
	attributes = append(attributes, formattedDensityCode)

	initializationCount := cmAttrNew("Initialization count", 0x0007, 2, typeBinary, "err")
	attributes = append(attributes, initializationCount)

	identifier := cmAttrNew("Identifier (deprecated)", 0x0008, 32, typeASCII, "err")
	attributes = append(attributes, identifier)

	volumeChangeReference := cmAttrNew("Volume change reference", 0x0009, 4, typeBinary, "err")
	attributes = append(attributes, volumeChangeReference)

	deviceAtLoadN0 := cmAttrNew("Device Vendor/Serial at current load",
		0x020A, 40, typeASCII, "FAKEVENDMODEL012345678901234567890123456")
	attributes = append(attributes, deviceAtLoadN0)

	deviceAtLoadN1 := cmAttrNew("Device Vendor/Serial at load N-1",
		0x020B, 40, typeASCII, "FAKEVEND   MODEL12345")
	attributes = append(attributes, deviceAtLoadN1)

	deviceAtLoadN2 := cmAttrNew("Device Vendor/Serial at load N-2",
		0x020C, 40, typeASCII, "ACMEINC \u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000")
	attributes = append(attributes, deviceAtLoadN2)

	deviceAtLoadN3 := cmAttrNew("Device Vendor/Serial at load N-3", 0x020D, 40, typeASCII, "FAKEVEND   MODEL34567")
	attributes = append(attributes, deviceAtLoadN3)

	totalWritten := cmAttrNew("Total MiB written", 0x0220, 8, typeBinary, 17476)
	attributes = append(attributes, totalWritten)

	totalRead := cmAttrNew("Total MiB read", 0x0221, 8, typeBinary, 15827)
	attributes = append(attributes, totalRead)

	totalWrittenSession := cmAttrNew("Total MiB written in current load", 0x0222, 8, typeBinary, 0)
	attributes = append(attributes, totalWrittenSession)

	totalReadSession := cmAttrNew("Total MiB Read in current load", 0x0223, 8, typeBinary, 139)
	attributes = append(attributes, totalReadSession)

	logicalPosFirstEncrypted := cmAttrNew("Logical pos. of 1st encrypted block", 0x0224, 8, typeBinary, "err")
	attributes = append(attributes, logicalPosFirstEncrypted)

	logicalPosFirstUnencrypted := cmAttrNew("Logical pos. of 1st unencrypted block after 1st encrypted block",
		0x0225, 8, typeBinary, "err")
	attributes = append(attributes, logicalPosFirstUnencrypted)

	usageHistory := cmAttrNew("Medium Usage History", 0x0340, 90, typeBinary, "err")
	attributes = append(attributes, usageHistory)

	partUsageHistory := cmAttrNew("Partition Usage History", 0x0341, 90, typeBinary, "err")
	attributes = append(attributes, partUsageHistory)

	manufacturer := cmAttrNew("Manufacturer", 0x0400, 8, typeASCII, "FAKMANUF")
	attributes = append(attributes, manufacturer)

	serialNo := cmAttrNew("Serial No", 0x0401, 32, typeASCII, "123456789")
	attributes = append(attributes, serialNo)

	length := cmAttrNew("Tape length meters", 0x0402, 4, typeBinary, 999)
	attributes = append(attributes, length)

	width := cmAttrNew("Tape width mm", 0x0403, 4, typeBinary, 111)
	attributes = append(attributes, width)

	assigningOrg := cmAttrNew("Assigning Organization", 0x0404, 8, typeASCII, "LTO-FAKE")
	attributes = append(attributes, assigningOrg)

	mediumDensity := cmAttrNew("Medium density code", 0x0405, 1, typeBinary, 0x42)
	attributes = append(attributes, mediumDensity)

	manufactureDate := cmAttrNew("Manufacture Date", 0x0406, 8, typeASCII, "20191231")
	attributes = append(attributes, manufactureDate)

	mAMCapacity := cmAttrNew("MAM Capacity bytes", 0x0407, 8, typeBinary, 4096)
	attributes = append(attributes, mAMCapacity)

	tapeType := cmAttrNew("Type", 0x0408, 1, typeBinary, 1)
	attributes = append(attributes, tapeType)

	typeInformation := cmAttrNew("Type Information", 0x0409, 2, typeBinary, 50)
	attributes = append(attributes, typeInformation)

	userText := cmAttrNew("User Medium Text Label", 0x0803, 160, typeASCII, "User Label")
	attributes = append(attributes, userText)

	dateTimeLastWritten := cmAttrNew("Date and Time Last Written", 0x0804, 12, typeASCII, "err")
	attributes = append(attributes, dateTimeLastWritten)

	textLocalizationID := cmAttrNew("Text Localization Identifier", 0x0805, 1, typeBinary, "err")
	attributes = append(attributes, textLocalizationID)

	barcode := cmAttrNew("Barcode", 0x0806, 12, typeASCII, "err")
	attributes = append(attributes, barcode)

	owningHostTextualName := cmAttrNew("Owning Host Textual Name", 0x0807, 80, typeASCII, "err")
	attributes = append(attributes, owningHostTextualName)

	mediaPool := cmAttrNew("Media Pool", 0x0808, 160, typeASCII, "err")
	attributes = append(attributes, mediaPool)

	applicationFormatVersion := cmAttrNew("Application Format Version", 0x080B, 16, typeASCII, "err")
	attributes = append(attributes, applicationFormatVersion)

	mediumGloballyUniqID := cmAttrNew("Medium Globally Unique Identifier", 0x0820, 36, typeASCII, "err")
	attributes = append(attributes, mediumGloballyUniqID)

	mediaPoolGloballyUniqID := cmAttrNew("Media Pool Globally Unique Identifier", 0x0821, 36, typeASCII, "err")
	attributes = append(attributes, mediaPoolGloballyUniqID)

	cm := &Cm{
		PartCapRemain:              partCapRemain,
		PartCapMax:                 partCapMax,
		TapeAlertFlags:             tapeAlertFlags,
		LoadCount:                  loadCount,
		MAMSpaceRemaining:          mAMSpaceRemaining,
		AssigningOrganization:      assigningOrganization,
		FormattedDensityCode:       formattedDensityCode,
		InitializationCount:        initializationCount,
		Identifier:                 identifier,
		VolumeChangeReference:      volumeChangeReference,
		DeviceAtLoadN0:             deviceAtLoadN0,
		DeviceAtLoadN1:             deviceAtLoadN1,
		DeviceAtLoadN2:             deviceAtLoadN2,
		DeviceAtLoadN3:             deviceAtLoadN3,
		TotalWritten:               totalWritten,
		TotalRead:                  totalRead,
		TotalWrittenSession:        totalWrittenSession,
		TotalReadSession:           totalReadSession,
		LogicalPosFirstEncrypted:   logicalPosFirstEncrypted,
		LogicalPosFirstUnencrypted: logicalPosFirstUnencrypted,
		UsageHistory:               usageHistory,
		PartUsageHistory:           partUsageHistory,
		Manufacturer:               manufacturer,
		SerialNo:                   serialNo,
		Length:                     length,
		Width:                      width,
		AssigningOrg:               assigningOrg,
		MediumDensity:              mediumDensity,
		ManufactureDate:            manufactureDate,
		MAMCapacity:                mAMCapacity,
		Type:                       tapeType,
		TypeInformation:            typeInformation,
		UserText:                   userText,
		DateTimeLastWritten:        dateTimeLastWritten,
		TextLocalizationID:         textLocalizationID,
		Barcode:                    barcode,
		OwningHostTextualName:      owningHostTextualName,
		MediaPool:                  mediaPool,
		ApplicationFormatVersion:   applicationFormatVersion,
		MediumGloballyUniqID:       mediumGloballyUniqID,
		MediaPoolGloballyUniqID:    mediaPoolGloballyUniqID,

		attributes: attributes,
	}
	cm.attributes = attributes

	return cm
}

func (cm *Cm) GetOptions() *model.TapeInfo {
	out := new(model.TapeInfo)

	cm.addCap(out)
	cm.addType(out)
	cm.addSpecs(out)
	cm.addOrg(out)
	cm.addSessions(out)

	return out
}

func (cm *Cm) addSessions(out *model.TapeInfo) {
	out.Sessions = make([]*model.SessionAttribute, 0, 4)
	for i, load := range []*CmAttr{cm.DeviceAtLoadN0, cm.DeviceAtLoadN1, cm.DeviceAtLoadN2, cm.DeviceAtLoadN3} {
		if load.IsValid {
			var devname, serial string
			if len(load.DataStr) > 8 {
				devname = strings.Trim(load.DataStr[:8], " \u0000")
				serial = strings.Trim(load.DataStr[8:], " \u0000")
			} else {
				devname = strings.Trim(load.DataStr, " \u0000")
			}
			out.Sessions = append(out.Sessions, &model.SessionAttribute{
				Number:  i,
				Devname: devname,
				Serial:  serial,
			})
		}
	}
}

func (cm *Cm) addCap(out *model.TapeInfo) {
	if cm.PartCapRemain.IsValid {
		out.PartCapRemain = &model.CapAttribute{
			Name:  cm.PartCapRemain.Name,
			Value: cm.PartCapRemain.DataInt,
		}
	}

	if cm.PartCapMax.IsValid {
		out.PartCapMax = &model.CapAttribute{
			Name:  cm.PartCapMax.Name,
			Value: cm.PartCapMax.DataInt,
		}
	}

	if cm.Length.IsValid {
		out.TapeLength = &model.CapAttribute{
			Name:  cm.Length.Name,
			Value: cm.Length.DataInt,
		}
	}

	if cm.Width.IsValid {
		out.TapeWidth = &model.CapAttribute{
			Name:  cm.Width.Name,
			Value: uint64(float32(cm.Width.DataInt) / 10),
		}
	}

	if cm.MAMCapacity.IsValid {
		if cm.MAMSpaceRemaining.IsValid {
			out.MAMCapacity = model.MAMCapacityAttribute{
				Name:           cm.MAMCapacity.Name,
				Capacity:       cm.MAMCapacity.DataInt,
				SpaceRemaining: &cm.MAMSpaceRemaining.DataInt,
			}
		} else {
			out.MAMCapacity = model.MAMCapacityAttribute{
				Name:           cm.MAMCapacity.Name,
				Capacity:       cm.MAMCapacity.DataInt,
				SpaceRemaining: nil,
			}
		}
	}

	if cm.LoadCount.IsValid {
		out.CartridgeLoadCount = cm.LoadCount.DataInt
	}

	if cm.TotalWritten.IsValid && cm.TotalRead.IsValid {
		out.TotalWritten = &model.CapAttribute{
			Name:  "data written MiB",
			Value: cm.TotalWritten.DataInt,
		}
		out.TotalRead = &model.CapAttribute{
			Name:  "data read MiB",
			Value: cm.TotalRead.DataInt,
		}
	}

	if cm.TotalWrittenSession.IsValid && cm.TotalReadSession.IsValid {
		out.WrittenSession = &model.CapAttribute{
			Name:  "data written session MiB",
			Value: cm.TotalWrittenSession.DataInt,
		}
		out.ReadSession = &model.CapAttribute{
			Name:  "data read session MiB",
			Value: cm.TotalReadSession.DataInt,
		}
	}
}

func (cm *Cm) addType(out *model.TapeInfo) {
	if cm.Type.IsValid {
		friendlyName := "Unknown"
		switch cm.Type.DataInt {
		case 0x00:
			friendlyName = "Data cartridge"
		case 0x01:
			friendlyName = "Cleaning cartridge"
			if cm.TypeInformation.IsValid {
				friendlyName = fmt.Sprintf("%s (%d cycles max)", friendlyName, cm.TypeInformation.DataInt)
			}
		case 0x80:
			friendlyName = "WORM (Write-once) cartridge"
		}
		out.CartridgeType = &model.StrAttribute{
			Name:  cm.Type.Name,
			Value: friendlyName,
		}
	}
}

func (cm *Cm) addSpecs(out *model.TapeInfo) {
	var specs SpecsType
	specs.IsValid = false
	if cm.MediumDensity.IsValid {
		var mediumformat string
		mediumformat, specs = cmDensityFriendly(int(cm.MediumDensity.DataInt))
		formattedAs, _ := cmDensityFriendly(int(cm.FormattedDensityCode.DataInt))

		out.MediumDensity = &model.MediumDensityAttribute{
			Name:         cm.MediumDensity.Name,
			Mediumformat: mediumformat,
			FormattedAs:  formattedAs,
		}
	}

	if specs.IsValid {
		out.Specs = &model.SpecsAttribute{
			Capacity: &model.SpecsCapacityAttribute{
				Name:           "Capacity GB native compressed",
				Native:         specs.NativeCap,
				Compressed:     specs.CompressedCap,
				CompressFactor: specs.CompressFactor,
			},
			Speed: &model.SpecsSpeedAttribute{
				Name:       "R/W Speed MB/s native compressed",
				Native:     specs.NativeSpeed,
				Compressed: specs.CompressedSpeed,
			},
			Partitions: &model.SpecsPartitionsAttribute{
				Name:            "max partitions supported",
				PartitionNumber: specs.PartitionNumber,
			},
			Phy: &model.SpecsPhyAttribute{
				Name:          "bands/tape wraps/band tracks/wrap total tracks",
				BandsPerTape:  specs.BandsPerTape,
				WrapsPerBand:  specs.WrapsPerBand,
				TracksPerWrap: specs.TracksPerWrap,
				Total:         specs.BandsPerTape * specs.WrapsPerBand * specs.TracksPerWrap,
			},
			Duration: &model.SpecsDurationAttribute{
				Name:            "Duration to fill tape",
				FullTapeMinutes: specs.FullTapeMinutes,
			},
		}
	}
}

func (cm *Cm) addOrg(out *model.TapeInfo) {
	if cm.AssigningOrg.IsValid {
		out.AssigningOrg = &model.StrAttribute{
			Name:  cm.AssigningOrg.Name,
			Value: cm.AssigningOrg.DataStr,
		}
	}
	if cm.Manufacturer.IsValid {
		out.Manufacturer = &model.StrAttribute{
			Name:  cm.Manufacturer.Name,
			Value: cm.Manufacturer.DataStr,
		}
	}

	manufactureDate := cm.ManufactureDate.DataStr
	if len(cm.ManufactureDate.DataStr) == 8 {
		if d, err := time.Parse("20060102", cm.ManufactureDate.DataStr); err == nil {
			years := time.Since(d).Hours() / 24.0 / 365.0
			manufactureDate = fmt.Sprintf("%s-%s-%s (roughly %.1f years ago)",
				cm.ManufactureDate.DataStr[0:4], cm.ManufactureDate.DataStr[4:6], cm.ManufactureDate.DataStr[6:8], years)
		} else {
			manufactureDate = fmt.Sprintf("%s-%s-%s",
				cm.ManufactureDate.DataStr[0:4], cm.ManufactureDate.DataStr[4:6], cm.ManufactureDate.DataStr[6:8])
		}
	}
	out.ManufactureDate = &model.StrAttribute{
		Name:  cm.ManufactureDate.Name,
		Value: manufactureDate,
	}
}
