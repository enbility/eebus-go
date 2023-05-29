# eebus-go

[![Build Status](https://github.com/enbility/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)](https://github.com/enbility/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/enbility/eebus-go)
[![Coverage Status](https://coveralls.io/repos/github/enbility/eebus-go/badge.svg?branch=dev)](https://coveralls.io/github/enbility/eebus-go?branch=dev)
[![Go report](https://goreportcard.com/badge/github.com/enbility/eebus-go)](https://goreportcard.com/report/github.com/enbility/eebus-go)

This library provides a complete foundation for implementing [EEBUS](https://eebus.org) use cases. The use cases define various functional scenarios for different device categories, e.g. energy management systems, charging stations, heat pumps, and more.

## Introduction

The supported functionality contains:

- Support for SHIP 1.0.1
- Support for big parts of SPINE 1.1.1
- (De-)serialization for EEBUS specific JSON format requirements
- Certificate handling
- mDNS Support, incl. avahi support (recommended)
- Connection (websocket) handling, including reconnection and double connections
- Support for handling pairing of devices

Basic understanding of the EEBUS concepts SHIP and SPINE to use this library is required. Please check the corresponding specifications on the [EEBUS downloads website](https://www.eebus.org/media-downloads/).

An open source SDK written in go providing the foundation to use EEBUS in your projects. Contains support for SHIP and SPINE communication.

## Usage

The included small demo applications do not implement any usecases and thus will end the connection once it reached exchanging usecase information.

Services with implemented use cases will be implemented in different repositories and are also early work in progress:

- [HEMS](https://github.com/enbility/cemd)

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

## SHIP implementation notes

- Double connection handling is not implemented according to SHIP 12.2.2. Instead the connection initiated by the higher SKI will be kept. Much simpler and always works
- PIN Verification is _NOT_ supported other than SHIP 13.4.4.3.5.1 _"none"_ PIN state is supported!
- Access Methods SHIP 13.4.6 only supports the most basic scenario and only works after PIN verification state is completed.
- Supported registration mechanisms (SHIP 5):
  - auto accept (without any interaction mechanism!)
  - user verification

This approach has been tested with:

- Elli Charger Connect
- Porsche Mobile Charger Connect
- SMA Home Energy Manager 2.0

## Roadmap - Spine specification implementation

### General request processing

- [X] Request and process full data
- [ ] Request partial data
  - [ ] Delete Selectors
  - [ ] Update Selectors
  - [ ] Elements
- [ ] Send
  - [X] Full data
  - [ ] Partial data
- [X] Process partial data
  - [X] Delete Selectors
  - [X] Update Selectors
  - [X] Elements
- [ ] Request types
  - [X] Read
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
      - [ ] Partial Delete
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
      - [X] Partial Delete
  - [X] Reply
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
      - [ ] Partial Delete
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
      - [X] Partial Delete
  - [X] Notify
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
      - [ ] Partial Delete
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
      - [X] Partial Delete
  - [X] Write
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
      - [ ] Partial Delete
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
      - [X] Partial Delete
- [X] Result message handling
  - [X] Handle incoming error results
  - [X] Handle incoming success results
  - [X] Respond with error result when processing failed
- [X] Acknowledgement support
  - [X] Request
  - [X] Respond
- [x] Use maximum response delay to timeout requests

### Node Management

- [ ] Detailed Discovery
  - [ ] Read Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
  - [ ] Reply Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
  - [ ] Notify Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
- [ ] Destination List
  - [ ] Request and process full data
  - [X] Response full data
  - [ ] Request and process partial data
  - [ ] Response partial data
  - [ ] Notify subscribers
- [ ] Binding
  - [ ] Send Requests
    - [X] Add Binding
    - [ ] Delete Binding
  - [X] Receive Requests
    - [X] Add Binding
    - [X] Delete Binding
- [ ] Subscription
  - [ ] Send Requests
    - [X] Add Subscription
    - [ ] Delete Subscription
  - [X] Receive Requests
    - [X] Add Subscription
    - [X] Delete Subscription
  - [X] Notify subscribers
- [ ] Use Case Discovery
  - [ ] Read Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
  - [ ] Reply Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request
  - [ ] Notify Messages
    - [ ] Send
      - [X] Full Request
      - [ ] Partial Request
    - [X] Receive
      - [X] Full Request
      - [X] Partial Request

### General feature implementation

- [ ] Hearbeat Support
  - [X] Send hearbeats
  - [ ] Receive hearbeats

### Partial, selector, elements support

All list types do support processing of incoming partial messages, including selectors and elements. Sending partial messages is possible but there is no special support implemented right now.

## Interfaces

### Verbose logging

Use `SetLogger` on `Service` to set the logger which needs to conform to the `logging.Logging` interface.

Example:

```go
configuration = service.NewConfiguration(...)
h.myService = service.NewEEBUSService(configuration, h)
h.myService.SetLogging(h)
```
