# FLOWFLUX

Connect small, reactive processes in every programming language that can read and write to Stdin and Stdout in a flow graph to form applications.

- Communicating processes are also known as actor-based concurrency
- Processes are interconnecting via a messaging-bus
- The form of programming is also known as event-driven or reactive
- The concurrency paradigm like actor based programming languages (Erlang) / messaging queues (RabbitMQ)

## Why?

- Easy to handle form of parallel/concurrent programming
- Easy to handle form of inter-process-communication between programs in different programming languages
- Faster than the same thing via sockets or network interface, also no port-management required
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

Messages/events flow through a process/actor via `stdin` and `stdout`. Stream-connectors are expressed as named pipes. The builtin actors "fork" and "merge" are used to splice and reunite the event streams. The builtin actors "pipe" is used to rewire any `stdout` to any `stdin` via a name pipe.

## Simple Wire Format

Messages/events are JSON encoded, because it's pervasively supported. To distinguish messages from another, an additional envelope is used: the message is bas64-encoded with the delimiter `\n---\n`. Th reasons for this are pervasiveness of base64, ease of use and the resulting robustness.

## Usage

First the `flow` command line tool is used to set up the infrastructure of named pipes. Usually automated by a shell-script. Second the actual actors/processes are started with their Stdin and Stdout connected to the named pipes.

### `flow` CLI usage

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

### Example:

Given the following flow-graph...

```ascii
        Client / Network

            |     ^
            V     |

         web-server.js  <---------------o
                                        |
               |                        |
               |                        |
     o---------o---------o              |
     |                   |              |
     V                   V              |
                                        |
read-file.js     find-media-type.js     |
                                        |
     |                   |              |
     o---------o---------o              |
               |                        |
               |                        |
               V                        |
                                        |
        merge-response.js  -------------o
```

The communications can be described in a flow.def files as follows:

```ascii
node web-server.js -> node read-file.js -> node merge-response.js
node web-server.js -> node find-media-type.js -> node merge-response.js
node merge-response.js -> node web-server.js
```

And run as follows:

```bash
$ ./flow flow.flow
```

Alternatively `flow` CLI can be used to manage a collection of named UNIX pipes. See [example/README.md](example/README.md) for details.

### Example folder

The example folder contains an actor-based static web-server implementation in NodeJS. Where each process follows a fire-and-forget-strategy. Messages/events are running in one direction only. Both things eliminate all problems and difficulties of concurrency and parallel-programming at once. 

The file [example/flow-messaging.js](example/flow-messaging.js) contains the only API necessary, to get the implementation DRY, in 36 LOCs. Porting it to other programming languages should be simple.