package serialize

import (
	"encoding/binary"
	"math"
	"reflect"
)

type primitiveReadWriter struct {
	fixed

	size   int
	reader func([]byte, reflect.Value)
	writer func([]byte, reflect.Value)
}

func (p *primitiveReadWriter) length() int {
	return p.size
}

func (p *primitiveReadWriter) readFixed(data []byte, v reflect.Value) {
	p.reader(data, v)
}

func (p *primitiveReadWriter) writeFixed(data []byte, v reflect.Value) {
	p.writer(data, v)
}

var primitiveIndex = map[reflect.Kind]*primitiveReadWriter{
	reflect.Int: {
		size: 0,
		reader: func(b []byte, v reflect.Value) {
			panic("only integers with an explicit length are supported (e.g. int8, int16, int32, int64)")
		},
		writer: func(b []byte, v reflect.Value) {
			panic("only integers with an explicit length are supported (e.g. int8, int16, int32, int64)")
		},
	},
	reflect.Uint: {
		size: 0,
		reader: func(b []byte, v reflect.Value) {
			panic("only integers with an explicit length are supported (e.g. int8, int16, int32, int64)")
		},
		writer: func(b []byte, v reflect.Value) {
			panic("only integers with an explicit length are supported (e.g. int8, int16, int32, int64)")
		},
	},
	reflect.Bool: {
		size: 1,
		reader: func(b []byte, v reflect.Value) {
			v.SetBool(b[0] != 0)
		},
		writer: func(b []byte, v reflect.Value) {
			if v.Bool() {
				b[0] = 1
			}
		},
	},
	reflect.Int8: {
		size: 1,
		reader: func(b []byte, v reflect.Value) {
			v.SetInt(int64(b[0]))
		},
		writer: func(b []byte, v reflect.Value) {
			b[0] = byte(v.Int())
		},
	},
	reflect.Int16: {
		size: 2,
		reader: func(b []byte, v reflect.Value) {
			v.SetInt(int64(binary.BigEndian.Uint16(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint16(b, uint16(v.Int()))
		},
	},
	reflect.Int32: {
		size: 4,
		reader: func(b []byte, v reflect.Value) {
			v.SetInt(int64(binary.BigEndian.Uint32(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint32(b, uint32(v.Int()))
		},
	},
	reflect.Int64: {
		size: 8,
		reader: func(b []byte, v reflect.Value) {
			v.SetInt(int64(binary.BigEndian.Uint64(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint64(b, uint64(v.Int()))
		},
	},
	reflect.Uint8: {
		size: 1,
		reader: func(b []byte, v reflect.Value) {
			v.SetUint(uint64(b[0]))
		},
		writer: func(b []byte, v reflect.Value) {
			b[0] = byte(v.Uint())
		},
	},
	reflect.Uint16: {
		size: 2,
		reader: func(b []byte, v reflect.Value) {
			v.SetUint(uint64(binary.BigEndian.Uint16(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint16(b, uint16(v.Uint()))
		},
	},
	reflect.Uint32: {
		size: 4,
		reader: func(b []byte, v reflect.Value) {
			v.SetUint(uint64(binary.BigEndian.Uint32(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint32(b, uint32(v.Uint()))
		},
	},
	reflect.Uint64: {
		size: 8,
		reader: func(b []byte, v reflect.Value) {
			v.SetUint(binary.BigEndian.Uint64(b))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint64(b, v.Uint())
		},
	},
	reflect.Float32: {
		size: 4,
		reader: func(b []byte, v reflect.Value) {
			v.SetFloat(float64(math.Float32frombits(binary.BigEndian.Uint32(b))))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint32(b, math.Float32bits(float32(v.Float())))
		},
	},
	reflect.Float64: {
		size: 8,
		reader: func(b []byte, v reflect.Value) {
			v.SetFloat(math.Float64frombits(binary.BigEndian.Uint64(b)))
		},
		writer: func(b []byte, v reflect.Value) {
			binary.BigEndian.PutUint64(b, math.Float64bits(v.Float()))
		},
	},
}
