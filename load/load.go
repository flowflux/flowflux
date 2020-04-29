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

// Sample ...
type Sample struct {
	Comm   string
	Utime  uint64
	Stime  uint64
	Cutime uint64
	Cstime uint64
}

func (s Sample) String() string {
	return fmt.Sprintf(
		"%s, Utime: %d, Stime: %d, Cutime: %d, Cstime: %d",
		s.Comm,
		s.Utime,
		s.Stime,
		s.Cutime,
		s.Cstime,
	)
}

// https://linux.die.net/man/5/proc

// Classes ...
const (
	CommIdx   = 1
	UtimeIdx  = 13
	StimeIdx  = 14
	CutimeIdx = 15
	CstimeIdx = 16
)

// StartSampling ...
func StartSampling(pid int, every time.Duration) <-chan Sample {
	channel := make(chan Sample)
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

func takeSample(pid int) (Sample, error) {
	procStatPath := fmt.Sprintf("/proc/%d/stat", pid)
	contents, err := ioutil.ReadFile(procStatPath)
	if err != nil {
		return Sample{}, err
	}
	fields := bytes.Fields(contents)
	sample := Sample{
		Comm:   parseStringValue(CommIdx, fields),
		Utime:  parseUint64Value(UtimeIdx, fields, pid),
		Stime:  parseUint64Value(StimeIdx, fields, pid),
		Cutime: parseUint64Value(CutimeIdx, fields, pid),
		Cstime: parseUint64Value(CstimeIdx, fields, pid),
	}
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
