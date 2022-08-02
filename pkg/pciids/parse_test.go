package pciids_test

import (
	"bufio"
	"bytes"
	"sync"
	"testing"

	"github.com/hertg/go-pciids/pkg/pciids"
	"github.com/stretchr/testify/assert"
)

const (
	pciids_example = `
# vendors
0001  ACME Corporation
	7a00  Generic Microcontroller
		0001 018a  Controller 1337
		0001 018b  Next-Gen Controller 3000
		71bf 001c  Foo's Microcontroller
		4200 6969  Bar Incredible Microcontroller
	7a01  Crazy Network Interface
		4200 6660  Crazy Funny Network Interface
71bf  Foo Technology LLC
	1337  Foo Reader
4200  Bar Electronics Ltd.
	6969  Funny Number Device
# classes
C 01  Mass storage controller
	00  SCSI storage controller
	08  Non-Volatile memory controller
		01  NVMHCI
		02  NVM Express
C 03  Display controller
	00  VGA compatible controller
		00  VGA controller
		01  8514 controller
`
)

var db *pciids.DB
var once sync.Once

func testDB() *pciids.DB {
	once.Do(func() {
		b := []byte(pciids_example)
		r := bytes.NewReader(b)
		s := bufio.NewScanner(r)
		db, _ = pciids.NewDB(s)
	})
	return db
}

func TestEntryCounts(t *testing.T) {
	db = testDB()
	assert.Equal(t, 3, len(db.Vendors))
	assert.Equal(t, 4, len(db.Devices))
	assert.Equal(t, 5, len(db.Subsystems))
	assert.Equal(t, 2, len(db.Classes))
}

func TestVendorParsing(t *testing.T) {
	acmeCorp, ok := db.Vendors[0x0001]
	assert.True(t, ok)
	assert.Equal(t, uint16(0x0001), acmeCorp.ID)
	assert.Equal(t, "ACME Corporation", acmeCorp.Label)
	assert.Equal(t, 2, len(acmeCorp.Devices))
	assert.Equal(t, 2, len(acmeCorp.Subsystems))

	fooLLC, ok := db.Vendors[0x71bf]
	assert.True(t, ok)
	assert.Equal(t, uint16(0x71bf), fooLLC.ID)
	assert.Equal(t, "Foo Technology LLC", fooLLC.Label)
	assert.Equal(t, 1, len(fooLLC.Devices))
	assert.Equal(t, 1, len(fooLLC.Subsystems))

	barLtd, ok := db.Vendors[0x4200]
	assert.True(t, ok)
	assert.Equal(t, uint16(0x4200), barLtd.ID)
	assert.Equal(t, "Bar Electronics Ltd.", barLtd.Label)
	assert.Equal(t, 1, len(barLtd.Devices))
	assert.Equal(t, 2, len(barLtd.Subsystems))
}

func TestClassParse(t *testing.T) {
	db = testDB()
	storage, ok := db.Classes[0x01]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x01), storage.ID)
	assert.Equal(t, "Mass storage controller", storage.Label)
	assert.Equal(t, 2, len(storage.Subclasses))

	scsi, ok := storage.Subclasses[0x00]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x00), scsi.ID)
	assert.Equal(t, "SCSI storage controller", scsi.Label)

	nvm, ok := storage.Subclasses[0x08]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x08), nvm.ID)
	assert.Equal(t, "Non-Volatile memory controller", nvm.Label)
	assert.Equal(t, 2, len(nvm.ProgrammingInterfaces))

	nvmhci, ok := nvm.ProgrammingInterfaces[0x01]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x01), nvmhci.ID)
	assert.Equal(t, "NVMHCI", nvmhci.Label)

	nvme, ok := nvm.ProgrammingInterfaces[0x02]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x02), nvme.ID)
	assert.Equal(t, "NVM Express", nvme.Label)

	display, ok := db.Classes[0x03]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x03), display.ID)
	assert.Equal(t, "Display controller", display.Label)
	assert.Equal(t, 1, len(display.Subclasses))

	vgacomp, ok := display.Subclasses[0x00]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x00), vgacomp.ID)
	assert.Equal(t, "VGA compatible controller", vgacomp.Label)
	assert.Equal(t, 2, len(vgacomp.ProgrammingInterfaces))

	vga, ok := vgacomp.ProgrammingInterfaces[0x00]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x00), vga.ID)
	assert.Equal(t, "VGA controller", vga.Label)

	num, ok := vgacomp.ProgrammingInterfaces[0x01]
	assert.True(t, ok)
	assert.Equal(t, uint8(0x01), num.ID)
	assert.Equal(t, "8514 controller", num.Label)
}
