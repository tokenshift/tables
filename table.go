package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tokenshift/tables/dice"
)

type Table struct {
	Columns Columns
	Rows    Rows
}

type Columns []Column

type Rows []Row

type Column struct {
	Name    string
	Options map[string]string
}

type Row []string

type RowMap map[string]string

func (cols Columns) Names() []string {
	result := make([]string, len(cols))

	for c, col := range cols {
		result[c] = col.Name
	}

	return result
}

var rxColumnName = regexp.MustCompile(`^(.*?)(?:\[(.*)\])?$`)

func NewColumn(name string) Column {
	return Column{
		Name:    strings.TrimSpace(name),
		Options: make(map[string]string),
	}
}

func ParseColumnName(input string) (c Column) {
	c = NewColumn(input)

	match := rxColumnName.FindStringSubmatch(c.Name)
	if match == nil {
		return
	}

	c.Name = strings.TrimSpace(match[1])

	for _, option := range strings.Split(match[2], ",") {
		if option != "" {
			split := strings.SplitN(option, "=", 2)
			if len(split) > 1 {
				c.Options[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
			} else {
				c.Options[strings.TrimSpace(split[0])] = ""
			}
		}
	}

	return
}

func LoadFile(filename string) (t Table, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return t, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return t, err
	}

	t.Columns = make([]Column, len(data[0]))

	for c, name := range data[0] {
		t.Columns[c] = ParseColumnName(name)
	}

	t.Rows = make([]Row, len(data[1:]))

	for r, row := range data[1:] {
		// "Extend" the column headers to match the number of cells in the row, in
		// case the input CSV was poorly formed (i.e. doesn't have the same number
		// of cells in each row, or row lengths don't match the number of column
		// headers)
		for c := len(t.Columns); c < len(row); c++ {
			t.Columns = append(t.Columns, NewColumn(fmt.Sprintf("Column %d", c+1)))
		}

		// "Extend" the row to match the number of column headers. See above.
		for i := len(row); i < len(t.Columns); i++ {
			row = append(row, "")
		}

		t.Rows[r] = row
	}

	return t, nil
}

func (orig Table) Filter(filters map[string]string) (copy Table) {
	copy.Columns = orig.Columns
	copy.Rows = make([]Row, 0, len(orig.Rows))

	for _, row := range orig.Rows {
		if orig.RowMap(row).Matches(filters) {
			copy.Rows = append(copy.Rows, row)
		}
	}

	return copy
}

func (table Table) RandomRow() Row {
	rowNum := Rand.Intn(len(table.Rows))

	orig := table.Rows[rowNum]
	result := make([]string, len(orig))

	for c, val := range orig {
		if ind, ok := table.Columns[c].Options["Independent"]; ok && strings.ToLower(ind) != "false" {
			result[c] = table.RandomColumnValue(c)
		} else {
			result[c] = val
		}
	}

	return result
}

func (table Table) ColumnValues(col int) []string {
	vals := make([]string, 0, len(table.Rows))

	for _, row := range table.Rows {
		val := row[col]
		if val != "" {
			vals = append(vals, val)
		}
	}

	return vals
}

func (table Table) RandomColumnValue(col int) string {
	if percent, ok := table.Columns[col].Options["Percent"]; ok {
		p, err := strconv.ParseInt(percent, 10, 0)
		if err != nil {
			return fmt.Sprintf("ERROR: Percentage=%s", percent)
		}

		if Rand.Intn(100) >= int(p) {
			return ""
		}
	}

	vals := table.ColumnValues(col)
	rowNum := Rand.Intn(len(vals))
	return vals[rowNum]
}

func (t Table) RowMap(row Row) RowMap {
	result := make(map[string]string)

	for c, val := range row {
		result[t.Columns[c].Name] = val
	}

	return result
}

func (r Row) String() string {
	var b bytes.Buffer

	for i, val := range r {
		if i > 0 && val != "" {
			fmt.Fprint(&b, " ")
		}

		fmt.Fprint(&b, val)
	}

	return b.String()
}

func (row RowMap) Matches(filters map[string]string) bool {
	include := true

	for key, vals := range filters {
		matched := false

		for _, val := range strings.Split(vals, ",") {
			if row[key] == val {
				matched = true
				break
			}
		}

		include = include && matched
	}

	return include
}

var rxDieRolls = regexp.MustCompile(`\[.*?\]`)

func (orig Row) ApplyDieRolls(r *rand.Rand) Row {
	copy := make([]string, len(orig))

	for i, val := range orig {
		copy[i] = applyDieRolls(val)
	}

	return copy
}

func applyDieRolls(input string) string {
	matches := rxDieRolls.FindAllStringSubmatchIndex(input, -1)

	if matches == nil {
		return input
	}

	var result bytes.Buffer

	offset := 0
	for i, match := range matches {
		// Check if the match was escaped with a backslash.
		if match[0] > 0 && input[match[0]-1] == '\\' {
			// Append everything up to the start of the match.
			fmt.Fprint(&result, input[offset:match[0]-1])

			// Append the match itself, minus the backslashes.
			if input[match[1]-2] == '\\' {
				fmt.Fprintf(&result, "[%s]", input[match[0]+1:match[1]-2])
			} else {
				fmt.Fprintf(&result, "[%s]", input[match[0]+1:match[1]-1])
			}
		} else {
			// Append everything up to the start of the match.
			fmt.Fprint(&result, input[offset:match[0]])

			// Get and validate the dice spec (between square brackets).
			// Append new square brackets with the result of the die roll,
			// or an error message if it couldn't be parsed.
			spec := input[match[0]+1 : match[1]-1]
			parsed, err := dice.Parse(spec)

			if err == nil {
				fmt.Fprintf(
					&result,
					"[%s = %d]",
					parsed,
					parsed.Roll(Rand).Sum())
			} else {
				fmt.Fprintf(
					&result,
					"[ERR:%s]",
					spec)
			}
		}

		// Append everything after the end of the last match.
		if i == len(matches)-1 {
			fmt.Fprint(&result, input[match[1]:])
		}

		offset = match[1]
	}

	return result.String()
}

func (orig Table) Omit(columnNames []string) (copy Table) {
	include := make(map[int]bool, len(orig.Columns))

	for c, col := range orig.Columns {
		include[c] = true

		for _, omitted := range columnNames {
			include[c] = include[c] && strings.TrimSpace(col.Name) != strings.TrimSpace(omitted)
		}
	}

	copy.Columns = make([]Column, 0, len(orig.Columns))
	copy.Rows = make([]Row, len(orig.Rows))

	for c, col := range orig.Columns {
		if include[c] {
			copy.Columns = append(copy.Columns, col)
		}
	}

	for r, row := range orig.Rows {
		newRow := make([]string, 0, len(copy.Columns))

		for c, val := range row {
			if include[c] {
				newRow = append(newRow, val)
			}
		}

		copy.Rows[r] = newRow
	}

	return copy
}
