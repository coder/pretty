package pretty

import (
	"fmt"
	"io"
	"strings"
)

// Text is a linked-list structure that represents an in-progress text string.
// Most formatters work by prepending and appending to text, so this structure
// is far more efficient than manipulating strings directly.
type Text struct {
	S    string
	Next *Text
	Prev *Text
}

func (t *Text) Head() *Text {
	for t.Prev != nil {
		t = t.Prev
	}
	return t
}

func (t *Text) Tail() *Text {
	for t.Next != nil {
		t = t.Next
	}
	return t
}

// Split splits the current text into two parts at the given index.
// It returns the new node.
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
	oldPrev.Next = tt
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

// String allocates a new string for the entire text.
func (t *Text) String() string {
	var sb strings.Builder
	sb.Grow(t.Len())
	t.WriteTo(&sb)
	return sb.String()
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

// Append appends a string to the end of the text and returns the new tail.
func (t *Text) Append(s string) *Text {
	oldTail := t.Tail()
	newTail := &Text{S: s, Prev: oldTail}
	oldTail.Next = newTail
	return newTail
}

// Prepend prepends a string to the beginning of the text and returns the new
// head.
func (t *Text) Prepend(s string) *Text {
	oldHead := t.Head()
	newHead := &Text{S: s, Next: oldHead}
	oldHead.Prev = newHead
	return newHead
}

// String returns a new Text object from a String.
func String(s string) *Text {
	return &Text{S: s}
}

// Formatter manipulates a Text object.
type Formatter interface {
	Format(*Text)
}

var _ Formatter = Style(nil)

// Style is a special formatter that applies multiple formatters to a text.
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
	return append(s, fs...)
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
func Fprintf(w io.Writer, f Formatter, format string, args ...interface{}) {
	txt := String(fmt.Sprintf(format, args...))
	f.Format(txt)
	txt.WriteTo(w)
}
