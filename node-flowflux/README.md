# FLOWFLUX

Connect small, reactive processes to form applications.

1. Small is beautiful.
2. Make each program do one thing well.
3. Build a prototype as soon as possible.
4. Choose portability over efficiency.
5. Store data in flat text files.
6. Use software leverage to your advantage.
7. Use shell scripts to increase leverage and portability.
8. Avoid captive user interfaces.
9. Make every program a filter.

*Mike Gancarz: The UNIX Philosoph*

## Nomenclature

### Flowflux executable

Reads a formation and spins up a flow-graph of actors.

### Actor

A process classified by:

- Parametrized by environment variables
- Listens for events from an input-streams
- Processes events and data
- Writes events to an output-streams
- As low complexity as possible
- Listen for events until it receives SIGINT

### Input, Output

Events flow through an actor via `stdin` and `stdout`. The builtin actors "fork" and "merge" are used to splice the event stream.

## Why

- To allow for designing applications as a flow-graph of stream actors.
- To design each of them as a low complexity stream process, that transforms input into output.
- To build multi-process-applications easily and with every programming language that can read/write to stdin/-out.
- To mix and match programming languages as comes handy.
- To choose the best library (or programming language) for a single problem.
- Not to be forced into compromises.

A flow-graph of stream actors is an event driven architecture. Listening for messages / or IPC via reading from `stdin` is among the most basic of tasks in the very most of programming languages. Also:
- Very good performance.
- Very good documentation.
- Maximum operating system support.

### Why not unix sockets?

Too difficult to reach in most programming languages. Too complex to use.
Sockets don't perform better than pipes.

### Why not http, websockets, long polling, ...?

Request-response-based network communication protocols are not a natural or effective choice for event driven architectures. Especially not if used on the same machine.
Communicating processes via network interface is not as ressource efficient as pipes.