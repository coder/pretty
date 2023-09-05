package bench

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/coder/pretty"
	"github.com/muesli/termenv"
)

const helloWorld = "Hello, World!"

func BenchmarkPretty(b *testing.B) {
	style := &pretty.Style{
		pretty.FgColor(termenv.RGBColor("#FF0000")),
		pretty.BgColor(termenv.RGBColor("#0000FF")),
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(helloWorld)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = style.Sprint(helloWorld)
	}
}

func BenchmarkLipgloss(b *testing.B) {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Background(lipgloss.Color("#0000FF"))

	b.ReportAllocs()
	b.SetBytes(int64(len(helloWorld)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = style.Render(helloWorld)
	}
}
