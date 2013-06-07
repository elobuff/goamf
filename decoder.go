package amf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type Decoder struct {
	refCache    []interface{}
	stringCache []interface{}
	objectCache []interface{}
	traitCache  []interface{}
}

func (d *Decoder) Decode(r io.Reader, v Version) (interface{}, error) {
	switch v {
	case 0:
		return d.DecodeAmf0(r)
	case 3:
		return d.DecodeAmf3(r)
	}

	return nil, errors.New(fmt.Sprintf("unsupported amf version %d", v))
}

func (d *Decoder) DecodeAmf0(r io.Reader) (interface{}, error) {
	marker, err := ReadMarker(r)
	if err != nil {
		return nil, err
	}

	switch marker {
	case AMF0_NUMBER_MARKER:
		return d.DecodeAmf0Number(r, false)
	case AMF0_BOOLEAN_MARKER:
		return d.DecodeAmf0Boolean(r, false)
	case AMF0_STRING_MARKER:
		return d.DecodeAmf0String(r, false)
	case AMF0_OBJECT_MARKER:
		return nil, errors.New("decode amf0: unsupported type object")
	case AMF0_MOVIECLIP_MARKER:
		return nil, errors.New("decode amf0: unsupported type movieclip")
	case AMF0_NULL_MARKER:
		return nil, errors.New("decode amf0: unsupported type null")
	case AMF0_UNDEFINED_MARKER:
		return nil, errors.New("decode amf0: unsupported type undefined")
	case AMF0_REFERENCE_MARKER:
		return nil, errors.New("decode amf0: unsupported type reference")
	case AMF0_ECMA_ARRAY_MARKER:
		return nil, errors.New("decode amf0: unsupported type ecma array")
	case AMF0_STRICT_ARRAY_MARKER:
		return nil, errors.New("decode amf0: unsupported type strict array")
	case AMF0_DATE_MARKER:
		return nil, errors.New("decode amf0: unsupported type date")
	case AMF0_LONG_STRING_MARKER:
		return nil, errors.New("decode amf0: unsupported type long string")
	case AMF0_UNSUPPORTED_MARKER:
		return nil, errors.New("decode amf0: unsupported type unsupported")
	case AMF0_RECORDSET_MARKER:
		return nil, errors.New("decode amf0: unsupported type recordset")
	case AMF0_XML_DOCUMENT_MARKER:
		return nil, errors.New("decode amf0: unsupported type xml document")
	case AMF0_TYPED_OBJECT_MARKER:
		return nil, errors.New("decode amf0: unsupported type typed object")
	case AMF0_ACMPLUS_OBJECT_MARKER:
		return nil, errors.New("decode amf0: unsupported type acm plus object")
	}

	return nil, nil
}

func (d *Decoder) DecodeAmf3(r io.Reader) (interface{}, error) {
	return nil, nil
}

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

	return false, errors.New(fmt.Sprintf("decode boolean failed: unexpected value %v", b))
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
