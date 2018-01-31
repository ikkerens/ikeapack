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

	ret := new(structWrapper)
	structIndex.Store(t.String(), ret)

	interfaceTest := reflect.New(t).Type()
	var (
		hasDeserializer = interfaceTest.Implements(deserializerInterface)
		hasSerializer   = interfaceTest.Implements(serializerInterface)
	)
	if hasDeserializer && hasSerializer {
		ret.wrapped = &customReadWriter{fallback: nil}
	} else if hasDeserializer || hasSerializer {
		ret.wrapped = &customReadWriter{fallback: scanStruct(t)}
	} else {
		ret.wrapped = scanStruct(t)
	}

	return ret.wrapped
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
			h = &compressionReadWriter{handler: h}
		}

		handlers = append(handlers, h)
	}

	if length != -1 {
		return &fixedStructReadWriter{size: length, handlers: handlers}
	}

	return &variableStructReadWriter{handlers: handlers}
}

type structWrapper struct {
	wrapped readWriter
}

func (s *structWrapper) isFixed() bool {
	return s.wrapped.isFixed()
}

func (s *structWrapper) readVariable(r io.Reader, v reflect.Value) error {
	return s.wrapped.(variableReadWriter).readVariable(r, v)
}

func (s *structWrapper) writeVariable(w io.Writer, v reflect.Value) error {
	return s.wrapped.(variableReadWriter).writeVariable(w, v)
}

func (s *structWrapper) length() int {
	return s.wrapped.(fixedReadWriter).length()
}

func (s *structWrapper) readFixed(b []byte, v reflect.Value) {
	s.wrapped.(fixedReadWriter).readFixed(b, v)
}

func (s *structWrapper) writeFixed(b []byte, v reflect.Value) {
	s.wrapped.(fixedReadWriter).writeFixed(b, v)
}

type fixedStructReadWriter struct {
	fixedImpl

	size     int
	handlers []readWriter
}

func (s *fixedStructReadWriter) length() int {
	return s.size
}

func (s *fixedStructReadWriter) readFixed(data []byte, v reflect.Value) {
	read := 0
	for i, handler := range s.handlers {
		r := handler.(fixedReadWriter)
		r.readFixed(data[read:read+r.length()], v.Field(i))
		read += r.length()
	}
}

func (s *fixedStructReadWriter) writeFixed(data []byte, v reflect.Value) {
	written := 0
	for i, handler := range s.handlers {
		w := handler.(fixedReadWriter)
		w.writeFixed(data[written:written+w.length()], v.Field(i))
		written += w.length()
	}
}

type variableStructReadWriter struct {
	variableImpl

	handlers []readWriter
}

func (h *variableStructReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	for i, handler := range h.handlers {
		if err := handleVariableReader(r, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func (h *variableStructReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	for i, handler := range h.handlers {
		if err := handleVariableWriter(w, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}
