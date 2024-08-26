# iox
Generic extension of Go's io pkg.

Index 
- [Errors](#errors)
- [Core interfaces](#core-interfaces)
- [Constructors](#constructors)
- [Modifiers](#modifiers)



## Errors
<details>
<summary> Expand/collapse section </summary>

This package does *not* define any new errors, it inherits them from the `io` package in the standard library.
```go
io.EOF              // Used by e.g iox.Reader: Stop reading/consuming
io.ErrClosedPipe    // Used by e.g iox.Writer: Stop writing/producing.
```

</details>



## Core interfaces
Core interfaces are the `iox.Reader` and `iox.Writer`listed below. They mirror `io.Reader` and `io.Writer`but differ in that they work with generic values instead. An extra addition is the use of `context.Context`, since io often involves program bounds (also it gives some added flexibility).

```go
type Reader[T any] interface {
	Read(context.Context) (T, error)
}
```

```go
type Writer[T any] interface {
	Write(context.Context, T) error
}
```

<details>
<summary> As with the io package from the standard library, iox readers and writers are combined with eachother and io.Closer. The full set of combinations can be seen by clicking here </summary>

```go
type Reader[T any] interface {
	Read(context.Context) (T, error)
}

type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

type Writer[T any] interface {
	Write(context.Context, T) error
}

type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}

type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
}
```
</details>

<br>
<details>
<summary>  There are also "impl" structs which let you implement core interfaces with functions, which reduces a lot of boilerplate. These can be seen by clicking on this section </summary>

<br>

Signatures are links to the Go playground (examples).
- [`type ReaderImpl[T any] struct`](https://go.dev/play/p/gkzrDGzLRtc)
- [`type ReadCloserImpl[T any] struct`](https://go.dev/play/p/SXA7OWQl5ee)
- [`type WriterImpl[T any] struct`](https://go.dev/play/p/796B8udkJKy)
- [`type WriteCloserImpl[T any] struct`](https://go.dev/play/p/UE0Bxls3D5D)
- [`type ReadWriterImpl[T, U any] struct`](https://go.dev/play/p/yl_e7ics0oY)
- [`type ReadWriteCloserImpl[T, U any] struct`](https://go.dev/play/p/RvmasSrtNo_c)

</details>



## Constructors

All links go to examples on the Go playground.

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



## Modifiers
All links go to examples on the Go playground.

Batching.
- [`func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T]`](
	https://go.dev/play/p/Mn3Cipq8-Gy
)
- [`func NewReaderWithUnbatching[T any](r Reader[[]T]) Reader[T]`](
	https://go.dev/play/p/zaLBILUnkgE
)
- [`func NewWriterWithBatching[T any](w Writer[[]T], size int) Writer[T]`](
	https://go.dev/play/p/sbOaajf3Jt8
)
- [`func NewWriterWithUnbatching[T any](w Writer[T]) Writer[[]T]`](
	https://go.dev/play/p/E-qP0CE8wV3
)

Filtering & mapping.
* [`func NewReaderWithFilterFn[T any](r Reader[T]) func(f func(v T) bool) Reader[T]`](
	https://go.dev/play/p/vYCJChGUKF_Y
)
* [`func NewReaderWithMapperFn[T, U any](r Reader[T]) func(f func(T) U) Reader[U]`](
	https://go.dev/play/p/CaB0N1N5nur
)
* [`func NewWriterWithFilterFn[T any](w Writer[T]) func(f func(T) bool) Writer[T]`](
	https://go.dev/play/p/BgKAgGVvJ7b
)
* [`func NewWriterWithMapperFn[T, U any](w Writer[U]) func(f func(T) U) Writer[T]`](
	https://go.dev/play/p/V3OvYkJS-mC
)