package iox

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

func assertEq[T any](subject string, a T, b T, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

// -----------------------------------------------------------------------------
// Encoder.
// -----------------------------------------------------------------------------

func TestEncoderImplEncodeIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := EncoderImpl{Impl: gob.NewEncoder(buf).Encode}

	err := enc.Encode("test")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test", string(buf.Bytes()[4:]), func(s string) { t.Fatal(s) })
}

func TestEncoderImplEncodeWithNilImpl(t *testing.T) {
	enc := EncoderImpl{}

	err := enc.Encode("test")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Decoder.
// -----------------------------------------------------------------------------

func TestDecoderImplDecodeIdeal(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`"test"`))
	dec := DecoderImpl{Impl: json.NewDecoder(buf).Decode}

	val := ""
	err := dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test", val, func(s string) { t.Fatal(s) })
}

func TestDecoderImplDecodeWithNilImpl(t *testing.T) {
	dec := DecoderImpl{}

	val := ""
	err := dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// readWriteCloserImpl
// -----------------------------------------------------------------------------

func TestIOReadWriteCloserImplIdeal(t *testing.T) {
	err := *new(error)
	rwc := readWriteCloserImpl{}
	rwc.ImplC = func() error { return nil }
	rwc.ImplR = func([]byte) (int, error) { return 0, nil }
	rwc.ImplW = func([]byte) (int, error) { return 0, nil }

	err = rwc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	_, err = rwc.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	_, err = rwc.Write(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestIOReadWriteCloserImplWithNilImpl(t *testing.T) {
	err := *new(error)
	rwc := readWriteCloserImpl{}

	err = rwc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	_, err = rwc.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rwc.Write(nil)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}
