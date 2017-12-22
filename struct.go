package serialize

import (
	"io"
	"reflect"
	"sync"
)

var structIndex = sync.Map{}

func getStructHandlerFromType(t reflect.Type) readWriter {
	infoV, found := structIndex.Load(t.String())
	if found {
		return infoV.(readWriter)
	}

	var ret readWriter

	interfaceTest := reflect.New(t).Type()
	var (
		hasDeserializer = interfaceTest.Implements(deserializerInterface)
		hasSerializer   = interfaceTest.Implements(serializerInterface)
	)
	if hasDeserializer && hasSerializer {
		ret = &customReadWriter{nil}
	} else if hasDeserializer || hasSerializer {
		ret = &customReadWriter{scanStruct(t)}
	} else {
		ret = scanStruct(t)
	}

	structIndex.Store(t.String(), ret)
	return ret
}

func scanStruct(t reflect.Type) readWriter {
	handlers := make([]readWriter, 0)

	length := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typ := field.Type
		h := getTypeHandler(typ)

		if f, ok := h.(fixedReadWriter); ok && length != -1 {
			length += f.length()
		} else {
			length = -1
		}

		if value, ok := field.Tag.Lookup("compressed"); ok && value == "true" {
			length = -1
			h = &compressionReadWriter{h}
		}

		handlers = append(handlers, h)
	}

	if length != -1 {
		return &fixedStructReadWriter{length, handlers}
	}

	return &variableStructReadWriter{handlers}
}

type fixedStructReadWriter struct {
	size     int
	handlers []readWriter
}

func (s *fixedStructReadWriter) length() int {
	return s.size
}

func (s *fixedStructReadWriter) read(data []byte, v reflect.Value) {
	read := 0
	for i, handler := range s.handlers {
		r := handler.(fixedReadWriter)
		r.read(data[read:read+r.length()], v.Field(i))
		read += r.length()
	}
}

func (s *fixedStructReadWriter) write(data []byte, v reflect.Value) {
	written := 0
	for i, handler := range s.handlers {
		w := handler.(fixedReadWriter)
		w.write(data[written:written+w.length()], v.Field(i))
		written += w.length()
	}
}

type variableStructReadWriter struct {
	handlers []readWriter
}

func (h *variableStructReadWriter) read(r io.Reader, v reflect.Value) error {
	for i, handler := range h.handlers {
		if err := handleVariableReader(r, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func (h *variableStructReadWriter) write(w io.Writer, v reflect.Value) error {
	for i, handler := range h.handlers {
		if err := handleVariableWriter(w, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}
