package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"
)

// FUTURE FEATURE
// To be able to spawn parallel actors if load of one is getting too high.
// Additional syntax needed for flow.def.

// https://linux.die.net/man/5/proc

// Process stat columns ...
const (
	procUtimeIdx  = 13
	procStimeIdx  = 14
	procCutimeIdx = 15
	procCstimeIdx = 16
)

// CPU stat columns ...
const (
	cpuUserIdx   = 1
	cpuNiceIdx   = 2
	cpuSystemIdx = 3
	cpuIdleIdx   = 4
)

// ProcessSample ...
type ProcessSample struct {
	PID  int
	Load uint64
}

// CPUSample ...
type CPUSample struct {
	Load uint64
	Idle uint64
}

var cpuCoreLineRe *regexp.Regexp
var fullCPULineRe *regexp.Regexp

func init() {
	cpuCoreLineRe = regexp.MustCompile(`^cpu\d+`)
	fullCPULineRe = regexp.MustCompile(`^cpu\s+`)
}

// CPUCount ...
func CPUCount() (uint, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}
	lines := bytes.Split(contents, []byte{'\n'})
	var count uint
	for _, line := range lines {
		if cpuCoreLineRe.Match(line) {
			count++
		}
	}
	return count, nil
}

// StartSamplingCPU ...
func StartSamplingCPU(every time.Duration, channel chan<- CPUSample) {
	for {
		load, idle, err := takeCPUSample()
		if err != nil {
			log.Fatalf("Error taking CPU load sample: %s", err)
		}
		channel <- CPUSample{Load: load, Idle: idle}
		time.Sleep(every)
	}
}

// StartSamplingProcess ...
func StartSamplingProcess(pid int, every time.Duration, channel chan<- ProcessSample) {
	for {
		load, err := takeProcessSample(pid)
		if err != nil {
			log.Fatalf("Error taking load sample for process %d: %s", pid, err)
		}
		channel <- ProcessSample{PID: pid, Load: load}
		time.Sleep(every)
	}
}

func takeProcessSample(pid int) (uint64, error) {
	procStatPath := fmt.Sprintf("/proc/%d/stat", pid)
	contents, err := ioutil.ReadFile(procStatPath)
	if err != nil {
		return 0, err
	}
	fields := bytes.Fields(contents)
	load := parseUint64Value(procUtimeIdx, fields, pid)
	load += parseUint64Value(procStimeIdx, fields, pid)
	load += parseUint64Value(procCutimeIdx, fields, pid)
	load += parseUint64Value(procCstimeIdx, fields, pid)
	return load, nil
}

func takeCPUSample() (uint64, uint64, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	lines := bytes.Split(contents, []byte{'\n'})
	for _, line := range lines {
		if fullCPULineRe.Match(line) {
			fields := bytes.Fields(line)
			load := parseUint64Value(cpuUserIdx, fields, -1)
			load += parseUint64Value(cpuNiceIdx, fields, -1)
			load += parseUint64Value(cpuSystemIdx, fields, -1)
			idle := parseUint64Value(cpuIdleIdx, fields, -1)
			return load, idle, nil
		}
	}
	err = fmt.Errorf("Overall CPU stats not found in /proc/stat")
	return 0, 0, err
}

func parseUint64Value(index int, fields [][]byte, pid int) uint64 {
	field := string(fields[index])
	val, err := strconv.ParseUint(field, 10, 64)
	if err != nil {
		if pid == -1 {
			log.Fatalf("Error parsing CPU load sample at index %d: %s", index, err)
		} else {
			log.Fatalf("Error parsing CPU usage sample at index %d for PID %d: %s", index, pid, err)
		}
	}
	return val
}

func parseStringValue(index int, fields [][]byte) string {
	return string(bytes.TrimSpace(fields[index]))
}
