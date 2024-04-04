package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

var states = []string{
	INIT_STATE: "INIT_STATE",

	OBJECT_START: "OBJECT_START",
	OBJECT_OPEN:  "OBJECT_OPEN",
	OBJECT_END:   "OBJECT_END",

	ARRAY_START: "ARRAY_START",
	ARRAY_END:   "ARRAY_END",
}

func (s State) String() string {
	return states[s]
}

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

func NewRoot() *Root {
	return &Root{
		Type: ROOT,
	}
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

type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	root      *Root
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer: l,
		root:  NewRoot(),
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
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("parseValue")
	}
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
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("parseObject")
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

		if os.Getenv("DEBUG") == "true" {
			fmt.Println("parseObject state: ", state)
		}
		switch state {
		case OBJECT_START:
			if p.curToken.Type != L_BRACE {
				return Object{}, fmt.Errorf("invalid start of object, expected '{' got '%s'", p.curToken.Literal)
			}
			switch p.peekToken.Type {
			case R_BRACE:
				p.nextToken()
				state = OBJECT_END
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

			switch p.peekToken.Type {
			case R_BRACE:
				p.nextToken()
				state = OBJECT_END
			case COMMA:
				p.nextToken()
				p.nextToken()
				state = OBJECT_OPEN
			default:
				return Object{}, fmt.Errorf("invalid object, invalid token '%s' of type %s", p.peekToken.Literal, p.peekToken.Type)
			}

		case OBJECT_END:
			p.nextToken()
			return object, nil
		default:
			panic("parsing object reached unknown state")
		}
	}
}

func (p *Parser) parseArray() (Array, error) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("parseArray")
	}
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
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("parseLiteral")
	}
	if p.curToken.Type == IDENT || p.curToken.Type == BOOLEAN || p.curToken.Type == NULL || p.curToken.Type == NUMBER {
		return Literal{Type: LITERAL, Value: p.curToken.Literal}, nil
	}
	return Literal{}, fmt.Errorf("cannot parse literal from token '%s' of type '%s'", p.curToken.Literal, p.curToken.Type)
}

func (p *Parser) parseProperty() (Property, error) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("parseProperty cur:", p.curToken.Literal, " next:", p.peekToken.Literal)
	}
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
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("func Parse()")
	}
	value, err := p.parseValue()
	if err != nil {
		return err
	}

	p.root.Value = value
	return nil
}

func (p *Parser) String() string {
	if p.root != nil {
		return p.printValue(p.root.Value)
	}
	return ""
}

// Get returns a json value corresponding to the given query
// A query is a string referring to a path in the underlying JSON data
// A query starts with a . (dot) representing the root node, and followed by either attribute names
// or a [i] to access the i-th element of an array
func (p *Parser) Get(query string) (string, error) {
	fields := strings.Split(query, ".")
	v := p.root.Value

	for _, field := range fields {
		if field == "" {
			continue
		}

		switch v.GetType() {
		case OBJECT:
			o := v.(Object)
			for _, prop := range o.Properties {
				if field == prop.Key.Value[1:len(prop.Key.Value)-1] {
					v = prop.Value
				}
			}
		case ARRAY:
			a := v.(Array)
			if end := strings.Index(field, "]"); end != -1 {
				idx, err := strconv.Atoi(field[1:end])
				if err != nil {
					return "", err
				}
				v = a.Elements[idx]
			}
		default:
			panic("cannot parse query: " + query)
		}
	}

	return p.printValue(v), nil
}

func (p *Parser) printValue(v Value) string {
	res := bytes.Buffer{}

	switch v.GetType() {
	case OBJECT:
		o := v.(Object)
		res.WriteString("{")
		for idx, prop := range o.Properties {
			res.WriteString(fmt.Sprintf("%s: %s", prop.Key.Value, p.printValue(prop.Value)))
			if idx != len(o.Properties)-1 {
				res.WriteString(",")
			}
		}
		res.WriteString("}")
	case LITERAL:
		l := v.(Literal)
		res.WriteString(l.Value)
	case ARRAY:
		a := v.(Array)
		res.WriteString("[")
		for idx, elem := range a.Elements {
			res.WriteString(p.printValue(elem))
			if idx != len(a.Elements)-1 {
				res.WriteString(",")
			}
		}
		res.WriteString("]")

	}

	return res.String()
}
