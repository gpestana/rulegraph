package rulegraph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
)

func parseToInt64(valueIf interface{}, rightSideStr string) (float64, float64, error) {
	value, ok := valueIf.(float64)
	if !ok {
		return 0, 0, fmt.Errorf("Error casting value (%s) to int64", valueIf)
	}
	rightSide, err := strconv.ParseFloat(rightSideStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing rule right side (%s) to float64", rightSideStr)
	}

	return value, rightSide, nil
}

func encodeBuffer(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

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
