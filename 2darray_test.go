package ikea

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

// Causes error: too much data in section SBSS / SNOPTRBSS
type HolderLarge struct {
	vals [50000000][10]uint16
	buffer    *bytes.Buffer `ikea:"-"`
}

func (h *HolderLarge) Init() {
	h.buffer = new(bytes.Buffer)
}

func (h *HolderLarge) compress() error {
	return Pack(h.buffer, h)
}

type HolderSmall struct {
	vals [10][10]uint16
	buffer    *bytes.Buffer `ikea:"-"`
}

func (h *HolderSmall) Init() {
	h.buffer = new(bytes.Buffer)
}

func (h *HolderSmall) compress() error {
	return Pack(h.buffer, h)
}


var _holderL HolderLarge
var _holderS HolderLarge

func init() {
	_holderL.Init()
	_holderS.Init()

	for i := range _holderL.vals {
		for j := range _holderL.vals[i] {
			v := rand.Uint32()
			_holderL.vals[i][j] = uint16(v)
		}
	}

	for i := range _holderS.vals {
		for j := range _holderS.vals[i] {
			v := rand.Uint32()
			_holderS.vals[i][j] = uint16(v)
		}
	}
}

func TestLarge2DArrayPackingAndUnpacking(t *testing.T) {
	if err := _holderL.compress(); err != nil {
		t.Fatal(err)
	}

	h2 := HolderLarge{}
	h2.Init()
	err := Unpack(_holderL.buffer, &h2)
	if err != nil {
		t.Fatal(err)
	}

	for i := range _holderL.vals {
		for j := range _holderL.vals[i] {
			if _holderL.vals[i][j] != h2.vals[i][j] {
				t.Fatal(fmt.Sprintf("diff values. At vals[%d][%d]", i, j))
			}
		}
	}
}

func TestSmall2DArrayPackingAndUnpacking(t *testing.T) {
	if err := _holderS.compress(); err != nil {
		t.Fatal(err)
	}

	h2 := HolderSmall{}
	h2.Init()
	err := Unpack(_holderS.buffer, &h2)
	if err != nil {
		t.Fatal(err)
	}

	for i := range _holderS.vals {
		for j := range _holderS.vals[i] {
			if _holderS.vals[i][j] != h2.vals[i][j] {
				t.Fatal(fmt.Sprintf("diff values. At vals[%d][%d]", i, j))
			}
		}
	}
}