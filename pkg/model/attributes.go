package model

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

	PartCapRemain *uint64 `json:"partCapRemain,omitempty"` // Свободное место в байтах
	PartCapMax    *uint64 `json:"partCapMax,omitempty"` // Всего места в байтах

	MediumDensity *MediumDensityAttribute `json:"mediumDensity,omitempty"`

	CartridgeType string `json:"cartridgeType,omitempty"`

	AssigningOrg    string `json:"assigningOrg,omitempty"`
	Manufacturer    string `json:"manufacturer,omitempty"`
	ManufactureDate string `json:"manufactureDate,omitempty"`

	TapeLength *uint64 `json:"tapeLength,omitempty"` // Длинна ленты в метрах
	TapeWidth  *uint64 `json:"tapeWidth,omitempty"` // Ширина ленты в милиметрах

	MAMCapacity *MAMCapacityAttribute `json:"mAMCapacity,omitempty"`

	Specs *SpecsAttribute `json:"specs,omitempty"`

	CartridgeLoadCount *uint64 `json:"cartridgeLoadCount,omitempty"`

	TotalWritten *uint64 `json:"totalWritten,omitempty"`
	TotalRead    *uint64 `json:"totalRead,omitempty"`

	WrittenSession *uint64 `json:"writtenSession,omitempty"`
	ReadSession    *uint64 `json:"readSession,omitempty"`

	Sessions []*SessionAttribute `json:"sessions,omitempty"`
}
