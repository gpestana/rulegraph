package rulegraph

import (
	uuid "github.com/satori/go.uuid"
)

// RuleGraph represents the graph of rules to evaluate
type RuleGraph struct {
	ruleNodes []RulesNode
}

// NewRuleGraph returns a new RuleGraph instance without any set of rules
func NewRuleGraph() *RuleGraph {
	return &RuleGraph{
		ruleNodes: []RulesNode{},
	}
}

// NewRuleGraphWith returns a new RuleGraph instance with a predefined ruleset
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

// LoadRulesFromString replaces the current rules with rules encoded as a
// string, or return an error if the input is malformed
func (rg *RuleGraph) LoadRulesFromString(rstring string) error {
	rules, err := RulesFromString(rstring)
	if err != nil {
		return err
	}

	rg.ruleNodes = rules

	return nil
}

// IsRulesetEmpty returns true if the rule set of the instance is empty
func (rg RuleGraph) IsRulesetEmpty() bool {
	return len(rg.ruleNodes) == 0
}
