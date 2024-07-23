package iox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// -----------------------------------------------------------------------------
// New Reader iface + impl.
// -----------------------------------------------------------------------------

// Reader reads T, it is intended as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}

// ReaderImpl lets you implement Reader with a function. Place it into "impl"
// and it will be called by the "Read" method.
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Read implements Reader by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx)
}

// -----------------------------------------------------------------------------
// New ReadCloser iface + impl.
// -----------------------------------------------------------------------------

// ReadCloser groups Reader with io.Closer.
type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

// ReadCloserImpl lets you implement ReadCloser with functions. This is similar
// to ReaderImpl but lets you implement io.Closer as well.
type ReadCloserImpl[T any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
}

// Read implements Closer by deferring to the internal "ImplC" func.
// If the internal "ImplC" func is nil, nothing will happen.
func (impl ReadCloserImpl[T]) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

// Read implements Reader by deferring to the internal "ImplR" func.
// If the internal "ImplR" is not set, an io.EOF will be returned.
func (impl ReadCloserImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

// NewReaderFrom returns a Reader which yields values from the given vals.
func NewReaderFrom[T any](vs ...T) Reader[T] {
	i := 0
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			if i >= len(vs) {
				return val, io.EOF
			}

			val = vs[i]
			i++
			return
		},
	}
}

// NewReaderFromBytes converts an io.Reader (bytes) into a iox.Reader (values).
// Nil 'r' returns an empty non-nil Reader; nil 'f' uses json.NewDecoder.
// Example:
//
//	// Used as io.Reader
//	b := bytes.NewBuffer(nil)
//
//	// Using json encoder, so the decoder has to be json in NewReaderFromBytes
//	json.NewEncoder(b).Encode("test1")
//	json.NewEncoder(b).Encode("test2")
//
//	r := NewReaderFromBytes[string](b)(
//		func(r io.Reader) Decoder {
//			return json.NewDecoder(r)
//		},
//	)
//
//	t.Log(r.Read(context.Background())) // "test1" <nil>
//	t.Log(r.Read(context.Background())) // "test2" <nil>
//	t.Log(r.Read(context.Background())) // "", io.EOF
func NewReaderFromBytes[T any](r io.Reader) func(f decoderFn) Reader[T] {
	return func(f func(io.Reader) Decoder) Reader[T] {
		if r == nil {
			return ReaderImpl[T]{}
		}

		var d Decoder = json.NewDecoder(r)
		if f != nil {
			if _d := f(r); _d != nil {
				d = _d
			}
		}

		return ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, err error) {
				err = d.Decode(&v)
				return
			},
		}
	}
}

// NewReaderFromValues converts an iox.Reader (values) into an io.Reader (bytes).
// Nil 'r' returns an empty non-nil Reader; nil 'f' uses json.NewEncoder.
// Example:
//
//	// Create the io.Reader from value Reader.
//	r := NewReaderFromValues(NewReaderFrom("test1", "test2"))(
//	    func(w io.Writer) Encoder {
//	        return json.NewEncoder(w)
//	    },
//	)
//
//	// Instantly pass it to a decoder just so we may log out the values.
//	dec := json.NewDecoder(r)
//	val := ""
//
//	t.Log(dec.Decode(&val), val) // <nil>, "test1"
//	t.Log(dec.Decode(&val), val) // <nil>, "test2"
//	t.Log(dec.Decode(&val), val) // EOF, "test2" <--- val is unchanged.
func NewReaderFromValues[T any](r Reader[T]) func(f encoderFn) io.Reader {
	return func(f func(io.Writer) Encoder) io.Reader {
		if r == nil {
			return readWriteCloserImpl{}
		}

		b := bytes.NewBuffer(nil)
		e := Encoder(json.NewEncoder(b))
		if f != nil {
			if _e := f(b); _e != nil {
				e = _e
			}
		}

		return readWriteCloserImpl{
			ImplR: func(p []byte) (n int, err error) {
				v, err := r.Read(context.Background())
				if err != nil {
					return 0, err
				}

				err = e.Encode(v)
				if err != nil {
					return 0, err
				}

				return b.Read(p)
			},
		}
	}
}
