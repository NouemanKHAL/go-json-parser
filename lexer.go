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
	Pos     Position
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

func NewToken(typ TokenType, lit string, pos Position) Token {
	return Token{
		Type:    typ,
		Literal: lit,
		Pos:     pos,
	}
}

type Position struct {
	line   int
	column int
}

func (p Position) String() string {
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

func (l *Lexer) Lex() Token {
	token := l.lex()
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("%-8s %-8s %-8s\n", token.Pos.String(), token.Type, token.Literal)
	}
	return token
}

func (l *Lexer) readRune() (rune, error) {
	cur, _, err := l.reader.ReadRune()
	return cur, err
}

func (l *Lexer) lex() Token {
	for {
		l.pos.column += 1
		cur, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return NewToken(EOF, "", l.pos)
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
			return NewToken(L_BRACE, "{", l.pos)
		case '}':
			return NewToken(R_BRACE, "}", l.pos)
		case '[':
			return NewToken(L_BRACKET, "[", l.pos)
		case ']':
			return NewToken(R_BRACKET, "]", l.pos)
		case ':':
			return NewToken(COLON, ":", l.pos)
		case ',':
			return NewToken(COMMA, ",", l.pos)
		case '"':
			token := l.lexString()

			next, err := l.readRune()
			if err != nil {
				handleError(err)
			}
			if next != '"' {
				err = fmt.Errorf("invalid json: unterminated literal expected '\"' got '%c' at position %d:%d", next, l.pos.line, l.pos.column)
				handleError(err)
			}
			token.Literal = "\"" + token.Literal + "\""
			return token
		default:
			if unicode.IsSpace(cur) {
				continue
			} else if unicode.IsDigit(cur) || cur == '-' || cur == '+' {
				l.rewind()
				return l.lexNumber()
			} else if unicode.IsLetter(cur) {
				l.rewind()
				return l.lexLiteral()
			} else {
				return NewToken(INVALID, string(cur), l.pos)
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

func (l *Lexer) lexLiteral() Token {
	token := l.lexString()
	if token.Literal == "true" || token.Literal == "false" {
		token.Type = BOOLEAN
	} else if token.Literal == "null" {
		token.Type = NULL
	} else {
		token.Type = INVALID
	}
	return token

}

// character
//
//	'0020' . '10FFFF' - '"' - '\'
//	'\' escape
func (l *Lexer) isValidCharacter(c rune) bool {
	if c < '\u0020' || c > unicode.MaxRune {
		return false
	}
	if c == '\\' || c == '"' || c == ',' {
		return false
	}

	return true
}

func (l *Lexer) lexString() Token {
	lit := ""
	startPos := l.pos
	for {
		cur, _, err := l.reader.ReadRune()
		if err != nil {
			handleError(err)
		}
		l.pos.column += 1

		if l.isValidCharacter(cur) {
			lit = lit + string(cur)
		} else if cur == '\\' {
			/*
				escape
				'"'
				'\'
				'/'
				'b'
				'f'
				'n'
				'r'
				't'
				'u' hex hex hex hex
			*/

			next, _, err := l.reader.ReadRune()
			if err != nil {
				handleError(err)
			}
			l.pos.column += 1

			if next == '"' || next == '\\' || next == '/' || next == 'b' || next == 'f' || next == 'n' || next == 'r' {
				lit = lit + string(cur) + string(next)
			} else {
				l.rewind()
				l.rewind()
				return NewToken(IDENT, lit, startPos)
			}

		} else {
			l.rewind()
			return NewToken(IDENT, lit, startPos)
		}
	}
}

func (l *Lexer) lexNumber() Token {
	lit := ""
	startPos := l.pos
	isDecimal := false
	firstTime := true

	for {
		cur, _, err := l.reader.ReadRune()
		if err != nil {
			handleError(err)
		}
		l.pos.column += 1

		if firstTime {
			firstTime = false
			if cur == '+' || cur == '-' {
				lit += string(cur)
				continue
			}
		}
		if unicode.IsDigit(cur) || (cur == '.' && !isDecimal) {
			lit += string(cur)
			if cur == '.' {
				isDecimal = true
			}
		} else {
			l.rewind()
			return NewToken(NUMBER, lit, startPos)
		}
	}

}
