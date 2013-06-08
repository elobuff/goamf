package amf

import (
	"io"
	"reflect"
)

type Encoder struct {
	stringRefs []string
}

func (e *Encoder) Encode(w io.Writer, val interface{}, ver Version) (int, error) {
	switch ver {
	case AMF0:
		return e.EncodeAmf0(w, val)
	case AMF3:
		return e.EncodeAmf3(w, val)
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
	case reflect.Array:
		length := v.Len()
		arr := make(StrictArray, length)
		for i := 0; i < length; i++ {
			arr[i] = v.Index(int(i)).Interface()
		}
		return e.EncodeAmf0StrictArray(w, arr, true)
	case reflect.Slice:
		return 0, Error("encode amf0: unsupported type slice")
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
	if val == nil {
		return e.EncodeAmf3Null(w, true)
	}

	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return e.EncodeAmf3Null(w, true)
	}

	switch v.Kind() {
	case reflect.String:
		return e.EncodeAmf3String(w, v.String(), true)
	case reflect.Bool:
		if v.Bool() {
			return e.EncodeAmf3True(w, true)
		} else {
			return e.EncodeAmf3False(w, true)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return e.EncodeAmf3Integer(w, uint32(v.Int()), true)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return e.EncodeAmf3Integer(w, uint32(v.Uint()), true)
	case reflect.Int64:
		return e.EncodeAmf3Double(w, float64(v.Int()), true)
	case reflect.Uint64:
		return e.EncodeAmf3Double(w, float64(v.Uint()), true)
	case reflect.Float32, reflect.Float64:
		return e.EncodeAmf3Double(w, float64(v.Float()), true)
	case reflect.Array, reflect.Slice:
		return 0, Error("encode amf3: unsupported type array")
	case reflect.Map:
		return 0, Error("encode amf3: unsupported type object")
	}

	if _, ok := val.(TypedObject); ok {
		return 0, Error("encode amf3: unsupported type typed object")
	}

	return 0, Error("encode amf3: unsupported type %s", v.Type())
}
