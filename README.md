# iox
Generic variant of Go's io pkg

Index 
- [Core interfaces](#core-interfaces)
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