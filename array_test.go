package ikea

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

type arr1DHolder struct {
	vals   [10]uint16
	buffer *bytes.Buffer `ikea:"-"`
}

func (h *arr1DHolder) Init() {
	h.buffer = new(bytes.Buffer)
}

func (h *arr1DHolder) compress() error {
	return Pack(h.buffer, h)
}

// var _holderL HolderLarge
var _arr1dHolder arr1DHolder

func init() {
	_arr1dHolder.Init()
	for i := range _arr1dHolder.vals {
		v := rand.Uint32()
		_arr1dHolder.vals[i] = uint16(v)
	}
}

func Test1DArrayPackingAndUnpacking(t *testing.T) {
	if err := _arr1dHolder.compress(); err != nil {
		t.Fatal(err)
	}

	h2 := arr1DHolder{}
	h2.Init()
	err := Unpack(_arr1dHolder.buffer, &h2)
	if err != nil {
		t.Fatal(err)
	}

	for i := range _arr1dHolder.vals {
		if _arr1dHolder.vals[i] != h2.vals[i] {
			t.Fatal(fmt.Sprintf("diff values. At vals[%d]", i))
		}
	}
}
