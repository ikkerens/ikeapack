package ikea

import (
	"compress/flate"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"
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
		hasUnpacker = interfaceTest.Implements(unpackerInterface)
		hasPacker   = interfaceTest.Implements(packerInterface)
	)
	if hasUnpacker && hasPacker {
		ret.r = &customReadWriter{fallback: nil}
	} else if hasUnpacker || hasPacker {
		ret.r = &customReadWriter{fallback: scanStruct(t)}
	} else {
		ret.r = scanStruct(t)
	}

	// Replace the original with the direct version (major performance boost)
	structIndexLock.Lock()
	structIndex[t.String()] = ret.r
	structIndexLock.Unlock()

	return ret.r
}

func scanStruct(t reflect.Type) readWriter {
	handlers := make([]readWriter, 0)

	length := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		r := rune(field.Name[0])
		if unicode.ToLower(r) == r {
			handlers = append(handlers, nil)
			continue // Ignore, unexported
		}

		tag, ok := field.Tag.Lookup("ikea")
		if !ok {
			tag = ""
		}

		if tag == "-" {
			handlers = append(handlers, nil)
			continue // Ignore, ignored
		}

		h := getTypeHandler(field.Type)
		if h.isFixed() && length != -1 {
			length += h.(fixedReadWriter).length()
		} else {
			length = -1
		}

		if strings.HasPrefix(tag, "compress") {
			length = -1

			var level = flate.BestCompression
			if strings.HasPrefix(tag, "compress:") {
				var err error
				level, err = strconv.Atoi(strings.TrimPrefix(tag, "compress:"))
				if err != nil {
					panic(err)
				}
			}

			h = &compressionReadWriter{handler: h, level: level}
		}

		handlers = append(handlers, h)
	}

	if length != -1 {
		return &fixedStructReadWriter{size: length, handlers: handlers}
	}

	return &variableStructReadWriter{handlers: handlers}
}

var _ variableReadWriter = (*structWrapper)(nil)

type structWrapper struct {
	sync.Mutex
	variable
	r readWriter
}

func (s *structWrapper) vLength(v reflect.Value) int {
	return s.r.(variableReadWriter).vLength(v)
}

func (s *structWrapper) readVariable(r io.Reader, v reflect.Value) error {
	return s.r.(variableReadWriter).readVariable(r, v)
}

func (s *structWrapper) writeVariable(w io.Writer, v reflect.Value) error {
	return s.r.(variableReadWriter).writeVariable(w, v)
}

var _ fixedReadWriter = (*fixedStructReadWriter)(nil)

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
		if handler == nil {
			continue
		}
		r := handler.(fixedReadWriter)
		r.readFixed(data[read:read+r.length()], v.Field(i))
		read += r.length()
	}
}

func (s *fixedStructReadWriter) writeFixed(data []byte, v reflect.Value) {
	written := 0
	for i, handler := range s.handlers {
		if handler == nil {
			continue
		}
		w := handler.(fixedReadWriter)
		w.writeFixed(data[written:written+w.length()], v.Field(i))
		written += w.length()
	}
}

var _ variableReadWriter = (*variableStructReadWriter)(nil)

type variableStructReadWriter struct {
	variable

	handlers []readWriter
}

func (h *variableStructReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	for i, handler := range h.handlers {
		if handler == nil {
			continue
		}
		if err := handleVariableReader(r, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func (h *variableStructReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	for i, handler := range h.handlers {
		if handler == nil {
			continue
		}
		if err := handleVariableWriter(w, handler, v.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func (h *variableStructReadWriter) vLength(v reflect.Value) int {
	size := 0

	for i, handler := range h.handlers {
		if handler == nil {
			continue
		}
		size += handleVariableLength(handler, v.Field(i))
	}

	return size
}
