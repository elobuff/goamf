package amf

import (
	"io"
	"reflect"
)

type Encoder struct {
}

func (e *Encoder) Encode(w io.Writer, val interface{}, ver Version) (int, error) {
	switch ver {
	case AMF0:
		return e.EncodeAmf0(w, val)
	case AMF3:
		return e.EncodeAmf0(w, val)
	}

	return 0, Error("encode amf: unsupported version %d", ver)
}

func (e *Encoder) EncodeAmf0(w io.Writer, val interface{}) (int, error) {
	if val == nil {
		return e.EncodeAmf0Null(w, true)
	}

	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return e.EncodeAmf0Null(w, true)
	}

	switch v.Kind() {
	case reflect.String:
		str := v.String()
		if len(str) <= AMF0_STRING_MAX {
			return e.EncodeAmf0String(w, str, true)
		} else {
			return e.EncodeAmf0LongString(w, str, true)
		}
	case reflect.Bool:
		return e.EncodeAmf0Boolean(w, v.Bool(), true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.EncodeAmf0Number(w, float64(v.Int()), true)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.EncodeAmf0Number(w, float64(v.Uint()), true)
	case reflect.Float32, reflect.Float64:
		return e.EncodeAmf0Number(w, float64(v.Float()), true)
	case reflect.Array, reflect.Slice:
		return 0, Error("encode amf0: unsupported type array")
	case reflect.Map:
		obj, ok := val.(Object)
		if ok != true {
			return 0, Error("encode amf0: unable to create object from map")
		}
		return e.EncodeAmf0Object(w, obj, true)
	}

	if _, ok := val.(TypedObject); ok {
		return 0, Error("encode amf0: unsupported type typed object")
	}

	return 0, Error("encode amf0: unsupported type %s", v.Type())
}

func (e *Encoder) EncodeAmf3(w io.Writer, val interface{}) (int, error) {
	return 0, nil
}
