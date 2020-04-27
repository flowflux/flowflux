Stdin and stdout support for flow-definitions started with the `flow` CLI is not yet fully functional.

The idea is that "open ends" in the definition file are reading from stdin and writing to stdout:

```ascii
-> node to-msg.js -> node compose.js ->
```