# About failing tests

`command_test.go` and `sequence_test.go` contain the heavy duty tests. Both of them can fail once in a while due to overload of system pipes or too high pressure on nodes.js pipes. Error might look something like that:

> panic: send on closed channel

A workaround is to throttle the producer to a lesser value: `NewThrottledProducer(..., <events-per-minute>)`.
