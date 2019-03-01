# IkeaPack [![Build Status](https://travis-ci.org/ikkerens/ikeapack.svg?branch=master)](https://travis-ci.org/ikkerens/ikeapack) [![Go Report Card](https://goreportcard.com/badge/github.com/ikkerens/ikeapack)](https://goreportcard.com/report/github.com/ikkerens/ikeapack) [![GoDoc](https://godoc.org/github.com/ikkerens/ikeapack?status.svg)](https://godoc.org/github.com/ikkerens/ikeapack)

> Named IkeaPack because it compacts structs in a very compact and packed manner, and if you don't know how to reassemble it, it may just look like a random blob of parts. Just like ikea products!

(If anyone from ikea doesn't like me using their name, just reach out to me and I'll change it, no copyright/trademark infringement intended.)

This is a packed struct serializer that is mostly meant for a private project but was released as it may be useful to someone else.

Originally this package was made as an extension to binary.Read and binary.Write, but I soon found those functions didn't match my use case as they offered no support for strings nor compression.

#### Features
* Caches types for faster calls to the same type
* Compression support
* Tread safe (the calls are, reading to the value is not)
* Easy to implement in other languages
* Supported types:
  * uint8 (and byte) up to uint64
  * int8 up to int64
  * float32 and float64
  * string
  * anything implementing the Packer/Unpacker interfaces
  * slices
  * structs

#### Format
* All primitives are stored in big endian format
* All slices are stored with a uint32 prefix indicating their length
* Strings are stored with a uint32 prefix indicating their length
* Compression blocks are stored using deflate (level 9) with a uint32 prefixing the size of the compressed data blob

#### Note about int/uint
The types `int` and `uint` are not supported because their actual sizes depend on the compiler architecture.  
Instead, be explicit and use int32/int64/uint32/uint64.

## Include in your project
```go
import "github.com/ikkerens/ikeapack"
```

## Usage
```go
package main

import (
	"bytes"
	"log"

	"github.com/ikkerens/ikeapack"
)

type myBlob struct {
	A uint64  // all fields have to be exported
	B []byte  `ikea:"compress:9"` // this field will be packed and compressed, with flate level 5
	C subBlob // If you omit the level `ikea:"compress"`, level 9 will be assumed.
	D int32
}

type subBlob struct {
	D string
}

func main() {
	b := new(bytes.Buffer)
	blob := &myBlob{A: 1, B: []byte{1, 2, 3, 4}, C: subBlob{D: "test message"}}

	// Pack
	if err := ikea.Pack(b, blob); err != nil { // Write does not need a pointer, but it is recommended
		log.Fatalln(err)
	}

	// Unpack
	newBlob := new(myBlob)
	if err := ikea.Unpack(b, newBlob); err != nil { // Read *needs* a pointer, or it will panic
		log.Fatalln(err)
	}

	log.Printf("Successfully unpacked: %+v", newBlob)
}
```

## Benchmarks
These benchmarks can be found in [alecthomas](https://github.com/alecthomas)'s [go serialization benchmarks](https://github.com/alecthomas/go_serialization_benchmarks).
While not all benchmarks are included since not all dependencies could resolve, these give a good overview of the performance of this lib vs the others.  
Note that this library does not have a focus on *being* the fastest in any way, as this was made to cover a specific use-case. But it does strive to be as fast as it can be.

These benchmarks were executed on a Dell laptop with an i7-8550U cpu and 16GB of ram.

```
BenchmarkIkeaMarshal-8                           3000000               505 ns/op              72 B/op          8 allocs/op
BenchmarkIkeaUnmarshal-8                         2000000               670 ns/op             160 B/op         11 allocs/op
BenchmarkJsonMarshal-8                            500000              3785 ns/op            1224 B/op          9 allocs/op
BenchmarkJsonUnmarshal-8                          300000              4412 ns/op             464 B/op          7 allocs/op
BenchmarkEasyJsonMarshal-8                       1000000              1559 ns/op             784 B/op          5 allocs/op
BenchmarkEasyJsonUnmarshal-8                     1000000              1363 ns/op             160 B/op          4 allocs/op
BenchmarkBsonMarshal-8                           1000000              1433 ns/op             392 B/op         10 allocs/op
BenchmarkBsonUnmarshal-8                         1000000              1928 ns/op             244 B/op         19 allocs/op
BenchmarkGobMarshal-8                            2000000               930 ns/op              48 B/op          2 allocs/op
BenchmarkGobUnmarshal-8                          2000000               943 ns/op             112 B/op          3 allocs/op
BenchmarkXdrMarshal-8                            1000000              1740 ns/op             456 B/op         21 allocs/op
BenchmarkXdrUnmarshal-8                          1000000              1449 ns/op             240 B/op         11 allocs/op
BenchmarkUgorjiCodecMsgpackMarshal-8             1000000              1141 ns/op             561 B/op          6 allocs/op
BenchmarkUgorjiCodecMsgpackUnmarshal-8           1000000              1349 ns/op             449 B/op          6 allocs/op
BenchmarkSerealMarshal-8                          500000              2680 ns/op             912 B/op         21 allocs/op
BenchmarkSerealUnmarshal-8                        500000              2943 ns/op            1008 B/op         34 allocs/op
BenchmarkBinaryMarshal-8                         1000000              1427 ns/op             334 B/op         20 allocs/op
BenchmarkBinaryUnmarshal-8                       1000000              1554 ns/op             336 B/op         22 allocs/op
BenchmarkHproseMarshal-8                         2000000               971 ns/op             479 B/op          8 allocs/op
BenchmarkHproseUnmarshal-8                       1000000              1140 ns/op             320 B/op         10 allocs/op
BenchmarkGoAvroMarshal-8                          500000              2561 ns/op            1030 B/op         31 allocs/op
BenchmarkGoAvroUnmarshal-8                        200000              6346 ns/op            3437 B/op         87 allocs/op
BenchmarkGoAvro2TextMarshal-8                     500000              2875 ns/op            1326 B/op         20 allocs/op
BenchmarkGoAvro2TextUnmarshal-8                   500000              2690 ns/op             807 B/op         34 allocs/op
BenchmarkGoAvro2BinaryMarshal-8                  2000000               916 ns/op             510 B/op         11 allocs/op
BenchmarkGoAvro2BinaryUnmarshal-8                2000000               979 ns/op             576 B/op         13 allocs/op
BenchmarkProtobufMarshal-8                       2000000               984 ns/op             200 B/op          7 allocs/op
BenchmarkProtobufUnmarshal-8                     2000000               827 ns/op             192 B/op         10 allocs/op
```

Below you will find some of the benchmarks that do not rely on Go's reflection at runtime, which creates a significant performance boost, but for diligence it's worth mentioning here regardless.
```
BenchmarkMsgpMarshal-8                          10000000               178 ns/op             128 B/op          1 allocs/op
BenchmarkMsgpUnmarshal-8                         5000000               340 ns/op             112 B/op          3 allocs/op
BenchmarkFlatBuffersMarshal-8                    5000000               341 ns/op               0 B/op          0 allocs/op
BenchmarkFlatBuffersUnmarshal-8                  5000000               249 ns/op             112 B/op          3 allocs/op
BenchmarkCapNProtoMarshal-8                      3000000               483 ns/op              56 B/op          2 allocs/op
BenchmarkCapNProtoUnmarshal-8                    3000000               438 ns/op             200 B/op          6 allocs/op
BenchmarkCapNProto2Marshal-8                     2000000               723 ns/op             244 B/op          3 allocs/op
BenchmarkCapNProto2Unmarshal-8                   1000000              1019 ns/op             320 B/op          6 allocs/op
BenchmarkGoprotobufMarshal-8                     3000000               396 ns/op              96 B/op          2 allocs/op
BenchmarkGoprotobufUnmarshal-8                   2000000               614 ns/op             200 B/op         10 allocs/op
BenchmarkGogoprotobufMarshal-8                  10000000               162 ns/op              64 B/op          1 allocs/op
BenchmarkGogoprotobufUnmarshal-8                10000000               223 ns/op              96 B/op          3 allocs/op
BenchmarkColferMarshal-8                        10000000               133 ns/op              64 B/op          1 allocs/op
BenchmarkColferUnmarshal-8                      10000000               187 ns/op             112 B/op          3 allocs/op
BenchmarkGencodeMarshal-8                       10000000               179 ns/op              80 B/op          2 allocs/op
BenchmarkGencodeUnmarshal-8                     10000000               202 ns/op             112 B/op          3 allocs/op
BenchmarkGencodeUnsafeMarshal-8                 20000000               102 ns/op              48 B/op          1 allocs/op
BenchmarkGencodeUnsafeUnmarshal-8               10000000               146 ns/op              96 B/op          3 allocs/op
BenchmarkXDR2Marshal-8                          10000000               162 ns/op              64 B/op          1 allocs/op
BenchmarkXDR2Unmarshal-8                        10000000               136 ns/op              32 B/op          2 allocs/op
```