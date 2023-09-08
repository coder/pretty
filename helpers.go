package pretty

import (
	"strings"
	"unicode"

	"github.com/muesli/termenv"
)

// FgColor returns a formatter that sets the foreground color.
// Example:
//
//	FgColor(termenv.RGBColor("#ff0000"))
//	FgColor(termenv.ANSI256Color(196))
//	FgColor(termenv.ANSIColor(31))
func FgColor(c termenv.Color) Formatter {
	seq := c.Sequence(false)
	return CSI(seq)
}

// BgColor returns a formatter that sets the background color.
// Example:
//
//	BgColor(termenv.RGBColor("#ff0000"))
//	BgColor(termenv.ANSI256Color(196))
//	BgColor(termenv.ANSIColor(31))
func BgColor(c termenv.Color) Formatter {
	seq := c.Sequence(true)
	return CSI(seq)
}

// CSI wraps the text in the given CSI (Control Sequence Introducer) sequence.
// Example:
//
//	CSI(termenv.BoldSeq)
//	CSI(termenv.UnderlineSeq)
//	CSI(termenv.ItalicSeq)
func CSI(seq string) Formatter {
	if seq == "" {
		return Nop
	}
	return Wrap(termenv.CSI+seq+"m", termenv.CSI+termenv.ResetSeq+"m")
}

// Bold returns a formatter that makes the text bold.
func Bold() Formatter {
	return CSI(termenv.BoldSeq)
}

// Italic returns a formatter that makes the text italic.
func Italic() Formatter {
	return CSI(termenv.ItalicSeq)
}

// Underline returns a formatter that underlines the text.
func Underline() Formatter {
	return CSI(termenv.UnderlineSeq)
}

// Wrap wraps the text in the given prefix and suffix.
// It is useful for wrapping text in ANSI sequences.
func Wrap(prefix, suffix string) Formatter {
	return formatterFunc(func(t *Text) {
		t.Prepend(prefix)
		t.Append(suffix)
	})
}

// XPad pads the text on the left and right.
func XPad(left, right int) Formatter {
	return formatterFunc(func(t *Text) {
		t.Prepend(strings.Repeat(" ", left))
		t.Append(strings.Repeat(" ", right))
	})
}

// LineWrap wraps the text at the given width.
// It breaks lines at word boundaries when possible. It will never break up
// a word so that URLs and other long strings present correctly.
func LineWrap(width int) Formatter {
	return formatterFunc(func(t *Text) {
		var col int

		for at := t.Head(); at != nil; at = at.Next {
			nlAt := strings.IndexByte(at.S, '\n')
			if nlAt < 0 {
				nlAt = len(at.S)
			}
			col += nlAt

			overflow := (width - col) * -1
			if overflow <= 0 {
				continue
			}

			spaceAt := strings.LastIndexFunc(at.S[:nlAt-overflow+1], unicode.IsSpace)
			if spaceAt < 0 {
				// Never break up a word.
				continue
			}

			next := at.Split(spaceAt)
			at.S = strings.TrimRight(at.S, " \t")
			next.S = strings.TrimLeft(next.S, " \t")
			next.Insert("\n")
			col = 0
		}
	})
}

// Nop is a no-op formatter.
var Nop Formatter = formatterFunc(func(t *Text) {})
