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
//
// Example (interactive):
//   - https://go.dev/play/p/gkzrDGzLRtc
//
// Example:
//
//	func myReader() Reader[int] {
//	    return ReaderImpl[int]{
//	        Impl: func(ctx context.Context) (int, error) {
//	            // Your implementation.
//	        },
//	    }
//	}
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
//
// Example (interactive):
//   - https://go.dev/play/p/SXA7OWQl5ee
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
//
// Example (interactive):
//   - https://go.dev/play/p/bP73PU1mQvf
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
//
// Example (interactive):
//   - https://go.dev/play/p/ltcwrgk41Gw
//
// Example:
//
//	// Used as io.Reader
//	b := bytes.NewBuffer(nil)
//
//	// Using json encoder, so the decoder has to be json in NewReaderFromBytes
//	json.NewEncoder(b).Encode("test1")
//
//	r := NewReaderFromBytes[string](b)(
//		func(r io.Reader) Decoder {
//			return json.NewDecoder(r)
//		},
//	)
//
//	t.Log(r.Read(context.Background())) // "test1" <nil>
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
//
// Example (interactive):
//   - https://go.dev/play/p/e9Sp5od3iE6
//
// Example:
//
//	// Create the io.Reader from value Reader.
//	r := NewReaderFromValues(NewReaderFrom("test1"))(
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
//	t.Log(dec.Decode(&val), val) // EOF, "test1" <--- val is unchanged.
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

// -----------------------------------------------------------------------------
// Modifiers.
// -----------------------------------------------------------------------------

// NewReaderWithBatching returns a reader which batches 'r' into slices with
// the given size. Nil 'r' returns an empty non-nil Reader, size <= 0 defaults
// to 8. Note, the last []T before an err (e.g io.EOF) may be smaller than 'size'.
//
// Example (interactive):
//   - https://go.dev/play/p/Mn3Cipq8-Gy
//
// Example:
//
//	vr := NewReadWriterFrom(1,2,3)
//	sr := NewReaderWithBatching(vr, 2)
//
//	t.Log(sr.Read(nil)) // [1, 2], nil
//	t.Log(sr.Read(nil)) // [3], nil
//	t.Log(sr.Read(nil)) // [], io.EOF
func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T] {
	if r == nil {
		return ReaderImpl[[]T]{}
	}

	if size <= 0 {
		size = 8
	}

	var errCache error
	return ReaderImpl[[]T]{
		Impl: func(ctx context.Context) (s []T, err error) {
			s = make([]T, 0, size)
			if errCache != nil {
				return s, errCache
			}

			var v T
			for i := 0; i < size; i++ {
				v, errCache = r.Read(ctx)
				if errCache != nil {
					break

				}

				s = append(s, v)
			}

			if errCache != nil && len(s) == 0 {
				return s, errCache
			}

			return s, err
		},
	}
}

// NewReaderWithUnbatching returns a reader of T from a reader of []T.
// Note that there is some internal buffering, so you may want to use this
// with caution as an unread buffer may cause value loss.
//
// Example (interactive):
//   - https://go.dev/play/p/zaLBILUnkgE
//
// Example:
//
//	sr := NewReaderFrom([]int{1, 2}, []int{3})
//	vr := NewReaderWithUnbatching(sr)
//
//	t.Log(vr.Read(nil)) // 1, nil
//	t.Log(vr.Read(nil)) // 2, nil
//	t.Log(vr.Read(nil)) // 3, nil
//	t.Log(vr.Read(nil)) // 0, io.EOF
func NewReaderWithUnbatching[T any](r Reader[[]T]) Reader[T] {
	if r == nil {
		return ReaderImpl[T]{}
	}

	var errCache error
	var buf []T
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			if len(buf) > 0 {
				val = buf[0]
				buf = buf[1:]
				return
			}

			if errCache != nil {
				err = errCache
				return
			}

			buf, err = r.Read(ctx)

			switch {
			case len(buf) == 0 && err != nil:
				return val, err
			case len(buf) == 0 && err == nil:
				return val, io.EOF
			case len(buf) != 0 && err != nil:
				errCache = err
				err = nil
			case len(buf) != 0 && err == nil:
			}

			val = buf[0]
			buf = buf[1:]
			return
		},
	}
}

// NewReaderWithFilterFn returns a reader of values from 'r', except for those
// filtered by 'f'. Nil 'r' returns an empty non-nil Reader; nil 'f' returns 'r'.
//
// Example (interactive):
//   - https://go.dev/play/p/vYCJChGUKF_Y
//
// Example:
//
//	r := NewReaderFrom(1, 2, 3)
//	r = NewReaderWithFilterFn(r)(
//		func(v int) bool {
//			return v > 1
//		},
//	)
//
//	t.Log(r.Read(nil)) // 2, nil
//	t.Log(r.Read(nil)) // 3, nil
//	t.Log(r.Read(nil)) // 0, io.EOF
func NewReaderWithFilterFn[T any](r Reader[T]) func(f func(v T) bool) Reader[T] {
	return func(f func(v T) bool) Reader[T] {
		if r == nil {
			return ReaderImpl[T]{}
		}
		if f == nil {
			return r
		}

		return ReaderImpl[T]{
			Impl: func(ctx context.Context) (val T, err error) {
				for val, err = r.Read(ctx); err == nil; val, err = r.Read(ctx) {
					if f(val) {
						return
					}
				}

				return
			},
		}
	}
}

// NewReaderWithMapperFn returns a reader of mapped values from 'r'.
// An empty non-nil Reader is returned if either 'r' or 'f' is nil.
//
// Example (interactive):
//   - https://go.dev/play/p/CaB0N1N5nur
//
// Example:
//
//	ri := NewReaderFrom(1, 2)
//	rs := NewReaderWithMapperFn[int, string](ri)(
//	    func(v int) string {
//	        return fmt.Sprint(v)
//	    },
//	)
//
//	t.Log(rs.Read(nil)) // "1", nil
//	t.Log(rs.Read(nil)) // "2", nil
//	t.Log(rs.Read(nil)) // "", io.EOF
func NewReaderWithMapperFn[T, U any](r Reader[T]) func(f func(T) U) Reader[U] {
	return func(f func(T) U) Reader[U] {
		if r == nil || f == nil {
			return ReaderImpl[U]{}
		}

		return ReaderImpl[U]{
			Impl: func(ctx context.Context) (valOut U, err error) {
				valIn, err := r.Read(ctx)
				if err != nil {
					return valOut, err
				}

				return f(valIn), err
			},
		}
	}
}
