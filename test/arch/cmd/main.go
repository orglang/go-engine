package main

import (
	"fmt"
	"os"

	"orglang/go-engine/test/arch"
)

func main() {
	graph, err := arch.Load("./...")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	violations := arch.Check(graph)
	if len(violations) == 0 {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "hexagon boundary violated in %d place(s):\n", len(violations))
	for _, v := range violations {
		_, _ = fmt.Fprintf(os.Stderr, "  %s\n", v)
	}
	_, _ = fmt.Fprintln(os.Stderr, "an agnostic core must not import a toolkit-specific adapter")
	os.Exit(1)
}
