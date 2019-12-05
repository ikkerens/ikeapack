package ikea

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

// Causes error: too much data in section SBSS / SNOPTRBSS
// this is a issue with Go: https://github.com/golang/go/issues/17378
// This test is currently removed, until go fixes this issue.
//type HolderLarge struct {
//	vals [100000000][10]uint16
//	buffer    *bytes.Buffer `ikea:"-"`
//}
//
//func (h *HolderLarge) Init() {
//	h.buffer = new(bytes.Buffer)
//}
//
//func (h *HolderLarge) compress() error {
//	return Pack(h.buffer, h)
//}

type arr2dHolderSmall struct {
	vals   [10][10]uint16
	buffer *bytes.Buffer `ikea:"-"`
}

func (h *arr2dHolderSmall) Init() {
	h.buffer = new(bytes.Buffer)
}

func (h *arr2dHolderSmall) compress() error {
	return Pack(h.buffer, h)
}

// var _holderL HolderLarge
var _arr2dHolderS arr2dHolderSmall

func init() {
	//_holderL.Init()
	_arr2dHolderS.Init()

	//for i := range _holderL.vals {
	//	for j := range _holderL.vals[i] {
	//		v := rand.Uint32()
	//		_holderL.vals[i][j] = uint16(v)
	//	}
	//}

	for i := range _arr2dHolderS.vals {
		for j := range _arr2dHolderS.vals[i] {
			v := rand.Uint32()
			_arr2dHolderS.vals[i][j] = uint16(v)
		}
	}
}

//func TestLarge2DArrayPackingAndUnpacking(t *testing.T) {
//	if err := _holderL.compress(); err != nil {
//		t.Fatal(err)
//	}
//
//	h2 := HolderLarge{}
//	h2.Init()
//	err := Unpack(_holderL.buffer, &h2)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for i := range _holderL.vals {
//		for j := range _holderL.vals[i] {
//			if _holderL.vals[i][j] != h2.vals[i][j] {
//				t.Fatal(fmt.Sprintf("diff values. At vals[%d][%d]", i, j))
//			}
//		}
//	}
//}

func TestSmall2DArrayPackingAndUnpacking(t *testing.T) {
	if err := _arr2dHolderS.compress(); err != nil {
		t.Fatal(err)
	}

	h2 := arr2dHolderSmall{}
	h2.Init()
	err := Unpack(_arr2dHolderS.buffer, &h2)
	if err != nil {
		t.Fatal(err)
	}

	for i := range _arr2dHolderS.vals {
		for j := range _arr2dHolderS.vals[i] {
			if _arr2dHolderS.vals[i][j] != h2.vals[i][j] {
				t.Fatal(fmt.Sprintf("diff values. At vals[%d][%d]", i, j))
			}
		}
	}
}
