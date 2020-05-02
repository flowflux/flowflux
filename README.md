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

Messages/events flow through a process/actor via `stdin` and `stdout`. Internally the processes are connector by named pipes or handled by Go channels (depending on which execution model you chosen). The builtin actors "fork" and "merge" are used to splice and reunite the streams. The builtin actor "pipe" is used to rewire any `stdout` to any `stdin`.

### Simple Wire Format

Flowflux suggests messages/events to be JSON encoded, because it's pervasively supported and fast. To distinguish messages from another, each string is required to be prefixed with it's length, like so:

```json
00000000000000078377{"id":"ck9pwb050001u26i021e96bed","cmd":"PROCESS_FILE_CHUNK","encoding":"base64","payload":"WQtbibrqaWDYlz...
```

#### Length-Prefixes

Length-prefixes are required to be unsigned 64bit integers rendered as UTF8-strings. Each needs to be exactly 20 digits long and 0-padded to the left. See [examples/flow-messaging](examples/flow-messaging/index.js) for an example of a NodeJS-implementation of this wire-format.

```js
function Envelope(msg) {
  const payload = Buffer.from(JSON.stringify(msg));
  const prefix = `${payload.length}`.padStart(20, '0');
  const length = Buffer.from(prefix);
  return Buffer.concat([length, payload]);
}
```

For an example of how to parse this messages, have a look at `ParseMessage` in [examples/flow-messaging](examples/flow-messaging/index.js).

Raw byte streams are supported too if flow is used as a communication hub with a `flow.def` file. A connection specified with `*->` instead of `->` defines a raw binary stream.

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

The communications can be described in a flow.def file as follows:

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