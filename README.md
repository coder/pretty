# pretty

`pretty` is a high performance Terminal pretty printer for Go. We created it
due to the significant performance issues encountered with `lipgloss`, which
seemed impossible to fix without a complete rewrite.

It is relatively low-level, meant to used in tandem with [termenv](https://pkg.go.dev/github.com/muesli/termenv).

