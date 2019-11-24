package logging

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

// TestChannelLogWithSingleReader ...
func TestChannelLogWithSingleReader(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)
	reader := bytes.NewReader(originalCombined)

	receiver := make(chan []byte)
	l := NewChannelLog(receiver)
	l.Add(reader)
	go l.Start()

	count := 0
	for b := range receiver {
		test.CheckTestSet(t, originalForID, b)
		count++
		if count == len(originalBytes) {
			close(receiver)
		}
	}

	if len(originals) != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}

func TestChannelLogWithMultipleReaders(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	receiver := make(chan []byte)
	l := NewChannelLog(receiver)

	for i := 0; i < 5; i++ {
		reader := bytes.NewReader(originalCombined)
		l.Add(reader)
	}

	go l.Start()
	goal := 5 * len(originals)
	count := 0

	for b := range receiver {
		test.CheckTestSet(t, originalForID, b)
		count++
		if count == goal {
			close(receiver)
		}
	}

	if goal != count {
		t.Fatal("Logged line count doesn't corespond to received ones")
	}
}
