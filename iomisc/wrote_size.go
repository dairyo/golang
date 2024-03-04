package iomisc

import "io"

var _ io.Writer = (*WroteSize[io.Writer])(nil)

// NewWroteSize is a constructor for WroteSize.
func NewWroteSize[T io.Writer](base T) *WroteSize[T] {
	return &WroteSize[T]{Base: base}
}

// WroteSize is a wrapper for io.Writer to memorize wrote size.
type WroteSize[T io.Writer] struct {
	// Base is an io.Writer to memorize wrote size.
	Base T

	// Size is memorized size. This must be use as read only.
	Size int
}

// Write write p to Base writer and memorize the wrote size.
func (l *WroteSize[T]) Write(p []byte) (int, error) {
	n, err := l.Base.Write(p)
	if err != nil {
		return 0, err
	}
	l.Size += n
	return n, nil
}
