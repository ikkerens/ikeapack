package serialize

import (
	"fmt"
	"io"
	"reflect"
)

type typeHandler struct {
	length  int
	handler interface{}
}

type fixedReadWriter interface {
	read([]byte, reflect.Value)

	write([]byte, reflect.Value)
}

type variableReadWriter interface {
	read(io.Reader, reflect.Value) error

	write(io.Writer, reflect.Value) error
}

func getTypeHandler(typ reflect.Type) *typeHandler {
	kind := typ.Kind()

	primitive, ok := primitiveIndex[kind]
	if ok {
		return primitive
	}

	switch kind {
	case reflect.String:
		return stringTypeHandler
	case reflect.Struct:
		return getStructHandlerFromType(typ)
	case reflect.Slice:
		return getSliceHandlerFromType(typ)
	default:
		panic(fmt.Sprintf("Cannot build type handler for type \"%s\" with kind nr. %d", typ.String(), typ.Kind()))
	}
}
