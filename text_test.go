package pretty

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireText(t *testing.T, txt *Text, s string) {
	if assert.Equal(t, s, txt.String()) {
		assert.Equal(t, len(s), txt.Len())
	}
	if t.Failed() {
		t.Logf("txt: %s", txt.debugString())
		t.FailNow()
	}
}

func TestText(t *testing.T) {
	t.Run("Len", func(t *testing.T) {
		txt := String("")
		requireText(t, txt, "")
		txt = txt.Append("a")
		txt = txt.Append("b")
		requireText(t, txt, "ab")
		txt = txt.Prepend("c")
		requireText(t, txt, "cab")
	})
	t.Run("PrependAppend", func(t *testing.T) {
		txt := String("")
		txt = txt.Append("a")
		txt = txt.Append("b")
		requireText(t, txt, "ab")
		txt = txt.Prepend("c")
		requireText(t, txt, "cab")
		txt = txt.Prepend("d")
		requireText(t, txt, "dcab")
	})
	t.Run("InsertEnd", func(t *testing.T) {
		txt := String("")
		txt = txt.Append("11")
		requireText(t, txt, "11")
		txt.Insert(("ef"))
		requireText(t, txt, "ef11")
	})
	t.Run("SplitEnd", func(t *testing.T) {
		txt := String("11")
		txt.Split(1)
		requireText(t, txt, "11")
	})

	t.Run("SplitMiddle", func(t *testing.T) {
		txt := String("123456")
		txt.Append("789")
		txt.Head().Split(3)
		require.Equal(t, "123", txt.Head().S)
		requireText(t, txt, "123456789")
	})
}

func TestFgColor(t *testing.T) {
	txt := String("disgusting red on green")
	FgColor(termenv.RGBColor("#ff0000")).Format(txt)
	BgColor(termenv.RGBColor("#00ff00")).Format(txt)
	t.Logf("txt: %v", txt.debugString())
	t.Logf("txt: %s", txt)
}

func TestLineWrap(t *testing.T) {
	t.Parallel()

	t.Run("None", func(t *testing.T) {
		txt := String(
			"The crazy fox jumped",
		)

		LineWrap(100).Format(txt)

		requireText(t, txt, "The crazy fox jumped")
	})
	t.Run("Basic", func(t *testing.T) {
		txt := String(
			"The crazy fox jumped",
		)

		LineWrap(10).Format(txt)

		requireText(t, txt, "The crazy\nfox jumped")
	})
	t.Run("WordBoundary", func(t *testing.T) {
		txt := String(
			"The crazy_fox_jumped",
		)

		LineWrap(10).Format(txt)

		requireText(t, txt, "The\ncrazy_fox_jumped")
	})
	t.Run("MultiLine", func(t *testing.T) {
		txt := String(
			"aabb cc dd ee ff",
		)

		LineWrap(4).Format(txt)

		requireText(t, txt, "aabb\ncc\ndd\nee\nff")
	})
}

func TestStyle(t *testing.T) {
	errorStyle := Style{
		FgColor(termenv.RGBColor("#ff0000")),
		BgColor(termenv.RGBColor("#000000")),
		WrapCSI(termenv.BoldSeq),
	}

	t.Logf("%s", errorStyle.Sprint("SOME ERROR"))
}
