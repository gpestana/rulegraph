package rulegraph

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// TestObj is a struct used as test input
type TestObj struct {
	User  TestUser  `json:"user"`
	House TestHouse `json:"house"`
}

// TestUser is for testing purposes only
type TestUser struct {
	Name         string       `json:"name"`
	Age          int          `json:"age"`
	Permission   bool         `json:"permission"`
	RegisteredAt sql.NullTime `json:"registered_at"`
	LeftAt       time.Time    `json:"left_at"`
}

// TestHouse is for testing purposes only
type TestHouse struct {
	ID      uuid.UUID `json:"id"`
	Size    int       `json:"size"`
	Price   int       `json:"price"`
	BuiltAt time.Time `json:"built_at"`
}

func objFactory(test string) ([]byte, error) {
	obj := TestObj{}

	switch test {
	case "testCase1":
		regAt, _ := time.Parse(time.RFC3339, "2010-11-01T02:04:05Z")
		builtAt, _ := time.Parse(time.RFC3339, "2011-01-02T15:04:05Z")

		obj = TestObj{
			User: TestUser{
				Name:         "Pete",
				Age:          13,
				Permission:   true,
				RegisteredAt: sql.NullTime{Time: regAt, Valid: true},
			},
			House: TestHouse{
				ID:      uuid.NewV4(),
				Size:    15,
				Price:   100,
				BuiltAt: builtAt,
			},
		}

	case "testCase2":
		obj = TestObj{
			User: TestUser{
				Name:       "Andrew",
				Age:        13,
				Permission: true,
			},
			House: TestHouse{
				ID:    uuid.NewV4(),
				Size:  15,
				Price: 100,
			},
		}

	case "testCase3":
		builtAt, _ := time.Parse(time.RFC3339, "2022-01-02T15:04:05Z")

		obj = TestObj{
			User: TestUser{},
			House: TestHouse{
				BuiltAt: builtAt,
			},
		}

	case "testCase4":
		leftAt := time.Date(2011, 01, 2, 15, 4, 5, 0, time.UTC).UTC()

		obj = TestObj{
			User: TestUser{
				LeftAt: leftAt,
			},
			House: TestHouse{
				Size: 10,
			},
		}

	}

	return json.Marshal(&obj)
}

func TestRuleGraphSimpleCases(t *testing.T) {
	rulesNode1Str := `{ 
	"id": "1117b810-9dad-11d1-80b4-00c04fd11111",
	"rules": [
		{ "operation": "equal", "left_side": "user.name", "right_side": "Pete" }
	]
}
`

	rulesNode2Str := `{
	"id": "2227b810-9dad-11d1-80b4-00c04fd22222",
	"rules": [
		{ "operation": "greater_than", "left_side": "user.age", "right_side": "10" },
		{ "operation": "lower_than", "left_side": "house.size", "right_side": "30" }
	]
}
`

	// Test 1
	testObj, err := objFactory("testCase1")
	assert.NoError(t, err)

	expectedMatchIDs := []uuid.UUID{
		uuid.FromStringOrNil("1117b810-9dad-11d1-80b4-00c04fd11111"),
		uuid.FromStringOrNil("2227b810-9dad-11d1-80b4-00c04fd22222"),
	}

	var rulesNode1 RulesNode
	err = json.Unmarshal([]byte(rulesNode1Str), &rulesNode1)
	if err != nil {
		t.Error("Error: ", err)
	}

	var rulesNode2 RulesNode
	err = json.Unmarshal([]byte(rulesNode2Str), &rulesNode2)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode1, rulesNode2})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval test1 matches both rules")

	// Test 2
	testObj, err = objFactory("testCase2")
	assert.NoError(t, err)

	expectedMatchIDs = []uuid.UUID{
		uuid.FromStringOrNil("2227b810-9dad-11d1-80b4-00c04fd22222"),
	}

	err = json.Unmarshal([]byte(rulesNode1Str), &rulesNode1)
	if err != nil {
		t.Error("Error: ", err)
	}

	err = json.Unmarshal([]byte(rulesNode2Str), &rulesNode2)
	if err != nil {
		t.Error("Error: ", err)
	}

	matchIDs, err = rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval test2 matches both rules")
}

func TestRuleGraphInequalities(t *testing.T) {
	rulesNodeStr := `{
	"id": "2227b810-9dad-11d1-80b4-00c04fd33333",
	"rules": [
		{ "operation": "greater_than", "left_side": "user.age", "right_side": "10" },
		{ "operation": "lower_than", "left_side": "user.age", "right_side": "30" }
	]
}
`
	// Test 1
	testObj, err := objFactory("testCase1")
	assert.NoError(t, err)

	expectedMatchIDs := []uuid.UUID{
		uuid.FromStringOrNil("2227b810-9dad-11d1-80b4-00c04fd33333"),
	}

	var rulesNode RulesNode
	err = json.Unmarshal([]byte(rulesNodeStr), &rulesNode)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval inequality should match struct state")

	// Test 2
	rulesNodeStr = `{
	"id": "2227b810-9dad-11d1-80b4-00c04fd33333",
	"rules": [
		{ "operation": "greater_than", "left_side": "user.age", "right_side": "30" },
		{ "operation": "lower_than", "left_side": "user.age", "right_side": "50" }
	]
}`

	testObj, err = objFactory("testCase1")
	assert.NoError(t, err)

	expectedMatchIDs = []uuid.UUID{}

	err = json.Unmarshal([]byte(rulesNodeStr), &rulesNode)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph = NewRuleGraphWith([]RulesNode{rulesNode})

	matchIDs, err = rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval inequality should not match struct state")
}

func TestRuleGraphAnd(t *testing.T) {
	rulesNodeStr := `{
	"id": "2227b810-9dad-11d1-80b4-00c04fd33333",
	"rules": [
		{ "operation": "lower_than", "left_side": "user.age", "right_side": "30" },
		{ "operation": "equal", "left_side": "user.age", "right_side": "15" }
	]
}
`
	// Test 1
	testObj, err := objFactory("testCase1")
	assert.NoError(t, err)

	expectedMatchIDs := []uuid.UUID{}

	var rulesNode RulesNode
	err = json.Unmarshal([]byte(rulesNodeStr), &rulesNode)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval should not match struct state")
}

func TestRuleInequalityTimestamp(t *testing.T) {
	rulesNodeStr := `{
	"id": "2227b810-9dad-11d1-80b4-00c04fd33333",
	"rules": [
		{ "operation": "greater_than", "left_side": "house.built_at", "right_side": "2021-01-02T15:04:05Z" },
		{ "operation": "lower_than", "left_side": "house.built_at", "right_side": "2023-01-02T15:04:05Z" }
	]
}
`
	// Test 3 (built_at "2022-01-02T15:04:05Z)
	testObj, err := objFactory("testCase3")
	assert.NoError(t, err)

	matchID, _ := uuid.FromString("2227b810-9dad-11d1-80b4-00c04fd33333")
	expectedMatchIDs := []uuid.UUID{matchID}

	var rulesNode RulesNode
	err = json.Unmarshal([]byte(rulesNodeStr), &rulesNode)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval should not match struct state")

	// Test with failing built_at timestamp ("2011-01-02T15:04:05Z")
	testObj, err = objFactory("testCase1")
	assert.NoError(t, err)

	// no match for testCase1
	expectedMatchIDs = []uuid.UUID{}

	matchIDs, err = rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval should not match struct state")
}

func TestRuleGraphNullTime(t *testing.T) {
	rulesNode1Str := `{ 
	"id": "1117b810-9dad-11d1-80b4-00c04fd11111",
	"rules": [
		{ "operation": "equal", "left_side": "user.registered_at", "right_side": "2010-11-01T02:04:05Z" }
	]
}
`
	// Test 1
	testObj, err := objFactory("testCase1")
	assert.NoError(t, err)

	expectedMatchIDs := []uuid.UUID{
		uuid.FromStringOrNil("1117b810-9dad-11d1-80b4-00c04fd11111"),
	}

	var rulesNode1 RulesNode
	err = json.Unmarshal([]byte(rulesNode1Str), &rulesNode1)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode1})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval test1 should match rule")
}

func TestRuleTimeInequality(t *testing.T) {
	rulesNode1Str := `{ 
	"id": "1117b810-9dad-11d1-80b4-00c04fd11112",
	"rules": [
		{ "operation": "equal", "left_side": "house.size", "right_side": "10" },
		{ "operation": "lower_than", "left_side": "user.left_at", "right_side": "2020-01-02T15:04:05Z" }
	]
}
`
	// Test 4
	testObj, err := objFactory("testCase4")
	assert.NoError(t, err)

	fmt.Println("Test Object::")
	fmt.Println(string(testObj))
	fmt.Println("Rules::")
	fmt.Println(rulesNode1Str)

	expectedMatchIDs := []uuid.UUID{
		uuid.FromStringOrNil("1117b810-9dad-11d1-80b4-00c04fd11112"),
	}

	var rulesNode1 RulesNode
	err = json.Unmarshal([]byte(rulesNode1Str), &rulesNode1)
	if err != nil {
		t.Error("Error: ", err)
	}

	rgraph := NewRuleGraphWith([]RulesNode{rulesNode1})

	matchIDs, err := rgraph.Evaluate(testObj)
	assert.NoError(t, err)

	assert.Equal(t, expectedMatchIDs, matchIDs, "Result of eval test4 should match rule")
}
