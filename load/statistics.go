package load

import (
	"time"
)

const maxSampleSize = 1000

// Statistics ...
type Statistics struct {
	cpuChannel     chan CPUSample
	processChannel chan ProcessSample
	cpuSamples     []CPUSample
	cpuIndex       uint
	procSamples    map[int][]ProcessSample
	procIndexes    map[int]uint
}

// SampleProcess ...
type SampleProcess interface {
	TakeLoadSamples(time.Duration, chan<- ProcessSample)
}

// NewStatistics ...
func NewStatistics() Statistics {
	return Statistics{
		cpuChannel:     make(chan CPUSample),
		processChannel: make(chan ProcessSample),
		cpuSamples:     make([]CPUSample, maxSampleSize),
		cpuIndex:       0,
		procSamples:    make(map[int][]ProcessSample),
		procIndexes:    make(map[int]uint),
	}
}

// CPUChannel ...
func (s Statistics) CPUChannel() chan<- CPUSample {
	return s.cpuChannel
}

// ProcessChannel ...
func (s Statistics) ProcessChannel() chan<- ProcessSample {
	return s.processChannel
}

// Start ...
func (s Statistics) Start(procs []SampleProcess) {
	every := 1 * time.Second
	go StartSamplingCPU(every, s.cpuChannel)
	for _, p := range procs {
		p.TakeLoadSamples(every, s.processChannel)
	}
	for {
		select {
		case cpu := <-s.cpuChannel:
			s.addCPUSample(cpu)
		case proc := <-s.processChannel:
			s.addProcSample(proc)
		}
	}
}

func (s Statistics) addCPUSample(cpu CPUSample) {
	s.cpuSamples[s.cpuIndex] = cpu
	s.cpuIndex++
	if s.cpuIndex == maxSampleSize {
		s.cpuIndex = 0
	}
}

func (s Statistics) addProcSample(proc ProcessSample) {
	procSamples, procOK := s.procSamples[proc.PID]

	if !procOK {
		procSamples = make([]ProcessSample, maxSampleSize)
		s.procIndexes[proc.PID] = 0
	}

	procIndex := s.procIndexes[proc.PID]
	procSamples[procIndex] = proc

	procIndex++
	if procIndex == maxSampleSize {
		procIndex = 0
	}

	s.procSamples[proc.PID] = procSamples
	s.procIndexes[proc.PID] = procIndex
}
