package rulegraph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func extractValueNotNull(mp map[string]interface{}) interface{} {
	delete(mp, "Valid")
	keys := reflect.ValueOf(mp).MapKeys()
	return mp[keys[0].String()]
}

func parseToString(valueIf interface{}, rightSideStr string) (string, string, error) {
	valueStr, ok := valueIf.(string)
	if ok {
		return valueStr, rightSideStr, nil
	}
	valueStr = strconv.FormatFloat(valueIf.(float64), 'f', -1, 64)

	return valueStr, rightSideStr, nil
}

func parseToFloat64(valueIf interface{}, rightSideStr string) (float64, float64, error) {

	v, r, err := tryParseNumber(valueIf, rightSideStr)
	if err == nil {
		// return values if parsing was successful
		return v, r, nil
	}

	v, r, err = tryParseDate(valueIf, rightSideStr)
	if err == nil {
		// return values if parsing was successful
		return v, r, nil
	}

	return tryParseString(valueIf, rightSideStr)
}

func tryParseNumber(valueIf interface{}, rightSideStr string) (float64, float64, error) {
	value, ok := valueIf.(float64)
	if !ok {
		return 0, 0, fmt.Errorf("Error parsing value (%s) to number", valueIf)
	}

	rightSide, err := strconv.ParseFloat(rightSideStr, 64)

	return value, rightSide, err
}

func tryParseDate(valueIf interface{}, rightSideStr string) (float64, float64, error) {
	valueDate, ok := valueIf.(time.Time)
	if !ok {
		v, err := time.Parse(time.RFC3339, valueIf.(string))
		if err != nil {
			return 0, 0, fmt.Errorf("Error parsing value (%s) to Date", valueIf)
		}
		valueDate = v
	}

	rightSideDate, err := time.Parse(time.RFC3339, rightSideStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing rule right side (%s) to Date", rightSideStr)
	}

	return float64(valueDate.UnixNano()), float64(rightSideDate.UnixNano()), nil
}

func tryParseString(valueIf interface{}, rightSideStr string) (float64, float64, error) {
	value, err := strconv.ParseFloat(valueIf.(string), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing value (%s) to string", valueIf)
	}

	rightSide, err := strconv.ParseFloat(rightSideStr, 64)

	return value, rightSide, err
}

//nolint
func encodeBuffer(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//nolint
func bufferIsEqual(x []byte, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i, b := range x {
		if b != y[i] {
			return false
		}
	}
	return true
}

func randomNumberGenerator() float32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float32()
}
