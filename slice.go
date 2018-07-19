package ikea

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"sync"
)

var (
	sliceIndex     = make(map[string]*sliceReadWriter)
	sliceIndexLock sync.RWMutex
)

func getSliceHandlerFromType(t reflect.Type) readWriter {
	sliceIndexLock.RLock()
	infoV, found := sliceIndex[t.String()]
	sliceIndexLock.RUnlock()
	if found {
		return infoV
	}

	info := &sliceReadWriter{typ: t, handler: getTypeHandler(t.Elem())}
	sliceIndexLock.Lock()
	sliceIndex[t.String()] = info
	sliceIndexLock.Unlock()

	return info
}

type sliceReadWriter struct {
	variable
	typ     reflect.Type
	handler readWriter
}

func (s *sliceReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}

	ul := binary.BigEndian.Uint32(b)
	if ul > math.MaxInt32 {
		return fmt.Errorf("transmitted slice size too large (%d>%d)", ul, math.MaxInt32)
	}
	l := int(ul)

	slice := reflect.MakeSlice(s.typ, l, l)

	if s.handler.isFixed() {
		hr := s.handler.(fixedReadWriter)
		sb := make([]byte, l*hr.length())
		if _, err := io.ReadFull(r, sb); err != nil {
			return err
		}

		for i := 0; i < l; i++ {
			idx := i * hr.length()
			hr.readFixed(sb[idx:idx+hr.length()], slice.Index(i))
		}
	} else {
		hr := s.handler.(variableReadWriter)
		for i := 0; i < l; i++ {
			if err := hr.readVariable(r, slice.Index(i)); err != nil {
				return err
			}
		}
	}

	v.Set(slice)
	return nil
}

func (s *sliceReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v.Len()))

	if _, err := w.Write(b); err != nil {
		return err
	}

	if s.handler.isFixed() {
		hw := s.handler.(fixedReadWriter)
		sb := make([]byte, v.Len()*hw.length())

		for i := 0; i < v.Len(); i++ {
			idx := i * hw.length()
			hw.writeFixed(sb[idx:idx+hw.length()], v.Index(i))
		}

		if _, err := w.Write(sb); err != nil {
			return err
		}
	} else {
		hw := s.handler.(variableReadWriter)
		for i := 0; i < v.Len(); i++ {
			if err := hw.writeVariable(w, v.Index(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *sliceReadWriter) vLength(v reflect.Value) (int, error) {
	if s.handler.isFixed() {
		return 4 + (v.Len() * s.handler.(fixedReadWriter).length()), nil
	}

	// variable
	size := 4
	h := s.handler.(variableReadWriter)
	for i := 0; i < v.Len(); i++ {
		l, err := h.vLength(v.Index(i))
		if err != nil {
			return 0, err
		}
		size += l
	}
	return size, nil
}
