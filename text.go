package pretty

import "strings"

// Text is a linked-list structure that represents an in-progress text string.
// Most formatters work by prepending and appending to text, so this structure
// is far more efficient than manipulating strings directly.
type Text struct {
	S    string
	Next *Text
	Prev *Text
}

func (t *Text) head() *Text {
	for t.Prev != nil {
		t = t.Prev
	}
	return t
}

func (t *Text) tail() *Text {
	for t.Next != nil {
		t = t.Next
	}
	return t
}

func (t *Text) debugString() string {
	var sb strings.Builder
	sb.Grow(t.Len())
	for at := t.head(); at != nil; at = at.Next {
		sb.WriteString(at.S)
		sb.WriteString("->")
	}
	sb.WriteString("▫️")
	return sb.String()
}

// String allocates a new string for the entire text.
func (t *Text) String() string {
	var sb strings.Builder
	sb.Grow(t.Len())
	for at := t.head(); at != nil; at = at.Next {
		sb.WriteString(at.S)
	}
	return sb.String()
}

// Len returns the length of the text.
func (t *Text) Len() int {
	l := 0
	at := t.head()
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
	oldTail := t.tail()
	newTail := &Text{S: s, Prev: oldTail}
	oldTail.Next = newTail
	return newTail
}

// Prepend prepends a string to the beginning of the text and returns the new
// head.
func (t *Text) Prepend(s string) *Text {
	oldHead := t.head()
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
