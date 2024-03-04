package iomisc

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}

type wantStateBase struct {
	n      int
	size   int
	isOver bool
	err    error
}

func checkBase[T io.Writer](t *testing.T, f *FirstChunk[T], p []byte, want *wantStateBase) {
	t.Helper()
	n, err := f.Write(p)
	if err != nil {
		t.Fatalf("must not return error but %s", err.Error())
	}
	if n != want.n {
		t.Fatalf("return value mismatch: want=%d, got=%d", want.n, n)
	}
	if f.Size != want.size {
		t.Fatalf("f.Size mismatch: want=%d, got=%d", want.size, f.Size)
	}
	if f.IsOver != want.isOver {
		t.Fatalf("f.IsOver mismatch: want=%t, got=%t", want.isOver, f.IsOver)
	}
	if f.Err != want.err {
		t.Fatalf("f.Err mismatch: want=%s, got=%s", want.err, f.Err)
	}

}

func TestFirstChunk(t *testing.T) {
	type contentState struct {
		wantStateBase
		content []byte
	}

	equalBytes := func(a []byte, b []byte) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	check := func(t *testing.T, f *FirstChunk[*bytes.Buffer], p []byte, want *contentState) {
		t.Helper()
		checkBase(t, f, p, &want.wantStateBase)

		if !equalBytes(f.Base.Bytes(), want.content) {
			t.Fatalf("f.Bytes() mismatch: want=%v, got=%v", want.content, f.Base.Bytes())
		}
	}

	f1 := NewFirstChunk(&bytes.Buffer{}, 16)
	if f1.Limit != 16 {
		t.Fatalf("limit must be 16 but %d", f1.Limit)
	}
	check(t, f1, []byte("foo"), &contentState{wantStateBase{3, 3, false, nil}, []byte("foo")})
	check(t, f1, []byte("bar"), &contentState{wantStateBase{3, 6, false, nil}, []byte("foobar")})
	check(t, f1, []byte("buz"), &contentState{wantStateBase{3, 9, false, nil}, []byte("foobarbuz")})
	check(t, f1, []byte("qux"), &contentState{wantStateBase{3, 12, false, nil}, []byte("foobarbuzqux")})
	check(t, f1, []byte("qux"), &contentState{wantStateBase{3, 15, false, nil}, []byte("foobarbuzquxqux")})
	check(t, f1, []byte("1"), &contentState{wantStateBase{1, 16, false, nil}, []byte("foobarbuzquxqux1")})
	check(t, f1, []byte("2"), &contentState{wantStateBase{1, 16, true, nil}, []byte("foobarbuzquxqux1")})
	check(t, f1, []byte("quux"), &contentState{wantStateBase{4, 16, true, nil}, []byte("foobarbuzquxqux1")})

	f2 := NewFirstChunk(&bytes.Buffer{}, 5)
	if f2.Limit != 5 {
		t.Fatalf("limit must be 5 but %d", f2.Limit)
	}
	check(t, f2, []byte("foobarbuz"), &contentState{wantStateBase{9, 5, true, nil}, []byte("fooba")})

	want := errors.New("foo")
	f3 := NewFirstChunk(writerFunc(func([]byte) (int, error) { return 0, want }), 16)
	checkBase(t, f3, []byte("foo"), &wantStateBase{3, 0, false, want})
	checkBase(t, f3, []byte("foo"), &wantStateBase{3, 0, false, want})

	f4 := NewFirstChunk(writerFunc(func(p []byte) (int, error) { return len(p) / 2, nil }), 16)
	checkBase(t, f4, []byte("foobar"), &wantStateBase{6, 3, false, ErrFailToWriteAll})
}
