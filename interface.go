package serialize

import (
	"io"
	"reflect"
)

var (
	deserializerInterface = reflect.TypeOf((*Deserializer)(nil)).Elem()
	serializerInterface   = reflect.TypeOf((*Serializer)(nil)).Elem()
)

// Deserializer allows you to implement a custom deserialization strategy for a type
type Deserializer interface {
	Deserialize(r io.Reader) error
}

// Serializer allows you to implement a custom serialization strategy for a type
type Serializer interface {
	Serialize(w io.Writer) error
}

type customReadWriter struct {
	fallback *typeHandler
}

func (c *customReadWriter) read(r io.Reader, v reflect.Value) error {
	var err error
	if d, ok := v.Addr().Interface().(Deserializer); ok {
		err = d.Deserialize(r)
	} else {
		err = handleVariableReader(r, c.fallback, v)
	}
	return err
}

func (c *customReadWriter) write(w io.Writer, v reflect.Value) error {
	var err error
	if s, ok := v.Addr().Interface().(Serializer); ok {
		err = s.Serialize(w)
	} else {
		err = handleVariableWriter(w, c.fallback, v)
	}
	return err
}
