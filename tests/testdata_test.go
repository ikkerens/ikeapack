package tests

import (
	"encoding/hex"
	"io"

	"github.com/ikkerens/ikeapack"
)

/* Test types */

type testStruct struct {
	TestBool          bool
	TestByte          byte
	TestUint8         uint8
	TestUint16        uint16
	TestUint32        uint32
	TestUint64        uint64
	TestInt8          int8
	TestInt16         int16
	TestInt32         int32
	TestInt64         int64
	TestFloat32       float32
	TestFloat64       float64
	TestString        string
	TestSubStruct     testSubStruct
	TestPackerOnly    testPackerOnly
	TestUnpackerOnly  testUnpackerOnly
	TestInterface     testInterface
	TestFixedPtr      *uint8
	TestVariablePtr   *string
	TestSlice         []byte
	TestVariableSlice []string
	TestCompression   []byte `ikea:"compress:9"`
	Padding           uint16 // Maps randomise iteration order, we can't verify this string, so we split using this
	TestMap           map[string]string
}

type testSubStruct struct {
	A byte
	B []testSubStruct
}

type testInterface struct {
	A int64
}

func (t *testInterface) Unpack(r io.Reader) error {
	var temp int64
	if err := ikea.Unpack(r, &temp); err != nil {
		return err
	}

	t.A = temp - 10

	return nil
}

func (t *testInterface) Pack(w io.Writer) error {
	return ikea.Pack(w, t.A+10)
}

/* Test data */
var source = &testStruct{
	TestBool:          true,
	TestByte:          0x11,
	TestUint8:         0x88,
	TestUint16:        0x1616,
	TestUint32:        0x32323232,
	TestUint64:        0x6464646464646464,
	TestInt8:          0x12,
	TestInt16:         0x1234,
	TestInt32:         0x12345678,
	TestInt64:         0x1234567812345678,
	TestFloat32:       0.12345678,
	TestFloat64:       0.12345678901234567890,
	TestString:        "amazing serialization lib",
	TestSubStruct:     testSubStruct{A: 0x42, B: []testSubStruct{{A: 0x24}, {A: 0x42}}},
	TestPackerOnly:    testPackerOnly{A: 0x24},
	TestUnpackerOnly:  testUnpackerOnly{A: 0x42},
	TestInterface:     testInterface{A: 0x24},
	TestFixedPtr:      makeIntPtr(),
	TestVariablePtr:   makeStringPtr(),
	TestSlice:         make([]byte, 100),
	TestVariableSlice: []string{"a", "bc", "def"},
	TestCompression:   make([]byte, 10000),
	Padding:           0x4242,
	TestMap: map[string]string{
		"keynr1":     "valuenr1",
		"anotherkey": "anothervalue",
	},
}
var testData, _ = hex.DecodeString("01118816163232323264646464646464641212341234567812345678123456783dfcd6e93fbf9add3746f65f00000019616d617a696e672073657269616c697a6174696f6e206c69624200000002240000000042000000002442000000000000002e660000000e49276d206120706f696e74657221000000640001081b407dd85802dbeb38c69dc23c1044dee55f51c1b63646ec3016a4e1d380ed2223f6a32f9ffa478aca0e5ab526b15e323367d48173b03f25680f1f9e9304f56f7611451992b78e1d697953fc7cd7153a4d5455565d7095d22ead5731418d1cf2180000000300000001610000000262630000000364656600000024ecc0811000000803c021649147fe5079ecfe939d03000000000000000080021f0000ffff424200000002000000066b65796e72310000000876616c75656e72310000000a616e6f746865726b65790000000c616e6f7468657276616c7565")

func init() {
	for i := range source.TestSlice {
		source.TestSlice[i] = byte((i * i * i) % 0xFF)
	}
	for i := range source.TestCompression {
		source.TestCompression[i] = 0x42
	}
}

func makeIntPtr() *uint8 {
	i := uint8(0x66)
	return &i
}

func makeStringPtr() *string {
	s := "I'm a pointer!"
	return &s
}
