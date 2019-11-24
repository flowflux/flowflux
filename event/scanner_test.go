package event

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

// TestScanner ...
func TestScanner(t *testing.T) {
	// t.SkipNow()
	originals := test.LoadTestData() // [:10]
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	scanner := NewScanner(reader)
	count := 0

	for scanner.Scan() {
		b := scanner.Bytes()
		test.CheckTestSet(t, originalForID, b)
		count++
	}

	if len(originals) != count {
		t.FailNow()
	}
}
