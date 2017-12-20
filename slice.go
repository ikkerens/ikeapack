package serialize

import (
	"encoding/binary"
	"io"
	"reflect"
	"sync"
)

var sliceIndex = sync.Map{}

func getSliceHandlerFromType(t reflect.Type) *typeHandler {
	if t.Kind() != reflect.Slice {
		panic("passed value is not a slice")
	}

	infoV, found := sliceIndex.Load(t.String())
	if found {
		return infoV.(*typeHandler)
	}

	info := &typeHandler{
		length:  -1,
		handler: &sliceReadWriter{t, getTypeHandler(t.Elem())},
	}

	sliceIndex.Store(t.String(), info)

	return info
}

type sliceReadWriter struct {
	typ     reflect.Type
	handler *typeHandler
}

func (s *sliceReadWriter) read(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	l := int(binary.BigEndian.Uint32(b))
	slice := reflect.MakeSlice(s.typ, l, l)

	switch hr := s.handler.handler.(type) {
	case fixedReadWriter:
		sb := make([]byte, l*s.handler.length)
		if _, err := io.ReadFull(r, sb); err != nil {
			return err
		}

		for i := 0; i < l; i++ {
			idx := i * s.handler.length
			hr.read(sb[idx:idx+s.handler.length], slice.Index(i))
		}
	case variableReadWriter:
		for i := 0; i < l; i++ {
			if err := hr.read(r, slice.Index(i)); err != nil {
				return err
			}
		}
	}

	v.Set(slice)
	return nil
}

func (s *sliceReadWriter) write(w io.Writer, v reflect.Value) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v.Len()))

	if _, err := w.Write(b); err != nil {
		return err
	}

	switch hw := s.handler.handler.(type) {
	case fixedReadWriter:
		sb := make([]byte, v.Len()*s.handler.length)

		for i := 0; i < v.Len(); i++ {
			idx := i * s.handler.length
			hw.write(sb[idx:idx+s.handler.length], v.Index(i))
		}

		if _, err := w.Write(sb); err != nil {
			return err
		}
	case variableReadWriter:
		for i := 0; i < v.Len(); i++ {
			if err := hw.write(w, v.Index(i)); err != nil {
				return err
			}
		}
	}

	return nil
}
