package main

import (
	"fmt"
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
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			handleError(err)
		}

		l := NewLexer(f)
		p := NewParser(l)
		err = p.Parse()
		if err != nil {
			handleError(err)
		}

		fmt.Println(p)

	} else {
		handleError(fmt.Errorf("expected input file in first argument"))
	}

}
