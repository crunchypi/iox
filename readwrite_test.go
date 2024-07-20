package iox

import (
	"context"
	"io"
	"testing"
)

// -----------------------------------------------------------------------------
// ReadWriter impl.
// -----------------------------------------------------------------------------

func TestReadWriterImplReadIdeal(t *testing.T) {
	rw := ReadWriterImpl[int, int]{}
	rw.ImplR = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReadWriterImplReadWithNilImpl(t *testing.T) {
	rw := ReadWriterImpl[int, int]{}

	val, err := rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadWriterImplWriteIdeal(t *testing.T) {
	err := *new(error)
	val := 0

	rw := ReadWriterImpl[int, int]{}
	rw.ImplW = func(ctx context.Context, v int) error { val = v; return nil }

	err = rw.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })
}

func TestReadWriterImplWriteWithNilImpl(t *testing.T) {
	rw := ReadWriterImpl[int, int]{}

	err := rw.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// ReadWriteCloser impl.
// -----------------------------------------------------------------------------

func TestReadWriteCloserImplReadIdeal(t *testing.T) {
	rwc := ReadWriteCloserImpl[int, int]{}
	rwc.ImplR = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := rwc.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReadWriteCloserImplReadWithNilImpl(t *testing.T) {
	rwc := ReadWriteCloserImpl[int, int]{}

	val, err := rwc.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadWriteCloserImplWriteIdeal(t *testing.T) {
	err := *new(error)
	val := 0

	rwc := ReadWriteCloserImpl[int, int]{}
	rwc.ImplW = func(ctx context.Context, v int) error { val = v; return nil }

	err = rwc.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })
}

func TestReadWriteCloserImplWriteWithNilImpl(t *testing.T) {
	err := ReadWriteCloserImpl[int, int]{}.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestReadWriteCloserImplCloseIdeal(t *testing.T) {
	rwc := ReadWriteCloserImpl[int, int]{ImplC: func() error { return nil }}
	assertEq("err", *new(error), rwc.Close(), func(s string) { t.Fatal(s) })
}

func TestReadWriteCloserImplCloseWithNilImpl(t *testing.T) {
	rwc := ReadWriteCloserImpl[int, int]{}
	assertEq("err", *new(error), rwc.Close(), func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

func TestNewReadWriterFromIdeal(t *testing.T) {
	rw := NewReadWriterFrom(1, 2)

	val := 0
	err := *new(error)

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })

	err = rw.Write(nil, 3)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })
}
