package amf

import (
	"encoding/binary"
	"io"
	"math"
)

// marker: 1 byte 0x00
// format: 8 byte big endian float64
func (e *Encoder) EncodeAmf0Number(w io.Writer, val float64, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_NUMBER_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	num := math.Float64bits(val)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, num)

	m, err = w.Write(buf)
	if err != nil {
		return
	}
	n += m

	return
}

// marker: 1 byte 0x01
// format: 1 byte, 0x00 = false, 0x01 = true
func (e *Encoder) EncodeAmf0Boolean(w io.Writer, val bool, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_BOOLEAN_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	buf := make([]byte, 1)
	if val {
		buf[0] = AMF0_BOOLEAN_TRUE
	} else {
		buf[0] = AMF0_BOOLEAN_FALSE
	}

	m, err = w.Write(buf)
	if err != nil {
		return
	}
	n += m

	return
}

// marker: 1 byte 0x02
// format:
// - 2 byte big endian uint16 header to determine size
// - n (size) byte utf8 string
func (e *Encoder) EncodeAmf0String(w io.Writer, val string, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_STRING_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	length := uint16(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return
	}
	n += 2

	m, err = w.Write([]byte(val))
	if err != nil {
		return
	}
	n += m

	return
}

// marker: 1 byte 0x03
// format:
// - loop encoded string followed by encoded value
// - terminated with empty string followed by 1 byte 0x09
func (e *Encoder) EncodeAmf0Object(w io.Writer, val Object, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_OBJECT_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	for k, v := range val {
		m, err = e.EncodeAmf0String(w, k, false)
		if err != nil {
			return n, Error("encode amf0: unable to encode object key: %s", err)
		}
		n += m

		m, err = e.EncodeAmf0(w, v)
		if err != nil {
			return n, Error("encode amf0: unable to encode object value: %s", err)
		}
		n += m
	}

	m, err = e.EncodeAmf0String(w, "", false)
	if err != nil {
		return n, Error("encode amf0: unable to encode object empty string: %s", err)
	}
	n += m

	err = WriteMarker(w, AMF0_OBJECT_END_MARKER)
	if err != nil {
		return n, Error("encode amf0: unable to object end marker: %s", err)
	}
	n += 1

	return
}

// marker: 1 byte 0x05
// no additional data
func (e *Encoder) EncodeAmf0Null(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_NULL_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x06
// no additional data
func (e *Encoder) EncodeAmf0Undefined(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_UNDEFINED_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x0c
// format:
// - 4 byte big endian uint32 header to determine size
// - n (size) byte utf8 string
func (e *Encoder) EncodeAmf0LongString(w io.Writer, val string, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_LONG_STRING_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	length := uint32(len(val))
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return n, Error("encode amf0: unable to encode long string length: %s", err)
	}
	n += 4

	m, err = w.Write([]byte(val))
	if err != nil {
		return
	}
	n += m

	return
}

// marker: 1 byte 0x0d
// no additional data
func (e *Encoder) EncodeAmf0Unsupported(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_UNSUPPORTED_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}
