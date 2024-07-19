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
