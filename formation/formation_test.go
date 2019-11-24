package formation

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/logging"
	"github.com/jaqmol/approx/test"
)

// TestSimpleActorFormation ...
func TestSimpleActorFormation(t *testing.T) {
	// t.SkipNow()
	originals := test.LoadTestData() // [:100]
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)
	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	inputReader := bytes.NewReader(originalCombined)
	outputWriter := test.NewWriter()
	logChannel := make(chan []byte)
	logger := logging.NewChannelLog(logChannel)

	projDir, err := filepath.Abs("../test/gamma-test-proj") // 1st name only
	if err != nil {
		t.Fatal(err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	testArgs := []string{origArgs[0], projDir}
	os.Args = testArgs

	form, err := NewFormation(inputReader, outputWriter, logger)
	if err != nil {
		t.Fatal(err)
	}

	finished := form.Start()
	outCounter, logCounter := 0, 0

	loop := true
	for loop {
		select {
		case outMsg := <-outputWriter.Lines:
			err = test.CheckUpperExtraction(outMsg, originalForID)
			catchToFatal(t, err)
			outCounter++
		case logMsg := <-logChannel:
			test.CheckCmdLogMsg(t, "Did extract", logMsg)
			logCounter++
		case <-finished:
			loop = false
		}
	}

	if outCounter != len(originals) {
		t.Fatalf("Expected %v outputs, but got %v", len(originals), outCounter)
	}
	if logCounter != len(originals) {
		t.Fatalf("Expected %v log messages, but got %v", len(originals), logCounter)
	}
}

// TestComplexActorFormation ...
func TestComplexActorFormation(t *testing.T) {
	// t.SkipNow()
	originals := test.LoadTestData() // [:10]
	originalForID := test.MakePersonForIDMap(originals)
	originalBytes := test.MarshalPeople(originals)
	originalCombined := bytes.Join(originalBytes, config.EvntEndBytes)
	originalCombined = append(originalCombined, config.EvntEndBytes...)

	inputReader := bytes.NewReader(originalCombined)
	outputWriter := test.NewWriter()
	logChannel := make(chan []byte)
	logger := logging.NewChannelLog(logChannel)

	projDir, err := filepath.Abs("../test/beta-test-proj") // 1st and last name
	if err != nil {
		t.Fatal(err)
	}

	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	testArgs := []string{origArgs[0], projDir}
	os.Args = testArgs

	form, err := NewFormation(inputReader, outputWriter, logger)
	if err != nil {
		t.Fatal(err)
	}

	finished := form.Start()
	outCounter, logCounter := 0, 0

	loop := true
	for loop {
		select {
		case outMsg := <-outputWriter.Lines:
			err = test.CheckUpperExtraction(outMsg, originalForID)
			catchToFatal(t, err)
			outCounter++
		case logMsg := <-logChannel:
			test.CheckCmdLogMsg(t, "Did extract", logMsg)
			logCounter++
		case <-finished:
			loop = false
		}
	}

	goal := len(originals) * 2
	if outCounter != goal {
		t.Fatalf("Expected %v outputs, but got %v", goal, outCounter)
	}
	if logCounter != goal {
		t.Fatalf("Expected %v log messages, but got %v", goal, logCounter)
	}
}

func catchToFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
