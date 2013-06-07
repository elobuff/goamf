package amf

import (
	"io"
)

type Encoder struct {
}

func (e *Encoder) Encode(w io.Writer, val interface{}, ver Version) (int, error) {
	switch ver {
	case AMF0:
		return e.EncodeAmf0(w, val)
	case AMF3:
		return e.EncodeAmf0(w, val)
	}

	return 0, Error("encode amf: unsupported version %d", ver)
}

func (e *Encoder) EncodeAmf0(w io.Writer, val interface{}) (int, error) {
	return 0, Error("encode amf0: unsupported type %s", v.Type())
}

func (e *Encoder) EncodeAmf3(w io.Writer, val interface{}) (int, error) {
	return 0, nil
}
