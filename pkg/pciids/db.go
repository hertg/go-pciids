package pciids

type Vendor struct {
	ID         uint16
	Label      string
	Devices    map[uint32]*Device
	Subsystems map[uint32]*Subsystem
}

type Device struct {
	ID         uint32
	Label      string
	Subsystems map[uint32]*Subsystem
}

type Subsystem struct {
	ID    uint32
	Label string
}

type Class struct {
	ID         uint8
	Label      string
	Subclasses map[uint8]Subclass
}

type Subclass struct {
	ID                    uint8
	Label                 string
	ProgrammingInterfaces map[uint8]ProgrammingInterface
}

type ProgrammingInterface struct {
	ID    uint8
	Label string
}

type DB struct {
	Vendors    map[uint16]*Vendor
	Devices    map[uint32]*Device
	Subsystems map[uint32]*Subsystem
	Classes    map[uint8]*Class
}
