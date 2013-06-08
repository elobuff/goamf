package amf

import (
	"io"
)

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
