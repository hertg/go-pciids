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

func NewDB(scanner *bufio.Scanner) *DB {
	db := &DB{
		Vendors:    make(map[VendorID]Vendor),
		Devices:    make(map[DeviceID]Device),
		Classes:    make(map[uint8]Class),
		Subsystems: make(map[SubsystemID]Subsystem),
	}
	var vcur1 *VendorID // vendor cursor 1
	var vcur2 *DeviceID // vendor cursor 2
	var ccur1 *uint8    // class cursor 1
	var ccur2 *uint8    // class cursor 2
	subsystemCache := make(map[VendorID][]*Subsystem)

	for scanner.Scan() {
		line := scanner.Bytes()
		length := len(line)

		if length == 0 {
			continue
		} else if line[0] == HASH_SIGN {
			continue
		} else if line[0] == UPPER_C {
			id := uint8(conv.ParseByteNum(line[2:4]))
			class := Class{
				ID:         id,
				Label:      string(line[6:]),
				Subclasses: make(map[uint8]Subclass, 16),
			}
			db.Classes[id] = class
			ccur1 = &class.ID
			ccur2 = nil
			vcur1 = nil
			vcur2 = nil
		} else if line[0] >= 0x30 && line[0] <= 0x39 || line[0] >= 0x61 && line[0] <= 0x66 {
			id := VendorID(conv.ParseByteNum(line[0:4]))
			vendor := Vendor{
				ID:         id,
				Label:      string(line[6:]),
				Devices:    make(map[DeviceID]*Device),
				Subsystems: make(map[SubsystemID]*Subsystem),
			}
			if subsystems, ok := subsystemCache[id]; ok {
				for _, subsystem := range subsystems {
					vendor.Subsystems[subsystem.ID] = subsystem
				}
			}
			db.Vendors[id] = vendor
			vcur1 = &vendor.ID
			vcur2 = nil
			ccur1 = nil
			ccur2 = nil
		} else if line[0] == TAB {
			if line[1] != TAB {
				if ccur1 != nil {
					id := uint8(conv.ParseByteNum(line[1:3]))
					subclass := Subclass{
						ID:                    id,
						Label:                 string(line[5:]),
						ProgrammingInterfaces: make(map[uint8]ProgrammingInterface),
					}
					db.Classes[*ccur1].Subclasses[id] = subclass
					ccur2 = &id
				} else if vcur1 != nil {
					id := DeviceID(uint(*vcur1)<<16 | conv.ParseByteNum(line[1:5]))
					device := Device{
						ID:         id,
						Label:      string(line[7:]),
						Subsystems: make(map[SubsystemID]*Subsystem),
					}
					db.Devices[id] = device
					db.Vendors[*vcur1].Devices[id] = &device
					vcur2 = &id
				} else {
					panic("got no cursor for tabbed line")
				}
			} else {
				if ccur2 != nil {
					id := uint8(conv.ParseByteNum(line[2:4]))
					progif := ProgrammingInterface{
						ID:    id,
						Label: string(line[6:]),
					}
					db.Classes[*ccur1].Subclasses[*ccur2].ProgrammingInterfaces[id] = progif
				} else if vcur2 != nil {
					svid := conv.ParseByteNum(line[2:6])
					id := SubsystemID(svid<<16 | conv.ParseByteNum(line[7:11]))
					subsystem := Subsystem{
						ID:    id,
						Label: string(line[13:]),
					}
					db.Subsystems[id] = subsystem
					db.Devices[*vcur2].Subsystems[id] = &subsystem
					if vendor, ok := db.Vendors[VendorID(svid)]; ok {
						vendor.Subsystems[id] = &subsystem
					} else {
						subsystemCache[VendorID(svid)] = append(subsystemCache[VendorID(svid)], &subsystem)
					}
				} else {
					panic("got no cursor for double tabbed line")
				}
			}
		} else {
			continue // unexpected
		}
	}
	return db
}
