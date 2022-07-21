# rulegraph


![example workflow](https://github.com/gpestana/rulegraph/actions/workflows/main.yml/badge.svg)
[![GoDoc](http://godoc.org/github.com/gpestana/rulegraph?status.svg)](http://godoc.org/github.com/gpestana/rulegraph)

`rulegraph` is a DSL for defining and evaluating rule sets over arbitrary struct states. The rule sets are encoded as a directional graph of rule nodes that are evaluated against any data structure.

For example, given the following rulegraph configuration:

```json
[
  {
    "id": "d8212113-3732-42de-bc40-be63b89d7864",
    "skip_probability": 0.3,
    "rules": [
      {
       "operation": "greater_than",
       "left_side": "person.born_at",
       "right_side": "2012-01-01",
      },
      {
       "operation": "equal",
       "left_side": "house.size",
       "right_side": "100",
      },
    ],
  },
  {
    "id": "14f68d22-0e79-4129-b775-a49d25dd025c",
    "skip_probability": 0.6,
    "rules": [
      {
       "operation": "lower_than",
       "left_side": "person.height",
       "right_side": "50",
      },
    ]
  }
]
```
Example 1. Ruleset in JSON format

The above rulegraph will evaluated to `true` IFF the input object has `person.born_at > 2012-01-01` AND `house.size > 100` OR `person.height < 50`. In addition, each ruleset match will be skipped (i.e. considered `false`, even if the rule matched) with a probability of `skip_probability`.

### Rulegraph structure 

![rulegraph](https://github.com/gpestana/rulegraph/blob/main/_docs/rulegraph.jpg)

Figure 1. `rulegraph` tree example

Figure 1. shows the visual representation of a rule graph. The rule graph contains 2 `RuleNodes`, each with three `Rules`. The result of the evaluation of the `rulegraph` represented in the Figure 1. is the following:

```
 match = RuleNodeA OR RuleNodeB <=> 
   (RuleNodeA_R1 AND RuleNodeA_R2 AND RuleNodeA_R3) OR (RuleNodeB_R1 AND RuleNodeB_R2 AND RuleNodeB_R3)
``` 

When calling the `rulegraph.Evaluate(obj)`, the rules loaded in the rule graph instance will be evaluated against the `object` state and. The evaluate function returns a list of `RuleNode` IDs that have been matched.

### API example

```go
package main

import (
  "encoding/json"

  "github.com/gpestana/rulegraph"
  uuid "github.com/satori/go.uuid"
)

func main() {
 ruleset := []RulesNode{}

 ruleID := uuid.NewV4()
 skipProbability := 0.4
 rules := `[{
  "operation": "greater_than",
  "left_side": "person.born_at",
  "right_side": "2012-01-01"
 },{
  "operation": "equal",
  "left_side": "person.born_at",
  "right_side": "2012-01-01"
 }]
`

 // create rule nodes
 r1 := rulegraph.NewRulesNode(ruleID, rules, skipProbability)
 // r2 := NewRulesNode(...)

 // init rulegraph with rulesets
 rg := rulegraph.NewRuleGraphWith([]RuleNodes{r1, r2})

 // get object input to evaluate against the rules set
 obj := CreateAndPopulateObj()
 objJSON, _ := json.Marshal(obj)

 matchedIDs, err := rg.Evaluate(objJSON)

 // matchedIDs will contain all the RuleNode IDs that matched with
 // the object state, given the skipProbability defined
}
```

### Operations

Operations are applied between the value associated with the input `left_side` value and the rule `right_side` parameter. Currently `rulegraph` supports the following operations over any type: 

- `equal`: Implements a `A == B` comparison
- `not_equal` Implement a `A !== B` comparison
- `greater_than`: Implements an inequality `A > B` comparison
- `lower_than` Implements an inequality `A < B` comparison

### Contributing
Fork and PR. Issues for discussion and requests for PR/help.
