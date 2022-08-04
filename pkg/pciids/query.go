package pciids

// FindClassLabel Get the label for the provided class.
// May return nil if no matching label can be found.
func (db *DB) FindClassLabel(class uint8) *string {
	if c, exists := db.Classes[class]; exists {
		return &c.Label
	}
	return nil
}

// FindSubclassLabel Get the label for the subclass by providing
// the class and the subclass. If there exists a label
// for the subclass, it will be returned directly. If there
// exists no label for the subclass, but one for the class,
// then the label of the class is returned.
// May return nil if no matching label can be found.
func (db *DB) FindSubclassLabel(class uint8, subclass uint8) *string {
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			return &s.Label
		}
		return &c.Label
	}
	return nil
}

// FindSubclassLabelAlt Similar to the FindClassLabel method,
// but allows passing a 'combined' number of class and subclass.
// Eg. You can pass in (0x0300) here, instead of passing (0x03, 0x00)
// into FindSubclassLabel
// func (db *DB) FindSubclassLabelAlt(cls uint16) *string {
// 	class := uint8(cls >> 8 & 0b111111)
// 	subclass := uint8(cls & 0b111111)
// 	return db.FindSubclassLabel(class, subclass)
// }

func (db *DB) FindVendorLabel(vendor uint16) *string {
	if v := db.FindVendor(vendor); v != nil {
		return &v.Label
	}
	return nil
}

func (db *DB) FindDeviceLabel(vendor uint16, device uint16) *string {
	if d := db.FindDevice(vendor, device); d != nil {
		return &d.Label
	}
	return nil
}

func (db *DB) FindSubsystemLabel(vendor uint16, subsystem uint16) *string {
	if sd := db.FindSubsystem(vendor, subsystem); sd != nil {
		return &sd.Label
	}
	return nil
}

func (db *DB) FindClass(class uint8) *Class {
	if c, exists := db.Classes[class]; exists {
		return c
	}
	return nil
}

func (db *DB) FindSubclass(class uint8, subclass uint8) *Subclass {
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			return &s
		}
	}
	return nil
}

func (db *DB) FindVendor(vendor uint16) *Vendor {
	if v, exists := db.Vendors[vendor]; exists {
		return v
	}
	return nil
}

func (db *DB) FindDevice(vendor uint16, device uint16) *Device {
	id := uint32(vendor)<<16 | uint32(device)
	if d, exists := db.Devices[id]; exists {
		return d
	}
	return nil
}

func (db *DB) FindSubsystem(subvendor uint16, subdevice uint16) *Subsystem {
	id := uint32(subvendor)<<16 | uint32(subdevice)
	if sd, exists := db.Subsystems[id]; exists {
		return sd
	}
	return nil
}
