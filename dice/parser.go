package dice

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/regex"
)

type DiceSpec struct {
	Head *MaybeOpClause `@@`
	Tail []*OpClause    `@@*`
}

func (spec DiceSpec) Clauses() []*OpClause {
	if spec.Head.Operator == "" {
		head := OpClause{
			Operator: "+",
			Clause:   spec.Head.Clause,
		}

		return append([]*OpClause{&head}, spec.Tail...)
	} else {
		head := OpClause{
			Operator: spec.Head.Operator,
			Clause:   spec.Head.Clause,
		}

		return append([]*OpClause{&head}, spec.Tail...)
	}
}

func (spec DiceSpec) String() string {
	var b bytes.Buffer
	fmt.Fprint(&b, spec.Head)

	for _, t := range spec.Tail {
		fmt.Fprintf(&b, " %s", t)
	}

	return b.String()
}

type MaybeOpClause struct {
	Operator string `@Op?`
	Clause   Clause `@@`
}

func (c MaybeOpClause) String() string {
	if c.Operator == "" {
		return fmt.Sprint(c.Clause)
	} else {
		return fmt.Sprintf("%s %s", c.Operator, c.Clause)
	}
}

type OpClause struct {
	Operator string `@Op`
	Clause   Clause `@@`
}

func (c OpClause) String() string {
	return fmt.Sprintf("%s %s", c.Operator, c.Clause)
}

type Clause struct {
	Count uint `@Int`
	Die   uint `(D @Int)?`
}

func (c Clause) String() string {
	if c.Die < 1 {
		return fmt.Sprint(c.Count)
	} else {
		return fmt.Sprintf("%dd%d", c.Count, c.Die)
	}
}

func Parse(input string) (DiceSpec, error) {
	lexer, err := regex.New(`
Int = \d+
D = d
Op = \+|\-
Whitespace = \s+
`)

	if err != nil {
		return DiceSpec{}, err
	}

	parser, err := participle.Build(
		&DiceSpec{},
		participle.Lexer(lexer),
		participle.Elide("Whitespace"))
	if err != nil {
		return DiceSpec{}, err
	}

	var spec DiceSpec
	err = parser.ParseString(input, &spec)
	if err != nil {
		return DiceSpec{}, err
	}

	return spec, nil
}
