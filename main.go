package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL

	L_BRACE
	R_BRACE

	DBL_QUOTE
	IDENT
	COLON
	COMMA
)

var tokens []string = []string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	L_BRACE:   "{",
	R_BRACE:   "}",
	DBL_QUOTE: "\"",
	IDENT:     "IDENT",
	COLON:     ":",
	COMMA:     ",",
}

type Position struct {
	line   int
	column int
}

func (p *Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (t Token) String() string {
	return tokens[t]
}

type Lexer struct {
	reader *bufio.Reader
	pos    Position
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		pos:    Position{line: 1, column: 0},
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		l.pos.column += 1
		cur, _, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return l.pos, EOF, ""
			}
			err = fmt.Errorf("invalid json: %s", err.Error())
			handleError(err)
		}
		switch cur {
		case '\n':
			l.pos.line += 1
			l.pos.column = 0
			continue
		case '{':
			return l.pos, L_BRACE, "{"
		case '}':
			return l.pos, R_BRACE, "}"
		case ':':
			return l.pos, COLON, ":"
		case ',':
			return l.pos, COMMA, ","
		case '"':
			return l.LexIdent()
		default:
			if unicode.IsSpace(cur) {
				continue
			} else {
				return l.pos, ILLEGAL, string(cur)
			}
		}
	}
}

func (l *Lexer) rewind() {
	err := l.reader.UnreadRune()
	if err != nil {
		handleError(err)
	}
	l.pos.column -= 1
}

func (l *Lexer) LexIdent() (Position, Token, string) {
	lit := "\""
	startPos := l.pos
	for {
		cur, _, err := l.reader.ReadRune()
		if err != nil {
			handleError(err)
		}
		l.pos.column += 1

		if unicode.IsLetter(cur) || unicode.IsDigit(cur) {
			lit = lit + string(cur)
		} else {
			if cur == '"' {
				return startPos, IDENT, lit + "\""
			}
			l.rewind()
			return l.pos, IDENT, lit
		}
	}
}

func handleError(err error) {
	if err == nil {
		return
	}
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func main() {
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			handleError(err)
		}

		l := NewLexer(f)

		for {
			pos, token, lit := l.Lex()
			if token == EOF {
				break
			}
			fmt.Printf("%-8s %-8s %-8s\n", pos.String(), token, lit)
		}
	} else {
		handleError(fmt.Errorf("expected input file in first argument"))
	}

}
