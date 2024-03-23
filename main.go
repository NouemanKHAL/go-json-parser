package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func handleError(err error) {
	if err == nil {
		return
	}
	os.Stderr.WriteString("error:\n\t" + err.Error() + "\n")
	os.Exit(1)
}

func parseObject(r io.RuneReader) error {
	cur, _, err := r.ReadRune()
	if err != nil {
		return err
	}
	if cur == '}' {
		return nil
	}

	return fmt.Errorf("invalid json: missing '}' at position 1")
}

func ParseJSON(r io.RuneReader) error {
	cur, _, err := r.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = fmt.Errorf("invalid json: empty file")
		}
		handleError(err)
	}

	if cur == '{' {
		err := parseObject(r)
		if err != nil {
			handleError(err)
		}

	} else {
		err := fmt.Errorf("invalid json: missing '{' at position 0")
		if err != nil {
			handleError(err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			handleError(err)
		}

		r := bufio.NewReader(f)
		err = ParseJSON(r)
		if err != nil {
			handleError(err)
		}
	} else {
		handleError(fmt.Errorf("expected input file in first argument"))
	}

}
