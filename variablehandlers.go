package ikea

import (
	"io"
	"reflect"
)

func handleVariableReader(r io.Reader, h readWriter, v reflect.Value) error {
	if h.isFixed() {
		hr := h.(fixedReadWriter)
		b := make([]byte, hr.length())
		if _, err := io.ReadFull(r, b); err != nil {
			return err
		}

		hr.readFixed(b, v)
	} else {
		hr := h.(variableReadWriter)
		if err := hr.readVariable(r, v); err != nil {
			return err
		}
	}

	return nil
}

func handleVariableWriter(w io.Writer, h readWriter, v reflect.Value) error {
	if h.isFixed() {
		hw := h.(fixedReadWriter)
		b := make([]byte, hw.length())
		hw.writeFixed(b, v)

		if _, err := w.Write(b); err != nil {
			return err
		}
	} else {
		hw := h.(variableReadWriter)
		if err := hw.writeVariable(w, v); err != nil {
			return err
		}
	}

	return nil
}

func handleVariableLength(h readWriter, v reflect.Value) (int, error) {
	if h.isFixed() {
		hl := h.(fixedReadWriter)
		return hl.length(), nil
	}

	// variable
	hl := h.(variableReadWriter)
	return hl.vLength(v)
}
