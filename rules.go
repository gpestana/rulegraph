package rulegraph

import (
	uuid "github.com/satori/go.uuid"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

// RulesNode represents node with rules
type RulesNode struct {
	ID    uuid.UUID `json:"id"`
	Rules []Rule    `json:"rules"`
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

	case "larger_than":
		value, rightSide, err := parseToInt64(valueIf, r.RightSide)
		if err != nil {
			return false, err
		}

		return value > rightSide, nil

	case "lower_than":
		value, rightSide, err := parseToInt64(valueIf, r.RightSide)
		if err != nil {
			return false, err
		}

		return value < rightSide, nil
	}

	return true, nil
}
