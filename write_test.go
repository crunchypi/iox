package iox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
)

// -----------------------------------------------------------------------------
// Test utils.
// -----------------------------------------------------------------------------

func newSliceWriter[T any](s *[]T) Writer[T] {
	return WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			*s = append(*s, v)
			return nil
		},
	}
}

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

func TestNewWriterFromBytesIdeal(t *testing.T) {
	s := make([]int, 0, 3)
	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	w := NewWriterFromBytes(newSliceWriter(&s))(f)

	json.NewEncoder(w).Encode(2)
	json.NewEncoder(w).Encode(3)

	assertEq("s", []int{2, 3}, s, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromBytesWithNilWriter(t *testing.T) {
	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	w := NewWriterFromBytes[int](nil)(f)

	err := json.NewEncoder(w).Encode(2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromBytesWithNilDecoder(t *testing.T) {
	s := make([]int, 0, 3)
	w := NewWriterFromBytes(newSliceWriter(&s))(nil)

	json.NewEncoder(w).Encode(2)
	json.NewEncoder(w).Encode(3)

	assertEq("s", []int{2, 3}, s, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromBytesWithDecodeErr(t *testing.T) {
	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	w := NewWriterFromBytes(WriterImpl[int]{})(f)

	_, err := w.Write([]byte("["))
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewWriterFromBytesWithWriteErr(t *testing.T) {
	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	w := NewWriterFromBytes(WriterImpl[int]{})(f)

	want := io.EOF
	have := json.NewEncoder(w).Encode(1)
	assertEq("err", want, have, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Modifiers.
// -----------------------------------------------------------------------------

func TestWriterWithBatchingIdeal(t *testing.T) {
	s := make([][]int, 0, 2)
	w := NewWriterWithBatching(newSliceWriter(&s), 2)

	assertEq("err", *new(error), w.Write(nil, 2), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 3), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 4), func(s string) { t.Fatal(s) })

	assertEq("len", 1, len(s), func(s string) { t.Fatal(s) })
	assertEq("val", []int{2, 3}, s[0], func(s string) { t.Fatal(s) })

	assertEq("err", *new(error), w.Write(nil, 5), func(s string) { t.Fatal(s) })
	assertEq("len", 2, len(s), func(s string) { t.Fatal(s) })
	assertEq("val", []int{4, 5}, s[1], func(s string) { t.Fatal(s) })
}

func TestWriterWithBatchingWithNilWriter(t *testing.T) {
	w := NewWriterWithBatching[int](nil, 2)

	err := w.Write(nil, 2)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestWriterWithUnbatchingIdeal(t *testing.T) {
	s := make([]int, 0, 4)
	w := NewWriterWithUnbatching(newSliceWriter(&s))

	assertEq("err", *new(error), w.Write(nil, []int{1, 2}), func(s string) { t.Fatal(s) })
	assertEq("val", []int{1, 2}, s, func(s string) { t.Fatal(s) })
}

func TestWriterWithUnbatchingWithNilReader(t *testing.T) {
	w := NewWriterWithUnbatching[int](nil)
	assertEq("err", io.EOF, w.Write(nil, []int{1, 2}), func(s string) { t.Fatal(s) })
}

func TestWriterWithUnbatchingWithCustomErrr(t *testing.T) {
	vw := WriterImpl[int]{}
	vw.Impl = func(ctx context.Context, i int) error { return io.ErrClosedPipe }

	sw := NewWriterWithUnbatching(vw)
	err := sw.Write(nil, []int{1, 2})
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestNewWriterWithFilterFnIdeal(t *testing.T) {
	s := make([]int, 0, 2)
	w := NewWriterWithFilterFn(newSliceWriter(&s))(func(v int) bool { return v%2 != 0 })

	assertEq("err", *new(error), w.Write(nil, 1), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 2), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 3), func(s string) { t.Fatal(s) })

	assertEq("val", []int{1, 3}, s, func(s string) { t.Fatal(s) })
}

func TestNewWriterWithFilterFnWithNilWriter(t *testing.T) {
	w := NewWriterWithFilterFn[int](nil)(func(v int) bool { return v%2 != 0 })

	assertEq("err", io.ErrClosedPipe, w.Write(nil, 1), func(s string) { t.Fatal(s) })
}

func TestNewWriterWithFilterFnWithNilFilter(t *testing.T) {
	s := make([]int, 0, 2)
	w := NewWriterWithFilterFn(newSliceWriter(&s))(nil)

	assertEq("err", *new(error), w.Write(nil, 1), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 2), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 3), func(s string) { t.Fatal(s) })

	assertEq("val", []int{1, 2, 3}, s, func(s string) { t.Fatal(s) })
}

func TestNewWriterWithMapperFnIdeal(t *testing.T) {
	s := make([]int, 0, 3)
	w := newSliceWriter(&s)
	w = NewWriterWithMapperFn[int](w)(func(v int) int { return v + 1 })

	assertEq("err", *new(error), w.Write(nil, 1), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 2), func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), w.Write(nil, 3), func(s string) { t.Fatal(s) })

	assertEq("val", []int{2, 3, 4}, s, func(s string) { t.Fatal(s) })
}

func TestNewWriterWithMapperFnWithNilWriter(t *testing.T) {
	w := NewWriterWithMapperFn[int, int](nil)(func(v int) int { return v + 1 })

	assertEq("err", io.ErrClosedPipe, w.Write(nil, 1), func(s string) { t.Fatal(s) })
}

func TestNewWriterWithMapperFnWithNilMapper(t *testing.T) {
	s := make([]int, 0, 3)
	w := newSliceWriter(&s)
	w = NewWriterWithMapperFn[int](w)(nil)

	assertEq("err", io.ErrClosedPipe, w.Write(nil, 1), func(s string) { t.Fatal(s) })
}
