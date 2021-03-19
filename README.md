# Kafka Test Application

Pipeline: [![CircleCI](https://circleci.com/gh/ftamas88/kafka-test.svg?style=svg)](https://circleci.com/gh/ftamas88/kafka-test)

## Table of Contents
- [About](#about)
- [Quickstart](#quickstart)
- [Presentation](#presentation)
- [Learnings](#learnings)

## About
The service reads AVRO payload from the filesystem and publishes it to the internal and cloud Kafka clusters.

**Task:**
- [x] Reading files from the filesystem when they arrive
- [x] Parsing the data in the files
- [x] Serializing that data in preparation for publishing to Kafka clusters
- [x] Creating a Kafka client
- [x] Publishing to multiple internal clusters
- [x] Publishing to a Kafka cluster running on a cloud platform
- [x] Removing the files from the filesystem when their contents have been published
- [ ] Providing an SDK for use by other teams in your startup for some of the key, high-level functionality

- [x] Presentation



## Quickstart
### Run the service
```
make run
```

### Generate mocked AVRO payload
```
make generate-mock
```
Sample:
```json
{
  "Id" : "fC30BffaBE04",
  "Reference" : "azDUa",
  "Amount" : {
    "Value" : -231.51184,
    "Currency" : "USD"
  },
  "BookedTime" : "1609911988"
}
```
This will generate a file in the `/tmo/ingest/` folder

### Formatting, linting and tests
```
make fmt lint test
```

### Tests in CI container (for integration)
```
make test-ci
```

### List all make targets
```
make help
```

## Environment Variables
| Name | Description |
|:-----|-------------|
| `HTTP_PORT` | HTTP port to serve from |
| `KAFKA_SERVERS` | List of Kafka brokers (comma separated) |
| `KAFKA_SCHEMA_SERVERS` | List of Kafka Schema servers. (comma separated) |
| `KAFKA_CLOUD_SERVERS` | List of CCloud servers |
| `KAFKA_CLOUD_SCHEMA_SERVERS` | List of CCloud Schema Registry servers |
| `KAFKA_CLOUD_KEY` | CCloud authentication key |
| `KAFKA_CLOUD_SECRET` | CCloud authentication secret |

## Presentation
The files located in the `/docs` folder

- [PDF](docs/Kafka%20test.pdf)
- [PPTX](docs/Kafka%20test.pptx)

## Learnings
These are my notes I made during my Kafka research
Available in the `/docs` folder, or [click here](docs/etc.MD)