package logging

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

// TestWriterLogWithSingleReader ...
func TestWriterLogWithSingleReader(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	writer := test.NewWriter()
	l := NewWriterLog(writer)
	l.Add(reader)
	go l.Start()

	count := 0
	for b := range writer.Lines {
		test.CheckTestSet(t, originalForID, b)
		count++
		writer.Stop(count == len(originalBytes))
	}

	if len(originals) != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}

func TestWriterLogWithMultipleReaders(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	writer := test.NewWriter()
	l := NewWriterLog(writer)

	for i := 0; i < 5; i++ {
		reader := bytes.NewReader(originalCombined)
		l.Add(reader)
	}

	go l.Start()
	goal := 5 * len(originals)
	count := 0

	for b := range writer.Lines {
		test.CheckTestSet(t, originalForID, b)
		count++
		writer.Stop(count == goal)
	}

	if goal != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}
