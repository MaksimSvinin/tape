package model

type Operation int

func (o Operation) String() string {
	return [...]string{"Unknown", "Info", "Write", "Read", "Erase", "Rewind", "Eject"}[o]
}

const (
	Unknown Operation = iota
	Info
	Write
	Read
	Erase
	Rewind
	Eject
)
