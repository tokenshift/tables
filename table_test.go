package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyDieRolls(t *testing.T) {
	assert.Equal(t,
		"testing",
		applyDieRolls("testing"))

	assert.Regexp(t,
		`^\[1d6 = \d+\]$`,
		applyDieRolls("[1d6]"))

	assert.Regexp(t,
		`^\[1d6 \+ 5 = \d+\]$`,
		applyDieRolls("[1d6+5]"))

	assert.Regexp(t,
		`^\[1d6 = \d+\]$`,
		applyDieRolls("[  1d6   ]"))

	assert.Regexp(t,
		`^\[ERR:  1d6  a \]$`,
		applyDieRolls("[  1d6  a ]"))

	assert.Regexp(t,
		`^Testing \[1d6 = \d+\] something \[2d8 \+ 5 = \d+\] here \[ERR:wat\]$`,
		applyDieRolls("Testing [1d6] something [ 2d8 + 5 ] here [wat]"))

	assert.Regexp(t,
		`Escaped \[1d6\] Unescaped \[1d6 = \d+\] Partially escaped \[1d7\]`,
		applyDieRolls(`Escaped \[1d6\] Unescaped [1d6] Partially escaped \[1d7]`))
}

func TestParseColumnName(t *testing.T) {
	col := ParseColumnName("Testing Something")
	assert.Equal(t, col.Name, "Testing Something")
	assert.Len(t, col.Options, 0)

	col = ParseColumnName("Testing [Foo]")
	assert.Equal(t, col.Name, "Testing")
	assert.Len(t, col.Options, 1)
	assert.Contains(t, col.Options, "Foo")

	col = ParseColumnName("Testing [Foo,Bar]")
	assert.Equal(t, col.Name, "Testing")
	assert.Len(t, col.Options, 2)
	assert.Contains(t, col.Options, "Foo")
	assert.Contains(t, col.Options, "Bar")

	col = ParseColumnName("Testing [Foo,Bar=50]")
	assert.Equal(t, col.Name, "Testing")
	assert.Len(t, col.Options, 2)
	assert.Contains(t, col.Options, "Foo")
	assert.Contains(t, col.Options, "Bar")
	assert.Equal(t, col.Options["Bar"], "50")
}
