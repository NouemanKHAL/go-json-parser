package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token struct {
	Type    TokenType
	Literal string
}
type TokenType int

const (
	INVALID TokenType = iota

	EOF

	L_BRACE
	R_BRACE

	L_BRACKET
	R_BRACKET

	IDENT

	BOOLEAN
	NUMBER
	NULL

	COLON
	COMMA
)

var tokens = []string{
	INVALID:   "INVALID",
	EOF:       "EOF",
	L_BRACE:   "L_BRACE",
	R_BRACE:   "R_BRACE",
	L_BRACKET: "L_BRACKET",
	R_BRACKET: "R_BRACKET",
	IDENT:     "IDENT",
	BOOLEAN:   "BOOLEAN",
	NUMBER:    "NUMBER",
	NULL:      "NULL",
	COLON:     "COLON",
	COMMA:     "COMMA",
}

func (t TokenType) String() string {
	return tokens[t]
}

func NewToken(typ TokenType, lit string) Token {
	return Token{
		Type:    typ,
		Literal: lit,
	}
}

type Position struct {
	line   int
	column int
}

func (p *Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
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

func (l *Lexer) Lex() (Position, Token) {
	pos, token := l.lex()
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("%-8s %-8s %-8s\n", pos.String(), token.Type, token.Literal)
	}
	return pos, token
}

func (l *Lexer) readRune() (rune, error) {
	cur, _, err := l.reader.ReadRune()
	return cur, err
}

func (l *Lexer) lex() (Position, Token) {
	for {
		l.pos.column += 1
		cur, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return l.pos, NewToken(EOF, "")
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
			return l.pos, NewToken(L_BRACE, "{")
		case '}':
			return l.pos, NewToken(R_BRACE, "}")
		case ':':
			return l.pos, NewToken(COLON, ":")
		case ',':
			return l.pos, NewToken(COMMA, ",")
		case '"':
			p, token := l.lexString()
			next, err := l.readRune()
			if err != nil || next != '"' {
				err = fmt.Errorf("invalid json: unterminated literal expected '\"' got %s: %s", next, err)
				handleError(err)
			}
			token.Literal = "\"" + token.Literal + "\""
			return p, token
		default:
			if unicode.IsSpace(cur) {
				continue
			} else if unicode.IsDigit(cur) {
				l.rewind()
				return l.lexNumber()
			} else if unicode.IsLetter(cur) {
				l.rewind()
				return l.lexLiteral()
			} else {
				return l.pos, NewToken(INVALID, string(cur))
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

func (l *Lexer) lexLiteral() (Position, Token) {
	pos, token := l.lexString()
	if token.Literal == "true" || token.Literal == "false" {
		token.Type = BOOLEAN
	} else if token.Literal == "null" {
		token.Type = NULL
	} else {
		token.Type = INVALID
	}
	return pos, token

}

func (l *Lexer) lexString() (Position, Token) {
	lit := ""
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
			l.rewind()
			return startPos, NewToken(IDENT, lit)
		}
	}
}

func (l *Lexer) lexNumber() (Position, Token) {
	lit := ""
	startPos := l.pos
	for {
		cur, _, err := l.reader.ReadRune()
		if err != nil {
			handleError(err)
		}
		l.pos.column += 1

		if unicode.IsDigit(cur) {
			lit += string(cur)
		} else {
			l.rewind()
			return startPos, NewToken(NUMBER, lit)
		}
	}

}
