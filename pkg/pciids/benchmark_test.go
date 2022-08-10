package pciids_test

var benchFileName = "/usr/share/hwdata/pci.ids"

// func BenchmarkGoPCIIDS(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		f, err := os.Open(benchFileName)
// 		if err != nil {
// 			b.Skipf("unable to open pci.ids file '%s' for benchmark: %s", benchFileName, err)
// 		}
// 		s := bufio.NewScanner(f)
// 		pciids.NewDB(s)
// 	}
// }

// func BenchmarkPCIDB(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, err := pcidb.New(&pcidb.WithOption{Path: &benchFileName})
// 		if err != nil {
// 			b.Skipf("unable to initialize pcidb: %s", err)
// 		}
// 	}
// }
