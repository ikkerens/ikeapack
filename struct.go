package serialize

import (
	"io"
	"reflect"
	"sync"
)

var structIndex = sync.Map{}

func getStructHandlerFromType(t reflect.Type) *typeHandler {
	if t.Kind() != reflect.Struct {
		panic("passed value is not a struct")
	}

	infoV, found := structIndex.Load(t.String())
	if found {
		return infoV.(*typeHandler)
	}
	info := new(typeHandler)

	handlers := make([]*typeHandler, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typ := field.Type
		h := getTypeHandler(typ)

		if h.length == -1 {
			info.length = -1
		} else if info.length != -1 {
			info.length += h.length
		}

		if value, ok := field.Tag.Lookup("compressed"); ok && value == "true" {
			h = makeCompressionHandler(h)
		}

		handlers = append(handlers, h)
	}

	if info.length != -1 {
		info.handler = &fixedStructReadWriter{handlers: handlers}
	} else {
		info.handler = &variableStructReadWriter{handlers: handlers}
	}

	testT := reflect.New(t).Type()
	if testT.Implements(deserializerInterface) || testT.Implements(serializerInterface) {
		info = &typeHandler{
			length:  -1,
			handler: &customReadWriter{info},
		}
	}

	structIndex.Store(t.String(), info)

	return info
}

type fixedStructReadWriter struct {
	handlers []*typeHandler
}

func (s *fixedStructReadWriter) read(data []byte, v reflect.Value) {
	read := 0
	for i, handler := range s.handlers {
		handler.handler.(fixedReadWriter).read(data[read:read+handler.length], v.Field(i))
		read += handler.length
	}
}

func (s *fixedStructReadWriter) write(data []byte, v reflect.Value) {
	written := 0
	for i, handler := range s.handlers {
		handler.handler.(fixedReadWriter).write(data[written:written+handler.length], v.Field(i))
		written += handler.length
	}
}

type variableStructReadWriter struct {
	handlers []*typeHandler
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
