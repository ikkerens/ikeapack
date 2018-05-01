package serialize

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

/* Tests */

func TestBool(t *testing.T) {
	i := rand.Int()%2 == 1
	typeTest(t, "TestBool", &i, i)
}

func TestByte(t *testing.T) {
	i := byte(rand.Int() & 0xFF)
	typeTest(t, "TestByte", &i, i)
}

func TestUint8(t *testing.T) {
	i := uint8(rand.Int() & 0xFF)
	typeTest(t, "TestUint8", &i, i)
}

func TestUint16(t *testing.T) {
	i := uint16(rand.Int() & 0xFFFF)
	typeTest(t, "TestUint16", &i, i)
}

func TestUint32(t *testing.T) {
	i := rand.Uint32()
	typeTest(t, "TestUint32", &i, i)
}

func TestUint64(t *testing.T) {
	i := rand.Uint64()
	typeTest(t, "TestUint64", &i, i)
}

func TestInt8(t *testing.T) {
	i := int8(rand.Int() & 0xFF)
	typeTest(t, "TestInt8", &i, i)
}

func TestInt16(t *testing.T) {
	i := int16(rand.Int() & 0xFFFF)
	typeTest(t, "TestInt16", &i, i)
}

func TestInt32(t *testing.T) {
	i := rand.Int31()
	typeTest(t, "TestInt32", &i, i)
}

func TestInt64(t *testing.T) {
	i := rand.Int63()
	typeTest(t, "TestInt64", &i, i)
}

func TestString(t *testing.T) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, rand.Intn(30))
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	s := string(b)

	typeTest(t, "TestString", &s, s)
}

func typeTest(t *testing.T, typ string, value, compare interface{}) {
	var b bytes.Buffer

	if err := Write(&b, value); err != nil {
		fmt.Fprintf(os.Stderr, "Failing %s, could not write value: %s\n", typ, err.Error())
		t.FailNow()
	}

	if err := Read(&b, value); err != nil {
		fmt.Fprintf(os.Stderr, "Failing %s, could not read value: %s\n", typ, err.Error())
		t.FailNow()
	}

	dereference := reflect.Indirect(reflect.ValueOf(value)).Interface()
	if dereference != compare {
		fmt.Fprintf(os.Stderr, "Failing %s, value %+v does not match original %+v\n", typ, dereference, compare)
		t.FailNow()
	}
}

func TestOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Write(buf, source); err != nil {
		t.Error(err)
		return
	}

	result := buf.Bytes()
	if len(result) != len(testData) {
		fmt.Printf("Failing TestWrite, result \"%s\" length (%d) does not match test data length (%d)\n", hex.EncodeToString(result), len(result), len(testData))
		t.FailNow()
	}

	// In struct creation we've added 0x4242 padding to we can split out the map.
	// We need to do this because golangs map iteration is always random.
	// However, for the sake of this test, I'm working around it.
	originalParts := strings.Split(hex.EncodeToString(testData), "4242")
	resultParts := strings.Split(hex.EncodeToString(result), "4242")

	if originalParts[0] != resultParts[0] {
		fmt.Printf("Failing TestWrite, hex output \"%s\" does not match test data slice\n", hex.EncodeToString(result))
		t.FailNow()
	}
}

func TestCompleteRead(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(testData)

	tst := new(testStruct)
	if err := Read(buf, tst); err != nil {
		t.Error(err)
		return
	}

	compare(t, "TestBool", tst.TestBool, source.TestBool)
	compare(t, "TestByte", tst.TestByte, source.TestByte)
	compare(t, "TestUint8", tst.TestUint8, source.TestUint8)
	compare(t, "TestUint16", tst.TestUint16, source.TestUint16)
	compare(t, "TestUint32", tst.TestUint32, source.TestUint32)
	compare(t, "TestUint64", tst.TestUint64, source.TestUint64)
	compare(t, "TestInt8", tst.TestInt8, source.TestInt8)
	compare(t, "TestInt16", tst.TestInt16, source.TestInt16)
	compare(t, "TestInt32", tst.TestInt32, source.TestInt32)
	compare(t, "TestInt64", tst.TestInt64, source.TestInt64)
	compare(t, "TestString", tst.TestString, source.TestString)
	compare(t, "TestSubStruct", tst.TestSubStruct.A, source.TestSubStruct.A)
	compare(t, "TestInterface", tst.TestInterface.A, source.TestInterface.A)
	for i := range source.TestSlice {
		compare(t, "TestSlice["+strconv.FormatInt(int64(i), 10)+"]", tst.TestSlice[i], source.TestSlice[i])
	}
	for i := range source.TestCompression {
		compare(t, "TestCompression["+strconv.FormatInt(int64(i), 10)+"]", tst.TestCompression[i], source.TestCompression[i])
	}
	for k := range source.TestMap {
		compare(t, "TestMap["+k+"]", tst.TestMap[k], source.TestMap[k])
	}
}

func compare(t *testing.T, field string, value1, value2 interface{}) {
	if value1 != value2 {
		fmt.Printf("Failing TestCompleteRead, decoded data field '%s' with value '%v' does not match '%v'", field, value1, value2)
		t.FailNow()
	}
}

/* Test types */

type testStruct struct {
	TestBool        bool
	TestByte        byte
	TestUint8       uint8
	TestUint16      uint16
	TestUint32      uint32
	TestUint64      uint64
	TestInt8        int8
	TestInt16       int16
	TestInt32       int32
	TestInt64       int64
	TestString      string
	TestSubStruct   testSubStruct
	TestInterface   testInterface
	TestSlice       []byte
	TestCompression []byte `compressed:"true"`
	Padding         uint16 // Maps randomise iteration order, we can't verify this string, so we split using this
	TestMap         map[string]string
}

type testSubStruct struct {
	A byte
}

type testInterface struct {
	A int64
}

func (t *testInterface) Deserialize(r io.Reader) error {
	var temp int64
	if err := Read(r, &temp); err != nil {
		return err
	}

	t.A = temp - 10

	return nil
}

func (t *testInterface) Serialize(w io.Writer) error {
	return Write(w, t.A+10)
}

/* Test data */
var source = &testStruct{
	TestBool:        true,
	TestByte:        0x11,
	TestUint8:       0x88,
	TestUint16:      0x1616,
	TestUint32:      0x32323232,
	TestUint64:      0x6464646464646464,
	TestInt8:        0x12,
	TestInt16:       0x1234,
	TestInt32:       0x12345678,
	TestInt64:       0x1234567812345678,
	TestString:      "amazing serialization lib",
	TestSubStruct:   testSubStruct{A: 0x42},
	TestInterface:   testInterface{A: 0x24},
	TestSlice:       make([]byte, 100),
	TestCompression: make([]byte, 10000),
	Padding:         0x4242,
	TestMap: map[string]string{
		"keynr1":     "valuenr1",
		"anotherkey": "anothervalue",
	},
}
var testData, _ = hex.DecodeString("011188161632323232646464646464646412123412345678123456781234567800000019616d617a696e672073657269616c697a6174696f6e206c696242000000000000002e000000640001081b407dd85802dbeb38c69dc23c1044dee55f51c1b63646ec3016a4e1d380ed2223f6a32f9ffa478aca0e5ab526b15e323367d48173b03f25680f1f9e9304f56f7611451992b78e1d697953fc7cd7153a4d5455565d7095d22ead5731418d1cf21800000024ecc0811000000803c021649147fe5079ecfe939d03000000000000000080021f0000ffff424200000002000000066b65796e72310000000876616c75656e72310000000a616e6f746865726b65790000000c616e6f7468657276616c7565")

func init() {
	for i := range source.TestSlice {
		source.TestSlice[i] = byte((i * i * i) % 0xFF)
	}
	for i := range source.TestCompression {
		source.TestCompression[i] = 0x42
	}
}
