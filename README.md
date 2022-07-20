# rulegraph


![example workflow](https://github.com/gpestana/rulegraph/actions/workflows/main.yml/badge.svg)
[![GoDoc](http://godoc.org/github.com/gpestana/rulegraph?status.svg)](http://godoc.org/github.com/gpestana/rulegraph)

`rulegraph` is a DSL for defining and evaluating rule sets over arbitrary struct states. The rule sets are encoded as a directional graph of rule nodes that are evaluated against any data structure.

For example, given the following rulegraph:

```json
[
  {
    "id": "d8212113-3732-42de-bc40-be63b89d7864",
    "skip_probability": 0.3,
    "rules": [
      {
       "operation": "larger_than",
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
_WIP_

### API example
_WIP_

```go
package main

import (
  "github.com/gpestana/rulegraph"
)

type Person struct {
  BornAt  time.Time  `json:"born_at"`
  Height  int        `json:"height"`
}

type House struct {
  Size int `json:"size"`
}

type Box struct {
  Person
  House
}

func main() {

 instance := Box{ 
  Person: Person{
    BornAt: time.Now(),
    Height: 120,
   },
   House: House{
    Size: 100,
   },
  }
  
  rg := rulegraph.
  
}
```

### Operations
_WIP_

### Contributing
Fork and PR. Issues for discussion and requests for PR/help.
