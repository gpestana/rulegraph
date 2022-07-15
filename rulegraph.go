package rulegraph

import (
	uuid "github.com/satori/go.uuid"
)

// RuleGraph represents the graph of rules to evaluate
type RuleGraph struct {
	ruleNodes []RulesNode
}

// NewRuleGraphWith returns a new RuleGraph with a predefined ruleset
func NewRuleGraphWith(ruleset []RulesNode) *RuleGraph {
	return &RuleGraph{
		ruleset,
	}
}

// Evaluate runs all the graph rules against the input and returns a list of
// rule IDs that evaluated as true against the provided input
func (rg *RuleGraph) Evaluate(json []byte) ([]uuid.UUID, error) {
	idsEvalTrue := []uuid.UUID{}

	// evaluate rule for all nodes
	for _, rnode := range rg.ruleNodes {
		result, err := rnode.Evaluate(json)
		if err != nil {
			return idsEvalTrue, err
		}

		if result {
			idsEvalTrue = append(idsEvalTrue, rnode.ID)
		}
	}

	return idsEvalTrue, nil
}
