package pciids

// unstable
func (db *DB) FindSubclass(cls uint16) *string {
	class := uint8(cls >> 8 & 0b111111)
	subclass := uint8(cls & 0b111111)
	if c, exists := db.Classes[class]; exists {
		if s, exists := c.Subclasses[subclass]; exists {
			return &s.Label
		}
	}
	return nil
}
