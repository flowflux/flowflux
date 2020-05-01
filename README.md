# FLOWFLUX

Actor-model concurrent computation utility. Connect small, reactive processes in every programming language in a flow graph to form more complex applications.

## Domain Key Words

- Communicating sequential processes
- Inter-process communication
- Actor-model based concurrency (Erlang)
- Event-driven
- Reactive
- Fire-and-forget
- Messaging-bus or event-bus (RabbitMQ, Kafka)

That means that flowflux can solve problems in the domains of concurrency and parallelism (actor pattern) as well as messaging queues (RabbitMQ).

## Why?

- Easy to handle form of parallel/concurrent programming
- Easy to handle form of inter-process-communication between programs in different programming languages
- Faster than the same thing via network interface, no port-management required
- Unix pipes are cool and everything that makes them usable for more problems is highly welcome 
- Communication between processes can be visualized

### Also 

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

## Input, Output

Messages/events flow through a process/actor via `stdin` and `stdout`. Stream-connectors are expressed as named pipes or handled by Go channels (depending on which execution model is chosen). The builtin actors "fork" and "merge" are used to splice and reunite the event streams. The builtin actor "pipe" is used to rewire any `stdout` to any `stdin`.

### Simple Wire Format

When development on this idea started, it became clear quite fast, that data transfer between different programming languages and how they handle binary data streams isn't easy to get reliable. Many different solutions where tried (length-prefixed streams, character-separated messages. etc.pp.). The following combination turned out to be a rock solid solution, though pure binary streams that are not segmented automatically at message-boundaries are supported too.

Messages/events themselves are JSON encoded, because it's pervasively supported. To distinguish messages from another, an additional envelope layer is used: The message is base64-encoded and suffixed with the delimiter `\n---\n`. This seems unnecessary at first, but a pure JSON wire-format would require 100% control over how the JSON is generated and how it's data is escaped. An additional base64 envelope is an equally pervasive solution that gives communication robustness. The delimiter `\n---\n` is nicely recognizable in between any form of generated base64.

As mentioned before, raw byte streams are supported too if flow is used as a communication hub with a `flow.def` file. A connection specified with `*->` instead of `->` stands for a raw binary stream.

## Example:

Given the following flow-graph...

```ascii
           Client / Network

               |     ^
               V     |

            web-server.js  <------------------o
                                              |
                  |                           |
                  |                           |
     o--o--o--o---o-------------o             |
     |  |  |  |                 |             |
     V  V  V  V                 V             |
                                              |
 read-file.js (x4)     find-media-type.js     |
                                              |
     |  |  |  |                 |             |
     o--o--o--o---o-------------o             |
                  |                           |
                  |                           |
                  V                           |
                                              |
          merge-response.js  -----------------o
```

The communications can be described in a flow.def files as follows:

```ascii
node web-server.js -> node read-file.js (x4) -> node merge-response.js
node web-server.js -> node find-media-type.js -> node merge-response.js
node merge-response.js -> node web-server.js
```

This flow.def describes a pipeline of communicating sequential processes, where the messaging flow is forked into 5 parallel processes, from which 4 handle file reading `node read-file.js (x4)`, and 1 handles finding the correct content type `node find-media-type.js`.

And run as follows:

```bash
$ ./flow flow.def
```

Alternatively `flow` CLI can be used to manage a collection of named UNIX pipes. See [examples/README.md](examples/README.md) for details.

### Example folder

The example folder contains an actor-based static web-server implementation in NodeJS. Where each process follows a fire-and-forget-strategy. Messages/events are running in one direction only. Both things eliminate all problems and difficulties of concurrency and parallel-programming at once. 

The file [examples/flow-messaging/index.js](examples/flow-messaging/index.js) contains the only API necessary, to get the implementation DRY, in 36 LOCs. Porting it to other programming languages should be simple.

## Usage

Flow allows for 2 different ways to run applications:

1. With flow as communication hub, the most comfortable one
2. With flow managing a collection of named UNIX pipes

Please see the example folder. The web-server is usable in both ways.

### `flow` CLI

```bash
$ ./flow
flow
Utility to build messaging systems by composing command line processes
pipe <name>
  Pipe message stream from <name>.wr to <name>.rd
fork <wr-name> <rd-name-1> <rd-name-2> <...>
  Fork message stream from wr-fifo into all provided rd-fifos
merge <wr-name-1> <wr-name-2> <...< <rd-name>
  Merge message stream from all provided wr-fifos into rd-fifo
input <name>
  Input JSON messages to stream them to <name>
cleanup <directory>
  Cleanup directory from fifos (wr & rd)
```