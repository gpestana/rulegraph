package rulegraph

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRulesNodeMarshalingFromString(t *testing.T) {
	rules := `[{
	"id": "d8212113-3732-42de-bc40-be63b89d7864",
	"skip_probability": 0.3,
	"rules": [
		{
			"operation": "greater_than",
			"left_side": "person.born_at",
			"right_side": "2012-01-01"
		}]
	}
]
`
	// create list of rule nodes from marshaled box
	ruleset, err := RulesFromString(rules)
	assert.NoError(t, err)

	assert.Equal(t, len(ruleset[0].Rules), 1, "number of rules does not match with the expected")
}

func TestRulesNodeMarshalingFromObject(t *testing.T) {
	type Box struct {
		ID              uuid.UUID `json:"id"`
		Rules           string    `json:"rules"`
		SkipProbability float32   `json:"skip_probability"`
		CreatedAt       time.Time `json:"created_at"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	id := uuid.NewV4()
	skipProbability := float32(0.4)
	createdAt := time.Now()
	updatedAt := time.Now()
	rules := `[{
		"operation": "greater_than",
		"left_side": "person.born_at",
		"right_side": "2012-01-01"
	},
	{
		"operation": "equal",
		"left_side": "person.born_at",
		"right_side": "2022-11-10 23:00:00 +0000 UTC"
	}]
	`

	box := Box{id, rules, skipProbability, createdAt, updatedAt}

	rn, err := NewRulesNode(box.ID, box.Rules, box.SkipProbability)
	assert.NoError(t, err)

	assert.Equal(t, rn.ID, id, "IDs do not match")
	assert.Equal(t, rn.SkipProbability, skipProbability, "skip probabilities do not match")
	assert.Equal(t, len(rn.Rules), 2, "number of rules do not match")

	assert.Equal(
		t,
		rn.Rules[0].Operation,
		"greater_than",
		"rule operation does not match the expected",
	)
	assert.Equal(
		t,
		rn.Rules[0].LeftSide,
		"person.born_at",
		"left side does not match the expected",
	)
	assert.Equal(
		t,
		rn.Rules[0].RightSide,
		"2012-01-01",
		"right side does not match the expected",
	)
}
