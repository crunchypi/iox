package iox

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// -----------------------------------------------------------------------------
// New Writer iface + impl.
// -----------------------------------------------------------------------------

// Writer writes T, it is intended as a generic variant of io.Writer.
// Use io.ErrClosedPipe as a signal for when writing should stop.
type Writer[T any] interface {
	Write(context.Context, T) error
}

// WriterImpl lets you implement Writer with a function. Place it into "impl"
// and it will be called by the "Write" method.
//
// Example:
//
//	func myWriter() Writer[int] {
//	    return WriterImpl[int]{
//	        Impl: func(ctx context.Context, v int) error {
//	            // Your implementation.
//	        },
//	    }
//	}
type WriterImpl[T any] struct {
	Impl func(context.Context, T) error
}

// Write implements Writer by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.ErrClosedPipe will be returned.
func (impl WriterImpl[T]) Write(ctx context.Context, v T) (err error) {
	if impl.Impl == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.Impl(ctx, v)
}

// -----------------------------------------------------------------------------
// New WriteCloser iface + impl.
// -----------------------------------------------------------------------------

// WriteCloser groups Writer with io.Closer.
type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

// WriteCloserImpl lets you implement WriteCloser with functions. This is
// similar to WriterImpl but lets you implement io.Closer as well.
type WriteCloserImpl[T any] struct {
	ImplC func() error
	ImplW func(context.Context, T) error
}

// Close implements io.Closer by deferring to the internal ImplC func.
// If the internal ImplC func is nil, nothing will happen.
func (impl WriteCloserImpl[T]) Close() error {
	if impl.ImplC == nil {
		return nil
	}

	return impl.ImplC()
}

// Write implements Writer by deferring to the internal "ImplW" func.
// If the internal "ImplW" is not set, an io.ErrClosedPipe will be returned.
func (impl WriteCloserImpl[T]) Write(ctx context.Context, v T) (err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

// NewWriterFromValues creates a Writer (vals) which writes into 'w'.
// Nil 'w' returns an empty non-nil Writer; nil 'f' uses json.NewEncoder.
//
// Example (interactive):
//   - https://go.dev/play/p/5arKiC4ZxRt
//
// Example:
//
//	// Defining our io.Writer to rcv the data + encoding method.
//	b := bytes.NewBuffer(nil)
//	f := func(w io.Writer) Encoder { return json.NewEncoder(w) }
//	w := NewWriterFromValues[int](b)(f)
//
//	// Write values, they are encoded and passed to 'b'. Err handling ignored.
//	w.Write(nil, 2)
//
//	// We'll use these to read what's in 'b'.
//	dec := json.NewDecoder(b)
//	val := 0
//
//	t.Log(dec.Decode(&val), val) // <nil> 2
//	t.Log(dec.Decode(&val), val) // EOF 2
func NewWriterFromValues[T any](w io.Writer) func(f encoderFn) Writer[T] {
	return func(f func(io.Writer) Encoder) Writer[T] {
		if w == nil {
			return WriterImpl[T]{}
		}

		b := bytes.NewBuffer(nil)
		e := func(w io.Writer) Encoder { return json.NewEncoder(w) }(b)

		if f != nil {
			if _e := f(b); _e != nil {
				e = _e
			}
		}

		return WriterImpl[T]{
			Impl: func(ctx context.Context, v T) error {
				err := e.Encode(v)
				if err != nil {
					return err
				}

				_, err = b.WriteTo(w)
				return err
			},
		}
	}
}

// NewWriterFromBytes creates an io.Writer (bytes) which writes into 'w'.
// Nil 'w' returns an empty non-nil Writer; nil 'f' uses json.NewDecoder.
//
// Example (interactive):
//   - https://go.dev/play/p/yhaEWVIMoxw
//
// Example:
//
//	// Writes simply logs values.
//	vw := WriterImpl[int]{
//		Impl: func(ctx context.Context, v int) error {
//			t.Log(v)
//			return nil
//		},
//	}
//
//	// io.Writer
//	bw := NewWriterFromBytes(vw)(
//		func(r io.Reader) Decoder {
//			return json.NewDecoder(r)
//		},
//	)
//
//	// Logs "9"
//	json.NewEncoder(bw).Encode(9)
func NewWriterFromBytes[T any](w Writer[T]) func(f decoderFn) io.Writer {
	return func(f decoderFn) io.Writer {
		if w == nil {
			return readWriteCloserImpl{}
		}

		b := bytes.NewBuffer(nil)
		d := func(r io.Reader) Decoder { return json.NewDecoder(r) }(b)

		if f != nil {
			if _d := f(b); _d != nil {
				d = _d
			}
		}

		return readWriteCloserImpl{
			ImplW: func(p []byte) (n int, err error) {
				n, err = b.Write(p)
				if err != nil {
					return
				}

				var v T
				err = d.Decode(&v)

				if err != nil {
					return
				}

				err = w.Write(nil, v)
				if err != nil {
					return
				}

				return
			},
		}
	}
}

// -----------------------------------------------------------------------------
// Modifiers.
// -----------------------------------------------------------------------------

// NewWriterWithBatching returns a Writer which writes into a buffer of a given
// size. When the buffer is full, it is written into 'w'. Note that this should
// be used with caution due to the internal buffer, as there may be value loss
// if the process exits before the buffer is filled and written to 'w', e.g
// if 'size' is 10 but the process exits after only writing 9 times.
//
// Example (interactive):
//   - https://go.dev/play/p/0O4QR_en9h1
//
// Example:
//
//	// Writes which logs values through 't.Log'.
//	logWriter := WriterImpl[[]int]{}
//	logWriter.Impl = func(_ context.Context, v []int) error { t.Log(v); return nil }
//
//	w := NewWriterWithBatching(logWriter, 2)
//	w.Write(nil, 1)
//	w.Write(nil, 2) // Logger logs: '[1, 2]'
//	w.Write(nil, 3)
func NewWriterWithBatching[T any](w Writer[[]T], size int) Writer[T] {
	if w == nil {
		return WriterImpl[T]{}

	}

	buf := make([]T, 0, size)
	return WriterImpl[T]{
		Impl: func(ctx context.Context, val T) (err error) {
			buf = append(buf, val)

			if len(buf) >= size {
				err = w.Write(ctx, buf)
				buf = make([]T, 0, size)
			}

			return err
		},
	}
}

// NewWriterWithUnbatching returns a Writer which accepts []T on a Write call,
// then iterates through the slice and writes each value to 'w'.
//
// Example (interactive):
//   - https://go.dev/play/p/Z31KN0C2Q-Z
//
// Example:
//
//	// Writes which logs values through 't.Log'.
//	logWriter := WriterImpl[int]{}
//	logWriter.Impl = func(_ context.Context, v int) error { t.Log(v); return nil }
//
//	w := NewWriterWithUnbatching(logWriter)
//	w.Write(nil, []int{1, 2})
//	// ^ logWriter logs the following lines:
//	//  1
//	//  2
func NewWriterWithUnbatching[T any](w Writer[T]) Writer[[]T] {
	if w == nil {
		return WriterImpl[[]T]{}
	}

	return WriterImpl[[]T]{
		Impl: func(ctx context.Context, vs []T) (err error) {
			for _, v := range vs {
				err = w.Write(ctx, v)
				if err != nil {
					return
				}
			}

			return
		},
	}
}

// NewWriterWithFilterFn returns a writer which writes values into 'w', except
// those filtered by 'f'. Nil 'w' returns an empty Writer; nil 'f' returns 'w'.
//
// Example (interactive):
//   - https://go.dev/play/p/LM-XNzSmSNV
//
// Example:
//
//	// Writes which logs values through 't.Log'.
//	logWriter := WriterImpl[int]{}
//	logWriter.Impl = func(_ context.Context, v int) error { t.Log(v); return nil }
//
//	w := NewWriterWithFilterFn(logWriter)(
//		func(v int) bool {
//			return v > 1
//		},
//	)
//
//	w.Write(nil, 1) // Logs: nothing
//	w.Write(nil, 2) // Logs: 2
//	w.Write(nil, 3) // Logs: 3
func NewWriterWithFilterFn[T any](w Writer[T]) func(f func(T) bool) Writer[T] {
	return func(f func(T) bool) Writer[T] {
		if w == nil {
			return WriterImpl[T]{}
		}
		if f == nil {
			return w
		}

		return WriterImpl[T]{
			Impl: func(ctx context.Context, v T) error {
				if !f(v) {
					return nil
				}

				return w.Write(ctx, v)
			},
		}
	}
}
