package actor_test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/test"
)

func TestSimpleMerge(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	producerAlpha := actor.NewProducer(10)
	producerBeta := actor.NewProducer(10)
	merge := actor.NewMerge(10, "merge", 2)
	collector := actor.NewCollector(10)
	receiver := make(chan []byte, 10)

	producerAlpha.Next(merge)
	producerBeta.Next(merge)
	merge.Next(collector)

	counter := 0
	expectedLen := len(originals) * 2

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	merge.Start()
	startProducingTestMessages(t, producerAlpha, originalCombined)
	startProducingTestMessages(t, producerBeta, originalCombined)

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}

func TestMultipleMerge(t *testing.T) {
	originals := test.LoadTestData()
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	mergeAlpha := actor.NewMerge(10, "merge-alpha", 2)
	mergeBeta := actor.NewMerge(10, "merge-beta", 3)
	mergeGamma := actor.NewMerge(10, "merge-gamma", 2)
	producerAlpha := actor.NewProducer(10)
	producerBeta := actor.NewProducer(10)
	producerGamma := actor.NewProducer(10)
	producerDelta := actor.NewProducer(10)
	producerEpsilon := actor.NewProducer(10)
	collector := actor.NewCollector(10)
	receiver := make(chan []byte, 10)

	producerAlpha.Next(mergeAlpha)
	producerBeta.Next(mergeAlpha)
	producerGamma.Next(mergeBeta)
	producerDelta.Next(mergeBeta)
	producerEpsilon.Next(mergeBeta)
	mergeAlpha.Next(mergeGamma)
	mergeBeta.Next(mergeGamma)
	mergeGamma.Next(collector)

	startCollectingTestMessages(t, collector, receiver, func() {
		close(receiver)
	})
	mergeAlpha.Start()
	mergeBeta.Start()
	mergeGamma.Start()
	startProducingTestMessages(t, producerAlpha, originalCombined)
	startProducingTestMessages(t, producerBeta, originalCombined)
	startProducingTestMessages(t, producerGamma, originalCombined)
	startProducingTestMessages(t, producerDelta, originalCombined)
	startProducingTestMessages(t, producerEpsilon, originalCombined)

	counter := 0
	expectedLen := len(originals) * 5

	for message := range receiver {
		test.CheckTestSet(t, originalForID, message)
		counter++
	}

	if counter != expectedLen {
		t.Fatalf("%v messages expected to be produced, but got %v", expectedLen, counter)
	}
}
