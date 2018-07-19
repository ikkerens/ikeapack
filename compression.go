package ikea

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type compressionReadWriter struct {
	variable
	handler readWriter
	level   int
}

func (c *compressionReadWriter) readVariable(r io.Reader, v reflect.Value) error {
	lb := make([]byte, 4)
	if _, err := io.ReadFull(r, lb); err != nil {
		return err
	}

	ul := binary.BigEndian.Uint32(lb)
	if ul > math.MaxInt32 {
		return fmt.Errorf("transmitted compressed blob too large (%d>%d)", ul, math.MaxInt32)
	}
	l := int(ul)

	cb := make([]byte, l)
	if _, err := io.ReadFull(r, cb); err != nil {
		return err
	}

	z := flate.NewReader(bytes.NewBuffer(cb))
	defer z.Close()

	return handleVariableReader(z, c.handler, v)
}

func (c *compressionReadWriter) writeVariable(w io.Writer, v reflect.Value) error {
	var b bytes.Buffer

	z, err := flate.NewWriter(&b, c.level)
	if err != nil {
		return err
	}

	if err := handleVariableWriter(z, c.handler, v); err != nil {
		return err
	}
	z.Close()

	lb := make([]byte, 4)
	binary.BigEndian.PutUint32(lb, uint32(b.Len()))
	if _, err := w.Write(lb); err != nil {
		return err
	}
	if _, err = w.Write(b.Bytes()); err != nil {
		return err
	}

	return nil
}

func (c *compressionReadWriter) vLength(v reflect.Value) (int, error) {
	var b bytes.Buffer
	if err := c.writeVariable(&b, v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}
