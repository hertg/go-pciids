package pciids

import (
	"bufio"
	"fmt"

	"github.com/hertg/go-pciids/internal/conv"
)

const (
	TAB       byte = 0x09
	HASH_SIGN byte = 0x23
	UPPER_C   byte = 0x43
)

const (
	SIZE_HINT_VENDORS    = 2300  // actual length at 2022-08-01: 2289
	SIZE_HINT_DEVICES    = 18000 // actual length at 2022-08-01: 17207
	SIZE_HINT_CLASSES    = 22    // actual length at 2022-08-01: 22
	SIZE_HINT_SUBSYSTEMS = 10000 // actual length at 2022-08-01: 9741
)

type cursor struct {
	vendor      *uint16
	device      *uint32
	class       *uint8
	subclass    *uint8
	subsysCache map[uint16][]*Subsystem
}

// NewDB Parse a pci.ids file with the provided scanner and return
// a pointer to pciids.DB which can be used for querying.
// May return an error if an unexpected line is encountered while parsing.
func NewDB(scanner *bufio.Scanner) (*DB, error) {
	db := &DB{
		Vendors:    make(map[uint16]*Vendor, SIZE_HINT_VENDORS),
		Devices:    make(map[uint32]*Device, SIZE_HINT_DEVICES),
		Classes:    make(map[uint8]*Class, SIZE_HINT_CLASSES),
		Subsystems: make(map[uint32]*Subsystem, SIZE_HINT_SUBSYSTEMS),
	}
	cur := &cursor{
		subsysCache: make(map[uint16][]*Subsystem),
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || line[0] == HASH_SIGN {
			continue
		} else if line[0] == UPPER_C {
			parseClassLine(line, db, cur)
		} else if line[0] >= 0x30 && line[0] <= 0x39 || line[0] >= 0x61 && line[0] <= 0x66 {
			parseVendorLine(line, db, cur)
		} else if line[0] == TAB {
			if line[1] != TAB {
				if cur.class != nil {
					parseSubclassLine(line, db, cur)
				} else if cur.vendor != nil {
					parseDeviceLine(line, db, cur)
				} else {
					return nil, fmt.Errorf("parsing error: cursor is on a tabbed line, but without any cursor context info")
				}
			} else {
				if cur.subclass != nil {
					parseProgifLine(line, db, cur)
				} else if cur.device != nil {
					parseSubsystemLine(line, db, cur)
				} else {
					return nil, fmt.Errorf("parsing error: cursor is on a double tabbed line, but without any cursor context info")
				}
			}
		} else {
			return nil, fmt.Errorf("unexpected beginning of line: %x", line[0])
		}
	}
	return db, nil
}

func parseClassLine(line []byte, db *DB, cur *cursor) {
	id := uint8(conv.ParseByteNum(line[2:4]))
	class := &Class{
		ID:         id,
		Label:      string(line[6:]),
		Subclasses: make(map[uint8]Subclass, 16),
	}
	db.Classes[id] = class
	cur.vendor = nil
	cur.device = nil
	cur.class = &class.ID
	cur.subclass = nil
}

func parseVendorLine(line []byte, db *DB, cur *cursor) {
	id := uint16(conv.ParseByteNum(line[0:4]))
	vendor := &Vendor{
		ID:         id,
		Label:      string(line[6:]),
		Devices:    make(map[uint32]*Device),
		Subsystems: make(map[uint32]*Subsystem),
	}
	if subsystems, ok := cur.subsysCache[id]; ok {
		for _, subsystem := range subsystems {
			vendor.Subsystems[subsystem.ID] = subsystem
		}
	}
	db.Vendors[id] = vendor
	cur.vendor = &vendor.ID
	cur.device = nil
	cur.class = nil
	cur.subclass = nil
}

func parseSubclassLine(line []byte, db *DB, cur *cursor) {
	id := uint8(conv.ParseByteNum(line[1:3]))
	subclass := Subclass{
		ID:                    id,
		Label:                 string(line[5:]),
		ProgrammingInterfaces: make(map[uint8]ProgrammingInterface),
	}
	db.Classes[*cur.class].Subclasses[id] = subclass
	cur.subclass = &id
}

func parseDeviceLine(line []byte, db *DB, cur *cursor) {
	id := uint32(*cur.vendor)<<16 | uint32(conv.ParseByteNum(line[1:5]))
	device := &Device{
		ID:         id,
		Label:      string(line[7:]),
		Subsystems: make(map[uint32]*Subsystem),
	}
	db.Devices[id] = device
	db.Vendors[*cur.vendor].Devices[id] = device
	cur.device = &id
}

func parseProgifLine(line []byte, db *DB, cur *cursor) {
	id := uint8(conv.ParseByteNum(line[2:4]))
	progif := ProgrammingInterface{
		ID:    id,
		Label: string(line[6:]),
	}
	db.Classes[*cur.class].Subclasses[*cur.subclass].ProgrammingInterfaces[id] = progif
}

func parseSubsystemLine(line []byte, db *DB, cur *cursor) {
	svid := uint16(conv.ParseByteNum(line[2:6]))
	id := uint32(svid)<<16 | uint32(conv.ParseByteNum(line[7:11]))
	subsystem := &Subsystem{
		ID:    id,
		Label: string(line[13:]),
	}
	db.Subsystems[id] = subsystem
	db.Devices[*cur.device].Subsystems[id] = subsystem
	if vendor, ok := db.Vendors[svid]; ok {
		vendor.Subsystems[id] = subsystem
	} else {
		cur.subsysCache[svid] = append(cur.subsysCache[svid], subsystem)
	}
}
