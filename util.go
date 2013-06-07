package amf

import (
	"errors"
	"fmt"
	"github.com/jcoene/gologger"
	"io"
)

var log logger.Logger = *logger.NewLogger(logger.LOG_LEVEL_DEBUG, "amf")

func Error(f string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(f, v...))
}

func ReadByte(r io.Reader) (byte, error) {
	bytes, err := ReadBytes(r, 1)
	if err != nil {
		return 0x00, err
	}

	return bytes[0], nil
}

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)

	m, err := r.Read(bytes)
	if err != nil {
		return bytes, err
	}

	if m != n {
		return bytes, Error("decode read bytes failed: expected %d got %d", m, n)
	}

	return bytes, nil
}

func ReadMarker(r io.Reader) (byte, error) {
	return ReadByte(r)
}

func AssertMarker(r io.Reader, x bool, m byte) error {
	if x == false {
		return nil
	}

	marker, err := ReadMarker(r)
	if err != nil {
		return err
	}

	if marker != m {
		return Error("decode assert marker failed: expected %v got %v", m, marker)
	}

	return nil
}
