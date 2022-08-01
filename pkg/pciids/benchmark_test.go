package pciids_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/hertg/go-pciids/pkg/pciids"
	"github.com/jaypipes/pcidb"
)

var benchFileName = "/usr/share/hwdata/pci.ids"

func init() {
	// prepare benchmark file
	f, _ := os.CreateTemp("", "pciids")
	defer f.Close()
	b := []byte(pciids_example)
	benchFileName = f.Name()
	f.Write(b)
}

func BenchmarkGoPCIIDS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open(benchFileName)
		s := bufio.NewScanner(f)
		pciids.NewDB(s)
	}
}

func BenchmarkPCIDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pcidb.New(&pcidb.WithOption{Path: &benchFileName})
	}
}
