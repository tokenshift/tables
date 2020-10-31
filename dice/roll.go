package dice

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Results []Result

type Result struct {
	Operator string
	Values   []uint
}

func (spec DiceSpec) Roll(r *rand.Rand) Results {
	clauses := spec.Clauses()
	results := make([]Result, len(clauses))

	for i, clause := range clauses {
		results[i] = clause.Roll(r)
	}

	return results
}

func (clause OpClause) Roll(r *rand.Rand) Result {
	if clause.Clause.Die < 1 {
		return Result{
			Operator: clause.Operator,
			Values:   []uint{clause.Clause.Count},
		}
	} else {
		rolls := make([]uint, clause.Clause.Count)
		for i := uint(0); i < clause.Clause.Count; i++ {
			rolls[i] = uint(r.Intn(int(clause.Clause.Die))) + 1
		}
		return Result{
			Operator: clause.Operator,
			Values:   rolls,
		}
	}
}

func (results Results) Sum() int {
	sum := 0
	for _, result := range results {
		for _, r := range result.Values {
			if result.Operator == "+" {
				sum += int(r)
			} else {
				sum -= int(r)
			}
		}
	}
	return sum
}

func (results Results) String() string {
	var b bytes.Buffer

	for i, result := range results {
		if i > 0 {
			fmt.Fprintf(&b, " %s ", result.Operator)
		}

		fmt.Fprint(&b, "(")

		for j, r := range result.Values {
			if j > 0 {
				fmt.Fprint(&b, " + ")
			}

			fmt.Fprint(&b, r)
		}

		fmt.Fprint(&b, ")")
	}

	return b.String()
}
