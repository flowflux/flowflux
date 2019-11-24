package test

import "testing"

// CheckTestSet ...
func CheckTestSet(t *testing.T, originalForID map[string]Person, b []byte) *Person {
	parsed, err := UnmarshalPerson(b)
	if err != nil {
		t.Fatalf("Couldn't unmarshall person from: \"%v\" -> %v\n", string(b), err.Error())
	}
	original := originalForID[parsed.ID]
	if !original.Equals(parsed) {
		t.Fatal("Parsed data doesn't conform to original")
	}
	return parsed
}
