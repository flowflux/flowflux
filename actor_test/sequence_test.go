package actor_test

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/logging"
	"github.com/jaqmol/approx/test"
)

func TestSimpleCommandSequence(t *testing.T) {
	// t.SkipNow()
	originals := test.LoadTestData() // [:10]
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)

	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	conf := test.MakeSimpleSequenceConfig()
	producer := actor.NewThrottledProducer(10, 5000)

	fork := actor.NewForkFromConf(10, &conf.Fork)
	producer.Next(fork)

	firstNameExtractCmd, err := actor.NewCommandFromConf(10, &conf.FirstNameExtract)
	catchToFatal(t, err)
	lastNameExtractCmd, err := actor.NewCommandFromConf(10, &conf.LastNameExtract)
	catchToFatal(t, err)

	receiver := make(chan unifiedMessage, 10)

	// Logger
	logReceiver := make(chan []byte, 10)
	logger := logging.NewChannelLog(logReceiver)
	funnelIntoUnifiedLogMessages(logReceiver, receiver)
	logger.Add(firstNameExtractCmd.Logging())
	logger.Add(lastNameExtractCmd.Logging())
	// /Logger

	fork.Next(firstNameExtractCmd, lastNameExtractCmd)

	merge := actor.NewMergeFromConf(10, &conf.Merge)
	firstNameExtractCmd.Next(merge)
	lastNameExtractCmd.Next(merge)

	collector := actor.NewCollector(10)
	merge.Next(collector)

	startCollectingUnifiedDataMessages(t, collector, receiver, func() {
		close(receiver)
	})

	fork.Start()
	firstNameExtractCmd.Start()
	lastNameExtractCmd.Start()
	merge.Start()
	go logger.Start()

	startProducingTestMessages(t, producer, originalCombined)

	counter := 0
	for message := range receiver {
		if message.messageType == unifiedMsgDataType {
			err = test.CheckUpperExtraction(message.data, originalForID)
			catchToFatal(t, err)
			counter++
		} else if message.messageType == unifiedMsgLogType {
			test.CheckCmdLogMsg(t, "Did extract", message.data)
		}
	}

	goal := len(originals) * 2
	if counter != goal {
		t.Fatalf("%v messages expected to be produced, but got %v", goal, counter)
	}
}

func catchToFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
