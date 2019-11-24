package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// CheckUpperExtraction ...
func CheckUpperExtraction(ob []byte, originalForID map[string]Person) error {
	var extraction map[string]string
	err := json.Unmarshal(ob, &extraction)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshall person from: \"%v\" -> %v", string(ob), err.Error())
	}

	original := originalForID[extraction["id"]]
	return checkExtractedPerson(original, extraction)
}

func catchToFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func checkExtractedPerson(original Person, extraction map[string]string) (err error) {
	extractedValue, ok := extraction["first_name"]
	if ok {
		upperValue := strings.ToUpper(original.FirstName)
		if upperValue != extractedValue {
			err = fmt.Errorf("Extracted value %v not as expected: %v", extractedValue, upperValue)
		}
	} else {
		extractedValue, ok = extraction["last_name"]
		if ok {
			upperValue := strings.ToUpper(original.LastName)
			if upperValue != extractedValue {
				err = fmt.Errorf("Extracted value %v not as expected: %v", extractedValue, upperValue)
			}
		} else {
			err = fmt.Errorf("Extraction not as expected: %v", extraction)
		}
	}
	return
}
