package pciids_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/hertg/go-pciids/pkg/pciids"
	"github.com/jaypipes/pcidb"
)

var benchFileName = "/usr/share/hwdata/pci.ids"

func BenchmarkGoPCIIDS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open(benchFileName)
		if err != nil {
			b.Skipf("unable to open pci.ids file '%s' for benchmark: %s", benchFileName, err)
		}
		s := bufio.NewScanner(f)
		pciids.NewDB(s)
	}
}

func BenchmarkPCIDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pcidb.New(&pcidb.WithOption{Path: &benchFileName})
		if err != nil {
			b.Skipf("unable to initialize pcidb: %s", err)
		}
	}
}

// func BenchmarkCustom(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		f, _ := os.Open(benchFileName)
// 		s := bufio.NewScanner(f)

// 		db := &pciids.DB{
// 			Vendors:    make(map[pciids.VendorID]pciids.Vendor),
// 			Devices:    make(map[pciids.DeviceID]pciids.Device),
// 			Classes:    make(map[uint8]pciids.Class),
// 			Subsystems: make(map[pciids.SubsystemID]pciids.Subsystem),
// 		}
// 		subsystemCache := make(map[pciids.VendorID][]*pciids.Subsystem)

// 		for s.Scan() {
// 			line := s.Bytes()
// 			_ = db
// 			_ = subsystemCache
// 			_ = line
// 		}
// 	}
// }
