# Learnings

## Kafka

### Why Kafka?

- **Open Source:**
    - 2010: Initial development by LinkedIn
    - 2011: Made Open Source
    - 2012: Apache Licence
    - 2015-: Used by Netflix, Uber, Spotify, etc
    - 2017: Apache Kafka 1.0
    - 2019: Apache Kafka 2.0
- **Scala/Java:**
  Initialy the development started with Scala since using functional languages seemed to be the trend however Scala is outnumbered by Java in developer count by far. Since most of the users are Java programmers it just made sense to create the client API in Java. While the server code is written in Scala. In general Scala programmers can read Java but not viceversa. Furthermore the learning curve for Scala is steeper than the Java one.
- High throughput
    - No serialization/deserialisation in transit
    - All incoming and outgoing data will be in binary format
    - Zero Copy (Instead of JVM Heap (Memory) can write to the disk directly) for non-TLS connections
- Not just a messaging platform -> Distributed Streaming platform
    - Messaging system (pub/sub)
    - Distributed storage (store data in a fault-tolerant, durable way)
    - Data processing (process events as they occur)

### Components

- Broker
- Zookeeper

> ... a centralised service for maintaining configuration information, naming, providing distributed synchronisation, and providing group services.

- Topic
    - Retention
    - Compacting topics
    - Partitions / Replicated partitions
- Producer
    - Producer record
    - Key/Value serialiser
- Consumer
    - Constant polling
    - Consumer record
    - Key/Value deserialiser

## Decisions / thought process

### Event-driven architecture

> A Software architecture pattern promoting the production, detection, consumption of, and reaction to events - Wikipedia -

- Microservices
- Serverless
- FaaS
- Event Sourcing
- CQRS
- Streaming

---
**Pros**:

- Decoupling
- Encapsulation
- Optimisation (speed)
- Scalability

**Cons**:

- Steep learning curve
    - What is the source of truth?
- Complexity
- Loss of transactionality
- Lineage (harder tracking/debugging)

### Kafka Streams

Listening on a topic, process it and sends it to a different topic.  
It can abscract the concept of producers/consumers and makes it easy to use.

Source/Consumer -> Processor x -> Processor Y -> Sink/Producer

> Topology: acyclic graph of sources, processors, and sinks

#### Duality of Streams

- Regular stream
    - Processing independent events (delete topics)
    - Can be transfered via aggregate/reduce/count
- Table
    - Processing evolving events (compaction topics)

#### Processors

- Stateless
    - Branch
    - Filter
    - Inverse Filter
    - Map
    - FlatMap
    - Foreach
    - Peek
    - GroupBy
    - Merge
- State store
    - Aggregations
    - Count
    - Joins
    - Windowing
    - Custom processors

#### KSQL

- Why?
    - (+) Easy syntax
    - (-) Less capabilities
- How?
    - KSQL Cli -> Rest API -> KSQL Server
- When?
    - Streaming operations (data analytics, monitoring, iot, etc)
    - Viewing Data (show the content of a topic)
    - Manipulating Data

```sql
  CREATE STREAM high_value_transactions AS
    SELECT category, amount, description
    FROM transactions
    WHERE amount > 10 000;
```

VS

```java
.stream("transaction")
.filter(transaction -> 
  transaction.getAmount() > 10 000
)
.mapValues(transactions -> 
  transaction.getCategory() + "," +
  transaction.getAmount() + "," +
  transaction.getDescription()
)
.to("high_value_transactions")
```


## Optimisations

### Config

- On lower memory systems the `batch_size` can be configured to lower the items in the buffer before the "send" thread will pick it up
- `linger.ms` (by default 0) then the messages will only sent after x ms. This value should be considered based on *Latency* vs *Throughput* requirements
- Enabling compressing by `compression.type` (e.g. gzip) -> reduce disk space, network traffic
- `heartbeat.interval.ms` must be lower than `session.timeout.ms` in the consumer
- `max_poll.interval.ms` to prevent consumer to get stuck during longer processing
- Duplicate messages: Relying on `Exactly-once` could be problematic, it comes with a significant performance cost. Using the `At-least-once guarantee` is recommended
- `min.insync.replicas` e.g. 2 -> to prevent loses if the leader goes down after ack but not fetched by replicas yet who will be elected as a new leader then the message can be lost.

### Data retention

- Disk space is limited, segment configuration should be set accordingly:
    - `log.segment.bytes` once the segment reaches a certain sizes, it will get closed, and a new one will be created
    - `log.segment.ms` after the given ms the current segment will be closed, and a new one will be created
- Retention period can be customised:
    - (bytes/ms) behavior same as above
    - `log.retention.bytes`
    - `log.retention.ms`

### Cluster Sizing

- The load should be estimated before deploying the cluster
- Usual bottlenecks:
    - Disk size/throughput
    - Network throughput

### Serialisation

> The process of translating **data structures** or **object states** into a format that can be **stored** or **transmitted** and reconstructed later, possibly in a **different** computer environment. - Wikipedia -

Choosing a serialisation method can be tricky.
The cost of serialisation/deserialisation is NOT negligible.

Ideally we are looking for something which is:

- Fast
- Compact
- Language-agnostic
- Easily extendable schema

#### Formats

| Name     | Binary | Human readable | Schema / IDL |
|----------|:------:|:--------------:|:------------:|
| JSON     |        | X              |              |
| XML      |        | X              | X            |
| YAML     |        | X              |              |
| AVRO     | X      |                | X            |
| Protobuf | X      |                | X            |
| Thrift   | X      |                | X            |

> **IDL**: Interface Description Language

#### JSON

The problem with JSON in a big data system is the lot of wasted bandwith traffic caused by the keys that have to sent in every event.
It is also more prone to error due to changing/evolving structure.

#### Serialisation frameworks

Binary formats improve efficiency and simplify schema management.

- Protocol buffers (Google)

    ```proto
        message Ads {
            required int32 ad_id = 1;
            required string type = 2;
        }
    ```

  > Using `protoc` compiler it can generate code for various languages.

- Apache Avro (Confluence)

  `transaction.avsc`

  ```json
    {
      "type": "record",
      "namespace": "enrichment",
      "name": "Transaction",
      "fields": [
        {
          "name": "transactionId",
          "type": "string"
        },
        {
          "name": "description",
          "type": "string"
        }
      ]
    }
  ```

  > Files can be also generated but not required.

- Facebook Thrift

#### Schema Registry

*Producers* can register schemas which then can be discovered by the *consumers*.
**Dynamic messages** can help when we want to consume different message types.
For example: if we are streaming ads but there are different sub-types. Maybe at the first stage we only care about the main type and based on this we will route the message accordingly.

It uses a Kafka topic to store the schemas.
Messages/Records doesn't contain the schema just a refrence to it.

Resolving the schemas are done in the Producer/Consumers via an API call.

There are multiple strategies available:

- Topic name strategy
- Record name strategy
- Topic Record name strategy

> **Note:** There is some performance cost associated with using dynamic messages.

### Notes

- `rack` value on the brokers can be used to group it by different regions, ensuring the system is fault-tolerant when using it with the minimum InSyncReplicas
