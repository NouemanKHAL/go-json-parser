package main

import (
	"bytes"
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

func PrintValue(v Value) string {
	res := bytes.Buffer{}

	switch v.GetType() {
	case OBJECT:
		o := v.(Object)
		res.WriteString("{")
		for idx, prop := range o.Properties {
			res.WriteString(fmt.Sprintf("%s: %s", prop.Key.Value, PrintValue(prop.Value)))
			if idx != len(o.Properties)-1 {
				res.WriteString(",")
			}
		}
		res.WriteString("}")
	case LITERAL:
		l := v.(Literal)
		res.WriteString(fmt.Sprintf("%s", l.Value))
	case ARRAY:
		a := v.(Array)
		res.WriteString("[")
		for idx, elem := range a.Elements {
			res.WriteString(fmt.Sprintf("%s", PrintValue(elem)))
			if idx != len(a.Elements)-1 {
				res.WriteString(",")
			}
		}
		res.WriteString("]")

	}

	return res.String()
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

		cur := p.tree.root
		fmt.Printf("%s\n", PrintValue(cur.Value))

	} else {
		handleError(fmt.Errorf("expected input file in first argument"))
	}

}
