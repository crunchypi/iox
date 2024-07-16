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



## Errors
This package inherits errors from the `io` package and only uses:
- `io.EOF`: Used with `iox.Reader[T]` and `iox.Decoder`
- `io.ErrClosedPipe`: 



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