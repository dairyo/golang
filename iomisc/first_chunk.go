package iomisc

import (
	"errors"
	"io"
)

var _ io.Writer = (*FirstChunk[io.Writer])(nil)

var (
	ErrFailToWriteAll = errors.New("fail to write all")
)

func NewFirstChunk[T io.Writer](base T, limit int) *FirstChunk[T] {
	return &FirstChunk[T]{
		WroteSize: *NewWroteSize(base),
		Limit:     limit,
	}
}

type FirstChunk[T io.Writer] struct {
	WroteSize[T]
	Limit  int
	IsOver bool
	Err    error
}

func (f *FirstChunk[T]) Write(p []byte) (int, error) {
	l := len(p)
	if f.IsOver {
		return l, nil
	}
	if f.Err != nil {
		return l, nil
	}
	rest := f.Limit - f.WroteSize.Size
	if rest <= 0 {
		f.IsOver = true
		return l, nil
	}
	if rest <= len(p) {
		if rest < len(p) {
			f.IsOver = true
		}
		p = p[:rest]
	}
	n, err := f.WroteSize.Write(p)
	if err != nil {
		f.Err = err
		return l, nil
	}
	if n != len(p) {
		f.Err = ErrFailToWriteAll
		// if truncated, always returns original p size.
		return l, nil
	}
	return l, nil
}
