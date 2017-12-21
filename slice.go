package serialize

import (
	"encoding/binary"
	"io"
	"reflect"
	"sync"
)

var sliceIndex = sync.Map{}

func getSliceHandlerFromType(t reflect.Type) readWriter {
	infoV, found := sliceIndex.Load(t.String())
	if found {
		return infoV.(readWriter)
	}

	info := &sliceReadWriter{t, getTypeHandler(t.Elem())}
	sliceIndex.Store(t.String(), info)

	return info
}

type sliceReadWriter struct {
	typ     reflect.Type
	handler readWriter
}

func (s *sliceReadWriter) read(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	l := int(binary.BigEndian.Uint32(b))
	slice := reflect.MakeSlice(s.typ, l, l)

	switch hr := s.handler.(type) {
	case fixedReadWriter:
		sb := make([]byte, l*hr.length())
		if _, err := io.ReadFull(r, sb); err != nil {
			return err
		}

		for i := 0; i < l; i++ {
			idx := i * hr.length()
			hr.read(sb[idx:idx+hr.length()], slice.Index(i))
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

	switch hw := s.handler.(type) {
	case fixedReadWriter:
		sb := make([]byte, v.Len()*hw.length())

		for i := 0; i < v.Len(); i++ {
			idx := i * hw.length()
			hw.write(sb[idx:idx+hw.length()], v.Index(i))
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
