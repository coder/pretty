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
		require.Equal(t, "ab", txt.String())
		txt = txt.Prepend("c")
		require.Equal(t, "cab", txt.String())
		txt = txt.Prepend("d")
		require.Equal(t, "dcab", txt.String())
	})
}
