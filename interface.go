package ikea

import (
	"bytes"
	"io"
	"reflect"
)

var (
	unpackerInterface = reflect.TypeOf((*Unpacker)(nil)).Elem()
	packerInterface   = reflect.TypeOf((*Packer)(nil)).Elem()
)

// Unpacker allows you to implement a custom unpacking strategy for a type
type Unpacker interface {
	Unpack(r io.Reader) error
}

// Packer allows you to implement a custom packing strategy for a type
type Packer interface {
	Pack(w io.Writer) error
}

type customReadWriter struct {
	variable
	fallback readWriter
}

func (c *customReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	var err error
	if d, ok := v.Addr().Interface().(Unpacker); ok {
		err = d.Unpack(r)
	} else {
		err = handleVariableReader(r, c.fallback, v)
	}
	return err
}

func (c *customReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	var err error
	if s, ok := v.Addr().Interface().(Packer); ok {
		err = s.Pack(w)
	} else {
		err = handleVariableWriter(w, c.fallback, v)
	}
	return err
}

func (c *customReadWriter) vLength(v reflect.Value) (int, error) {
	var b bytes.Buffer
	if err := c.writeVariable(&b, v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}
