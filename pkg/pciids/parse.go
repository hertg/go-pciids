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
	cursor := struct {
		vendor      *uint16
		device      *uint32
		class       *uint8
		subclass    *uint8
		subsysCache map[uint16][]*Subsystem
	}{subsysCache: make(map[uint16][]*Subsystem)}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		} else if line[0] == HASH_SIGN {
			continue
		} else if line[0] == UPPER_C {
			id := uint8(conv.ParseByteNum(line[2:4]))
			class := &Class{
				ID:         id,
				Label:      string(line[6:]),
				Subclasses: make(map[uint8]Subclass, 16),
			}
			db.Classes[id] = class
			cursor.vendor = nil
			cursor.device = nil
			cursor.class = &class.ID
			cursor.subclass = nil
		} else if line[0] >= 0x30 && line[0] <= 0x39 || line[0] >= 0x61 && line[0] <= 0x66 {
			id := uint16(conv.ParseByteNum(line[0:4]))
			vendor := &Vendor{
				ID:         id,
				Label:      string(line[6:]),
				Devices:    make(map[uint32]*Device),
				Subsystems: make(map[uint32]*Subsystem),
			}
			if subsystems, ok := cursor.subsysCache[id]; ok {
				for _, subsystem := range subsystems {
					vendor.Subsystems[subsystem.ID] = subsystem
				}
			}
			db.Vendors[id] = vendor
			cursor.vendor = &vendor.ID
			cursor.device = nil
			cursor.class = nil
			cursor.subclass = nil
		} else if line[0] == TAB {
			if line[1] != TAB {
				if cursor.class != nil {
					id := uint8(conv.ParseByteNum(line[1:3]))
					subclass := Subclass{
						ID:                    id,
						Label:                 string(line[5:]),
						ProgrammingInterfaces: make(map[uint8]ProgrammingInterface),
					}
					db.Classes[*cursor.class].Subclasses[id] = subclass
					cursor.subclass = &id
				} else if cursor.vendor != nil {
					id := uint32(*cursor.vendor)<<16 | uint32(conv.ParseByteNum(line[1:5]))
					device := &Device{
						ID:         id,
						Label:      string(line[7:]),
						Subsystems: make(map[uint32]*Subsystem),
					}
					db.Devices[id] = device
					db.Vendors[*cursor.vendor].Devices[id] = device
					cursor.device = &id
				} else {
					return nil, fmt.Errorf("parsing error: cursor is on a tabbed line, but without any cursor context info")
				}
			} else {
				if cursor.subclass != nil {
					id := uint8(conv.ParseByteNum(line[2:4]))
					progif := ProgrammingInterface{
						ID:    id,
						Label: string(line[6:]),
					}
					db.Classes[*cursor.class].Subclasses[*cursor.subclass].ProgrammingInterfaces[id] = progif
				} else if cursor.device != nil {
					svid := uint16(conv.ParseByteNum(line[2:6]))
					id := uint32(svid)<<16 | uint32(conv.ParseByteNum(line[7:11]))
					subsystem := &Subsystem{
						ID:    id,
						Label: string(line[13:]),
					}
					db.Subsystems[id] = subsystem
					db.Devices[*cursor.device].Subsystems[id] = subsystem
					if vendor, ok := db.Vendors[svid]; ok {
						vendor.Subsystems[id] = subsystem
					} else {
						cursor.subsysCache[svid] = append(cursor.subsysCache[svid], subsystem)
					}
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
