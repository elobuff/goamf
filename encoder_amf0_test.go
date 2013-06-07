package amf

import (
	"bytes"
	"testing"
)

func TestEncodeAmf0Number(t *testing.T) {
	buf := new(bytes.Buffer)
	expect := []byte{0x00, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33}

	enc := new(Encoder)

	n, err := enc.EncodeAmf0(buf, float64(1.2))
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 9 {
		t.Errorf("expected to write 9 bytes, actual %d", n)
	}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("expected buffer: %+v, got: %+v", expect, buf.Bytes())
	}
}

func TestEncodeAmf0BooleanTrue(t *testing.T) {
	buf := new(bytes.Buffer)
	expect := []byte{0x01, 0x01}

	enc := new(Encoder)

	n, err := enc.EncodeAmf0(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 2 {
		t.Errorf("expected to write 2 bytes, actual %d", n)
	}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("expected buffer: %+v, got: %+v", expect, buf.Bytes())
	}
}

func TestEncodeAmf0BooleanFalse(t *testing.T) {
	buf := new(bytes.Buffer)
	expect := []byte{0x01, 0x00}

	enc := new(Encoder)

	n, err := enc.EncodeAmf0(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 2 {
		t.Errorf("expected to write 2 bytes, actual %d", n)
	}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("expected buffer: %+v, got: %+v", expect, buf.Bytes())
	}
}

func TestEncodeAmf0Null(t *testing.T) {
	buf := new(bytes.Buffer)
	expect := []byte{0x05}

	enc := new(Encoder)

	n, err := enc.EncodeAmf0(buf, nil)
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 1 {
		t.Errorf("expected to write 1 byte, actual %d", n)
	}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("expected buffer: %+v, got: %+v", expect, buf.Bytes())
	}
}
