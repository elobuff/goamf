package amf

import (
	"errors"
	"fmt"
	"io"
)

type Decoder struct {
	refCache    []interface{}
	stringCache []interface{}
	objectCache []interface{}
	traitCache  []interface{}
}

func (d *Decoder) Decode(r io.Reader, ver Version) (interface{}, error) {
	switch ver {
	case 0:
		return d.DecodeAmf0(r)
	case 3:
		return d.DecodeAmf3(r)
	}

	return nil, errors.New(fmt.Sprintf("decode amf: unsupported version %d", ver))
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
		return d.DecodeAmf0Object(r, false)
	case AMF0_MOVIECLIP_MARKER:
		return nil, Error("decode amf0: unsupported type movieclip")
	case AMF0_NULL_MARKER:
		return d.DecodeAmf0Null(r, false)
	case AMF0_UNDEFINED_MARKER:
		return d.DecodeAmf0Undefined(r, false)
	case AMF0_REFERENCE_MARKER:
		return nil, Error("decode amf0: unsupported type reference")
	case AMF0_ECMA_ARRAY_MARKER:
		return d.DecodeAmf0EcmaArray(r, false)
	case AMF0_STRICT_ARRAY_MARKER:
		return d.DecodeAmf0StrictArray(r, false)
	case AMF0_DATE_MARKER:
		return d.DecodeAmf0Date(r, false)
	case AMF0_LONG_STRING_MARKER:
		return d.DecodeAmf0LongString(r, false)
	case AMF0_UNSUPPORTED_MARKER:
		return d.DecodeAmf0Unsupported(r, false)
	case AMF0_RECORDSET_MARKER:
		return nil, Error("decode amf0: unsupported type recordset")
	case AMF0_XML_DOCUMENT_MARKER:
		return d.DecodeAmf0XmlDocument(r, false)
	case AMF0_TYPED_OBJECT_MARKER:
		return d.DecodeAmf0TypedObject(r, false)
	case AMF0_ACMPLUS_OBJECT_MARKER:
		return d.DecodeAmf3(r)
	}

	return nil, Error("decode amf0: unsupported type %d", marker)
}

func (d *Decoder) DecodeAmf3(r io.Reader) (interface{}, error) {
	marker, err := ReadMarker(r)
	if err != nil {
		return nil, err
	}

	return nil, Error("decode amf3: unsupported type %d", marker)
}
