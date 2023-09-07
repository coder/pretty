package pretty

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestFgColor(t *testing.T) {
	txt := String("disgusting red on green")
	FgColor(termenv.RGBColor("#ff0000")).Format(txt)
	BgColor(termenv.RGBColor("#00ff00")).Format(txt)
	t.Logf("txt: %v", txt.debugString())
	t.Logf("txt: %s", txt)
}

func TestLineWrap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "None",
			input:    "The crazy fox jumped",
			width:    100,
			expected: "The crazy fox jumped",
		},
		{
			name:     "Basic",
			input:    "The crazy fox jumped",
			width:    10,
			expected: "The crazy\nfox jumped",
		},
		{
			name:     "WordBoundary",
			input:    "The crazy_fox_jumped",
			width:    10,
			expected: "The\ncrazy_fox_jumped",
		},
		{
			name:     "MultiLine",
			input:    "aabb cc dd ee ff",
			width:    4,
			expected: "aabb\ncc\ndd\nee\nff",
		},
		{
			name:     "EmptyString",
			input:    "",
			width:    10,
			expected: "",
		},
		{
			name:     "SingleWordLongerThanWrap",
			input:    "supercalifragilisticexpialidocious",
			width:    10,
			expected: "supercalifragilisticexpialidocious",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			txt := String(tc.input)
			LineWrap(tc.width).Format(txt)
			requireText(t, txt, tc.expected)
		})
	}
}

func TestXPad(t *testing.T) {
	txt := String("a")
	XPad(1, 2).Format(txt)
	requireText(t, txt, " a  ")
}

func TestStyle(t *testing.T) {
	errorStyle := Style{
		FgColor(termenv.RGBColor("#ff0000")),
		BgColor(termenv.RGBColor("#000000")),
		CSI(termenv.BoldSeq),
	}

	t.Logf("%s", Sprint(errorStyle, "SOME ERROR"))
}
