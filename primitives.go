package serialize

import (
	"encoding/binary"
	"math"
	"reflect"
)

type primitiveReadWriter struct {
	reader func([]byte, reflect.Value)
	writer func([]byte, reflect.Value)
}

func (p *primitiveReadWriter) read(data []byte, v reflect.Value) {
	p.reader(data, v)
}

func (p *primitiveReadWriter) write(data []byte, v reflect.Value) {
	p.writer(data, v)
}

var primitiveIndex = map[reflect.Kind]*typeHandler{
	reflect.Bool: {
		1,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetBool(b[0] != 0)
			},
			func(b []byte, v reflect.Value) {
				if v.Bool() {
					b[0] = 1
				}
			},
		},
	},
	reflect.Int8: {
		1,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetInt(int64(b[0]))
			},
			func(b []byte, v reflect.Value) {
				b[0] = byte(v.Int())
			},
		},
	},
	reflect.Int16: {
		2,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetInt(int64(binary.BigEndian.Uint16(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint16(b, uint16(v.Int()))
			},
		},
	},
	reflect.Int32: {
		4,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetInt(int64(binary.BigEndian.Uint32(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint32(b, uint32(v.Int()))
			},
		},
	},
	reflect.Int64: {
		8,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetInt(int64(binary.BigEndian.Uint64(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint64(b, uint64(v.Int()))
			},
		},
	},
	reflect.Uint8: {
		1,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetUint(uint64(b[0]))
			},
			func(b []byte, v reflect.Value) {
				b[0] = byte(v.Uint())
			},
		},
	},
	reflect.Uint16: {
		2,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetUint(uint64(binary.BigEndian.Uint16(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint16(b, uint16(v.Uint()))
			},
		},
	},
	reflect.Uint32: {
		4,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetUint(uint64(binary.BigEndian.Uint32(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint32(b, uint32(v.Uint()))
			},
		},
	},
	reflect.Uint64: {
		8,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetUint(binary.BigEndian.Uint64(b))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint64(b, v.Uint())
			},
		},
	},
	reflect.Float32: {
		4,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetFloat(float64(math.Float32frombits(binary.BigEndian.Uint32(b))))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint32(b, math.Float32bits(float32(v.Float())))
			},
		},
	},
	reflect.Float64: {
		8,
		&primitiveReadWriter{
			func(b []byte, v reflect.Value) {
				v.SetFloat(math.Float64frombits(binary.BigEndian.Uint64(b)))
			},
			func(b []byte, v reflect.Value) {
				binary.BigEndian.PutUint64(b, math.Float64bits(v.Float()))
			},
		},
	},
}
