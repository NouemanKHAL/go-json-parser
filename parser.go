package main

import "fmt"

/*
JSON structure:
	root: value

	value : object | array | literal

	object: '{' [ property [, property]* ] '}'
	array: '[' [ value [, value]* ]  ']'
	literal:  string | number | "true" | "false" | "null"

	property: literal ':' value

	string: '"' [a-zA-Z0-9_] '"" // TODO: fix this expression
	number: [0-9]+
*/

type State int

const (
	INIT_STATE State = iota

	OBJECT_START
	OBJECT_OPEN
	OBJECT_END

	ARRAY_START
	ARRAY_END
)

type NodeType int

const (
	ROOT NodeType = iota
	OBJECT
	ARRAY
	LITERAL
	PROPERTY
)

type Value interface {
	GetType() NodeType
}

type Root struct {
	Type  NodeType
	Value Value
}

func (r Root) GetType() NodeType {
	return r.Type
}

type Object struct {
	Type       NodeType
	Properties []Property
}

func (o Object) GetType() NodeType {
	return o.Type
}

type Array struct {
	Type     NodeType
	Elements []Value
}

func (a Array) GetType() NodeType {
	return a.Type
}

type Property struct {
	Type  NodeType
	Key   Literal
	Value Value
}

func (p Property) GetType() NodeType {
	return p.Type
}

type Literal struct {
	Type  NodeType
	Value string
}

func (l Literal) GetType() NodeType {
	return l.Type
}

type JsonTree struct {
	root *Root
}

func NewJsonTree() *JsonTree {
	return &JsonTree{}

}

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	tree      *JsonTree
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer: l,
		tree:  NewJsonTree(),
	}
	// fill curToken & peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) readToken() Token {
	_, Token := p.lexer.Lex()
	return Token
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.readToken()
}

func (p *Parser) parseValue() (Value, error) {
	var value Value
	var err error
	switch p.curToken.Type {
	case L_BRACE:
		value, err = p.parseObject()
	case L_BRACKET:
		value, err = p.parseArray()
	case IDENT:
		value, err = p.parseLiteral()
	case NUMBER:
		value, err = p.parseLiteral()
	case BOOLEAN:
		value, err = p.parseLiteral()
	default:
		err := fmt.Errorf("cannot parse value, got token '%s'", p.curToken.Literal)
		return nil, err
	}
	return value, err
}

func (p *Parser) parseObject() (Object, error) {
	if p.curToken.Type != L_BRACE {
		return Object{}, fmt.Errorf("invalid start of object, expected '{' got '%s'", p.curToken.Literal)
	}

	state := OBJECT_START

	object := Object{
		Type:       OBJECT,
		Properties: make([]Property, 0),
	}

	for {
		if p.peekToken.Type == EOF {
			if state != OBJECT_OPEN {
				return object, nil

			}
			return Object{}, fmt.Errorf("invalid object, EOF reached before '{'")
		}

		switch state {
		case OBJECT_START:
			switch p.peekToken.Type {
			case R_BRACE:
				p.nextToken()
				state = OBJECT_END
			case COMMA:
				p.nextToken()
				p.nextToken()
				state = OBJECT_OPEN
			case IDENT:
				p.nextToken()
				state = OBJECT_OPEN
			default:
				return Object{}, fmt.Errorf("invalid object, invalid token '%s' of type %s", p.peekToken.Literal, p.peekToken.Type)
			}
		case OBJECT_OPEN:
			prop, err := p.parseProperty()
			if err != nil {
				return Object{}, err
			}
			object.Properties = append(object.Properties, prop)
			state = OBJECT_START
		case OBJECT_END:
			p.nextToken()
			return object, nil
		default:
			panic("parsing object reached unknown state")
		}
	}
}

func (p *Parser) parseArray() (Array, error) {
	if p.curToken.Type != L_BRACKET {
		return Array{}, fmt.Errorf("invalid start of array, expected '[' got '%s'", p.curToken.Literal)
	}
	Arr := Array{Type: ARRAY, Elements: make([]Value, 0)}
	state := ARRAY_START
	for {
		if p.peekToken.Type == EOF {
			if state == ARRAY_END {
				return Arr, nil

			}
			return Arr, fmt.Errorf("invalid array, EOF reached before '['")
		}
		switch state {
		case ARRAY_START:
			p.nextToken()
			if p.curToken.Type == R_BRACKET {
				p.nextToken()
				return Arr, nil
			}
			v, err := p.parseValue()
			handleError(err)

			Arr.Elements = append(Arr.Elements, v)

			if p.peekToken.Type == R_BRACKET {
				p.nextToken()
				return Arr, nil
			}

			if p.peekToken.Type == COMMA {
				p.nextToken()
			}
		case ARRAY_END:
			p.nextToken()
			return Arr, nil
		}
	}
}

func (p *Parser) parseLiteral() (Literal, error) {
	if p.curToken.Type == IDENT || p.curToken.Type == BOOLEAN || p.curToken.Type == NULL || p.curToken.Type == NUMBER {
		return Literal{Type: LITERAL, Value: p.curToken.Literal}, nil
	}
	return Literal{}, fmt.Errorf("cannot parse literal from token '%s' of type '%s'", p.curToken.Literal, p.curToken.Type)
}

func (p *Parser) parseProperty() (Property, error) {
	lit, err := p.parseLiteral()
	if err != nil {
		return Property{}, err
	}

	if p.peekToken.Type != COLON {
		return Property{}, fmt.Errorf("invalid property, expected ':' got '%s'", p.curToken.Literal)
	}
	p.nextToken()
	p.nextToken()

	value, err := p.parseValue()
	if err != nil {
		return Property{}, err
	}
	return Property{Type: PROPERTY, Key: lit, Value: value}, nil
}

func (p *Parser) Parse() error {
	value, err := p.parseValue()
	if err != nil {
		return err
	}

	p.tree.root = &Root{
		Type:  ROOT,
		Value: value,
	}

	return nil
}
