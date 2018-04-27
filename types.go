package serialize

import (
	"fmt"
	"io"
	"reflect"
)

type readWriter interface {
	isFixed() bool
}

type fixedReadWriter interface {
	readWriter

	length() int

	readFixed([]byte, reflect.Value)

	writeFixed([]byte, reflect.Value)
}

type variableReadWriter interface {
	readWriter

	readVariable(io.Reader, reflect.Value) error

	writeVariable(io.Writer, reflect.Value) error
}

func getTypeHandler(typ reflect.Type) readWriter {
	kind := typ.Kind()

	if primitive, ok := primitiveIndex[kind]; ok {
		return primitive
	}

	switch kind {
	case reflect.String:
		return stringTypeHandler
	case reflect.Struct:
		return getStructHandlerFromType(typ)
	case reflect.Slice:
		return getSliceHandlerFromType(typ)
	case reflect.Map:
		return getMapHandlerFromType(typ)
	case reflect.Uint:
		fallthrough
	case reflect.Int:
		panic("types uint and int are not supported, as their actual size is dependant on compiler architecture and" +
			" could cause data inconsistencies. use uint32/uint64/int32/int64 instead")
	default:
		panic(fmt.Sprintf("Cannot build type handler for type \"%s\" with kind nr. %d", typ.String(), typ.Kind()))
	}
}

type fixed struct{}

func (fixed) isFixed() bool {
	return true
}

type variable struct{}

func (variable) isFixed() bool {
	return false
}
