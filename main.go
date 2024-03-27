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

func PrintValue(value Value) string {
	res := bytes.Buffer{}

	q := []Value{value}

	for len(q) > 0 {
		levelSize := len(q)

		newQ := []Value{}

		for i := 0; i < levelSize; i += 1 {
			v := q[i]
			switch v.GetType() {
			case OBJECT:
				o := v.(Object)
				for idx, prop := range o.Properties {
					res.WriteString(fmt.Sprintf("{ \"Property\": {\"Key\": %s, \"Value\":%s} } ", prop.Key.Value, PrintValue(prop.Value)))
					if idx != len(o.Properties)-1 {
						res.WriteString(",")
					}
				}
			case LITERAL:
				l := v.(Literal)
				res.WriteString(fmt.Sprintf("{\"Literal\": %s} ", l.Value))
			}
		}

		q = newQ
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
		fmt.Printf("[%s]\n", PrintValue(cur.Value))

	} else {
		handleError(fmt.Errorf("expected input file in first argument"))
	}

}
