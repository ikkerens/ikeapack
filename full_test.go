package ikea

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

/* Tests */

func TestOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Pack(buf, source); err != nil {
		t.Error(err)
		return
	}

	result := buf.Bytes()
	if len(result) != len(testData) {
		t.Errorf("Failing TestWrite, result \"%s\" length (%d) does not match test data length (%d)\n", hex.EncodeToString(result), len(result), len(testData))
		return
	}

	// In struct creation we've added 0x4242 padding to we can split out the map.
	// We need to do this because golangs map iteration is always random.
	// So we'll deal with that one later in this test.
	originalParts := strings.Split(hex.EncodeToString(testData), "4242")
	resultParts := strings.Split(hex.EncodeToString(result), "4242")

	if originalParts[0] != resultParts[0] {
		t.Errorf("Failing TestWrite, hex output \"%s\" does not match test data slice\n", hex.EncodeToString(result))
		return
	}

	// Instead, we treat the map a bit differently, we put it back into a buffer
	buf.Reset()
	data, err := hex.DecodeString(resultParts[1])
	if err != nil {
		t.Errorf("Failing TestWrite, hex output \"%s\" is not a valid hex string: %s\n", resultParts[1], err.Error())
		return
	}
	buf.Write(data)

	// Unpack it
	var test map[string]string
	if err := Unpack(buf, &test); err != nil {
		t.Errorf("Failing TestWrite, could not unpack map: %s\n", err.Error())
		return
	}

	// And then compare it using DeepEqual
	if !reflect.DeepEqual(source.TestMap, test) {
		t.Errorf("Failing TestWrite, resulting map is not equal\n")
	}
}

func TestCompleteRead(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(testData)

	tst := new(testStruct)
	if err := Unpack(buf, tst); err != nil {
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
	compare(t, "TestFloat32", tst.TestFloat32, source.TestFloat32)
	compare(t, "TestFloat64", tst.TestFloat64, source.TestFloat64)
	compare(t, "TestString", tst.TestString, source.TestString)
	compare(t, "TestSubStruct", tst.TestSubStruct.A, source.TestSubStruct.A)
	compare(t, "TestInterface", tst.TestInterface.A, source.TestInterface.A)
	compare(t, "TestFixedPtr", *tst.TestFixedPtr, *source.TestFixedPtr)
	compare(t, "TestVariablePtr", *tst.TestVariablePtr, *source.TestVariablePtr)
	for i := range source.TestSlice {
		compare(t, "TestSlice["+strconv.Itoa(i)+"]", tst.TestSlice[i], source.TestSlice[i])
	}
	for i := range source.TestCompression {
		compare(t, "TestCompression["+strconv.Itoa(i)+"]", tst.TestCompression[i], source.TestCompression[i])
	}
	for k := range source.TestMap {
		compare(t, "TestMap["+k+"]", tst.TestMap[k], source.TestMap[k])
	}
}

func TestLen(t *testing.T) {
	if l := Len(source); l != len(testData) {
		t.Errorf("Failing TestLen, Len reported an incorrect value %d, should be %d", l, len(testData))
	}
}

func compare(t *testing.T, field string, value1, value2 interface{}) {
	if value1 != value2 {
		t.Errorf("Failing TestCompleteRead, decoded data field '%s' with value '%v' does not match '%v'", field, value1, value2)
	}
}

type testPackerOnly struct {
	A       uint8
	Ignored uint8 `ikea:"-"`
}

func (p *testPackerOnly) Pack(w io.Writer) error {
	return Pack(w, &p.A)
}

type testUnpackerOnly struct {
	A       uint8
	Ignored uint8 `ikea:"-"`
}

func (p *testUnpackerOnly) Unpack(r io.Reader) error {
	return Unpack(r, &p.A)
}
