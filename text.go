package pretty

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"
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
func (t *Text) Split(n int) *Text {
	if n > len(t.S) {
		panic("split index out of bounds")
	}

	if len(t.S) == n || len(t.S) == 0 {
		return t
	}

	// Split the string.
	right := t.S[n:]
	t.S = t.S[:n]

	if t.Next == nil {
		t.Next = &Text{S: right, Prev: t}
		return t.Next
	}
	t.Next.Insert(right)
	return t.Next
}

// Insert inserts the given text before the current text.
func (t *Text) Insert(s string) {
	tt := &Text{S: s}
	oldPrev := t.Prev
	oldPrev.Next = tt
	tt.Prev = oldPrev
	tt.Next = t
	t.Prev = tt
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

// Style is a list of formatters, which are applied in order.
type Style []Formatter

// Format applies all formatters in the style to the text and
// returns the modified text.
//
// When performance is a concern, use WriteTo instead of String
// on the returned text.
func (s Style) Format(t *Text) *Text {
	for _, f := range s {
		f.Format(t)
	}
	return t
}

// With returns a new style with the given formatters appended.
func (s Style) With(fs ...Formatter) Style {
	return append(s, fs...)
}

// Sprintf formats the given string with the style.
func (s Style) Sprintf(format string, args ...interface{}) string {
	return s.Format(String(fmt.Sprintf(format, args...))).String()
}

func (s Style) Sprint(args ...interface{}) string {
	return s.Format(String(fmt.Sprint(args...))).String()
}

func (s Style) Printf(format string, args ...interface{}) {
	str := s.Format(String(fmt.Sprintf(format, args...))).String()
	fmt.Print(str)
}

type formatterFunc func(*Text)

func (f formatterFunc) Format(t *Text) {
	f(t)
}

// FgColor returns a formatter that sets the foreground color.
// Example:
//
//	FgColor(termenv.RGBColor("#ff0000"))
//	FgColor(termenv.ANSI256Color(196))
//	FgColor(termenv.ANSIColor(31))
func FgColor(c termenv.Color) Formatter {
	seq := c.Sequence(false)
	return WrapCSI(seq)
}

// BgColor returns a formatter that sets the background color.
// Example:
//
//	BgColor(termenv.RGBColor("#ff0000"))
//	BgColor(termenv.ANSI256Color(196))
//	BgColor(termenv.ANSIColor(31))
func BgColor(c termenv.Color) Formatter {
	seq := c.Sequence(true)
	return WrapCSI(seq)
}

// WrapCSI wraps the text in the given CSI (Control Sequence Introducer) sequence.
// Example:
//
//	WrapCSI(termenv.BoldSeq)
//	WrapCSI(termenv.UnderlineSeq)
//	WrapCSI(termenv.ItalicSeq)
func WrapCSI(seq string) Formatter {
	return Wrap(termenv.CSI+seq+"m", termenv.CSI+termenv.ResetSeq+"m")
}

// Wrap wraps the text in the given prefix and suffix.
// It is useful for wrapping text in ANSI sequences.
func Wrap(prefix, suffix string) Formatter {
	return formatterFunc(func(t *Text) {
		t.Prepend(prefix)
		t.Append(suffix)
	})
}
