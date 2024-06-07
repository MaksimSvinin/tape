package model

type CapAttribute struct {
	Name  string
	Value uint64
}

type StrAttribute struct {
	Name  string
	Value string
}

type MediumDensityAttribute struct {
	Name         string
	Mediumformat string
	FormattedAs  string
}

type MAMCapacityAttribute struct {
	Name           string
	Capacity       uint64
	SpaceRemaining *uint64
}

type SpecsCapacityAttribute struct {
	Name           string
	Native         int
	Compressed     int
	CompressFactor string
}

type SpecsSpeedAttribute struct {
	Name       string
	Native     int
	Compressed int
}

type SpecsPartitionsAttribute struct {
	Name            string
	PartitionNumber int
}

type SpecsPhyAttribute struct {
	Name          string
	BandsPerTape  int
	WrapsPerBand  int
	TracksPerWrap int
	Total         int
}

type SpecsDurationAttribute struct {
	Name            string
	FullTapeMinutes int
}

type SpecsAttribute struct {
	Capacity   *SpecsCapacityAttribute
	Speed      *SpecsSpeedAttribute
	Partitions *SpecsPartitionsAttribute
	Phy        *SpecsPhyAttribute
	Duration   *SpecsDurationAttribute
}

type SessionAttribute struct {
	Number  int
	Devname string
	Serial  string
}

type TapeInfo struct {
	PartCapRemain *CapAttribute
	PartCapMax    *CapAttribute

	MediumDensity *MediumDensityAttribute

	CartridgeType *StrAttribute

	AssigningOrg    *StrAttribute
	Manufacturer    *StrAttribute
	ManufactureDate *StrAttribute

	TapeLength *CapAttribute
	TapeWidth  *CapAttribute

	MAMCapacity MAMCapacityAttribute

	Specs *SpecsAttribute

	CartridgeLoadCount uint64

	TotalWritten *CapAttribute
	TotalRead    *CapAttribute

	WrittenSession *CapAttribute
	ReadSession    *CapAttribute

	Sessions []*SessionAttribute
}
