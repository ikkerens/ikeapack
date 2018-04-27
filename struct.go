package serialize

import (
	"io"
	"reflect"
	"sync"
)

var (
	structIndex     = make(map[string]readWriter)
	structIndexLock sync.RWMutex
)

func getStructHandlerFromType(t reflect.Type) readWriter {
	structIndexLock.RLock()
	infoV, found := structIndex[t.String()]
	structIndexLock.RUnlock()
	if found {
		return infoV
	}

	ret := new(structWrapper)
	ret.Lock()
	defer ret.Unlock()

	// For now, insert the wrapper, so recursive struct calls won't cause an infinite stack
	structIndexLock.Lock()
	structIndex[t.String()] = ret
	structIndexLock.Unlock()

	interfaceTest := reflect.New(t).Type()
	var (
		hasDeserializer = interfaceTest.Implements(deserializerInterface)
		hasSerializer   = interfaceTest.Implements(serializerInterface)
	)
	if hasDeserializer && hasSerializer {
		ret.readWriter = &customReadWriter{fallback: nil}
	} else if hasDeserializer || hasSerializer {
		ret.readWriter = &customReadWriter{fallback: scanStruct(t)}
	} else {
		ret.readWriter = scanStruct(t)
	}

	// Replace the original with the direct version (major performance boost)
	structIndexLock.Lock()
	structIndex[t.String()] = ret.readWriter
	structIndexLock.Unlock()

	return ret.readWriter
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
	sync.Mutex
	readWriter
}

func (s *structWrapper) isFixed() bool {
	s.Lock()
	defer s.Unlock()

	return s.readWriter.isFixed()
}

func (s *structWrapper) readVariable(r io.Reader, v reflect.Value) error {
	return s.readWriter.(variableReadWriter).readVariable(r, v)
}

func (s *structWrapper) writeVariable(w io.Writer, v reflect.Value) error {
	return s.readWriter.(variableReadWriter).writeVariable(w, v)
}

func (s *structWrapper) length() int {
	return s.readWriter.(fixedReadWriter).length()
}

func (s *structWrapper) readFixed(b []byte, v reflect.Value) {
	s.readWriter.(fixedReadWriter).readFixed(b, v)
}

func (s *structWrapper) writeFixed(b []byte, v reflect.Value) {
	s.readWriter.(fixedReadWriter).writeFixed(b, v)
}

type fixedStructReadWriter struct {
	fixed

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
	variable

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
