package tests

import (
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/ikkerens/ikeapack"
)

func TestReadPointer(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	var i int
	_ = ikea.Unpack(nil, i)
}

func TestUseInt(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	var i int
	_ = ikea.Unpack(nil, &i) // int is not supported
}

func TestUseUint(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	var ui uint
	_ = ikea.Unpack(nil, &ui) // uint is not supported
}

func TestUnsupportedType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	var ui complex64
	_ = ikea.Unpack(nil, &ui) // uint is not supported
}

func TestVariableLengthOverflow(t *testing.T) {
	overflow := new(bytes.Buffer)
	_ = ikea.Pack(overflow, uint32(math.MaxInt32+1))
	var t1 struct {
		Data struct{} `ikea:"compress:9"`
	}
	if err := ikea.Unpack(overflow, &t1); err == nil {
		t.FailNow()
	}

	overflow.Reset()
	_ = ikea.Pack(overflow, uint32(math.MaxInt32+1))
	var t2 map[string]struct{}
	if err := ikea.Unpack(overflow, &t2); err == nil {
		t.FailNow()
	}

	overflow.Reset()
	_ = ikea.Pack(overflow, uint32(math.MaxInt32+1))
	var t3 []struct{}
	if err := ikea.Unpack(overflow, &t3); err == nil {
		t.FailNow()
	}

	overflow.Reset()
	_ = ikea.Pack(overflow, uint32(math.MaxInt32+1))
	var t4 string
	if err := ikea.Unpack(overflow, &t4); err == nil {
		t.FailNow()
	}
}

func TestCompressionInitError(t *testing.T) {
	s1 := struct {
		Data []byte `ikea:"compress:10"`
	}{make([]byte, 10)}
	if err := ikea.Pack(new(bytes.Buffer), &s1); err == nil {
		t.Fail()
	}

	s2 := struct {
		Data []byte `ikea:"compress:a"`
	}{make([]byte, 10)}
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()
	_ = ikea.Pack(new(bytes.Buffer), &s2)
}

func TestInvalidUTF8(t *testing.T) {
	var invalid string
	b := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x01, 0xF1})
	if ikea.Unpack(b, &invalid) == nil {
		t.FailNow()
	}
}

func TestPackFixedNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	s := struct {
		A *uint32
	}{}
	_ = ikea.Pack(new(bytes.Buffer), &s)
}

func TestPackVariableNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	s := struct {
		A *string
	}{}
	_ = ikea.Pack(new(bytes.Buffer), &s)
}

func TestLenFixedNil(t *testing.T) {
	s := struct {
		A *uint32
	}{}
	ikea.Len(&s) // Unlike all other nil values, this should succeed
}

func TestLenVariableNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	s := struct {
		A *string
	}{}
	ikea.Len(&s)
}

func TestReadErrors(t *testing.T) {
	e := new(errorStream)

	// Reading error tests
	tst := new(testStruct)
	err := errors.New("start of errors")
	for err != nil {
		err = ikea.Unpack(e, tst)
		if err == nil && e.pass != len(testData) {
			t.FailNow()
		}
		e.Reset()
	}
}

func TestWriteErrors(t *testing.T) {
	e := new(errorStream)

	// Reading error tests
	err := errors.New("start of errors")
	for err != nil {
		err = ikea.Pack(e, source)
		if err == nil && e.pass != len(testData) {
			t.FailNow()
		}
		e.Reset()
	}
}

type errorStream struct {
	pointer int
	pass    int
}

func (s *errorStream) Read(p []byte) (n int, err error) {
	if s.pointer+len(p) > s.pass {
		s.pointer += len(p)
		return 0, errors.New("test error")
	}

	copy(p, testData[s.pointer:s.pointer+len(p)])
	s.pointer += len(p)
	return len(p), nil
}

func (s *errorStream) Write(p []byte) (n int, err error) {
	if s.pointer+len(p) > s.pass {
		s.pointer += len(p)
		return 0, errors.New("test error")
	}

	s.pointer += len(p)
	return len(p), nil
}

func (s *errorStream) Reset() {
	s.pass = s.pointer
	s.pointer = 0
}
