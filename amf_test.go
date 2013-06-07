package amf

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func EncodeAndDecode(val interface{}, ver Version) (result interface{}, err error) {
	enc := new(Encoder)
	dec := new(Decoder)

	buf := new(bytes.Buffer)

	_, err = enc.Encode(buf, val, ver)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in encode: %s", err))
	}

	result, err = dec.Decode(buf, ver)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in decode: %s", err))
	}

	return
}

func Compare(val interface{}, ver Version, name string, t *testing.T) {
	result, err := EncodeAndDecode(val, ver)
	if err != nil {
		t.Errorf("%s: %s", name, err)
	}
	if val != result {
		t.Errorf("%s: comparison failed between %+v and %+v", name, val, result)
	}
}

func TestAmf0Number(t *testing.T) {
	Compare(float64(6), 0, "amf0 number uint32", t)
	Compare(float64(1245), 0, "amf0 number int32", t)
	Compare(float64(12345.678), 0, "amf0 number float64", t)
}

func TestAmf0String(t *testing.T) {
	Compare("a pup!", 0, "amf0 string simple", t)
	Compare("日本語", 0, "amf0 string utf8", t)
}

func TestAmf0Boolean(t *testing.T) {
	Compare(true, 0, "amf0 boolean true", t)
	Compare(false, 0, "amf0 boolean false", t)
}

func TestAmf0Null(t *testing.T) {
	Compare(nil, 0, "amf0 boolean nil", t)
}

func TestAmf0Object(t *testing.T) {
	obj := make(Object)
	obj["dog"] = "alfie"
	obj["coffee"] = true
	obj["drugs"] = false
	obj["pi"] = 3.14159

	res, err := EncodeAndDecode(obj, 0)
	if err != nil {
		t.Errorf("amf0 object: %s", err)
	}

	result, ok := res.(Object)
	if ok != true {
		t.Errorf("amf0 object conversion failed")
	}

	if result["dog"] != "alfie" {
		t.Errorf("amf0 object string: comparison failed")
	}

	if result["coffee"] != true {
		t.Errorf("amf0 object true: comparison failed")
	}

	if result["drugs"] != false {
		t.Errorf("amf0 object false: comparison failed")
	}

	if result["pi"] != float64(3.14159) {
		t.Errorf("amf0 object float: comparison failed")
	}
}
