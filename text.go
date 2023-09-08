package pretty

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Text is a linked-list structure that represents an in-progress text string.
// Most formatters work by prepending and appending to text, so this structure
// is far more efficient than manipulating strings directly.
//
// The pointer is instrinsicly a cursor that points to the current text segment.
// So, subsequent appends and prepends are O(1) since the cursor is already at
// the tail or head respectively.
type Text struct {
	S    string
	Next *Text
	Prev *Text
}

// Head returns the absolute head of the text.
// It adjusts the pointer to the head of the text.
func (t *Text) Head() *Text {
	for t.Prev != nil {
		t = t.Prev
	}
	return t
}

// Tail returns the absolute tail of the text.
// It adjusts the pointer to the tail of the text
func (t *Text) Tail() *Text {
	for t.Next != nil {
		t = t.Next
	}
	return t
}

// Split splits the current text into two parts at the given index. The current
// node contains the first part, and the new node contains the second part.
// It returns the new node, and does not adjust the pointer.
func (t *Text) Split(n int) *Text {
	if n > len(t.S) {
		panic(fmt.Sprintf("split index %d > len(t.S) (%v) ", n, len(t.S)))
	}

	if n <= 0 {
		panic("split index must be > 0")
	}

	if len(t.S) == n || len(t.S) == 0 {
		return t
	}

	// Split the string.
	nextStr := t.S[n:]
	t.S = t.S[:n]

	if t.Next == nil {
		t.Next = &Text{S: nextStr, Prev: t}
		return t.Next
	}
	t.Next.Insert(nextStr)
	return t.Next
}

// Insert inserts the given text before the current text.
// It returns the new node.
func (t *Text) Insert(s string) *Text {
	tt := &Text{S: s}
	oldPrev := t.Prev
	if oldPrev != nil {
		oldPrev.Next = tt
	}
	tt.Prev = oldPrev
	tt.Next = t
	t.Prev = tt
	return tt
}

func (t *Text) debugString() string {
	var sb strings.Builder
	sb.Grow(t.Len())
	for at := t.Head(); at != nil; at = at.Next {
		fmt.Fprintf(&sb, "%q ->", at.S)
	}
	sb.WriteString("▫️")
	return sb.String()
}

// String allocates a new string containing the entire text.
func (t *Text) String() string {
	var sb strings.Builder
	sb.Grow(t.Len())
	t.WriteTo(&sb)
	return sb.String()
}

// Bytes allocates a new byte slice containing the entire text.
// It uses the given buffer if it is large enough.
func (t *Text) Bytes(b []byte) []byte {
	buf := bytes.NewBuffer(b[:0])
	buf.Grow(t.Len())
	t.WriteTo(buf)
	return buf.Bytes()
}

// WriteTo writes the text to the given writer, avoiding
// string allocations.
func (t *Text) WriteTo(w io.Writer) (int64, error) {
	var n int64
	for at := t.Head(); at != nil; at = at.Next {
		nn, err := io.WriteString(w, at.S)
		n += int64(nn)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// Len returns the length of the text.
func (t *Text) Len() int {
	l := 0
	at := t.Head()
	for {
		l += len(at.S)
		at = at.Next
		if at == nil {
			return l
		}
	}
}

// Append appends strings to the end of the text
// in order.
// Example:
//
//	txt := String("a")
//	txt = txt.Append("b", "c")
//	fmt.Println(txt.String())
//	// Output: abc
func (t *Text) Append(ss ...string) *Text {
	for _, s := range ss {
		t = t.appendOne(s)
	}
	return t
}

// appendOne appends a string to the end of the text and returns the new tail.
func (t *Text) appendOne(s string) *Text {
	oldTail := t.Tail()
	newTail := &Text{S: s, Prev: oldTail}
	oldTail.Next = newTail
	return newTail
}

// Prepend prepends strings to the beginning of the text
// in order.
// Example:
//
//	txt := String("c")
//	txt = txt.Prepend("a", "b")
//	fmt.Println(txt.String())
//	// Output: abc
func (t *Text) Prepend(ss ...string) *Text {
	for i := len(ss) - 1; i >= 0; i-- {
		t = t.prependOne(ss[i])
	}
	return t
}

func (t *Text) prependOne(s string) *Text {
	oldHead := t.Head()
	newHead := &Text{S: s, Next: oldHead}
	oldHead.Prev = newHead
	return newHead
}

// String creates a new Text object from the given strings.
func String(s ...string) *Text {
	if len(s) == 0 {
		return &Text{}
	}
	txt := &Text{S: s[0]}
	for _, s := range s[1:] {
		txt = txt.appendOne(s)
	}
	return txt
}

// Formatter manipulates Text.
type Formatter interface {
	Format(*Text)
}

var _ Formatter = Style(nil)

// Style is a special Formatter that applies multiple Formatters to a text
// in order.
type Style []Formatter

// Format applies all formatters in the style to the text and
// returns the modified text.
//
// When performance is a concern, use WriteTo instead of String
// on the returned text.
func (s Style) Format(t *Text) {
	for _, f := range s {
		f.Format(t)
	}
}

// With returns a new style with the given formatters appended.
func (s Style) With(fs ...Formatter) Style {
	// Force a copy of the slice to avoid multiple calls to With
	// interfering with each other.
	return append(s[:len(s):len(s)], fs...)
}

type formatterFunc func(*Text)

func (f formatterFunc) Format(t *Text) {
	f(t)
}

// Printf formats the given string with the formatter and prints it to stdout.
func Printf(f Formatter, format string, args ...interface{}) {
	txt := String(fmt.Sprintf(format, args...))
	f.Format(txt)
	fmt.Print(txt.String())
}

// Sprintf formats the given string with the formatter.
func Sprintf(f Formatter, format string, args ...interface{}) string {
	txt := String(fmt.Sprintf(format, args...))
	f.Format(txt)
	return txt.String()
}

// Sprint formats the given string with the formatter.
func Sprint(f Formatter, args ...interface{}) string {
	txt := String(fmt.Sprint(args...))
	f.Format(txt)
	return txt.String()
}

// Fprintf formats the given string with the formatter and writes it to the
// given writer.
func Fprintf(w io.Writer, f Formatter, format string, args ...interface{}) {
	txt := String(fmt.Sprintf(format, args...))
	f.Format(txt)
	txt.WriteTo(w)
}

// Fprint formats the given string with the formatter and writes it to the
// given writer.
func Fprint(w io.Writer, f Formatter, args ...interface{}) {
	txt := String(fmt.Sprint(args...))
	f.Format(txt)
	txt.WriteTo(w)
}
