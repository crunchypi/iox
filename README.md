# iox
Generic variant of Go's io pkg

Index 
- [Core interfaces](#core-interfaces)
- [Constructors/Factories](#constructorsfactories)
- [Errors](errors)
- [Impl pattern](#impl-pattern)



## Core interfaces

#### Reader
```go
// Reader reads T, it is intended as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}
```

#### ReadCloser
```go
// ReadCloser groups Reader with io.Closer.
type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}
```

#### Writer
```go
// Writer writes T, it is intended as a generic variant of io.Writer.
type Writer[T any] interface {
	Write(context.Context, T) error
}
```

#### WriteCloser
```go
// WriteCloser groups Writer with io.Closer.
type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}
```

#### ReadWriter
```go
// ReadWriter groups Reader[T] and Writer[U].
type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}
```

#### ReadWriteCloser
```go
// ReadWriteCloser groups Reader[T] and Writer[U] with io.Closer.
type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
}
```



## Constructors/Factories


```go
// NewReaderFrom returns a Reader which yields values from the given vals.
func NewReaderFrom[T any](vs ...T) Reader[T]
```

```go
// NewReaderFromBytes creates a new T reader from an io.Reader and Decoder.
// It simply reads bytes from 'r', decodes them, and passes them along to the
// caller. As such, the decoder must match the encoder used to create the bytes.
// If 'r' is nil, an empty Reader is returned; if 'f' is nil, the decoder is set
// to json.NewDecoder.
func NewReaderFromBytes[T any](r io.Reader) func(f decoderFn) Reader[T]
```

```go
// NewReaderFromValues creates an io.Reader from a Reader and Encoder.
// It simply reads values from 'r', encodes them, and passes them along to the
// caller. As such, when decoding values from the returned io.Reader one should
// use a decoder which matches the encoder passed here. If 'r' is nil, an
// empty (not nil) io.Reader is returned; if 'f' is nil, the encoder is set to
// json.NewEncoder. 
func NewReaderFromValues[T any](r Reader[T]) func(f encoderFn) io.Reader
```

```go
// NewWriterFromValues returns a Writer which accepts values, encodes them
// using the given encoder, and then writes them to 'w'. If 'w' is nil, an empty
// Writer is returned; if 'f' is nil, the encoder is set to json.NewEncoder.
func NewWriterFromValues[T any](w io.Writer) func(f encoderFn) Writer[T]
```

```go
// NewReaderFromValues creates an io.Reader from a Reader and Encoder.
// It simply reads values from 'r', encodes them, and passes them along to the
// caller. As such, when decoding values from the returned io.Reader one should
// use a decoder which matches the encoder passed here. If 'r' is nil, an
// empty (not nil) io.Reader is returned; if 'f' is nil, the encoder is set to
// json.NewEncoder. 
func NewReaderFromValues[T any](r Reader[T]) func(f encoderFn) io.Reader
```



## Errors
This package inherits errors from the `io` package and only uses:
- `io.EOF`: Used with `iox.Reader[T]` and `iox.Decoder`
- `io.ErrClosedPipe`: Used with `iox.Writer[T]` and `iox.Encoder`



## Impl pattern
The impl pattern allows you to implement an interface in a functional way, avoiding the tedium of defining structs which implement small interfaces. You simply define the function and place it inside an impl struct.

#### Impl for Reader
```go
// ReaderImpl implements Reader with it's Read method by deferring to 'Impl'.
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Read implements Reader by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error)
```

#### Impl for ReadCloser
```go
// ReadCloserImpl implements Reader and io.Closer with its methods by deferring
// to ImplC (closer) and ImplR (reader). 
type ReadCloserImpl[T any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
}

// Read implements Closer by deferring to the internal "ImplC" func.
func (impl ReadCloserImpl[T]) Close() (err error)

// Read implements Reader by deferring to the internal "ImplR" func.
func (impl ReadCloserImpl[T]) Read(ctx context.Context) (r T, err error)
```

#### Impl for Writer.
```go
// WriterImpl implements Writer with its Write method by deferring to 'Impl'.
type WriterImpl[T any] struct {
	Impl func(context.Context, T) error
}

// Write implements Writer by deferring to the internal "Impl" func.
func (impl WriterImpl[T]) Write(ctx context.Context, v T) (err error) 
```

#### Impl for WriteCloser.
```go
// WriteCloserImpl implements Writer and io.Closer with its methods by deferring
// to ImplC (closer) and ImplW (writer). 
type WriteCloserImpl[T any] struct {
	ImplC func() error
	ImplW func(context.Context, T) error
}

// Close implements io.Closer by deferring to the internal ImplC func.
func (impl WriteCloserImpl[T]) Close() error

// Write implements Writer by deferring to the internal "ImplW" func.
func (impl WriteCloserImpl[T]) Write(ctx context.Context, v T) (err error)
```

```go
// ReadWriterImpl implements ReadWriter[T, U] with its Read and Write methods,
// their logic is deferred to the internal ImplR and ImplW fields (funcs).
type ReadWriterImpl[T, U any] struct {
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Read implements the Reader[T] part of ReadWriter[T, U] by calling ImplR.
func (impl ReadWriterImpl[T, U]) Read(ctx context.Context) (r T, err error)

// Write implements the Writer[U] part of ReadWriter[T, U] by calling ImplW.
func (impl ReadWriterImpl[T, U]) Write(ctx context.Context, v U) (err error)
```

#### Impl for ReadWriteCloser.
```go
// ReadWriteCloserImpl implements ReadWriteCloser with its methods Read, Write
// and Close, their logic is deferred to the internal ImplR, ImplW and ImplC.
type ReadWriteCloserImpl[T, U any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Close implements io.Close by deferring to the internal ImplC func.
func (impl ReadWriteCloserImpl[T, U]) Close() (err error)

// Read implements Reader[T] by deferring logic to the internal ImplR func.
func (impl ReadWriteCloserImpl[T, U]) Read(ctx context.Context) (r T, err error)

// Write implements Writer[U] by deferring logic to the internal ImplW func.
func (impl ReadWriteCloserImpl[T, U]) Write(ctx context.Context, v U) (err error)
```