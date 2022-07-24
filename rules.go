package rulegraph

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

const (
	minProbability = 0
	maxProbability = 1
)

// RulesNode represents node with rules
type RulesNode struct {
	ID              uuid.UUID `json:"id"`
	Rules           []Rule    `json:"rules"`
	SkipProbability float32   `json:"skip_probability"`
}

// NewRulesNode returns a new rules node given the explicit set of params as
// input (rules as a encoding string)
func NewRulesNode(id uuid.UUID, rulesJSON string, skipProbability float32) (RulesNode, error) {
	rulesNode := RulesNode{}

	rules := []Rule{}
	err := json.Unmarshal([]byte(rulesJSON), &rules)
	if err != nil {
		return rulesNode, err
	}

	rulesNode.ID = id
	rulesNode.Rules = rules
	rulesNode.SkipProbability = skipProbability

	return rulesNode, nil
}

// RulesFromString returns a new set of rules from a string or an error, if the string
// is malformed
func RulesFromString(rstring string) ([]RulesNode, error) {
	rules := []RulesNode{}

	err := json.Unmarshal([]byte(rstring), &rules)
	if err != nil {
		return rules, err
	}

	for _, rn := range rules {
		if rn.SkipProbability < minProbability || rn.SkipProbability > maxProbability {
			return []RulesNode{},
				fmt.Errorf("skip_probability should be within [0, 1], got %v", rn.SkipProbability)
		}
	}

	return rules, nil
}

// Evaluate evaluates the current rule against a provided input
func (rn *RulesNode) Evaluate(json []byte) (bool, error) {
	for _, rule := range rn.Rules {
		result, err := rule.isMatch(json)

		if err != nil {
			return false, err
		}

		if !result {
			return false, nil
		}
	}

	// flip the coin based on the rule's skip probability
	rand := randomNumberGenerator()
	if rn.SkipProbability > rand {
		return false, nil
	}

	return true, nil
}

// Rule represents a single rule
type Rule struct {
	Operation string `json:"operation"`
	LeftSide  string `json:"left_side"`
	RightSide string `json:"right_side"`
}

func (r *Rule) isMatch(json []byte) (bool, error) {
	jsonq := gojsonq.New().FromString(string(json))

	valueIf := jsonq.Find(r.LeftSide)

	// if left side value does not exist in the json representation, return false
	// since the rule does not match
	if valueIf == nil {
		return false, nil
	}

	switch r.Operation {
	case "equal":
		value, err := encodeBuffer(valueIf)
		if err != nil {
			return false, err
		}
		rightSide, err := encodeBuffer(r.RightSide)
		if err != nil {
			return false, err
		}

		return bufferIsEqual(value, rightSide), nil

	case "not_equal":
		value, err := encodeBuffer(valueIf)
		if err != nil {
			return false, err
		}

		rightSide, err := encodeBuffer(r.RightSide)
		if err != nil {
			return false, err
		}

		return !bufferIsEqual(value, rightSide), nil

	case "greater_than":
		value, rightSide, err := parseToFloat64(valueIf, r.RightSide)
		if err != nil {
			return false, err
		}

		return value > rightSide, nil

	case "lower_than":
		value, rightSide, err := parseToFloat64(valueIf, r.RightSide)
		if err != nil {
			return false, err
		}

		return value < rightSide, nil
	}

	return true, nil
}
