package amf

import (
	"bytes"
	"testing"
)

func TestEncodeAmf3Integer(t *testing.T) {
	enc := new(Encoder)

	for _, tc := range u29TestCases {
		buf := new(bytes.Buffer)
		_, err := enc.EncodeAmf3Integer(buf, tc.value, false)
		if err != nil {
			t.Errorf("EncodeAmf3Integer error: %s", err)
		}
		got := buf.Bytes()
		if !bytes.Equal(tc.expect, got) {
			t.Errorf("EncodeAmf3Integer expect n %x got %x", tc.value, got)
		}
	}

	buf := new(bytes.Buffer)
	expect := []byte{0x04, 0x80, 0xFF, 0xFF, 0xFF}

	n, err := enc.EncodeAmf3(buf, uint32(4194303))
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 5 {
		t.Errorf("expected to write 5 bytes, actual %d", n)
	}
	if bytes.Compare(buf.Bytes(), expect) != 0 {
		t.Errorf("expected buffer: %+v, got: %+v", expect, buf.Bytes())
	}
}
