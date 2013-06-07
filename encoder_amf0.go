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
