package serialize

import (
	"io"
	"reflect"
)

func Read(r io.Reader, data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	h := getTypeHandler(v.Type())

	return handleVariableReader(r, h, v)
}

func Write(w io.Writer, data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	h := getTypeHandler(v.Type())

	return handleVariableWriter(w, h, v)
}
