package pciids

import (
	"bufio"

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

func NewDB(scanner *bufio.Scanner) *DB {
	db := &DB{
		Vendors:    make(map[VendorID]*Vendor, SIZE_HINT_VENDORS),
		Devices:    make(map[DeviceID]*Device, SIZE_HINT_DEVICES),
		Classes:    make(map[uint8]*Class, SIZE_HINT_CLASSES),
		Subsystems: make(map[SubsystemID]*Subsystem, SIZE_HINT_SUBSYSTEMS),
	}
	cur := struct {
		vendor      *VendorID
		device      *DeviceID
		class       *uint8
		subclass    *uint8
		subsysCache map[uint][]*Subsystem
	}{
		vendor:      nil,
		device:      nil,
		class:       nil,
		subclass:    nil,
		subsysCache: make(map[uint][]*Subsystem),
	}

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
			cur.vendor = nil
			cur.device = nil
			cur.class = &class.ID
			cur.subclass = nil
		} else if line[0] >= 0x30 && line[0] <= 0x39 || line[0] >= 0x61 && line[0] <= 0x66 {
			id := conv.ParseByteNum(line[0:4])
			vendor := &Vendor{
				ID:         VendorID(id),
				Label:      string(line[6:]),
				Devices:    make(map[DeviceID]*Device),
				Subsystems: make(map[SubsystemID]*Subsystem),
			}
			if subsystems, ok := cur.subsysCache[id]; ok {
				for _, subsystem := range subsystems {
					vendor.Subsystems[subsystem.ID] = subsystem
				}
			}
			db.Vendors[VendorID(id)] = vendor
			cur.vendor = &vendor.ID
			cur.device = nil
			cur.class = nil
			cur.subclass = nil
		} else if line[0] == TAB {
			if line[1] != TAB {
				if cur.class != nil {
					id := uint8(conv.ParseByteNum(line[1:3]))
					subclass := Subclass{
						ID:                    id,
						Label:                 string(line[5:]),
						ProgrammingInterfaces: make(map[uint8]ProgrammingInterface),
					}
					db.Classes[*cur.class].Subclasses[id] = subclass
					cur.subclass = &id
				} else if cur.vendor != nil {
					id := DeviceID(uint(*cur.vendor)<<16 | conv.ParseByteNum(line[1:5]))
					device := &Device{
						ID:         id,
						Label:      string(line[7:]),
						Subsystems: make(map[SubsystemID]*Subsystem),
					}
					db.Devices[id] = device
					db.Vendors[*cur.vendor].Devices[id] = device
					cur.device = &id
				} else {
					panic("got no cursor for tabbed line")
				}
			} else {
				if cur.subclass != nil {
					id := uint8(conv.ParseByteNum(line[2:4]))
					progif := ProgrammingInterface{
						ID:    id,
						Label: string(line[6:]),
					}
					db.Classes[*cur.class].Subclasses[*cur.subclass].ProgrammingInterfaces[id] = progif
				} else if cur.device != nil {
					svid := conv.ParseByteNum(line[2:6])
					id := SubsystemID(svid<<16 | conv.ParseByteNum(line[7:11]))
					subsystem := &Subsystem{
						ID:    id,
						Label: string(line[13:]),
					}
					db.Subsystems[id] = subsystem
					db.Devices[*cur.device].Subsystems[id] = subsystem
					if vendor, ok := db.Vendors[VendorID(svid)]; ok {
						vendor.Subsystems[id] = subsystem
					} else {
						cur.subsysCache[svid] = append(cur.subsysCache[svid], subsystem)
					}
				} else {
					panic("got no cursor for double tabbed line")
				}
			}
		} else {
			panic("unexpected beginning of line")
			// continue // unexpected
		}
	}
	return db
}
