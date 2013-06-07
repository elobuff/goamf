package amf

import (
	"io"
)

type Encoder struct {
}

func (e *Encoder) Encode(w io.Writer, obj interface{}, v Version) (uint32, error) {
	switch v {
	case AMF0:
		return e.EncodeAmf0(w, obj)
	case AMF3:
		return e.EncodeAmf0(w, obj)
	}

	return 0, Error("encode amf: unsupported version %d", v)
}

func (e *Encoder) EncodeAmf0(w io.Writer, obj interface{}) (n uint32, err error) {
	return
}

func (e *Encoder) EncodeAmf3(w io.Writer, obj interface{}) (n uint32, err error) {
	return
}
