// iox is a generic extension of package io. See github for an overview:
// * https://github.com/crunchypi/iox
//
// iox is designed to be as simple as possible and does the following:
//   - Covers the most essential interfaces, Reader, Writer and variations.
//   - Defines interfaces for encoding and decoding.
//   - Inherits errors, namely io.EOF and io.ErrClosedPipe
//   - Defines converters for interoperability with io.
package iox

import "io"

// -----------------------------------------------------------------------------
// Encoder.
// -----------------------------------------------------------------------------

// Encoder encodes values into binary form. Some commonly used encoders are:
//   - json.NewEncoder(bytes.NewBuffer(nil))
//   - gob.NewEncoder(bytes.NewBuffer(nil))
type Encoder interface {
	Encode(e any) error
}

// EncoderImpl lets you implement Encoder with a function. Place it into "Impl"
// and it will be called by the "Encode" method.
type EncoderImpl struct {
	Impl func(e any) error
}

// Encode implements Encoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.ErrClosedPipe will be returned.
func (impl EncoderImpl) Encode(e any) error {
	if impl.Impl == nil {
		return io.ErrClosedPipe
	}

	return impl.Impl(e)
}

// -----------------------------------------------------------------------------
// Decoder.
// -----------------------------------------------------------------------------

// Decoder decodes values from binary form. Some commonly used encoders are:
//   - json.NewDecoder(bytes.NewBuffer(nil))
//   - gob.NewDecoder(bytes.NewBuffer(nil))
type Decoder interface {
	Decode(e any) error
}

// DecoderImpl lets you implement Decoder with a function. Place it into "Impl"
// and it will be called by the "Decode" method.
type DecoderImpl struct {
	Impl func(d any) error
}

// Decode implements Decoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl DecoderImpl) Decode(d any) error {
	if impl.Impl == nil {
		return io.EOF
	}

	return impl.Impl(d)
}

// -----------------------------------------------------------------------------
// Implementation io.Reader, io.Writer, io.ReadWriter and closer variants.
// -----------------------------------------------------------------------------

type readWriteCloserImpl struct {
	ImplC func() error
	ImplR func([]byte) (int, error)
	ImplW func([]byte) (int, error)
}

func (impl readWriteCloserImpl) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

func (impl readWriteCloserImpl) Read(p []byte) (n int, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(p)
}

func (impl readWriteCloserImpl) Write(p []byte) (n int, err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(p)
}

// -----------------------------------------------------------------------------
// Abbreviations.
// -----------------------------------------------------------------------------

type decoderFn = func(io.Reader) Decoder
type encoderFn = func(io.Writer) Encoder
