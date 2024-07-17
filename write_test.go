package iox

import (
	"context"
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
