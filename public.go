package serialize

import (
	"io"
	"reflect"
)

// Read will read exactly enough bytes from the specified Reader in order to fill the value passed to data.
// if data is not a pointer Read will panic
func Read(r io.Reader, data interface{}) error {
	pv := reflect.ValueOf(data)
	if pv.Kind() != reflect.Ptr {
		panic("passed data argument is not a pointer")
	}

	v := pv.Elem()
	h := getTypeHandler(v.Type())

	return handleVariableReader(r, h, v)
}

// Write will write the value passed in data to the specified Writer
func Write(w io.Writer, data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	h := getTypeHandler(v.Type())

	return handleVariableWriter(w, h, v)
}
