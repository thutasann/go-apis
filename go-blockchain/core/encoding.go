package core

import "io"

// Encoder[T] is a generic interface for types that can encode values of
// type T to an io.Writer. Implementations should write a deterministic
// binary or textual representation of v to w and return any error
// encountered during the write.
type Encoder[T any] interface {
	Encode(io.Writer, T) error
}

// Decoder[T] is a generic interface for types that can decode values of
// type T from an io.Reader. The implementation should read from r and
// populate the provided value (or return it) and return any error
// encountered during decoding.
type Decoder[T any] interface {
	Decode(io.Reader, T) error
}
