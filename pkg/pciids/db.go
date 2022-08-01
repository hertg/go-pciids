package pciids

type VendorID uint16
type Vendor struct {
	ID         VendorID
	Label      string
	Devices    map[DeviceID]*Device
	Subsystems map[SubsystemID]*Subsystem
}

type DeviceID uint32
type Device struct {
	ID         DeviceID
	Label      string
	Subsystems map[SubsystemID]*Subsystem
}

type SubsystemID uint32
type Subsystem struct {
	ID    SubsystemID
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
	Vendors    map[VendorID]*Vendor
	Devices    map[DeviceID]*Device
	Subsystems map[SubsystemID]*Subsystem
	Classes    map[uint8]*Class
}
