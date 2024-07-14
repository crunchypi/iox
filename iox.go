package iox

import "io"

// Encoder encodes values into binary form. Some commonly used encoders are:
//   - json.NewEncoder(bytes.NewBuffer(nil))
//   - gob.NewEncoder(bytes.NewBuffer(nil))
type Encoder interface {
	Encode(e any) error
}

// EncoderImpl implements Encoder with it's Encode method by deferring to 'Impl'.
// This is for convenience, as you may use functional implementation of Encoder
// without defining a new type (that's done for you here).
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
