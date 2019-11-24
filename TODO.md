# TODO

1. Buildup of flow is done via a list of processor configurations
2. At each step the processor pops the first element and configures itself
3. Then starts up the next step in the sequence and hands over the remaining processor configurations
4. The previously created processors are kept in a central place for access from all the others

## Implementation process

1. Implement 2 simple processors that can be tested against each other
2. Test them for high data throughput reliability
3. Continue with 1