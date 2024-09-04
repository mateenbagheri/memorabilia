# Challenges in Memorabilia Development

## Overview

Memorabilia is a personal project aimed at creating a Redis clone, initially written in Golang. The project’s goal is to build a fast, in-memory key-value store, with the potential to evolve into something more akin to TiKV, a distributed transactional key-value database. As the project progresses, various challenges have surfaced, impacting both the current development trajectory and the potential future pivot.

#### 1. Go Concurrency Model

Golang’s concurrency model, while powerful, presents its own set of challenges when developing a high-performance key-value store. Efficiently managing goroutines, channels, and ensuring thread safety while minimizing latency and maximizing throughput is a continuous balancing act.

### 2. Choosing Appropriate Communication Method for Interacting with Memorabilia Core

Redis utilizes TCP for client-server communication, which is a straightforward and well-established method for handling network communication in key-value stores. However, for Memorabilia, I have decided to implement gRPC as the communication protocol. 

gRPC offers several advantages, including support for strong typing through Protocol Buffers, efficient serialization, built-in support for bi-directional streaming, and automatic code generation for multiple programming languages. These features make gRPC an attractive option for modern, scalable applications, especially as the project may evolve into a distributed system like TiKV in the future.

The challenge lies in adapting gRPC to effectively mimic the simplicity and performance of Redis’s TCP-based communication while leveraging gRPC's strengths. This includes designing a gRPC interface that closely mirrors Redis commands, ensuring low-latency communication, and maintaining the flexibility to extend or modify the protocol as the project’s scope expands. Additionally, integrating gRPC's features into the core of Memorabilia requires careful planning to avoid introducing unnecessary complexity or overhead.

### 3. Choosing how data will be stored
As the project progresses, we may eventually opt to implement a custom data structure for memorabilia. Currently, we are using a Golang map to store our data in memory.

To ensure that any future changes to the data structure are seamless and maintainable, we have adopted the repository pattern. This approach decouples our data storage implementation from the rest of the application, promoting flexibility and modularity.

Another concern is maintaining data consistency during concurrent read and write operations on the map. To address this, we have implemented an RWMutex to prevent race conditions and ensure thread-safe access to the data.

Additionally, a key challenge that led us to consider designing a custom type system is the lack of a clean, built-in method in Go for defining a map that can store multiple types of data as values. Even with the use of generics, this issue persists, as initializing a map requires specifying a single data type, which does not align with our requirements.