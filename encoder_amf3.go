package amf

import (
	"io"
)

func (e *Encoder) EncodeAmf3Integer(w io.Writer, val uint32, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_INTEGER_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int

	if val <= 0x0000007F {
		err = WriteByte(w, byte(val))
		if err == nil {
			n += 1
		}
		return
	}

	if val <= 0x00003FFF {
		m, err = w.Write([]byte{byte(val>>7 | 0x80), byte(val & 0x7F)})
		n += m
		return
	}

	if val <= 0x001FFFFF {
		m, err = w.Write([]byte{byte(val>>14 | 0x80), byte(val>>7&0x7F | 0x80), byte(val & 0x7F)})
		n += m
		return
	}

	if val <= 0x1FFFFFFF {
		m, err = w.Write([]byte{byte(val>>22 | 0x80), byte(val>>15&0x7F | 0x80), byte(val>>8&0x7F | 0x80), byte(val)})
		n += m
		return
	}

	return n, Error("amf3 encode: cannot encode u29 (out of range)")
}
