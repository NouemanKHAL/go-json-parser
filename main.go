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
	var query string = "."

	if len(os.Args) == 1 {
		stat, err := os.Stdin.Stat()

		if (stat.Mode()&fs.ModeNamedPipe) == 0 || err != nil {
			fmt.Println("Usage:\n\tgojson [FILE] [PATH]")
			os.Exit(0)
		}
		r = bufio.NewReader(os.Stdin)
	} else if len(os.Args) >= 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			handleError(err)
		}
		r = bufio.NewReader(f)
		if len(os.Args) >= 3 {
			query = os.Args[2]
		}
	}

	l := NewLexer(r)
	p := NewParser(l)
	p.Parse()

	out, err := p.Get(query)
	handleError(err)

	fmt.Println(out)
}
