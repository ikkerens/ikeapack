package serialize

import (
	"io"
	"reflect"
)

func handleVariableReader(r io.Reader, h readWriter, v reflect.Value) error {
	switch hr := h.(type) {
	case fixedReadWriter:
		b := make([]byte, hr.length())
		if _, err := io.ReadFull(r, b); err != nil {
			return err
		}

		hr.read(b, v)
	case variableReadWriter:
		if err := hr.read(r, v); err != nil {
			return err
		}
	}

	return nil
}

func handleVariableWriter(w io.Writer, h readWriter, v reflect.Value) error {
	switch hw := h.(type) {
	case fixedReadWriter:
		b := make([]byte, hw.length())
		hw.write(b, v)

		if _, err := w.Write(b); err != nil {
			return err
		}
	case variableReadWriter:
		if err := hw.write(w, v); err != nil {
			return err
		}
	}

	return nil
}
