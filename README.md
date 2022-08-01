# go-pciids

A minimal, zero-dependency library to parse [pci.ids](https://pci-ids.ucw.cz/) files.

## Description

The library provides methods to parse the contents of a `pci.ids` file
and export methods to query the parsed data. It has been created
as a stripped down, but more performant alternative to [jaypipes/pcidb](https://github.com/jaypipes/pcidb).
Because of this minimalistic approach only functions for parsing and
querying are provided, without any extra features (see [Limitations](#limitations)).

## Limitations

The following features are **not** included:

- No caching
- No fetching of `pci.ids` via network
- No search for `pci.ids` files locally

## Usage

### Install

```shell
go get github.com/hertg/go-pciids
```

### Parse

```go
filepath := "/usr/share/hwdata/pci.ids"
file, _ := os.Open(filepath)
scanner := bufio.NewScanner(file)

// you need to provide a *bufio.Scanner for the pci.ids yourself
db, err := pciids.NewDB(scanner)
```

### Query
TODO

## Comparison

It parses the vendor, device, and class IDs directly to
unsigned integers. This prevents unnecessary string allocations
and significantly improves performance.

It has been found that this library parses at 2-3 times the speed
while using roughly half the amount memory in comparison.

<details>
	<summary>Open Benchmark</summary>

	```text
	goos: linux
	goarch: amd64
	pkg: github.com/hertg/go-pciids/pkg/pciids
	cpu: AMD Ryzen 9 5950X 16-Core Processor
	BenchmarkGoPCIIDS-32    	     100	  16775134 ns/op	 6400548 B/op	  116541 allocs/op
	BenchmarkPCIDB-32       	      37	  29394779 ns/op	11972032 B/op	  184512 allocs/op
	PASS
	ok  	github.com/hertg/go-pciids/pkg/pciids	2.813s
	```
</details>
