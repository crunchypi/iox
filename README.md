# iox
Generic extension of Go's io pkg.

Index 
- [Core interfaces](#core-interfaces)
- [Errors](#errors)
- [Constructors/Factories](#constructorsfactories)
- [Impl pattern](#impl-pattern)
- [Modifiers/Wrappers](#constructorsfactories)


## Core interfaces
Core interfaces are the `iox.Reader` and `iox.Writer`listed below. They mirror `io.Reader` and `io.Writer`but differ in that they work with generic values instead. An extra addition is the use of `context.Context`, since io often involves program bounds (also it gives some added flexibility).

#### Reader
```go
// Reader reads T, it is intended as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}
```

#### Writer
```go
// Writer writes T, it is intended as a generic variant of io.Writer.
type Writer[T any] interface {
	Write(context.Context, T) error
}
```

#### Permutations
As with the `io` package from the standard library, `iox` readers and writers can be combined with eachother and `io.Closer`. The full set of interfaces can be viewed by clicking belod.

<details>
<summary> Show all interfaces </summary>

```go
// Reader reads T, it is intended as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}

// ReadCloser groups Reader with io.Closer.
type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

// Writer writes T, it is intended as a generic variant of io.Writer.
type Writer[T any] interface {
	Write(context.Context, T) error
}

// WriteCloser groups Writer with io.Closer.
type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

// ReadWriter groups Reader[T] and Writer[U].
type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}

// ReadWriteCloser groups Reader[T] and Writer[U] with io.Closer.
type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
}
```
</details>



## Errors
This package does *not* define any new errors, it inherits them from the `io` package in the standard library.
```go
var io.EOF              // Used by iox.Reader and iox.Decoder
var io.ErrClosedPipe    // Used by iox.Writer and iox.Encoder
```



## Constructors/Factories

Here's an overview, all links go to the Go playground.

- [`func NewReaderFrom[T any](vs ...T) Reader[T]`](https://go.dev/play/p/bP73PU1mQvf)
- [`func NewReaderFromBytes[T any](r io.Reader) func(f decoderFn) Reader[T]`](https://go.dev/play/p/ltcwrgk41Gw)
- [`func NewReaderFromValues[T any](r Reader[T]) func(f encoderFn) io.Reader`](https://go.dev/play/p/e9Sp5od3iE6)
- [`func NewWriterFromValues[T any](w io.Writer) func(f encoderFn) Writer[T]`](https://go.dev/play/p/5arKiC4ZxRt)
- [`func NewWriterFromBytes[T any](w Writer[T]) func(f decoderFn) io.Writer`](https://go.dev/play/p/yhaEWVIMoxw)
- [`func NewReadWriterFrom[T any](vs ...T) ReadWriter[T, T]`](https://go.dev/play/p/tusGzivubiI)

<details>
<summary> Alternatively, you may see signatures and docs by clicking here</summary>


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
// NewWriterFromBytes returns an io.Writer which accepts bytes, decodes them
// using the given decoder, and then writes them to 'w'. If 'w' is nil, an emtpy
// io.Writer is returned; if 'f' is nil, the decoder is set to json.NewDecoder.
func NewWriterFromBytes[T any](w Writer[T]) func(f decoderFn) io.Writer 
```

```go
// NewReadWriterFrom returns a ReadWriter[T] which writes into- and read from
// an internal buffer. The buffer is initially populated with the given values.
// The buffer acts like a stack, and a read while the buf is empty returns io.EOF.
func NewReadWriterFrom[T any](vs ...T) ReadWriter[T, T]
```
</details>


## Impl pattern
The impl pattern allows you to implement an interface in a functional way, avoiding the tedium of defining structs which implement small interfaces. You simply define the function and place it inside an impl struct. There is an impl struct for all [Core interfaces](#core-interfaces), but I'll show the one associated with `iox.Reader` to make it clear:

```go
// Here's how it may be used to e.g implement a Reader mapper:
//	https://go.dev/play/p/JQY_1vQZ6pz.
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx)
}
```

With this pattern you may easily define e.g a mapper func for a `iox.Reader`, e.g [this go playground](https://go.dev/play/p/JQY_1vQZ6pz)

<details>
<summary>As mentioned, an impl struct exists for all core interfaces, you may see them by clicking here</summary>

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
</details>



## Modifiers/Wrappers

Some helpers are defined for convenience. Their signature are listed below (all links go to the Go playground).

- [`func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T]`](
	https://go.dev/play/p/SnGdMkV9PNE
)