package test

import (
	"testing"
)

// TestTestData ...
func TestTestData(t *testing.T) {
	// t.SkipNow()
	original := LoadTestData()
	originalForID := MakePersonForIDMap(original)
	originalBytes := MarshalPeople(original)
	parsed := UnmarshalPeople(originalBytes)
	parsedForID := MakePersonForIDMap(parsed)
	for id, person := range originalForID {
		readPerson, ok := parsedForID[id]
		if !ok {
			t.FailNow()
		}
		if !person.Equals(&readPerson) {
			t.FailNow()
		}
	}
}
