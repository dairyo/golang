package iomisc

import (
	"bytes"
	"errors"
	"testing"
)

func TestWroteSizeWrite(t *testing.T) {
	w := NewWroteSize(&bytes.Buffer{})
	n, err := w.Write([]byte("foo"))

	if n != 3 {
		t.Fatalf("ret val must be 3 but %d", n)
	}
	if err != nil {
		t.Fatalf("must not return error but %s", err)
	}
	if w.Size != 3 {
		t.Fatalf("size must be 3 but %d", w.Size)
	}
	if w.Base.String() != "foo" {
		t.Fatalf(`contents must be "foo" but %q`, w.Base.String())
	}

	n, err = w.Write([]byte("bar"))
	if n != 3 {
		t.Fatalf("ret val must be 3 but %d", n)
	}
	if err != nil {
		t.Fatalf("must not return error but %s", err)
	}
	if w.Size != 6 {
		t.Fatalf("size must be 3 but %d", w.Size)
	}
	if w.Base.String() != "foobar" {
		t.Fatalf(`contents must be "foo" but %q`, w.Base.String())
	}
}

func TestWroteSizeError(t *testing.T) {
	want := errors.New("foo")
	w := NewWroteSize(&errorWriter{want})
	n, got := w.Write([]byte("bar"))
	if n != 0 {
		t.Fatalf("ret val must be 0 but %d", n)
	}
	if got != want {
		t.Fatalf(`error mismatch: want=%s, got=%s`, want, got)
	}
	if w.Size != 0 {
		t.Fatalf("size must be 0 but %d", w.Size)
	}
}
