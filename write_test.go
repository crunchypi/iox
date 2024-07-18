package iox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
)

// -----------------------------------------------------------------------------
// Writer impl.
// -----------------------------------------------------------------------------

func TestWriterImplWriteIdeal(t *testing.T) {
	err := *new(error)
	val := 0

	w := WriterImpl[int]{}
	w.Impl = func(ctx context.Context, v int) error { val = v; return nil }

	err = w.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })
}

func TestWriterImplWriteWithNilImpl(t *testing.T) {
	err := *new(error)
	val := 0

	w := WriterImpl[int]{}

	err = w.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// WriteCloser impl.
// -----------------------------------------------------------------------------

func TestWriteCloserImplWriteIdeal(t *testing.T) {
	err := *new(error)
	val := 0

	wc := WriteCloserImpl[int]{}
	wc.ImplW = func(ctx context.Context, v int) error { val = v; return nil }

	err = wc.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })
}

func TestWriteCloserImplWriteWithNilImpl(t *testing.T) {
	err := *new(error)
	val := 0

	wc := WriteCloserImpl[int]{}

	err = wc.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestWriteCloserImplCloseIdeal(t *testing.T) {
	wc := WriteCloserImpl[int]{}
	wc.ImplC = func() error { return nil }

	err := wc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestWriteCloserImplCloseWithNilImpl(t *testing.T) {
	wc := WriteCloserImpl[int]{}

	err := wc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

func TestNewWriterFromValuesIdeal(t *testing.T) {
	b := bytes.NewBuffer(nil)
	f := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	w := NewWriterFromValues[int](b)(f)

	dec := json.NewDecoder(b)
	err := *new(error)
	val := 0

	err = w.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = w.Write(nil, 3)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromValuesWithNilWriter(t *testing.T) {
	f := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	w := NewWriterFromValues[int](nil)(f)

	err := w.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromValuesWithNilEncoder(t *testing.T) {
	b := bytes.NewBuffer(nil)
	w := NewWriterFromValues[int](b)(nil)

	dec := json.NewDecoder(b)
	err := *new(error)
	val := 0

	err = w.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = w.Write(nil, 3)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromValuesWithEncodeErr(t *testing.T) {
	b := bytes.NewBuffer(nil)
	f := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	w := NewWriterFromValues[chan int](b)(f)

	want := "json: unsupported type: chan int"
	have := w.Write(nil, make(chan int)).Error()
	assertEq("err", want, have, func(s string) { t.Fatal(s) })
}
