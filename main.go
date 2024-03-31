package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func handleError(err error) {
	if err == nil {
		return
	}
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func main() {
	var r io.Reader

	switch len(os.Args) {
	case 1:
		stat, err := os.Stdin.Stat()

		if (stat.Mode()&fs.ModeNamedPipe) == 0 || err != nil {
			fmt.Println("Usage:\n\tgojson [FILE] [PATH]")
			os.Exit(0)
		}
		r = bufio.NewReader(os.Stdin)
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			handleError(err)
		}
		r = bufio.NewReader(f)
	}

	l := NewLexer(r)
	p := NewParser(l)
	p.Parse()
	fmt.Println(p)
}
