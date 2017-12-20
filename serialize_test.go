package serialize

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

type test struct {
	A uint32
	B uint64
	C []uint16
	D string
	E []uint16 `compressed:"true"`
}

var compare []byte

const compressionConst = 42

var compressionValue []uint16

func TestMain(m *testing.M) {
	compare, _ = hex.DecodeString("000000010000000000000001000000040001000200030004000000047465737400000098ecc1311100000800a1b78efd033a5883a3a916000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e05d000000ffff")
	compressionValue = make([]uint16, 16*256*16)
	for i := range compressionValue {
		compressionValue[i] = compressionConst
	}

	os.Exit(m.Run())
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Write(buf, &test{1, 1, []uint16{1, 2, 3, 4}, "test", compressionValue}); err != nil {
		t.Error(err)
		return
	}

	result := buf.Bytes()
	if len(result) != len(compare) {
		fmt.Printf("Failing TestWrite, result \"%s\" length (%d) does not match compare length (%d)\n", hex.EncodeToString(result), len(result), len(compare))
		t.FailNow()
	}

	for i := 0; i < len(compare); i++ {
		if result[i] != compare[i] {
			fmt.Printf("Failing TestWrite, hex output \"%s\" does not match compare slice\n", hex.EncodeToString(result))
			t.FailNow()
		}
	}
}

func TestRead(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(compare)

	tst := new(test)
	if err := Read(buf, tst); err != nil {
		t.Error(err)
		return
	}

	if tst.A != 1 || tst.B != 1 ||
		len(tst.C) != 4 || tst.C[0] != 1 || tst.C[1] != 2 || tst.C[2] != 3 || tst.C[3] != 4 ||
		tst.D != "test" {
		fmt.Printf("Failing TestRead, result struct: %+v\n", tst)
		t.FailNow()
	}

	if len(tst.E) != len(compressionValue) {
		fmt.Printf("Failing TestRead, decompressed array length does not match (expected: %d, got %d)", len(compressionValue), len(tst.E))
		t.FailNow()
	}

	for _, val := range tst.E {
		if val != compressionConst {
			fmt.Printf("Failing TestRead, decompressed array entry does not match expected value")
			t.FailNow()
		}
	}
}
