package amf

import (
	"encoding/binary"
	"io"
	"math"
)

func (d *Decoder) DecodeAmf0Number(r io.Reader, x bool) (result float64, err error) {
	if err = AssertMarker(r, x, AMF0_NUMBER_MARKER); err != nil {
		return
	}

	var bytes []byte

	if bytes, err = ReadBytes(r, 8); err != nil {
		return
	}

	u64n := binary.BigEndian.Uint64(bytes)
	result = math.Float64frombits(u64n)

	return
}

func (d *Decoder) DecodeAmf0Boolean(r io.Reader, x bool) (result bool, err error) {
	if err = AssertMarker(r, x, AMF0_BOOLEAN_MARKER); err != nil {
		return
	}

	var b byte
	if b, err = ReadByte(r); err != nil {
		return
	}

	if b == AMF0_BOOLEAN_FALSE {
		return false, nil
	} else if b == AMF0_BOOLEAN_TRUE {
		return true, nil
	}

	return false, Error("decode amf0: unexpected value %v for boolean", b)
}

func (d *Decoder) DecodeAmf0String(r io.Reader, x bool) (result string, err error) {
	if err = AssertMarker(r, x, AMF0_STRING_MARKER); err != nil {
		return
	}

	var bytes []byte
	if bytes, err = ReadBytes(r, 2); err != nil {
		return
	}

	len := binary.BigEndian.Uint16(bytes)

	if bytes, err = ReadBytes(r, int(len)); err != nil {
		return
	}

	return string(bytes), nil
}

func (d *Decoder) DecodeAmf0Null(r io.Reader, x bool) (result interface{}, err error) {
	err = AssertMarker(r, x, AMF0_NULL_MARKER)
	return
}

func (d *Decoder) DecodeAmf0LongString(r io.Reader, x bool) (result string, err error) {
	if err = AssertMarker(r, x, AMF0_LONG_STRING_MARKER); err != nil {
		return
	}

	var bytes []byte
	if bytes, err = ReadBytes(r, 4); err != nil {
		return
	}

	len := binary.BigEndian.Uint32(bytes)

	if bytes, err = ReadBytes(r, int(len)); err != nil {
		return
	}

	return string(bytes), nil
}
