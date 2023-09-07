# pretty
[![Go Reference](https://pkg.go.dev/badge/github.com/coder/pretty.svg)](https://pkg.go.dev/github.com/coder/pretty)


`pretty` is a performant Terminal pretty printer for Go. We built it after
using lipgloss and experiencing significant performance issues.

`pretty` doesn't implement escape sequences and should be used alongside [termenv](https://pkg.go.dev/github.com/muesli/termenv).


## Basic Usage

```go
errorStyle := pretty.Style{
		pretty.FgColor(termenv.RGBColor("#ff0000")),
		pretty.BgColor(termenv.RGBColor("#000000")),
		pretty.WrapCSI(termenv.BoldSeq),
}

errorStyle.Printf("something bad")
```

## Color

You can use `termenv` to adapt the colors to the terminal's color palette:

```go
profile := termenv.NewOutput(os.Stdout, termenv.WithColorCache(true)).ColorProfile()
errorStyle := pretty.Style{
        pretty.FgColor(profile.Color("#ff0000")),
        pretty.BgColor(profile.Color("#000000")),
        pretty.WrapCSI(termenv.BoldSeq),
}
```

## Performance

```
$ go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/coder/pretty/bench
BenchmarkPretty-10               5142177               232.6 ns/op        55.88 MB/s         272 B/op          8 allocs/op
BenchmarkLipgloss-10              280276              4157 ns/op           3.13 MB/s         896 B/op         72 allocs/op
PASS
ok      github.com/coder/pretty/bench   2.921s
```

pretty remains fast even through dozens of transformations due to its linked-list
based intermediate representation of text. In general, operations scale with
the number of links rather than the length of the text. For example, coloring
a 1000 character string green is just as fast as wrapping a 1 character string.