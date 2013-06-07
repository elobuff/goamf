package amf

import (
	"bytes"
	"testing"
)

func TestDecodeAmf0Number(t *testing.T) {
	buf := bytes.NewReader([]byte{0x00, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33})
	expect := float64(1.2)

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test number interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0Number(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test number interface without marker
	buf.Seek(1, 0)
	got, err = dec.DecodeAmf0Number(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf0BooleanTrue(t *testing.T) {
	buf := bytes.NewReader([]byte{0x01, 0x01})
	expect := true

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test boolean interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0Boolean(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test boolean interface without marker
	buf.Seek(1, 0)
	got, err = dec.DecodeAmf0Boolean(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf0BooleanFalse(t *testing.T) {
	buf := bytes.NewReader([]byte{0x01, 0x00})
	expect := false

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test boolean interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0Boolean(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test boolean interface without marker
	buf.Seek(1, 0)
	got, err = dec.DecodeAmf0Boolean(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf0String(t *testing.T) {
	buf := bytes.NewReader([]byte{0x02, 0x00, 0x03, 0x66, 0x6f, 0x6f})
	expect := "foo"

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test string interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0String(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test string interface without marker
	buf.Seek(1, 0)
	got, err = dec.DecodeAmf0String(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf0Null(t *testing.T) {
	buf := bytes.NewReader([]byte{0x05})

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}

	// Test null interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0Null(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}
}

func TestDecodeAmf0LongString(t *testing.T) {
	buf := bytes.NewReader([]byte{0x0c, 0x00, 0x00, 0x00, 0x03, 0x66, 0x6f, 0x6f})
	expect := "foo"

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test string interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0LongString(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	// Test string interface without marker
	buf.Seek(1, 0)
	got, err = dec.DecodeAmf0LongString(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf0Unsupported(t *testing.T) {
	buf := bytes.NewReader([]byte{0x0d})

	dec := &Decoder{}

	// Test main interface
	got, err := dec.DecodeAmf0(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}

	// Test null interface with marker
	buf.Seek(0, 0)
	got, err = dec.DecodeAmf0Unsupported(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}
}
