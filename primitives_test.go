package ikea

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

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

func TestFloat32(t *testing.T) {
	i := rand.Float32()
	typeTest(t, "TestFloat32", &i, i)
}

func TestFloat64(t *testing.T) {
	i := rand.Float64()
	typeTest(t, "TestFloat32", &i, i)
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

	if err := Pack(&b, value); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failing %s, could not write value: %s\n", typ, err.Error())
		t.FailNow()
	}

	target := reflect.New(reflect.TypeOf(value).Elem())
	if err := Unpack(&b, target.Interface()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failing %s, could not read value: %s\n", typ, err.Error())
		t.FailNow()
	}

	dereference := target.Elem().Interface()
	if dereference != compare {
		_, _ = fmt.Fprintf(os.Stderr, "Failing %s, %T value %+v does not match original %T %+v\n", typ, dereference, dereference, compare, compare)
		t.FailNow()
	}
}
