package serialize

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"unicode/utf8"
)

var stringTypeHandler = new(stringReadWriter)

type stringReadWriter struct {
	variable
}

func (s *stringReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}

	ul := binary.BigEndian.Uint32(b)
	if ul > math.MaxInt32 {
		return fmt.Errorf("transmitted string size too large (%d>%d)", ul, math.MaxInt32)
	}
	l := int(ul)

	str := make([]byte, l)
	if _, err := io.ReadFull(r, str); err != nil {
		return err
	}

	if !utf8.Valid(str) {
		return errors.New("invalid utf8 string")
	}

	v.SetString(string(str))
	return nil
}

func (s *stringReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v.Len()))

	if _, err := w.Write(b); err != nil {
		return err
	}
	if _, err := w.Write([]byte(v.String())); err != nil {
		return err
	}

	return nil
}
