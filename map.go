package serialize

import (
	"encoding/binary"
	"io"
	"reflect"
	"sync"
)

var mapIndex = sync.Map{}

func getMapHandlerFromType(t reflect.Type) readWriter {
	infoV, found := mapIndex.Load(t.String())
	if found {
		return infoV.(readWriter)
	}

	t.Key()
	info := &mapReadWriter{t, t.Key(), t.Elem(), getTypeHandler(t.Key()), getTypeHandler(t.Elem())}
	mapIndex.Store(t.String(), info)

	return info
}

type mapReadWriter struct {
	mapType                  reflect.Type
	keyType, valueType       reflect.Type
	keyHandler, valueHandler readWriter
}

func (s *mapReadWriter) read(r io.Reader, v reflect.Value) error {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	l := int(binary.BigEndian.Uint32(b))
	mp := reflect.MakeMapWithSize(s.mapType, l)

	for i := 0; i < l; i++ {
		key := reflect.Indirect(reflect.New(s.keyType))
		val := reflect.Indirect(reflect.New(s.valueType))

		handleVariableReader(r, s.keyHandler, key)
		handleVariableReader(r, s.valueHandler, val)

		mp.SetMapIndex(key, val)
	}

	v.Set(mp)
	return nil
}

func (s *mapReadWriter) write(w io.Writer, v reflect.Value) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v.Len()))

	if _, err := w.Write(b); err != nil {
		return err
	}

	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)

		handleVariableWriter(w, s.keyHandler, key)
		handleVariableWriter(w, s.valueHandler, val)
	}

	return nil
}
