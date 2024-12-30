package shared

import (
	"unicode/utf8"
	"unsafe"
)

// A Builder is used to efficiently build a string using [Builder.Write] methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type StringBuilder struct {
	addr *StringBuilder // of receiver, to detect copies by value
	buf  []byte
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheck
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	//goland:noinspection GoVetUnsafePointer
	return unsafe.Pointer(x ^ 0)
}

func (b *StringBuilder) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*StringBuilder)(noescape(unsafe.Pointer(b)))
	} else if b.addr != b {
		panic("strings: illegal use of non-zero StringBuilder copied by value")
	}
}

// String returns the accumulated string.
func (b *StringBuilder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *StringBuilder) Len() int { return len(b.buf) }

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *StringBuilder) Cap() int { return cap(b.buf) }

// Reset resets the [Builder] to be empty.
func (b *StringBuilder) Reset() {
	b.addr = nil
	b.buf = nil
}

func (b *StringBuilder) Rewind(n int) {
	if n > len(b.buf) {
		n = len(b.buf)
	}
	b.buf = b.buf[:len(b.buf)-n]
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *StringBuilder) grow(n int) {
	buf := make([]byte, 2*cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *StringBuilder) Grow(n int) {
	b.copyCheck()
	if n < 0 {
		panic("strings.StringBuilder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *StringBuilder) Write(p []byte) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *StringBuilder) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *StringBuilder) WriteRune(r rune) (int, error) {
	b.copyCheck()
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *StringBuilder) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, s...)
	return len(s), nil
}
