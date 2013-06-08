package amf

import (
	"encoding/binary"
	"io"
)

// marker: 1 byte 0x00
// no additional data
func (e *Encoder) EncodeAmf3Undefined(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_UNDEFINED_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x01
// no additional data
func (e *Encoder) EncodeAmf3Null(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_NULL_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x02
// no additional data
func (e *Encoder) EncodeAmf3False(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_FALSE_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x03
// no additional data
func (e *Encoder) EncodeAmf3True(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_TRUE_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x04
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

// marker: 1 byte 0x05
func (e *Encoder) EncodeAmf3Double(w io.Writer, val float64, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_DOUBLE_MARKER); err != nil {
			return
		}
		n += 1
	}

	err = binary.Write(w, binary.BigEndian, &val)
	if err != nil {
		return
	}
	n += 8

	return
}
