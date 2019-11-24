package actor_test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

func TestSingleCollector(t *testing.T) {
	// Single collector is being tested in producer_test.go
}

func TestMultipleCollectors(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	producer := actor.NewProducer(10)
	collectorAlpha := actor.NewCollector(10)
	collectorBeta := actor.NewCollector(10)
	receiver := make(chan []byte, 10)

	producer.Next(collectorAlpha, collectorBeta)

	startCollectingTestMessages(t, collectorAlpha, receiver, func() {})
	startCollectingTestMessages(t, collectorBeta, receiver, func() {})
	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	expectedLen := len(originals) * 2

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
		if counter == expectedLen {
			close(receiver)
		}
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}

func startCollectingTestMessages(
	t *testing.T,
	collector *actor.Collector,
	receiver chan<- []byte,
	finished func(),
) {
	go func() {
		err := collector.Collect(func(message []byte) error {
			receiver <- message
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		finished()
	}()
}
