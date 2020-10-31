package dice

import (
	"math/rand"
	"testing"
	"time"
)

const TestRolls = 1000

func testRollResults(t *testing.T, diceSpec string, sumLow, sumHigh int, constraints ...[][]uint) bool {
	parsed, err := Parse(diceSpec)
	if err != nil {
		t.Error(err)
		return false
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for roll := 0; roll < TestRolls; roll++ {
		results := parsed.Roll(r)

		if results.Sum() < sumLow {
			t.Errorf("%s: Rolled %d (%s), which is < %d", diceSpec, results.Sum(), results, sumLow)
			return false
		}

		if results.Sum() > sumHigh {
			t.Errorf("%s: Rolled %d (%s), which is > %d", diceSpec, results.Sum(), results, sumHigh)
			return false
		}

		if len(constraints) == 0 {
			continue
		}

		if len(results) != len(constraints) {
			t.Errorf("Expected %d result blocks (%v), got %d (%v)", len(constraints), constraints, len(results), results)
			return false
		}

		for i, result := range results {
			if len(result.Values) != len(constraints[i]) {
				t.Errorf("Expected %d results (%v), got %d (%v)", len(constraints[i]), constraints[i], len(result.Values), result.Values)
				return false
			}

			for j, val := range result.Values {
				if len(constraints[i][j]) != 2 {
					t.Errorf("Expected a low/high pair, got %v", constraints[i][j])
					return false
				}

				if constraints[i][j][0] > constraints[i][j][1] {
					t.Errorf("`low` can't be greater than `high`")
					return false
				}

				if val < constraints[i][j][0] || val > constraints[i][j][1] {
					t.Errorf("%s: Expected %d <= %d <= %d", diceSpec, constraints[i][j][0], val, constraints[i][j][1])
					return false
				}
			}
		}
	}

	return true
}

func TestConstants(t *testing.T) {
	testRollResults(t, "1", 1, 1, [][]uint{[]uint{1, 1}})
	testRollResults(t, "6", 6, 6, [][]uint{[]uint{6, 6}})
	testRollResults(t, "1634615", 1634615, 1634615, [][]uint{[]uint{1634615, 1634615}})
}

func TestSimpleRolls(t *testing.T) {
	testRollResults(t, "1d6", 1, 6, [][]uint{[]uint{1, 6}})
	testRollResults(t, "2d6", 2, 12, [][]uint{[]uint{1, 6}, []uint{1, 6}})
	testRollResults(t, "3d6", 3, 18, [][]uint{[]uint{1, 6}, []uint{1, 6}, []uint{1, 6}})
	testRollResults(t, "4d6", 4, 24, [][]uint{[]uint{1, 6}, []uint{1, 6}, []uint{1, 6}, []uint{1, 6}})
	testRollResults(t, "5d6", 5, 30, [][]uint{[]uint{1, 6}, []uint{1, 6}, []uint{1, 6}, []uint{1, 6}, []uint{1, 6}})

	testRollResults(t, "1d20", 1, 20, [][]uint{[]uint{1, 20}})
	testRollResults(t, "2d20", 2, 40, [][]uint{[]uint{1, 20}, []uint{1, 20}})

	testRollResults(t, "15213d53115", 15213, 808038495)
}

func TestMultipleRolls(t *testing.T) {
	testRollResults(t,
		"1d20+2d6+4d4", 7, 20+6+6+4+4+4+4,
		[][]uint{[]uint{1, 20}},
		[][]uint{[]uint{1, 6}, []uint{1, 6}},
		[][]uint{[]uint{1, 4}, []uint{1, 4}, []uint{1, 4}, []uint{1, 4}})
}

func TestMultiplePlusAndMinusRolls(t *testing.T) {
	testRollResults(t, "1d20-2d6+4d4", 1-6-6+4, 20-1-1+4+4+4+4,
		[][]uint{[]uint{1, 20}},
		[][]uint{[]uint{1, 6}, []uint{1, 6}},
		[][]uint{[]uint{1, 4}, []uint{1, 4}, []uint{1, 4}, []uint{1, 4}})
}

func TestMultipleRollsWithWhitespace(t *testing.T) {
	testRollResults(t, "  1 d 20    - 2 d 6   + 4    d\t 4 ", 1-6-6+4, 20-1-1+4+4+4+4,
		[][]uint{[]uint{1, 20}},
		[][]uint{[]uint{1, 6}, []uint{1, 6}},
		[][]uint{[]uint{1, 4}, []uint{1, 4}, []uint{1, 4}, []uint{1, 4}})
}

func TestInvalidDiceSpecs(t *testing.T) {
	_, err := Parse("1dd6")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
