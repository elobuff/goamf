package amf

import (
	"encoding/binary"
	"io"
	"math"
)

// marker: 1 byte 0x00
// format: 8 byte big endian float64
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

// marker: 1 byte 0x01
// format: 1 byte, 0x00 = false, 0x01 = true
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

// marker: 1 byte 0x02
// format:
// - 2 byte big endian uint16 header to determine size
// - n (size) byte utf8 string
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

// marker: 1 byte 0x03
// format:
// - encoded string followed by encoded value
// - terminated with empty string followed by 1 byte 0x09
func (d *Decoder) DecodeAmf0Object(r io.Reader, x bool) (Object, error) {
	if err := AssertMarker(r, x, AMF0_OBJECT_MARKER); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	for {
		key, err := d.DecodeAmf0String(r, false)
		if err != nil {
			return nil, err
		}

		if key == "" {
			if err = AssertMarker(r, true, AMF0_OBJECT_END_MARKER); err != nil {
				return nil, Error("decode amf0: expected object end marker")
			}

			break
		}

		value, err := d.DecodeAmf0(r)
		if err != nil {
			return nil, Error("decode amf0: unable to decode object value: %s", err)
		}

		result[key] = value
	}

	return result, nil

}

// marker: 1 byte 0x05
// no additional data
func (d *Decoder) DecodeAmf0Null(r io.Reader, x bool) (result interface{}, err error) {
	err = AssertMarker(r, x, AMF0_NULL_MARKER)
	return
}

// marker: 1 byte 0x06
// no additional data
func (d *Decoder) DecodeAmf0Undefined(r io.Reader, x bool) (result interface{}, err error) {
	err = AssertMarker(r, x, AMF0_UNDEFINED_MARKER)
	return
}

// marker: 1 byte 0x0c
// format:
// - 4 byte big endian uint32 header to determine size
// - n (size) byte utf8 string
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

// marker: 1 byte 0x0d
// no additional data
func (d *Decoder) DecodeAmf0Unsupported(r io.Reader, x bool) (result interface{}, err error) {
	err = AssertMarker(r, x, AMF0_UNSUPPORTED_MARKER)
	return
}
