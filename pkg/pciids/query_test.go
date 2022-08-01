package pciids_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFullClass(t *testing.T) {
	db := testDB()
	label := db.FindSubclass(0x0300)
	assert.Equal(t, "VGA compatible controller", *label)
}
