package ikea

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

type Holder struct {
	vals [50000000][10]uint16
	buffer    *bytes.Buffer `ikea:"-"`
}

func (h *Holder) Init() {
	h.buffer = new(bytes.Buffer)
}

func (h *Holder) compress() error {
	return Pack(h.buffer, h)
}


var _holder Holder

func init() {
	_holder.Init()

	for i := range _holder.vals {
		for j := range _holder.vals[i] {
			v := rand.Uint32()
			_holder.vals[i][j] = uint16(v)
		}
	}
}

func Test2DArrayPackingAndUnpacking(t *testing.T) {
	if err := _holder.compress(); err != nil {
		panic(err)
	}

	h2 := Holder{}
	h2.Init()
	err := Unpack(_holder.buffer, &h2)
	if err != nil {
		panic(err)
	}

	for i := range _holder.vals {
		for j := range _holder.vals[i] {
			if _holder.vals[i][j] != h2.vals[i][j] {
				panic(fmt.Sprintf("diff values. At vals[%d][%d]", i, j))
			}
		}
	}
}