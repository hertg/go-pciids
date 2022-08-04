
<div align="center">
  <h1><strong>go-pciids</strong></h1>
  <p>
		<strong>A minimal, zero-dependency library to parse <a href="https://pci-ids.ucw.cz/">pci.ids</a> files.</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/hertg/go-pciids">
      <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hertg/go-pciids" />
    </a>
    <a href="#">
			<img alt="License Information" src="https://img.shields.io/github/license/hertg/go-pciids">
    </a>
  </p>
</div>

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
// example for creating a scanner
filepath := "/usr/share/hwdata/pci.ids"
file, _ := os.Open(filepath)
scanner := bufio.NewScanner(file)

// actual usage
db, err := pciids.NewDB(scanner)
```

### Query

Names and labels of vendors, devices, and classes can be easily
retrieved by using the `Find***Label` methods available inside the
`pciids` package.

```go
db.FindClassLabel(0x03) 							// -> 'Display controller'
db.FindSubclassLabel(0x03, 0x00) 			// -> 'VGA compatible controller'
db.FindVendorLabel(0x1002) 						// -> 'Advanced Micro Devices, Inc. [AMD/ATI]'
db.FindDeviceLabel(0x1002, 0x73bf) 		// -> 'Navi 21 [Radeon RX 6800/6800 XT / 6900 XT]'
db.FindSubsystemLabel(0x148c, 0x2408) // -> 'Red Devil AMD Radeon RX 6900 XT'
```

The `DB` can also be traversed manually, 
see a quick overview of the DB structure.

```text
DB
├─ vendors
│  ├─ devices
│  │  └─ subsystems
│  └─ subsystems 
├─ devices
│  └─ subsystems 
├─ subsystems
└─ classes
   └─ subclasses
      └─ progifs
```

## Comparison

It parses the vendor, device, and class IDs directly to
unsigned integers. This prevents unnecessary string allocations
and significantly improves performance.

It has been found that this library parses at 2-3 times the speed
while using roughly half the amount memory in comparison to [jaypipes/pcidb](https://github.com/jaypipes/pcidb).

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
