package pciids_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClassLabel(t *testing.T) {
	db := testDB()
	label := db.FindClassLabel(0x03)
	assert.Equal(t, "Display controller", *label)
}

func TestFindSubclassLabel(t *testing.T) {
	db := testDB()
	label := db.FindSubclassLabel(0x03, 0x00)
	assert.Equal(t, "VGA compatible controller", *label)

	// providing an inexistent subclass,
	// should return the label of the class instead
	label = db.FindSubclassLabel(0x03, 0x99)
	assert.Equal(t, "Display controller", *label)
}
