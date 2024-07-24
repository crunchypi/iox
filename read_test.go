package iox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
)

// -----------------------------------------------------------------------------
// Reader impl.
// -----------------------------------------------------------------------------

func TestReaderImplReadIdeal(t *testing.T) {
	r := ReaderImpl[int]{}
	r.Impl = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReaderImplReadWithNilImpl(t *testing.T) {
	r := ReaderImpl[int]{}

	val, err := r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// ReadCloser impl.
// -----------------------------------------------------------------------------

func TestReadCloserImplReadIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplR = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := rc.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplReadWithNilImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	val, err := rc.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplC = func() error { return nil }

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseWithNilImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

func TestNewReaderFromIdeal(t *testing.T) {
	r := NewReaderFrom(1, 2)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesIdeal(t *testing.T) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode("test1")
	json.NewEncoder(b).Encode("test2")

	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	r := NewReaderFromBytes[string](b)(f)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesWithNilReader(t *testing.T) {
	r := NewReaderFromBytes[string](nil)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesWithNilDecoder(t *testing.T) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode("test1")
	json.NewEncoder(b).Encode("test2")

	r := NewReaderFromBytes[string](b)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromValuesIdeal(t *testing.T) {
	fn := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	br := NewReaderFromValues(NewReaderFrom("test1", "test2"))(fn)

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromValuesWithNilReader(t *testing.T) {
	fn := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	br := NewReaderFromValues[int](nil)(fn)

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromValuesWithNilEncoder(t *testing.T) {
	vr := NewReaderFrom("test1", "test2")
	br := NewReaderFromValues(vr)(nil)

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromValuesWithEncodeError(t *testing.T) {
	fn := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	br := NewReaderFromValues(NewReaderFrom(make(chan int)))(fn)

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)

	want := "json: unsupported type: chan int"
	have := err.Error()
	assertEq("err", want, have, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Modifiers.
// -----------------------------------------------------------------------------

func TestNewReaderWithBatchingIdeal(t *testing.T) {
	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	vr := NewReaderFrom(vs...)
	sr := NewReaderWithBatching(vr, 0)

	s := []int{}
	err := *new(error)

	s, err = sr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", vs[0:8], s, func(s string) { t.Fatal(s) })

	s, err = sr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", vs[8:], s, func(s string) { t.Fatal(s) })

	s, err = sr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", []int{}, s, func(s string) { t.Fatal(s) })
}

func TestNewReaderWithBatchingWithNilReader(t *testing.T) {
	sr := NewReaderWithBatching[int](nil, 0)

	s, err := sr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", *new([]int), s, func(s string) { t.Fatal(s) })
}

