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

// FindProgifLabel Get the label for the programming interface by providing
// the class, subclass, and programming interface. If there exists a label
// for the progif, it will be returned directly. Otherwise the label of the
// subclass is returned. If that also doesn't exist, the label of the class is returned.
// exists no label for the subclass, but one for the class,
// then the label of the class is returned.
// May return nil if no matching label can be found.
func (db *DB) FindProgifLabel(class uint8, subclass uint8, progif uint8) *string {
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			if p, exists := s.ProgrammingInterfaces[progif]; exists {
				return &p.Label
			}
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

// FindVendorLabel Get the label for the given vendor.
// May return nil if no vendor name can be found for this vendor id.
func (db *DB) FindVendorLabel(vendor uint16) *string {
	if v := db.FindVendor(vendor); v != nil {
		return &v.Label
	}
	return nil
}

// FindDeviceLabel Get the label for the given vendor/device.
// May return nil if no vendor/device can be found for these parameters.
func (db *DB) FindDeviceLabel(vendor uint16, device uint16) *string {
	if d := db.FindDevice(vendor, device); d != nil {
		return &d.Label
	}
	return nil
}

// FindSubsystemLabel Get the label for the provided subsystem.
// May return nil if no subsystem can be found for these parameters.
func (db *DB) FindSubsystemLabel(subvendor uint16, subsystem uint16) *string {
	if sd := db.FindSubsystem(subvendor, subsystem); sd != nil {
		return &sd.Label
	}
	return nil
}

// FindClass Get the class for the provided class id.
// May return nil if no class can be found for this id.
func (db *DB) FindClass(class uint8) *Class {
	if c, exists := db.Classes[class]; exists {
		return c
	}
	return nil
}

// FindSubclass Get the subclass for the provided class, and subclass ids.
// May return nil if no subclass can be found for these parameters.
// Note: will always return nil if the subclass doesn't exist, even if the class does.
func (db *DB) FindSubclass(class uint8, subclass uint8) *Subclass {
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			return &s
		}
	}
	return nil
}

// FindProgrammingInterface Get the programming interface (progif) for the provided class, subclass, and progif ids.
// May return nil if no programming interface can be found for these parameters.
// Note: will always return nil if the progif doesn't exist, even if the class and subclass do.
func (db *DB) FindProgrammingInterface(class uint8, subclass uint8, progif uint8) *ProgrammingInterface {
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			if p, exists := s.ProgrammingInterfaces[progif]; exists {
				return &p
			}
		}
	}
	return nil
}

// FindVendor Get the vendor for the provided vendor id.
// May return nil if no vendor can be found for this id.
func (db *DB) FindVendor(vendor uint16) *Vendor {
	if v, exists := db.Vendors[vendor]; exists {
		return v
	}
	return nil
}

// FindDevice Get the device for the provided vendor/device id.
// May return nil if no device can be found for these parameters.
func (db *DB) FindDevice(vendor uint16, device uint16) *Device {
	id := uint32(vendor)<<16 | uint32(device)
	if d, exists := db.Devices[id]; exists {
		return d
	}
	return nil
}

// FindSubsystem Get the subsystem for the provided subvendor/subdevice ids.
// May return nil if no subsystem can be found for these parameters.
func (db *DB) FindSubsystem(subvendor uint16, subdevice uint16) *Subsystem {
	id := uint32(subvendor)<<16 | uint32(subdevice)
	if sd, exists := db.Subsystems[id]; exists {
		return sd
	}
	return nil
}
