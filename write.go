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

// WriterImpl implements Writer with its Write method by deferring to 'Impl'.
// This is for convenience, as you may use a functional implementation of Writer
// without defining a new type (that's done for you here).
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

// WriteCloserImpl implements Writer and io.Closer with its methods by deferring
// to ImplC (closer) and ImplW (writer). This is for convenience, as you may use
// a functional implementation of the interfaces without defining a new type.
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

// NewWriterFromValues returns a Writer which accepts values, encodes them
// using the given encoder, and then writes them to 'w'. If 'w' is nil, an empty
// Writer is returned; if 'f' is nil, the encoder is set to json.NewEncoder.
// Example:
//
//	// Defining our io.Writer to rcv the data + encoding method.
//	b := bytes.NewBuffer(nil)
//	f := func(w io.Writer) Encoder { return json.NewEncoder(w) }
//	w := NewWriterFromValues[int](b)(f)
//
//	// Write values, they are encoded and passed to 'b'. Err handling ignored.
//	w.Write(nil, 2)
//	w.Write(nil, 3)
//
//	// We'll use these to read what's in 'b'.
//	dec := json.NewDecoder(b)
//	val := 0
//
//	t.Log(dec.Decode(&val), val) // <nil> 2
//	t.Log(dec.Decode(&val), val) // <nil> 3
//	t.Log(dec.Decode(&val), val) // EOF 3
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
