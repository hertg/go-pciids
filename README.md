# go-pciids

## Alternatives

This library is an alternative to the more widespread
[jaypipes/pcidb](https://github.com/jaypipes/pcidb).

It has been created with the goal to use unsigned integers
instead of strings for the identifiers, and possibly be
more performant without the unnecessary string allocations.

```text
goos: linux
goarch: amd64
cpu: AMD Ryzen 9 5950X 16-Core Processor
BenchmarkGoPCIIDS-32    	      81	  21009679 ns/op	 9968074 B/op	  117197 allocs/op
BenchmarkPCIDB-32       	      44	  30055181 ns/op	11952707 B/op	  184159 allocs/op
PASS
```
