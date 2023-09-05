# pretty

`pretty` is a high performance Terminal pretty printer for Go. We created it
due to the significant performance issues encountered with `lipgloss`, which
seemed impossible to fix without a complete rewrite.

It is relatively low-level, meant to be used in tandem with [termenv](https://pkg.go.dev/github.com/muesli/termenv).


## Basic Usage

```go
// E.g:
errorStyle := pretty.Style{
		pretty.FgColor(termenv.RGBColor("#ff0000")),
		pretty.BgColor(termenv.RGBColor("#000000")),
		pretty.WrapCSI(termenv.BoldSeq),
}

errorStyle.Printf("something bad")
```

## Adaptive Color

You can use `termenv` to adapt the colors to the terminal's color palette:

```go
// E.g:
profile := termenv.NewOutput(os.Stdout, termenv.WithColorCache(true)).ColorProfile()
errorStyle := pretty.Style{
        pretty.FgColor(profile.Color("#ff0000")),
        pretty.BgColor(profile.Color("#000000")),
        pretty.WrapCSI(termenv.BoldSeq),
}
```

Note that termenv requires a round-trip between the terminal and the program
to determine the terminal's color palette. Even with caching, this can result
in noticeable unresponsiveness.