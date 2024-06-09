package model

type CapAttribute struct {
	Name  string `json:"name,omitempty"`
	Value uint64 `json:"value"`
}

type StrAttribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type MediumDensityAttribute struct {
	Name         string `json:"name,omitempty"`
	Mediumformat string `json:"mediumformat,omitempty"`
	FormattedAs  string `json:"formattedAs,omitempty"`
}

type MAMCapacityAttribute struct {
	Name           string  `json:"name,omitempty"`
	Capacity       uint64  `json:"capacity"`
	SpaceRemaining *uint64 `json:"spaceRemaining,omitempty"`
}

type SpecsCapacityAttribute struct {
	Name           string `json:"name,omitempty"`
	Native         int    `json:"native"`
	Compressed     int    `json:"compressed"`
	CompressFactor string `json:"compressFactor,omitempty"`
}

type SpecsSpeedAttribute struct {
	Name       string `json:"name,omitempty"`
	Native     int    `json:"native"`
	Compressed int    `json:"compressed"`
}

type SpecsPartitionsAttribute struct {
	Name            string `json:"name,omitempty"`
	PartitionNumber int    `json:"partitionNumber"`
}

type SpecsPhyAttribute struct {
	Name          string `json:"name,omitempty"`
	BandsPerTape  int    `json:"bandsPerTape"`
	WrapsPerBand  int    `json:"wrapsPerBand"`
	TracksPerWrap int    `json:"tracksPerWrap"`
	Total         int    `json:"total"`
}

type SpecsDurationAttribute struct {
	Name            string `json:"name,omitempty"`
	FullTapeMinutes int    `json:"fullTapeMinutes"`
}

type SpecsAttribute struct {
	Capacity   *SpecsCapacityAttribute   `json:"capacity,omitempty"`
	Speed      *SpecsSpeedAttribute      `json:"speed,omitempty"`
	Partitions *SpecsPartitionsAttribute `json:"partitions,omitempty"`
	Phy        *SpecsPhyAttribute        `json:"phy,omitempty"`
	Duration   *SpecsDurationAttribute   `json:"duration,omitempty"`
}

type SessionAttribute struct {
	Number  int    `json:"number"`
	Devname string `json:"devname,omitempty"`
	Serial  string `json:"serial,omitempty"`
}

type Attributes struct {
	SerialNo string `json:"serialNumber,omitempty"`

	PartCapRemain *CapAttribute `json:"partCapRemain,omitempty"`
	PartCapMax    *CapAttribute `json:"partCapMax,omitempty"`

	MediumDensity *MediumDensityAttribute `json:"mediumDensity,omitempty"`

	CartridgeType *StrAttribute `json:"cartridgeType,omitempty"`

	AssigningOrg    *StrAttribute `json:"assigningOrg,omitempty"`
	Manufacturer    *StrAttribute `json:"manufacturer,omitempty"`
	ManufactureDate *StrAttribute `json:"manufactureDate,omitempty"`

	TapeLength *CapAttribute `json:"tapeLength,omitempty"`
	TapeWidth  *CapAttribute `json:"tapeWidth,omitempty"`

	MAMCapacity *MAMCapacityAttribute `json:"mAMCapacity,omitempty"`

	Specs *SpecsAttribute `json:"specs,omitempty"`

	CartridgeLoadCount uint64 `json:"cartridgeLoadCount,omitempty"`

	TotalWritten *CapAttribute `json:"totalWritten,omitempty"`
	TotalRead    *CapAttribute `json:"totalRead,omitempty"`

	WrittenSession *CapAttribute `json:"writtenSession,omitempty"`
	ReadSession    *CapAttribute `json:"readSession,omitempty"`

	Sessions []*SessionAttribute `json:"sessions,omitempty"`
}
