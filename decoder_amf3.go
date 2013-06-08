package amf

import (
	"encoding/binary"
	"io"
)

// marker: 1 byte 0x00
// no additional data
func (d *Decoder) DecodeAmf3Undefined(r io.Reader, decodeMarker bool) (result interface{}, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_UNDEFINED_MARKER)
	return
}

// marker: 1 byte 0x01
// no additional data
func (d *Decoder) DecodeAmf3Null(r io.Reader, decodeMarker bool) (result interface{}, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_NULL_MARKER)
	return
}

// marker: 1 byte 0x02
// no additional data
func (d *Decoder) DecodeAmf3False(r io.Reader, decodeMarker bool) (result bool, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_FALSE_MARKER)
	result = false
	return
}

// marker: 1 byte 0x03
// no additional data
func (d *Decoder) DecodeAmf3True(r io.Reader, decodeMarker bool) (result bool, err error) {
	err = AssertMarker(r, decodeMarker, AMF3_TRUE_MARKER)
	result = true
	return
}

// marker: 1 byte 0x04
func (d *Decoder) DecodeAmf3Integer(r io.Reader, decodeMarker bool) (result uint32, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_INTEGER_MARKER); err != nil {
		return
	}

	var b byte

	for i := 0; i < 3; i++ {
		b, err = ReadByte(r)
		if err != nil {
			return
		}
		result = (result << 7) + uint32(b&0x7F)
		if (b & 0x80) == 0 {
			return
		}
	}
	b, err = ReadByte(r)
	if err != nil {
		return
	}

	return ((result << 8) + uint32(b)), nil
}

// marker: 1 byte 0x05
func (d *Decoder) DecodeAmf3Double(r io.Reader, decodeMarker bool) (result float64, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_DOUBLE_MARKER); err != nil {
		return
	}

	err = binary.Read(r, binary.BigEndian, &result)
	if err != nil {
		return float64(0), Error("amf3 decode: unable to read double: %s", err)
	}

	return
}

// marker: 1 byte 0x06
// format:
// - u29 reference int. if reference, no more data. if not reference,
//   length value of bytes to read to complete string.
func (d *Decoder) DecodeAmf3String(r io.Reader, decodeMarker bool) (result string, err error) {
	if err = AssertMarker(r, decodeMarker, AMF3_STRING_MARKER); err != nil {
		return
	}

	var ref bool
	var length uint32
	ref, length, err = d.decodeReferenceInt(r)
	if err != nil {
		return "", Error("amf3 decode: unable to decode reference and length: %s", err)
	}

	if ref {
		if length > uint32(len(d.stringRefs)) {
			return "", Error("amf3 decode: bad string reference")
		}

		result = d.stringRefs[length]
		return
	}

	buf := make([]byte, length)
	_, err = r.Read(buf)
	if err != nil {
		return "", Error("amf3 decode: unable to read string: %s", err)
	}

	result = string(buf)
	if result != "" {
		d.stringRefs = append(d.stringRefs, result)
	}

	return
}

func (d *Decoder) decodeReferenceInt(r io.Reader) (ref bool, val uint32, err error) {
	u29, err := d.DecodeAmf3Integer(r, false)
	if err != nil {
		return false, 0, Error("amf3 decode: unable to decode reference int: %s", err)
	}

	ref = u29&0x01 == 0
	val = u29 >> 1

	return
}
