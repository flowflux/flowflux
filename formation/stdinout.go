package formation

import (
	"fmt"
	"io"
	"log"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/event"
)

func newStdinActor(reader io.Reader) *actor.Producer {
	producer := actor.NewProducer(actorInboxSize)
	go func() {
		scanner := event.NewScanner(reader)
		err := producer.Produce(func() ([]byte, error) {
			if scanner.Scan() {
				raw := scanner.Bytes()
				return event.ScannedBytesCopy(raw), nil
			}
			return nil, io.EOF
		})
		if err != nil {
			log.Fatalln("Error processing <stdin> events:", err)
		}
	}()
	return producer
}

func newStdoutActor(writer io.Writer, finished chan<- bool) *actor.Collector {
	collector := actor.NewCollector(actorInboxSize)
	go func() {
		err := collector.Collect(func(message []byte) error {
			n, err := writer.Write(message)
			if err != nil {
				return err
			}
			if n < len(message) {
				return fmt.Errorf("Only %v bytes of %v could be written to <stdout>", n, len(message))
			}
			return nil
		})
		if err != nil && err != io.EOF {
			log.Fatalln("Error processing <stdout> events:", err)
		} else if err == nil {
			finished <- true
		}
	}()
	return collector
}
