# go-pciids

## Alternatives

This library is an alternative to the more widespread
[jaypipes/pcidb](https://github.com/jaypipes/pcidb).

It has been created with the goal to use unsigned integers
instead of strings for the identifiers, and to be more performant
by preventing unnecessary string allocations during parsing.

<details>
	<summary>Open Benchmark</summary>

	```text
	goos: linux
	goarch: amd64
	pkg: github.com/hertg/go-pciids/pkg/pciids
	cpu: AMD Ryzen 9 5950X 16-Core Processor
	BenchmarkGoPCIIDS
	BenchmarkGoPCIIDS-32    	     100	  16293188 ns/op	 7654091 B/op	  116536 allocs/op
	BenchmarkPCIDB
	BenchmarkPCIDB-32       	      36	  30413482 ns/op	11971077 B/op	  184507 allocs/op
	PASS
	ok  	github.com/hertg/go-pciids/pkg/pciids	2.779s
	```
</details>
