package pretty

import (
	"testing"

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
