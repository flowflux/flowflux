package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

// FUTURE FEATURE
// To be able to spawn parallel actors if load of one is getting too high.
// Additional syntax needed for flow.def.

// https://linux.die.net/man/5/proc

// Classes ...
const (
	UtimeIdx  = 13
	StimeIdx  = 14
	CutimeIdx = 15
	CstimeIdx = 16
)

// StartSampling ...
func StartSampling(pid int, every time.Duration) <-chan uint64 {
	channel := make(chan uint64)
	go func() {
		for {
			sample, err := takeSample(pid)
			if err != nil {
				log.Fatalf("Error taking CPU usage sample for PID %d: %s", pid, err)
			}
			channel <- sample
			time.Sleep(every) // 1 * time.Second
		}
	}()
	return channel
}

func takeSample(pid int) (uint64, error) {
	procStatPath := fmt.Sprintf("/proc/%d/stat", pid)
	contents, err := ioutil.ReadFile(procStatPath)
	if err != nil {
		return 0, err
	}
	fields := bytes.Fields(contents)
	sample := parseUint64Value(UtimeIdx, fields, pid)
	sample += parseUint64Value(StimeIdx, fields, pid)
	sample += parseUint64Value(CutimeIdx, fields, pid)
	sample += parseUint64Value(CstimeIdx, fields, pid)
	return sample, nil
}

func parseUint64Value(index int, fields [][]byte, pid int) uint64 {
	field := string(fields[index])
	val, err := strconv.ParseUint(field, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing CPU usage sample at index %d for PID %d: %s", index, pid, err)
	}
	return val
}

func parseStringValue(index int, fields [][]byte) string {
	return string(bytes.TrimSpace(fields[index]))
}
