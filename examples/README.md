# EXAMPLE: ACTOR-BASED WEB-SERVER IN NODE.JS

Below is a diagram that shows the message flow between the actors:

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

The example is a streaming parallel/concurrent web-server programmed in node.js that serves static files. The process read-files.js sends files chunk-wise via flow messages. Web-server.js receives requests and dispatches the read file-chunks to the clients as they come. Merge-response.js takes care that chunks have the correct content-type attached, which is handled by find-media-type.js.

Flowflux follows a fire-and-forget-strategy. Messages/events are running in one direction only. Both things eliminate all problems and difficulties of concurrency and parallel-programming at once.

Flowflux allows for 2 different ways to run this application.

1. With flow as communication hub
2. With flow managing a collection of named UNIX pipes

## Flowflux as communication hub

The messaging-flow between the actors is declared in a file, the file is then fed to the "flow" CLI:

**Definition file example: flow.def**

```ascii
node web-server.js -> node read-file.js (x4) -> node merge-response.js
node web-server.js -> node find-media-type.js -> node merge-response.js
node merge-response.js -> node web-server.js
```
This flow.def describes a pipeline of communicating sequential processes, where the messaging flow is forked into 5 parallel processes, from which 4 handle file reading `node read-file.js (x4)`, and 1 handles finding the correct content type `node find-media-type.js`.

### Startup procedure

```bash
$ ./flow flow.def
```

### Termination procedure

The usual `ctrl+c`.

## Flowflux managing a collection of named UNIX pipes

This way of running the same code does not feature parallel/concurrent execution, in order to keep the example comprehensible.

### Meaning of pipe characters:

- program stdout | program stdin
- program stdout > fifo/file 
- program stdin < fifo/file

### That means

- Example: `wc -l < input.txt > output.txt`
- Would be interpreted as: `(wc -l < input.txt) > output.txt`

### Expressing flow

- Goal: express data flow: `input.txt -> program -> output.txt`
- Solution: `program < input.txt > output.txt`

### Startup procedure

#### 1st start the named pipes

```bash
./flow pipe fifo-response-pipe &
./flow fork fifo-web-server-out fifo-read-file-in fifo-find-media-type-in &
./flow merge fifo-read-file-out fifo-find-media-type-out fifo-merge-response-in &
```

#### 2nd start the programs

> **Note:**
> `web-server.js, find-media-type.js, merge-response.js, read-file.js`
> Need to be executable (chmod +x)

```bash
./web-server.js < fifo-response-pipe.rd > fifo-web-server-out.wr &
./find-media-type.js < fifo-find-media-type-in.rd > fifo-find-media-type-out.wr &
./merge-response.js < fifo-merge-response-in.rd > fifo-response-pipe.wr &
./read-file.js < fifo-read-file-in.rd > fifo-read-file-out.wr &
```

### Termination procedure

When used like that, ever actor runs it's own process running the background, even the infrastructure nodes (named UNIX pipes). In order to terminate it, you need to bring every process in the foreground with `fg`, where it can be terminated with the usual `ctrl+c`.