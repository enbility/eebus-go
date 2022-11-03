# eebus-go

![Build Status](https://github.com/DerAndereAndi/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)

The goal is to provide a basic EEBUS implementation

## Introduction

This repository contains:

- adoptions of the SPINE and SHIP EEBUS model definitions, there are likely issues and some models are not 100% correct
- (De-)serialization for EEBUS specific JSON format requirements
- Certificate support incl. creating a compatible cert and key
- mDNS Support (announcement and connecting to an announced SKI) incl. avahi support if available
- ... work in progress

You need a basic understanding of the EEBUS concepts SHIP and SPINE to use this library. Please check the corresponding specification on the [EEBUS website](https://eebus.org).

## Usage

The included small demo applications do not implement any usecases and thus will end the connection once it reached exchanging usecase information.

Services with implemented use cases will be implemented in different repositories and are also early work in progress:

- [HEMS](https://github.com/DerAndereAndi/eebus-go-cem)
- [EVSE](https://github.com/DerAndereAndi/eebus-go-evse)

### HEMS

#### First Run

```sh
go run cmd/hems/main.go 4715
```

`4715` is the example server port that this process should use

The certificate and key and the local SKI will be generated and printed. You should then save the certificate and the key to a file.

#### General Usage

```sh
Usage: go run cmd/hems/main.go <serverport> <remoteski> <certfile> <keyfile>
```

- `remoteski` is the SKI of the remote device or service you want to connect to
- `certfile` is a local file containing the generated certificate in the first usage run
- `keyfile` is a local file containing the generated key in the first usage run

### EVSE

#### First Run

```sh
go run cmd/hems/main.go 4715
```

`4715` is the example server port that this process should use

The certificate and key and the local SKI will be generated and printed. You should then save the certificate and the key to a file.

#### General Usage

```sh
Usage: go run cmd/evse/main.go <serverport> <remoteski> <certfile> <keyfile>
```

- `remoteski` is the SKI of the remote device or service you want to connect to
- `certfile` is a local file containing the generated certificate in the first usage run
- `keyfile` is a local file containing the generated key in the first usage run

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.

## Roadmap - Spine specification implementation

### General request processing

- [x] Request / Handle acknowledgement
- [x] Use maximum response delay to timeout requests
- [X] Send error result when processing failed
- [X] Sending heartbeats

### Node Management

- [ ] Detailed Discovery
  - [X] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data
  - [ ] Notify subscribers
- [ ] Destination List
  - [ ] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data
  - [ ] Notify subscribers
- [ ] Binding
  - [X] Add Binding
  - [x] Delete Binding
  - [ ] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data
- [ ] Subscription
  - [X] Add subscription
  - [x] Delete subscription
  - [ ] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data
- [ ] Use Case Discovery
  - [X] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data

### General feature implementation

- [X] Request and process full data
- [X] Response full data
- [ ] Request partial data
- [ ] Process partial data
  - [X] Delete Selectors
  - [X] Update Selectors
  - [ ] Elements
- [ ] Response partial data
- [ ] Process write call
- [ ] Request subscription
- [X] Notify subscribers
- [ ] Handle incoming error results
- [X] Handle incoming success results

### Feature with partial data support

- `ElectricalConnection`
- `Measurement`

## Interfaces

### Verbose logging

Use `SetLogger` on `Service` to set the logger which needs to conform to the `logging.Logging` interface.

Example:

```go
h.myService = service.NewEEBUSService(serviceDescription, h)
h.myService.SetLogging(h)
```
