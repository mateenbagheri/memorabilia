# Challenges in Memorabilia Development

## Overview

Memorabilia is a personal project aimed at creating a Redis clone,
initially written in Golang. The project's goal is to build a fast,
in-memory key-value store, with the potential to evolve into something
more akin to TiKV, a distributed transactional key-value database.

As the project progresses, various challenges have surfaced, impacting
both the current development trajectory and the potential future pivot.

## 1. Go Concurrency Model

Golang's concurrency model, while powerful, presents its own set of
challenges when developing a high-performance key-value store.

Efficiently managing goroutines, channels, and ensuring thread safety
while minimizing latency and maximizing throughput is a continuous
balancing act.

## 2. Choosing an Appropriate Communication Method

Redis utilizes TCP for client-server communication, which is a
straightforward and well-established method for handling network
communication in key-value stores.

For Memorabilia, I decided to implement gRPC as the communication
protocol.

gRPC offers several advantages, including strong typing through
Protocol Buffers, efficient serialization, built-in support for
bi-directional streaming, and automatic code generation for multiple
programming languages.

These features make gRPC an attractive option for modern, scalable
applications, especially if the project evolves into a distributed
system similar to TiKV.

The challenge lies in adapting gRPC to effectively mimic the simplicity
and performance of Redis's TCP-based communication while leveraging
gRPC's strengths.

This includes designing a gRPC interface that closely mirrors Redis
commands, ensuring low-latency communication, and maintaining the
flexibility to extend or modify the protocol as the project's scope
expands.

Additionally, integrating gRPC features into the core of Memorabilia
requires careful planning to avoid introducing unnecessary complexity
or overhead.

## 3. Choosing How Data Will Be Stored

As the project progresses, we may eventually opt to implement a custom
data structure for Memorabilia. Currently, we are using a Golang map to
store data in memory.

To ensure future changes remain seamless and maintainable, we adopted
the repository pattern. This approach decouples the storage
implementation from the rest of the application, promoting flexibility
and modularity.

Another concern is maintaining data consistency during concurrent read
and write operations. To address this, we implemented an `RWMutex` to
prevent race conditions and ensure thread-safe access.

A further challenge is the lack of a clean built-in mechanism for
defining a map that stores multiple value types. Even with generics,
map initialization requires specifying a single value type, which does
not align with our requirements.

This limitation motivated us to consider designing a custom type system.

## 4. TTL

To design a TTL (Time-To-Live) implementation for Memorabilia, I
considered two approaches:

1. Launch a goroutine that continuously monitors TTL values for all
   keys.

2. Check TTL only when retrieving values. This approach could allow
   expired keys to accumulate in memory, so an additional cleanup task
   would be required.

I chose the second approach because it appeared less resource-intensive
than continuously scanning for expired values in real time.

To mitigate memory growth, users can configure a scheduled cleanup task
that periodically removes expired entries.

After evaluating whether to implement a scheduling utility from scratch
or use an existing package, I selected
`github.com/robfig/cron`.

While implementing a cron scheduler would be an interesting exercise,
it is not a priority during the current phase of development.

This feature was implemented behind a repository-style abstraction,
allowing the scheduler implementation to be replaced in the future,
whether with another package or a custom solution.

